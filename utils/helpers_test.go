package utils

import (
	"testing"
)

func TestGenerateRandomString_Length(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 6", 6},
		{"length 8", 8},
		{"length 10", 10},
		{"length 1", 1},
		{"length 0", 0},
		{"length 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateRandomString(tt.length)
			if len(result) != tt.length {
				t.Errorf("GenerateRandomString(%d) length = %d, want %d", tt.length, len(result), tt.length)
			}
		})
	}
}

func TestGenerateRandomString_CharacterSet(t *testing.T) {
	result := GenerateRandomString(100)
	for _, c := range result {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Errorf("GenerateRandomString() contains invalid character: %c", c)
		}
	}
}

func TestGenerateRandomString_Uniqueness(t *testing.T) {
	a := GenerateRandomString(16)
	b := GenerateRandomString(16)
	if a == b {
		t.Errorf("GenerateRandomString() produced same result twice: %s", a)
	}
}
