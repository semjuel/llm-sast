package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(s string) string {
	// Create a new SHA256 hash.
	hash := sha256.New()
	// Write the string as bytes.
	hash.Write([]byte(s))
	// Get the final hash sum as a byte slice.
	hashBytes := hash.Sum(nil)
	// Convert the hash bytes to a hexadecimal string.
	return hex.EncodeToString(hashBytes)
}
