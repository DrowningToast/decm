package logger

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

// AnonymizeIP anonymizes an IP address for PDPA/GDPR compliance.
// For IPv4: masks the last octet (192.168.1.xxx)
// For IPv6: masks the last 80 bits, keeping /48 prefix (2001:0db8:85a3::xxxx)
//
// This allows threat analysis while protecting user privacy.
// Uses proper IP parsing to handle compressed IPv6 addresses correctly.
func AnonymizeIP(ip string) string {
	if ip == "" {
		return "unknown"
	}

	// Parse IP address properly to handle both IPv4 and IPv6 (including compressed formats)
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// Invalid IP format - log warning as this may indicate proxy misconfiguration
		slog.Warn("Malformed IP address detected, falling back to hash",
			slog.String("ip", ip),
			slog.String("hint", "Check proxy configuration or IP extraction logic"))
		return HashIP(ip)
	}

	// Check if IPv4
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// Mask last octet: 192.168.1.xxx
		return fmt.Sprintf("%d.%d.%d.xxx", ipv4[0], ipv4[1], ipv4[2])
	}

	// Handle IPv6 - keep first 48 bits (3 hextets), mask the rest
	// Properly handles compressed formats like 2001::8a2e:370:7334
	if ipv6 := parsedIP.To16(); ipv6 != nil {
		// Create anonymized IPv6 with first 6 bytes (48 bits), rest zeros
		anonymized := make(net.IP, 16)
		copy(anonymized, ipv6[0:6]) // Copy first 48 bits
		// Rest is already zeros
		// Format as xxxx:xxxx:xxxx::xxxx for consistency
		return anonymized.String() + ":xxxx"
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
	// Order matters: Check most specific patterns first to handle overlapping cases.
	// - Edge user agents contain: "Edg/" + "Chrome" + "Safari"
	// - Chrome user agents contain: "Chrome" + "Safari"
	// - Safari user agents contain: "Safari" (but not "Chrome")
	if strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/") {
		// Edge (modern versions use "Edg/", older used "Edge/")
		browser = "edge"
	} else if strings.Contains(ua, "chrome") {
		// Chrome (but not Edge, since we checked Edge first)
		browser = "chrome"
	} else if strings.Contains(ua, "firefox") {
		browser = "firefox"
	} else if strings.Contains(ua, "safari") {
		// Safari (but not Chrome/Edge, since we checked them first)
		browser = "safari"
	} else if strings.Contains(ua, "curl") {
		browser = "curl"
	} else if strings.Contains(ua, "python") {
		browser = "python"
	} else if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") {
		browser = "bot"
	}

	return os + "/" + browser
}
