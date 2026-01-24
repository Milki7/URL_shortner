package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomCode creates a secure random string of a given length
func GenerateRandomCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		// crypto/rand is much safer for unique IDs than math/rand
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
