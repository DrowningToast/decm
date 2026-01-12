package security

import (
	"testing"
)

func TestSuspiciousPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		pattern  SuspiciousPattern
		testPath string
		want     bool
	}{
		// PREFIX matching tests
		{
			name:     "Prefix match - .env file",
			pattern:  SuspiciousPattern{Pattern: "/.env", MatchType: MatchPrefix},
			testPath: "/.env",
			want:     true,
		},
		{
			name:     "Prefix match - .env.production",
			pattern:  SuspiciousPattern{Pattern: "/.env", MatchType: MatchPrefix},
			testPath: "/.env.production",
			want:     true,
		},
		{
			name:     "Prefix match - env.js (no dot prefix)",
			pattern:  SuspiciousPattern{Pattern: "/env.", MatchType: MatchPrefix},
			testPath: "/env.js",
			want:     true,
		},
		{
			name:     "Prefix match - env.development.js",
			pattern:  SuspiciousPattern{Pattern: "/env.", MatchType: MatchPrefix},
			testPath: "/env.development.js",
			want:     true,
		},
		{
			name:     "Prefix no match - /environment",
			pattern:  SuspiciousPattern{Pattern: "/.env", MatchType: MatchPrefix},
			testPath: "/environment",
			want:     false,
		},
		{
			name:     "Prefix match - .git/config",
			pattern:  SuspiciousPattern{Pattern: "/.git", MatchType: MatchPrefix},
			testPath: "/.git/config",
			want:     true,
		},
		{
			name:     "Prefix match - .idea/workspace.xml",
			pattern:  SuspiciousPattern{Pattern: "/.idea", MatchType: MatchPrefix},
			testPath: "/.idea/workspace.xml",
			want:     true,
		},
		{
			name:     "Contains match - .well-known/openid-configuration (specific, not all .well-known)",
			pattern:  SuspiciousPattern{Pattern: "/.well-known/openid-configuration", MatchType: MatchContains},
			testPath: "/.well-known/openid-configuration",
			want:     true,
		},

		// SUFFIX matching tests
		{
			name:     "Suffix match - file.php",
			pattern:  SuspiciousPattern{Pattern: ".php", MatchType: MatchSuffix},
			testPath: "/index.php",
			want:     true,
		},
		{
			name:     "Suffix match - shell.aspx",
			pattern:  SuspiciousPattern{Pattern: ".aspx", MatchType: MatchSuffix},
			testPath: "/admin.aspx",
			want:     true,
		},
		{
			name:     "Suffix match - service.asmx",
			pattern:  SuspiciousPattern{Pattern: ".asmx", MatchType: MatchSuffix},
			testPath: "/WebCloud.asmx",
			want:     true,
		},
		{
			name:     "Suffix match - page.jspa",
			pattern:  SuspiciousPattern{Pattern: ".jspa", MatchType: MatchSuffix},
			testPath: "/secure/Signup.jspa",
			want:     true,
		},
		{
			name:     "Suffix match - page.jsf",
			pattern:  SuspiciousPattern{Pattern: ".jsf", MatchType: MatchSuffix},
			testPath: "/jbpm-console/app/tasks.jsf",
			want:     true,
		},
		{
			name:     "Suffix match - ldlogon.dll",
			pattern:  SuspiciousPattern{Pattern: ".dll", MatchType: MatchSuffix},
			testPath: "/ldlogon.dll",
			want:     true,
		},
		{
			name:     "Suffix match - mobile.log",
			pattern:  SuspiciousPattern{Pattern: ".log", MatchType: MatchSuffix},
			testPath: "/log/mobile.log",
			want:     true,
		},
		{
			name:     "Suffix match - win.ini",
			pattern:  SuspiciousPattern{Pattern: ".ini", MatchType: MatchSuffix},
			testPath: "/c:/windows/win.ini",
			want:     true,
		},
		{
			name:     "Suffix no match - .php in middle",
			pattern:  SuspiciousPattern{Pattern: ".php", MatchType: MatchSuffix},
			testPath: "/index.php.bak",
			want:     false,
		},

		// EXACT matching tests
		{
			name:     "Exact match - /admin",
			pattern:  SuspiciousPattern{Pattern: "/admin", MatchType: MatchExact},
			testPath: "/admin",
			want:     true,
		},
		{
			name:     "Exact no match - /admin/users",
			pattern:  SuspiciousPattern{Pattern: "/admin", MatchType: MatchExact},
			testPath: "/admin/users",
			want:     false,
		},
		{
			name:     "Exact match - /jwks",
			pattern:  SuspiciousPattern{Pattern: "/jwks", MatchType: MatchExact},
			testPath: "/jwks",
			want:     true,
		},
		{
			name:     "Exact no match - /jwks.json",
			pattern:  SuspiciousPattern{Pattern: "/jwks", MatchType: MatchExact},
			testPath: "/jwks.json",
			want:     false,
		},

		// CONTAINS matching tests
		{
			name:     "Contains match - /jmx-console/",
			pattern:  SuspiciousPattern{Pattern: "/jmx-console", MatchType: MatchContains},
			testPath: "/jmx-console/",
			want:     true,
		},
		{
			name:     "Contains match - /jbpm-console/app/tasks.jsf",
			pattern:  SuspiciousPattern{Pattern: "/jbpm-console", MatchType: MatchContains},
			testPath: "/jbpm-console/app/tasks.jsf",
			want:     true,
		},
		{
			name:     "Contains match - /jupyter/lab",
			pattern:  SuspiciousPattern{Pattern: "/jupyter", MatchType: MatchContains},
			testPath: "/jupyter/lab",
			want:     true,
		},
		{
			name:     "Contains match - /hub/login",
			pattern:  SuspiciousPattern{Pattern: "/hub/login", MatchType: MatchContains},
			testPath: "/hub/login?next=",
			want:     true,
		},
		{
			name:     "Contains match - CMSPages",
			pattern:  SuspiciousPattern{Pattern: "/CMSPages/", MatchType: MatchContains},
			testPath: "/CMSPages/Staging/SyncServer.asmx",
			want:     true,
		},
		{
			name:     "Contains no match - different path",
			pattern:  SuspiciousPattern{Pattern: "/jmx-console", MatchType: MatchContains},
			testPath: "/api/users",
			want:     false,
		},

		// Case insensitivity tests
		{
			name:     "Case insensitive - .PHP uppercase",
			pattern:  SuspiciousPattern{Pattern: ".php", MatchType: MatchSuffix},
			testPath: "/INDEX.PHP",
			want:     true,
		},
		{
			name:     "Case insensitive - /ADMIN uppercase",
			pattern:  SuspiciousPattern{Pattern: "/admin", MatchType: MatchExact},
			testPath: "/ADMIN",
			want:     true,
		},
		{
			name:     "Case insensitive - /JMX-CONSOLE uppercase",
			pattern:  SuspiciousPattern{Pattern: "/jmx-console", MatchType: MatchContains},
			testPath: "/JMX-CONSOLE/",
			want:     true,
		},

		// Edge cases
		{
			name:     "Empty path",
			pattern:  SuspiciousPattern{Pattern: "/admin", MatchType: MatchExact},
			testPath: "",
			want:     false,
		},
		{
			name:     "Root path",
			pattern:  SuspiciousPattern{Pattern: "/", MatchType: MatchExact},
			testPath: "/",
			want:     true,
		},
		{
			name:     "Query parameters preserved",
			pattern:  SuspiciousPattern{Pattern: "/login.html", MatchType: MatchContains},
			testPath: "/login.html?redirect=/admin",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPattern(tt.testPath, tt.pattern)
			if got != tt.want {
				t.Errorf("matchPattern(%q, %+v) = %v, want %v",
					tt.testPath, tt.pattern, got, tt.want)
			}
		})
	}
}

// matchPattern is a helper function that mimics the logic in the middleware
func matchPattern(path string, sp SuspiciousPattern) bool {
	pathLower := toLower(path)
	// Compute lowercase for test patterns that may not be from the global list
	patternLower := sp.LowerPattern
	if patternLower == "" {
		patternLower = toLower(sp.Pattern)
	}

	switch sp.MatchType {
	case MatchPrefix:
		return hasPrefix(pathLower, patternLower)
	case MatchSuffix:
		return hasSuffix(pathLower, patternLower)
	case MatchExact:
		return pathLower == patternLower
	case MatchContains:
		return contains(pathLower, patternLower)
	}
	return false
}

// Helper functions to avoid strings package import in test
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + ('a' - 'A')
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestRealWorldAttackPatterns(t *testing.T) {
	// Real attacks from production logs
	realAttacks := []struct {
		path        string
		description string
	}{
		{"/env.js", "Environment config file"},
		{"/env.development.js", "Development environment file"},
		{"/env.production.js", "Production environment file"},
		{"/.env", "Dotenv file"},
		{"/.env.production", "Production dotenv file"},
		{"/.well-known/openid-configuration/jwks.json", "OpenID JWKS"},
		{"/login.html", "HTML login page"},
		{"/CMSPages/Staging/SyncServer.asmx", "Kentico CMS ASMX"},
		{"/jmx-console/", "JBoss JMX Console"},
		{"/jbpm-console/app/tasks.jsf", "JBoss BPM JSF page"},
		{"/index.php", "PHP index"},
		{"/ldlogon.dll", "Windows DLL"},
		{"/log/mobile.log", "Mobile log file"},
		{"/jupyter/login", "Jupyter login"},
		{"/auth/admin", "Keycloak admin"},
		{"/secure/Signup.jspa", "JIRA signup"},
		{"/dana-na/auth/url_default/welcome.cgi", "Pulse Secure VPN"},
		{"/.idea/dataSources.xml", "IntelliJ IDEA config"},
		{"/karma.conf.js", "Karma config"},
		{"/WebCloud.asmx", "ASMX web service"},
	}

	for _, attack := range realAttacks {
		t.Run(attack.description, func(t *testing.T) {
			matched := false
			for _, pattern := range SuspiciousPathPatterns {
				if matchPattern(attack.path, pattern) {
					matched = true
					break
				}
			}
			if !matched {
				t.Errorf("Attack pattern %q (%s) was not detected by any pattern",
					attack.path, attack.description)
			}
		})
	}
}

func TestNoFalsePositives(t *testing.T) {
	// Legitimate paths that should NOT trigger
	legitimatePaths := []string{
		"/api/v1/events",
		"/api/v1/users/profile",
		"/api/v1/authentication/login",
		"/api/v1/certificates",
		"/",
		"/ready",
		"/health",
		"/metrics",
		"/swagger/index.html",
		"/.well-known/acme-challenge/token123", // Let's Encrypt SSL certificate validation (RFC 8615)
		"/.well-known/security.txt",            // Security disclosure (RFC 9116)
		"/.well-known/change-password",         // Password change endpoint (WICG spec)
	}

	for _, path := range legitimatePaths {
		t.Run(path, func(t *testing.T) {
			matched := false
			for _, pattern := range SuspiciousPathPatterns {
				if matchPattern(path, pattern) {
					matched = true
					t.Errorf("Legitimate path %q was incorrectly flagged by pattern %+v",
						path, pattern)
					break
				}
			}
			if matched {
				t.Errorf("Legitimate path %q should not match any suspicious pattern", path)
			}
		})
	}
}
