package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func GenerateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%s", now.Format("20060102150405"), GenerateRandomString(6))
}

func GenerateVoucherCode() string {
	return strings.ToUpper(GenerateRandomString(8))
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

func GenerateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%s-%s", now.Format("20060102"), GenerateRandomString(8))
}
