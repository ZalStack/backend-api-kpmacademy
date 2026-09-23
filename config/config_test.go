package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal string
		envVal     string
		want       string
	}{
		{
			name:       "env var set",
			key:        "TEST_GET_ENV_1",
			defaultVal: "default",
			envVal:     "custom",
			want:       "custom",
		},
		{
			name:       "env var not set",
			key:        "TEST_GET_ENV_2",
			defaultVal: "default",
			envVal:     "",
			want:       "default",
		},
		{
			name:       "env var with spaces",
			key:        "TEST_GET_ENV_3",
			defaultVal: "default",
			envVal:     "  trimmed  ",
			want:       "trimmed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.key, tt.envVal)
			}
			got := getEnv(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal bool
		envVal     string
		want       bool
	}{
		{"true string", "TEST_BOOL_1", false, "true", true},
		{"false string", "TEST_BOOL_2", true, "false", false},
		{"empty string", "TEST_BOOL_3", true, "", true},
		{"invalid string", "TEST_BOOL_4", true, "invalid", true},
		{"1", "TEST_BOOL_5", false, "1", true},
		{"0", "TEST_BOOL_6", true, "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.key, tt.envVal)
			}
			got := getEnvBool(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvBool(%q, %v) = %v, want %v", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal int
		envVal     string
		want       int
	}{
		{"valid int", "TEST_INT_1", 0, "42", 42},
		{"empty string", "TEST_INT_2", 10, "", 10},
		{"invalid int", "TEST_INT_3", 5, "abc", 5},
		{"negative int", "TEST_INT_4", 0, "-10", -10},
		{"zero", "TEST_INT_5", 100, "0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.key, tt.envVal)
			}
			got := getEnvInt(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvInt(%q, %d) = %d, want %d", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetEnvInt64(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		defaultVal int64
		envVal     string
		want       int64
	}{
		{"valid int64", "TEST_INT64_1", 0, "1234567890123", 1234567890123},
		{"empty string", "TEST_INT64_2", 100, "", 100},
		{"invalid int64", "TEST_INT64_3", 50, "abc", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				t.Setenv(tt.key, tt.envVal)
			}
			got := getEnvInt64(tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getEnvInt64(%q, %d) = %d, want %d", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	// Ensure JWT_SECRET is not set
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_REFRESH_SECRET")

	// Load should call log.Fatal when JWT_SECRET is missing
	// We can't easily test log.Fatal, but we can test that getEnv works
	// This is more of an integration test
	t.Skip("Skipping Load() test - requires log.Fatal handling")
}
