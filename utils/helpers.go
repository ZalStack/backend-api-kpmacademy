package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateRandomString(length int) string {
	b := make([]byte, (length+1)/2)
	_, _ = rand.Read(b) //nolint:errcheck // crypto/rand.Read never returns error on valid slice
	return strings.ToUpper(hex.EncodeToString(b))[:length]
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
