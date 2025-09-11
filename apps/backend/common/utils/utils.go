package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
)

// generateSecureRandomString creates a secure random string for session IDs
func GenerateSecureRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback to less secure method if crypto/rand fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
