package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateRandomString(length int) string {
	b := make([]byte, (length+1)/2)
	rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))[:length]
}
