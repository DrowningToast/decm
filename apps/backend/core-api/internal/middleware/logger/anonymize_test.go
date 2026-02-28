package logger

import (
	"strings"
	"testing"
)

func TestAnonymizeIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		wantType string // "masked", "hashed", or "unknown"
	}{
		{
			name:     "IPv4 - masks last octet",
			ip:       "192.168.1.100",
			wantType: "masked",
		},
		{
			name:     "IPv4 - different IP",
			ip:       "10.0.0.50",
			wantType: "masked",
		},
		{
			name:     "IPv6 - keeps first 3 segments",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			wantType: "masked",
		},
		{
			name:     "Empty IP",
			ip:       "",
			wantType: "unknown",
		},
		{
			name:     "Localhost IPv4",
			ip:       "127.0.0.1",
			wantType: "masked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnonymizeIP(tt.ip)

			switch tt.wantType {
			case "masked":
				if tt.ip != "" && strings.Contains(tt.ip, ".") {
					// IPv4: should end with .xxx
					if !strings.HasSuffix(got, ".xxx") {
						t.Errorf("AnonymizeIP(%q) = %q, expected IPv4 to end with .xxx", tt.ip, got)
					}
					// Should preserve first 3 octets
					parts := strings.Split(tt.ip, ".")
					if len(parts) == 4 {
						expectedPrefix := parts[0] + "." + parts[1] + "." + parts[2]
						if !strings.HasPrefix(got, expectedPrefix) {
							t.Errorf("AnonymizeIP(%q) = %q, expected prefix %q", tt.ip, got, expectedPrefix)
						}
					}
				} else if tt.ip != "" && strings.Contains(tt.ip, ":") {
					// IPv6: should end with ::xxxx
					if !strings.HasSuffix(got, "::xxxx") {
						t.Errorf("AnonymizeIP(%q) = %q, expected IPv6 to end with ::xxxx", tt.ip, got)
					}
				}
			case "unknown":
				if got != "unknown" {
					t.Errorf("AnonymizeIP(%q) = %q, want 'unknown'", tt.ip, got)
				}
			case "hashed":
				if !strings.HasPrefix(got, "hash_") {
					t.Errorf("AnonymizeIP(%q) = %q, expected hash prefix", tt.ip, got)
				}
			}
		})
	}
}

func TestAnonymizeIPConsistency(t *testing.T) {
	// Same IP should always produce same result
	ip := "192.168.1.100"
	result1 := AnonymizeIP(ip)
	result2 := AnonymizeIP(ip)

	if result1 != result2 {
		t.Errorf("AnonymizeIP not consistent: %q vs %q", result1, result2)
	}
}

func TestAnonymizeIPDifferentLastOctet(t *testing.T) {
	// Different IPs with same first 3 octets should produce same anonymized result
	ip1 := "192.168.1.100"
	ip2 := "192.168.1.200"

	result1 := AnonymizeIP(ip1)
	result2 := AnonymizeIP(ip2)

	if result1 != result2 {
		t.Errorf("Different last octets should produce same result: %q vs %q", result1, result2)
	}

	expected := "192.168.1.xxx"
	if result1 != expected {
		t.Errorf("AnonymizeIP(%q) = %q, want %q", ip1, result1, expected)
	}
}

func TestHashIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
	}{
		{"IPv4", "192.168.1.100"},
		{"IPv6", "2001:0db8:85a3::8a2e:0370:7334"},
		{"Localhost", "127.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashIP(tt.ip)

			// Should start with hash_
			if !strings.HasPrefix(hash, "hash_") {
				t.Errorf("HashIP(%q) = %q, should start with 'hash_'", tt.ip, hash)
			}

			// Should be consistent
			hash2 := HashIP(tt.ip)
			if hash != hash2 {
				t.Errorf("HashIP not consistent: %q vs %q", hash, hash2)
			}

			// Should be different for different IPs
			differentIP := "10.0.0.1"
			differentHash := HashIP(differentIP)
			if tt.ip != differentIP && hash == differentHash {
				t.Errorf("HashIP produced same hash for different IPs")
			}
		})
	}
}

func TestAnonymizeIP_MalformedAddresses(t *testing.T) {
	// Test that malformed IPs are handled gracefully and hashed
	tests := []struct {
		name string
		ip   string
	}{
		{"Invalid format", "not-an-ip"},
		{"Partial IPv4", "192.168.1"},
		{"Too many octets", "192.168.1.1.1"},
		{"Invalid characters", "192.168.1.abc"},
		{"Port included", "192.168.1.1:8080"},
		{"URL instead of IP", "http://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnonymizeIP(tt.ip)

			// Should fall back to hashing
			if !strings.HasPrefix(result, "hash_") {
				t.Errorf("AnonymizeIP(%q) = %q, expected hash prefix for malformed IP", tt.ip, result)
			}

			// Should be consistent
			result2 := AnonymizeIP(tt.ip)
			if result != result2 {
				t.Errorf("AnonymizeIP(%q) not consistent: %q vs %q", tt.ip, result, result2)
			}
		})
	}
}

func TestAnonymizeUserAgent(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want string
	}{
		{
			name: "Chrome on Windows",
			ua:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			want: "windows/chrome",
		},
		{
			name: "Firefox on macOS",
			ua:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
			want: "macos/firefox",
		},
		{
			name: "Safari on iOS",
			ua:   "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
			want: "ios/safari",
		},
		{
			name: "Chrome on Android",
			ua:   "Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
			want: "android/chrome",
		},
		{
			name: "curl",
			ua:   "curl/7.68.0",
			want: "unknown/curl",
		},
		{
			name: "Python requests",
			ua:   "python-requests/2.25.1",
			want: "unknown/python",
		},
		{
			name: "Bot",
			ua:   "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			want: "unknown/bot",
		},
		{
			name: "Empty user agent",
			ua:   "",
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnonymizeUserAgent(tt.ua)
			if got != tt.want {
				t.Errorf("AnonymizeUserAgent(%q) = %q, want %q", tt.ua, got, tt.want)
			}
		})
	}
}

func TestAnonymizeUserAgentRemovesVersions(t *testing.T) {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124"
	result := AnonymizeUserAgent(ua)

	// Should not contain version numbers
	if strings.Contains(result, "91.0") || strings.Contains(result, "10.0") {
		t.Errorf("AnonymizeUserAgent should remove version numbers, got: %q", result)
	}

	// Should be short and generic
	if len(result) > 20 {
		t.Errorf("AnonymizeUserAgent result too long: %q", result)
	}
}

func TestAnonymizeUserAgent_BrowserDetectionOrder(t *testing.T) {
	// Test that browser detection handles overlapping patterns correctly
	// Edge contains "Edg/" + "Chrome" + "Safari"
	// Chrome contains "Chrome" + "Safari"
	// Safari contains "Safari" only
	tests := []struct {
		name           string
		ua             string
		expectedBrowser string
		description    string
	}{
		{
			name:           "Edge (modern) - should detect as Edge, not Chrome",
			ua:             "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.59",
			expectedBrowser: "edge",
			description:    "Edge user agent contains 'Edg/', 'Chrome', and 'Safari'",
		},
		{
			name:           "Edge (older) - should detect as Edge, not Chrome",
			ua:             "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18362",
			expectedBrowser: "edge",
			description:    "Older Edge user agent contains 'Edge/'",
		},
		{
			name:           "Chrome - should detect as Chrome, not Safari",
			ua:             "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
			expectedBrowser: "chrome",
			description:    "Chrome user agent contains 'Chrome' and 'Safari'",
		},
		{
			name:           "Safari (macOS) - should detect as Safari",
			ua:             "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
			expectedBrowser: "safari",
			description:    "Safari user agent contains 'Safari' but not 'Chrome'",
		},
		{
			name:           "Safari (iOS) - should detect as Safari",
			ua:             "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
			expectedBrowser: "safari",
			description:    "iOS Safari contains 'Safari' but not 'Chrome'",
		},
		{
			name:           "Firefox - should detect as Firefox",
			ua:             "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
			expectedBrowser: "firefox",
			description:    "Firefox user agent is independent",
		},
		{
			name:           "Case insensitive - EDG uppercase",
			ua:             "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 EDG/91.0.864.59",
			expectedBrowser: "edge",
			description:    "Should handle case-insensitive detection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AnonymizeUserAgent(tt.ua)

			// Result format is "os/browser"
			parts := strings.Split(result, "/")
			if len(parts) != 2 {
				t.Fatalf("Expected result format 'os/browser', got: %q", result)
			}

			gotBrowser := parts[1]
			if gotBrowser != tt.expectedBrowser {
				t.Errorf("AnonymizeUserAgent(%q)\n  got browser: %q\n  want: %q\n  context: %s",
					tt.ua, gotBrowser, tt.expectedBrowser, tt.description)
			}
		})
	}
}
