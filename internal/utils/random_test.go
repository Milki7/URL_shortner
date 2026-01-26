package utils

import (
	"testing"
)

// TestGenerateRandomCodeLength ensures the output matches the requested length
func TestGenerateRandomCodeLength(t *testing.T) {
	lengths := []int{4, 6, 8, 10}

	for _, want := range lengths {
		got := GenerateRandomCode(want)
		if len(got) != want {
			t.Errorf("GenerateRandomCode(%d) = length %d; want %d", want, len(got), want)
		}
	}
}

// TestGenerateRandomCodeUniqueness checks if 1000 generated codes are unique
func TestGenerateRandomCodeUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		code := GenerateRandomCode(6)
		if seen[code] {
			t.Errorf("Collision detected at iteration %d: code %s already exists", i, code)
		}
		seen[code] = true
	}
}
