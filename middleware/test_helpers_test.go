package middleware

import (
	"testing"

	"backend-api-kpmacademy/config"
	"backend-api-kpmacademy/utils"

	"github.com/google/uuid"
)

func generateTestToken(t *testing.T, cfg *config.Config) string {
	t.Helper()
	userID := uuid.New()
	tokens, err := utils.GenerateTokenPair(userID, "test@example.com", "user", &cfg.JWT)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}
	return tokens.AccessToken
}
