package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// AnonymizeIP anonymizes an IP address for PDPA/GDPR compliance.
// For IPv4: masks the last octet (192.168.1.xxx)
// For IPv6: keeps only the first 3 segments
//
// This allows threat analysis while protecting user privacy.
func AnonymizeIP(ip string) string {
	if ip == "" {
		return "unknown"
	}

	// Handle IPv4
	if strings.Contains(ip, ".") && !strings.Contains(ip, ":") {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			// Mask last octet: 192.168.1.xxx
			return parts[0] + "." + parts[1] + "." + parts[2] + ".xxx"
		}
	}

	// Handle IPv6
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			// Keep first 3 segments: 2001:0db8:85a3::xxxx
			return parts[0] + ":" + parts[1] + ":" + parts[2] + "::xxxx"
		}
	}

	// Fallback: hash the IP
	return HashIP(ip)
}

// HashIP creates a one-way hash of an IP address.
// Useful for tracking unique attackers without storing raw IPs.
//
// Note: Not reversible, but same IP always produces same hash.
func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return "hash_" + hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter output
}

// AnonymizeUserAgent removes version numbers and detailed system info
// while keeping browser/OS type for security analysis.
func AnonymizeUserAgent(ua string) string {
	if ua == "" {
		return "unknown"
	}

	// Extract basic info: browser type and OS
	// Example: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124"
	//       -> "Windows/Chrome"

	ua = strings.ToLower(ua)
	browser := "unknown"
	os := "unknown"

	// Detect OS (order matters: check specific OS before generic ones)
	if strings.Contains(ua, "android") {
		os = "android"
	} else if strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		os = "ios"
	} else if strings.Contains(ua, "windows") {
		os = "windows"
	} else if strings.Contains(ua, "mac") || strings.Contains(ua, "darwin") {
		os = "macos"
	} else if strings.Contains(ua, "linux") {
		os = "linux"
	}

	// Detect Browser
	if strings.Contains(ua, "chrome") {
		browser = "chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "firefox"
	} else if strings.Contains(ua, "safari") {
		browser = "safari"
	} else if strings.Contains(ua, "edge") {
		browser = "edge"
	} else if strings.Contains(ua, "curl") {
		browser = "curl"
	} else if strings.Contains(ua, "python") {
		browser = "python"
	} else if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") {
		browser = "bot"
	}

	return os + "/" + browser
}
