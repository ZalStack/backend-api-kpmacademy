package utils

import (
	"testing"
	"time"

	"backend-api-kpmacademy/config"

	"github.com/google/uuid"
)

const testEmail = "test@example.com"

func TestGenerateTokenPair(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	userID := uuid.New()
	email := testEmail
	role := "user"

	tokens, err := GenerateTokenPair(userID, email, role, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("GenerateTokenPair() returned empty access token")
	}
	if tokens.RefreshToken == "" {
		t.Error("GenerateTokenPair() returned empty refresh token")
	}
	if tokens.ExpiresAt == 0 {
		t.Error("GenerateTokenPair() returned zero ExpiresAt")
	}
	if tokens.ExpiresAt <= time.Now().Unix() {
		t.Error("GenerateTokenPair() ExpiresAt is in the past")
	}
}

func TestValidateAccessToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	userID := uuid.New()
	email := testEmail
	role := "user"

	tokens, err := GenerateTokenPair(userID, email, role, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	claims, err := ValidateAccessToken(tokens.AccessToken, cfg)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateAccessToken() UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("ValidateAccessToken() Email = %v, want %v", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("ValidateAccessToken() Role = %v, want %v", claims.Role, role)
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	_, err := ValidateAccessToken("invalid-token", cfg)
	if err == nil {
		t.Error("ValidateAccessToken() expected error for invalid token")
	}
}

func TestValidateRefreshToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	userID := uuid.New()
	email := testEmail
	role := "admin"

	tokens, err := GenerateTokenPair(userID, email, role, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	claims, err := ValidateRefreshToken(tokens.RefreshToken, cfg)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("ValidateRefreshToken() UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("ValidateRefreshToken() Email = %v, want %v", claims.Email, email)
	}
	if claims.Role != role {
		t.Errorf("ValidateRefreshToken() Role = %v, want %v", claims.Role, role)
	}
}

func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	_, err := ValidateRefreshToken("invalid-token", cfg)
	if err == nil {
		t.Error("ValidateRefreshToken() expected error for invalid token")
	}
}

func TestValidateRefreshToken_WrongSecret(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:            "test-secret-key-for-jwt-signing",
		RefreshSecret:     "test-refresh-secret-key-for-jwt-signing",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	wrongCfg := &config.JWTConfig{
		Secret:            "wrong-secret",
		RefreshSecret:     "wrong-refresh-secret",
		Expiration:        3600,
		RefreshExpiration: 604800,
	}

	userID := uuid.New()
	tokens, err := GenerateTokenPair(userID, testEmail, "user", cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	_, err = ValidateRefreshToken(tokens.RefreshToken, wrongCfg)
	if err == nil {
		t.Error("ValidateRefreshToken() expected error with wrong secret")
	}
}
