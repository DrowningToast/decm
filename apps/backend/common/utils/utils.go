package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// DerefOrEmpty dereferences a string pointer, returning "" if nil.
func DerefOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// NormalizeToLower lowercases a non-nil, non-empty string pointer.
// Returns nil when the pointer is nil or points to an empty string.
func NormalizeToLower(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}
	lower := strings.ToLower(*s)
	return &lower
}

// GenerateSecureRandomString creates a secure random string for session IDs
func GenerateSecureRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Fallback to less secure method if crypto/rand fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
