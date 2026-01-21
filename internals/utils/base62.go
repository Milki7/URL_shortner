package utils

import (
	"errors"
	"math"
	"strings"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Encode converts a base-10 ID to a Base62 string
func Encode(number uint) string {
	if number == 0 {
		return string(alphabet[0])
	}

	var result strings.Builder
	base := uint(len(alphabet))

	for number > 0 {
		remainder := number % base
		result.WriteByte(alphabet[remainder])
		number = number / base
	}

	// Reverse the string
	s := result.String()
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Decode converts a Base62 string back to a base-10 ID
func Decode(encoded string) (uint, error) {
	var number uint
	base := uint(len(alphabet))

	for i, char := range encoded {
		pos := strings.IndexRune(alphabet, char)
		if pos == -1 {
			return 0, errors.New("invalid character in base62 string")
		}
		exponent := len(encoded) - 1 - i
		number += uint(pos) * uint(math.Pow(float64(base), float64(exponent)))
	}
	return number, nil
}
