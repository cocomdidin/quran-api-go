package repository

import (
	"testing"
)

// TestOffsetFormula memverifikasi formula offset pagination
func TestOffsetFormula(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"page 1 limit 20", 1, 20, 0},
		{"page 2 limit 20", 2, 20, 20},
		{"page 3 limit 20", 3, 20, 40},
		{"page 5 limit 10", 5, 10, 40},
		{"page 10 limit 50", 10, 50, 450},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Formula yang benar: (page - 1) * limit
			offset := (tt.page - 1) * tt.limit
			if offset != tt.expected {
				t.Errorf("offset = %d, want %d", offset, tt.expected)
			}
		})
	}
}
