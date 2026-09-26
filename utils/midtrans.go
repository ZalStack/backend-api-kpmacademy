package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type SnapTransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount int64   `json:"gross_amount"`
}

type SnapCustomerDetails struct {
	FirstName string `json:"first_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type SnapExpiry struct {
	Duration int    `json:"duration"`
	Unit     string `json:"unit"`
}

type SnapRequest struct {
	TransactionDetails SnapTransactionDetails `json:"transaction_details"`
	CustomerDetails    SnapCustomerDetails    `json:"customer_details,omitempty"`
	Expiry             SnapExpiry             `json:"expiry,omitempty"`
}

type SnapResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

type MidtransStatusResponse struct {
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	VANumbers         []struct {
		Bank     string `json:"bank"`
		VANumber string `json:"va_number"`
	} `json:"va_numbers,omitempty"`
	PermataVANumber string `json:"permata_va_number,omitempty"`
	BillKey         string `json:"bill_key,omitempty"`
	BillerCode      string `json:"biller_code,omitempty"`
	QRString        string `json:"qr_string,omitempty"`
}

func getMidtransSnapURL() string {
	if os.Getenv("MIDTRANS_IS_PRODUCTION") == "true" {
		return "https://app.midtrans.com/snap/v1/transactions"
	}
	return "https://app.sandbox.midtrans.com/snap/v1/transactions"
}

func getMidtransApiBaseURL() string {
	if os.Getenv("MIDTRANS_IS_PRODUCTION") == "true" {
		return "https://api.midtrans.com/v2"
	}
	return "https://api.sandbox.midtrans.com/v2"
}

// CreateSnapTransaction creates a new transaction token and redirect URL in Midtrans
func CreateSnapTransaction(orderNumber string, grossAmount float64, customerName, customerEmail, customerPhone string) (string, string, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	reqBody := SnapRequest{
		TransactionDetails: SnapTransactionDetails{
			OrderID:     orderNumber,
			GrossAmount: int64(grossAmount),
		},
		CustomerDetails: SnapCustomerDetails{
			FirstName: customerName,
			Email:     customerEmail,
			Phone:     customerPhone,
		},
		Expiry: SnapExpiry{
			Duration: 60,
			Unit:     "minutes",
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", getMidtransSnapURL(), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var snapResp SnapResponse
	if err := json.Unmarshal(body, &snapResp); err != nil {
		return "", "", err
	}

	if snapResp.Token == "" {
		return "", "", fmt.Errorf("midtrans error: %s", string(body))
	}

	return snapResp.Token, snapResp.RedirectURL, nil
}

// CheckMidtransTransactionStatus queries Midtrans directly for real-time status of an order
func CheckMidtransTransactionStatus(orderNumber string) (*MidtransStatusResponse, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	url := fmt.Sprintf("%s/%s/status", getMidtransApiBaseURL(), orderNumber)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(serverKey + ":"))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var statusResp MidtransStatusResponse
	if err := json.Unmarshal(body, &statusResp); err != nil {
		return nil, err
	}

	return &statusResp, nil
}
