package logger

import (
	"testing"
)

func TestIsSuspiciousPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		// Environment files
		{name: "Environment file - .env", path: "/.env", want: true},
		{name: "Environment file - .env.production", path: "/.env.production", want: true},
		{name: "Environment file - env.js", path: "/env.js", want: true},
		{name: "Environment file - env.development.js", path: "/env.development.js", want: true},
		{name: "Environment file - env.prod.js", path: "/env.prod.js", want: true},

		// Config files
		{name: "Git config - .git/config", path: "/.git/config", want: true},
		{name: "IDEA config - .idea/workspace.xml", path: "/.idea/workspace.xml", want: true},
		{name: "Config dir - .config/karma.conf.js", path: "/.config/karma.conf.js", want: true},
		{name: "Well-known - .well-known/jwks.json", path: "/.well-known/jwks.json", want: true},

		// PHP files
		{name: "PHP - index.php", path: "/index.php", want: true},
		{name: "PHP - config.php", path: "/config.php", want: true},
		{name: "PHP - shell.php", path: "/upload/shell.php", want: true},

		// ASP.NET files
		{name: "ASPX - admin.aspx", path: "/admin.aspx", want: true},
		{name: "ASMX - WebCloud.asmx", path: "/WebCloud.asmx", want: true},
		{name: "ASHX - handler.ashx", path: "/handler.ashx", want: true},

		// Java files
		{name: "JSP - index.jsp", path: "/index.jsp", want: true},
		{name: "JSPA - Signup.jspa", path: "/secure/Signup.jspa", want: true},
		{name: "JSF - tasks.jsf", path: "/jbpm-console/app/tasks.jsf", want: true},

		// Windows
		{name: "DLL - ldlogon.dll", path: "/ldlogon.dll", want: true},
		{name: "INI - win.ini", path: "/c:/windows/win.ini", want: true},

		// Log files
		{name: "Log - mobile.log", path: "/log/mobile.log", want: true},
		{name: "Log - firewall.log", path: "/log/firewall.log", want: true},

		// Java enterprise servers
		{name: "JMX Console", path: "/jmx-console/", want: true},
		{name: "JBoss BPM", path: "/jbpm-console/app/tasks.jsf", want: true},
		{name: "JBoss Web Services", path: "/jbossws/services", want: true},
		{name: "Jolokia", path: "/jolokia/search/*:test=test", want: true},
		{name: "JasperServer", path: "/jasperserver/login.html", want: true},
		{name: "Jenkins", path: "/jenkins/script", want: true},

		// Jupyter/Data science
		{name: "Jupyter lab", path: "/jupyter/lab", want: true},
		{name: "Jupyter login", path: "/jupyter/login", want: true},
		{name: "IPython", path: "/ipython/tree", want: true},
		{name: "JupyterHub", path: "/hub/login", want: true},

		// CMS
		{name: "Kentico CMS", path: "/CMSPages/logon.aspx", want: true},
		{name: "WordPress", path: "/wp-admin", want: true},

		// Auth/SSO
		{name: "Keycloak admin", path: "/auth/admin", want: true},
		{name: "Keycloak realms", path: "/auth/realms/master/.well-known/openid-configuration", want: true},
		{name: "JIRA signup", path: "/secure/Signup!default.jspa", want: true},

		// VPN
		{name: "Pulse Secure VPN", path: "/dana-na/auth/url_default/welcome.cgi", want: true},

		// Dashboards
		{name: "Dashboard", path: "/dashboard/", want: true},
		{name: "Webmail", path: "/webmail/login/", want: true},

		// Other enterprise
		{name: "Kafka UI", path: "/kafka/clusters", want: true},
		{name: "Kiali", path: "/kiali/api/status", want: true},

		// Legitimate paths (should NOT match)
		{name: "Legitimate - API events", path: "/api/v1/events", want: false},
		{name: "Legitimate - API users", path: "/api/v1/users", want: false},
		{name: "Legitimate - API auth", path: "/api/v1/authentication/login", want: false},
		{name: "Legitimate - root", path: "/", want: false},
		{name: "Legitimate - ready", path: "/ready", want: false},
		{name: "Legitimate - health", path: "/health", want: false},
		{name: "Legitimate - metrics", path: "/metrics", want: false},
		{name: "Legitimate - swagger", path: "/swagger/index.html", want: false},
		{name: "Legitimate - certificate", path: "/api/v1/certificates/123", want: false},

		// Edge cases
		{name: "Empty path", path: "", want: false},
		{name: "Query params", path: "/login.html?redirect=/admin", want: true},
		{name: "Case insensitive - uppercase PHP", path: "/INDEX.PHP", want: true},
		{name: "Case insensitive - uppercase JMX", path: "/JMX-CONSOLE/", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSuspiciousPath(tt.path)
			if got != tt.want {
				t.Errorf("isSuspiciousPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsSuspiciousPathRealAttacks(t *testing.T) {
	// Real attack paths from production logs
	realAttacks := []string{
		"/env.js",
		"/env.development.js",
		"/env.production.js",
		"/.env",
		"/.env.production",
		"/.well-known/openid-configuration/jwks.json",
		"/login.html",
		"/jwks.json",
		"/jwks",
		"/CMSPages/Staging/SyncServer.asmx",
		"/index.php",
		"/ldlogon.dll",
		"/log/mobile.log",
		"/jmx-console/",
		"/jbpm-console/app/tasks.jsf",
		"/jbossws/services",
		"/jolokia/",
		"/jasperserver/login.html",
		"/jeecg-boot/",
		"/jupyter/login",
		"/ipython/tree",
		"/hub/login",
		"/auth/admin",
		"/auth/realms/master/.well-known/openid-configuration",
		"/secure/Signup.jspa",
		"/dana-na/auth/url_default/welcome.cgi",
		"/.idea/dataSources.xml",
		"/.idea/WebServers.xml",
		"/karma.conf.js",
		"/WebCloud.asmx",
		"/kafka/clusters",
		"/kiali/api/status",
		"/dashboard/",
		"/webmail/login/",
		"/UtilServlet",
	}

	for _, attack := range realAttacks {
		t.Run(attack, func(t *testing.T) {
			if !isSuspiciousPath(attack) {
				t.Errorf("Real attack path %q was not detected as suspicious", attack)
			}
		})
	}
}

func TestIsSuspiciousPathNoFalsePositives(t *testing.T) {
	// Legitimate API paths that should NOT be flagged
	legitimatePaths := []string{
		"/api/v1/events",
		"/api/v1/events/123",
		"/api/v1/users",
		"/api/v1/users/profile",
		"/api/v1/authentication/login",
		"/api/v1/authentication/logout",
		"/api/v1/authentication/refresh",
		"/api/v1/certificates",
		"/api/v1/certificates/verify",
		"/api/v1/blockchain/status",
		"/",
		"/ready",
		"/health",
		"/metrics",
		"/swagger/index.html",
		"/swagger/doc.json",
		"/.well-known/acme-challenge/token123", // Let's Encrypt SSL validation (RFC 8615)
		"/.well-known/security.txt",            // Security disclosure (RFC 9116)
		"/.well-known/change-password",         // Password change endpoint (WICG spec)
	}

	for _, path := range legitimatePaths {
		t.Run(path, func(t *testing.T) {
			if isSuspiciousPath(path) {
				t.Errorf("Legitimate path %q was incorrectly flagged as suspicious", path)
			}
		})
	}
}

// Benchmark the pattern matching function
func BenchmarkIsSuspiciousPath(b *testing.B) {
	testPaths := []string{
		"/api/v1/events",       // Legitimate (should be fast)
		"/.env.production",     // Suspicious prefix match
		"/index.php",           // Suspicious suffix match
		"/jmx-console/",        // Suspicious contains match
		"/admin",               // Suspicious exact match
		"/very/long/path/that/is/legitimate/and/should/not/match",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			isSuspiciousPath(path)
		}
	}
}
