package utils

import (
	"testing"
)

func TestValidateStrongPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "valid strong password",
			password: "StrongPass123!",
			want:     true,
		},
		{
			name:     "missing uppercase",
			password: "strongpass123!",
			want:     false,
		},
		{
			name:     "missing lowercase",
			password: "STRONGPASS123!",
			want:     false,
		},
		{
			name:     "missing number",
			password: "StrongPass!",
			want:     false,
		},
		{
			name:     "missing special character",
			password: "StrongPass123",
			want:     false,
		},
		{
			name:     "too short",
			password: "Str1!",
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Password string `validate:"required,strong_password"`
			}
			s := TestStruct{Password: tt.password}
			errs := ValidateStruct(s)
			got := len(errs) == 0
			if got != tt.want {
				t.Errorf("validateStrongPassword() = %v, want %v for password: %s", got, tt.want, tt.password)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required,min=2"`
		Email string `validate:"required,email"`
	}

	tests := []struct {
		name       string
		input      TestStruct
		wantErrors bool
	}{
		{
			name:       "valid input",
			input:      TestStruct{Name: "John", Email: "john@example.com"},
			wantErrors: false,
		},
		{
			name:       "missing name",
			input:      TestStruct{Name: "", Email: "john@example.com"},
			wantErrors: true,
		},
		{
			name:       "invalid email",
			input:      TestStruct{Name: "John", Email: "invalid-email"},
			wantErrors: true,
		},
		{
			name:       "name too short",
			input:      TestStruct{Name: "J", Email: "john@example.com"},
			wantErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateStruct(tt.input)
			if (len(errs) > 0) != tt.wantErrors {
				t.Errorf("ValidateStruct() errors = %v, wantErrors %v", errs, tt.wantErrors)
			}
		})
	}
}
