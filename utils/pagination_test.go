package utils

import (
	"testing"
)

func TestCalcLastPage(t *testing.T) {
	tests := []struct {
		name    string
		total   int64
		perPage int
		want    int
	}{
		{
			name:    "10 items, 10 per page",
			total:   10,
			perPage: 10,
			want:    1,
		},
		{
			name:    "15 items, 10 per page",
			total:   15,
			perPage: 10,
			want:    2,
		},
		{
			name:    "0 items",
			total:   0,
			perPage: 10,
			want:    1,
		},
		{
			name:    "100 items, 10 per page",
			total:   100,
			perPage: 10,
			want:    10,
		},
		{
			name:    "1 item, 10 per page",
			total:   1,
			perPage: 10,
			want:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcLastPage(tt.total, tt.perPage); got != tt.want {
				t.Errorf("CalcLastPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		allowed []string
		want    bool
	}{
		{
			name:    "status in allowed list",
			status:  "pending",
			allowed: []string{"pending", "paid", "canceled"},
			want:    true,
		},
		{
			name:    "status not in allowed list",
			status:  "invalid",
			allowed: []string{"pending", "paid", "canceled"},
			want:    false,
		},
		{
			name:    "empty allowed list",
			status:  "pending",
			allowed: []string{},
			want:    false,
		},
		{
			name:    "empty status",
			status:  "",
			allowed: []string{"pending", "paid"},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStatus(tt.status, tt.allowed); got != tt.want {
				t.Errorf("ValidateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
