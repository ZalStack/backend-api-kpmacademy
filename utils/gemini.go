package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiSystemInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiRequest struct {
	SystemInstruction *GeminiSystemInstruction `json:"systemInstruction,omitempty"`
	Contents          []GeminiContent          `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

const defaultKPMPersona = `Anda adalah AI Master Tutor di KPM Academy (Klinik Pendidikan MIPA).
Keahlian Anda meliputi Matematika Nalaria Realistik (MNR), Sains Terpadu (Fisika, Kimia, Biologi), dan Penalaran Logika.
Pedoman Menjawab:
1. Bersikaplah ramah, antusias, memotivasi, dan gunakan gaya bahasa Indonesia yang mudah dipahami.
2. Jangan langsung berikan jawaban mentah tanpa penjelasan; jelaskan konsep dasar dan langkah bernalar (step-by-step) ala metode KPM.
3. Gunakan formatting rapi (bullet points, numbering, dan emoji pendukung).
4. Di akhir jawaban, berikan dorongan atau ajakan interaktif untuk bertanya lebih lanjut.`

// CallGeminiAI sends prompt to Google Gemini API
func CallGeminiAI(userPrompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-3.5-flash"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	reqBody := GeminiRequest{
		SystemInstruction: &GeminiSystemInstruction{
			Parts: []GeminiPart{
				{Text: defaultKPMPersona},
			},
		},
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: userPrompt},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	if geminiResp.Error != nil && geminiResp.Error.Message != "" {
		return "", fmt.Errorf("gemini api error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("empty response from gemini")
}
