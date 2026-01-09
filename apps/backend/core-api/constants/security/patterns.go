package security

// MatchType defines how a pattern should be matched against a request path
type MatchType int

const (
	// MatchContains checks if the pattern exists anywhere in the path
	MatchContains MatchType = iota
	// MatchPrefix checks if the path starts with the pattern
	MatchPrefix
	// MatchSuffix checks if the path ends with the pattern
	MatchSuffix
	// MatchExact checks if the path exactly matches the pattern
	MatchExact
)

// SuspiciousPattern defines a pattern and how it should be matched
type SuspiciousPattern struct {
	Pattern   string
	MatchType MatchType
	Comment   string // Optional description for documentation
}

// SuspiciousPathPatterns contains URL patterns that indicate security probes or attacks.
// These patterns are used by the logger middleware to identify and tag suspicious requests
// with category="security" for monitoring and alerting in Grafana/Loki.
//
// Usage in Grafana query:
//
//	{container="decm-backend-prod"} | json | category="security"
var SuspiciousPathPatterns = []SuspiciousPattern{
	// Configuration and credential files (PREFIX matching - more precise)
	{Pattern: "/.env", MatchType: MatchPrefix, Comment: "Environment files (dot prefix)"},
	{Pattern: "/env.", MatchType: MatchPrefix, Comment: "Environment files (no dot)"},
	{Pattern: "/.git", MatchType: MatchPrefix, Comment: "Git repository"},
	{Pattern: "/.aws", MatchType: MatchPrefix, Comment: "AWS credentials"},
	{Pattern: "/.ssh", MatchType: MatchPrefix, Comment: "SSH keys"},
	{Pattern: "/.npmrc", MatchType: MatchPrefix, Comment: "NPM config"},
	{Pattern: "/.dockerenv", MatchType: MatchPrefix, Comment: "Docker environment"},
	{Pattern: "/.htaccess", MatchType: MatchPrefix, Comment: "Apache config"},
	{Pattern: "/.idea", MatchType: MatchPrefix, Comment: "IntelliJ IDEA config"},
	{Pattern: "/.config", MatchType: MatchPrefix, Comment: "Config directory"},
	{Pattern: "/.well-known", MatchType: MatchPrefix, Comment: "Well-known URIs"},

	// Suspicious file extensions (SUFFIX matching)
	// Technologies we don't use (ASP.NET, JSP, CGI, Perl, PHP, etc.)
	{Pattern: ".aspx", MatchType: MatchSuffix, Comment: "ASP.NET pages"},
	{Pattern: ".asp", MatchType: MatchSuffix, Comment: "Classic ASP"},
	{Pattern: ".asmx", MatchType: MatchSuffix, Comment: "ASMX web services"},
	{Pattern: ".ashx", MatchType: MatchSuffix, Comment: "ASP.NET handlers"},
	{Pattern: ".jsp", MatchType: MatchSuffix, Comment: "Java Server Pages"},
	{Pattern: ".jspa", MatchType: MatchSuffix, Comment: "JIRA pages"},
	{Pattern: ".jsf", MatchType: MatchSuffix, Comment: "JavaServer Faces"},
	{Pattern: ".cgi", MatchType: MatchSuffix, Comment: "CGI scripts"},
	{Pattern: ".pl", MatchType: MatchSuffix, Comment: "Perl scripts"},
	{Pattern: ".php", MatchType: MatchSuffix, Comment: "PHP files"},
	{Pattern: ".dll", MatchType: MatchSuffix, Comment: "Windows DLL"},
	{Pattern: ".bak", MatchType: MatchSuffix, Comment: "Backup files"},
	{Pattern: ".sql", MatchType: MatchSuffix, Comment: "SQL dump files"},
	{Pattern: ".log", MatchType: MatchSuffix, Comment: "Log files"},
	{Pattern: ".zip", MatchType: MatchSuffix, Comment: "Archive files"},
	{Pattern: ".tar.gz", MatchType: MatchSuffix, Comment: "Compressed archives"},
	{Pattern: ".ini", MatchType: MatchSuffix, Comment: "INI config files"},

	// Exact path matches (EXACT matching - fastest for known endpoints)
	{Pattern: "/admin", MatchType: MatchExact, Comment: "Admin panel"},
	{Pattern: "/phpmyadmin", MatchType: MatchExact, Comment: "phpMyAdmin"},
	{Pattern: "/wp-admin", MatchType: MatchExact, Comment: "WordPress admin"},
	{Pattern: "/wp-login", MatchType: MatchExact, Comment: "WordPress login"},
	{Pattern: "/administrator", MatchType: MatchExact, Comment: "Joomla admin"},
	{Pattern: "/cpanel", MatchType: MatchExact, Comment: "cPanel"},
	{Pattern: "/jwks", MatchType: MatchExact, Comment: "JWKS endpoint"},
	{Pattern: "/jwks.json", MatchType: MatchExact, Comment: "JWKS JSON"},
	{Pattern: "/graphql", MatchType: MatchExact, Comment: "GraphQL endpoint"},
	{Pattern: "/server-status", MatchType: MatchExact, Comment: "Apache status"},
	{Pattern: "/server-info", MatchType: MatchExact, Comment: "Server info"},

	// CMS/Admin panels (CONTAINS matching - for path segments)
	{Pattern: "/CMSPages/", MatchType: MatchContains, Comment: "Kentico CMS"},
	{Pattern: "/Kenesto/", MatchType: MatchContains, Comment: "Kenesto platform"},
	{Pattern: "/wp-content", MatchType: MatchContains, Comment: "WordPress content"},
	{Pattern: "/wp-includes", MatchType: MatchContains, Comment: "WordPress includes"},
	{Pattern: "/cgi-bin", MatchType: MatchContains, Comment: "CGI directory"},
	{Pattern: "/manager/html", MatchType: MatchContains, Comment: "Tomcat manager"},
	{Pattern: "/dashboard", MatchType: MatchExact, Comment: "Generic dashboard root"},
	{Pattern: "/dashboard/", MatchType: MatchPrefix, Comment: "Generic dashboard paths"},
	{Pattern: "/webmail", MatchType: MatchContains, Comment: "Webmail interfaces"},

	// Java/Enterprise Application Servers
	{Pattern: "/jmx-console", MatchType: MatchContains, Comment: "JMX Console (JBoss)"},
	{Pattern: "/jbpm-console", MatchType: MatchContains, Comment: "JBoss BPM"},
	{Pattern: "/jbossws", MatchType: MatchContains, Comment: "JBoss Web Services"},
	{Pattern: "/web-console", MatchType: MatchContains, Comment: "JBoss Web Console"},
	{Pattern: "/jolokia", MatchType: MatchContains, Comment: "Jolokia JMX-HTTP"},
	{Pattern: "/jasperserver", MatchType: MatchContains, Comment: "JasperReports Server"},
	{Pattern: "/jeecg-boot", MatchType: MatchContains, Comment: "JeecgBoot framework"},
	{Pattern: "/jenkins", MatchType: MatchContains, Comment: "Jenkins CI/CD"},
	{Pattern: "/jexws", MatchType: MatchContains, Comment: "Java exploit shell"},
	{Pattern: "/jexinv", MatchType: MatchContains, Comment: "Java exploit invoker"},
	{Pattern: "/jbossass", MatchType: MatchContains, Comment: "JBoss assault"},

	// Jupyter/Data Science
	{Pattern: "/jupyter", MatchType: MatchContains, Comment: "Jupyter notebooks"},
	{Pattern: "/ipython", MatchType: MatchContains, Comment: "IPython"},
	{Pattern: "/hub/login", MatchType: MatchContains, Comment: "JupyterHub login"},
	{Pattern: "/terminals", MatchType: MatchContains, Comment: "Jupyter terminals"},

	// Authentication/SSO
	{Pattern: "/auth/admin", MatchType: MatchContains, Comment: "Keycloak admin"},
	{Pattern: "/auth/realms", MatchType: MatchContains, Comment: "Keycloak realms"},
	{Pattern: "/servicedesk/customer/user/signup", MatchType: MatchContains, Comment: "JIRA Service Desk"},
	{Pattern: "/secure/Signup", MatchType: MatchContains, Comment: "JIRA signup"},

	// VPN/Remote Access
	{Pattern: "/dana-na", MatchType: MatchContains, Comment: "Pulse Secure VPN"},

	// Other enterprise apps
	{Pattern: "/kafka/clusters", MatchType: MatchContains, Comment: "Kafka UI"},
	{Pattern: "/kiali", MatchType: MatchContains, Comment: "Kiali (Istio)"},
	{Pattern: "/jackett", MatchType: MatchContains, Comment: "Jackett"},
	{Pattern: "/UtilServlet", MatchType: MatchContains, Comment: "Util servlet exploits"},

	// OpenID/OAuth reconnaissance
	{Pattern: "/.well-known/openid-configuration", MatchType: MatchContains, Comment: "OpenID discovery"},
	{Pattern: "/oauth", MatchType: MatchContains, Comment: "OAuth endpoints"},

	// Common vulnerability paths and shells
	{Pattern: "/eval-stdin.php", MatchType: MatchContains, Comment: "PHP eval shell"},
	{Pattern: "/shell.php", MatchType: MatchContains, Comment: "Web shell"},
	{Pattern: "/upload.php", MatchType: MatchContains, Comment: "Upload endpoint"},
	{Pattern: "/filemanager", MatchType: MatchContains, Comment: "File manager"},
	{Pattern: "/elfinder", MatchType: MatchContains, Comment: "elFinder file manager"},
	{Pattern: "/tinymce", MatchType: MatchContains, Comment: "TinyMCE editor"},
	{Pattern: "/phpinfo.php", MatchType: MatchContains, Comment: "PHP info page"},
	{Pattern: "/config.php", MatchType: MatchContains, Comment: "PHP config"},
	{Pattern: "/configuration.php", MatchType: MatchContains, Comment: "PHP configuration"},
	{Pattern: "/web.config", MatchType: MatchContains, Comment: "IIS config"},
	{Pattern: "/database.yml", MatchType: MatchContains, Comment: "Database config"},

	// Legacy login pages
	{Pattern: "/login.html", MatchType: MatchContains, Comment: "HTML login page"},
	{Pattern: "/logon.aspx", MatchType: MatchContains, Comment: "ASP.NET login"},
	{Pattern: "/signin.html", MatchType: MatchContains, Comment: "HTML signin page"},

	// API abuse patterns (non-existent in our app)
	{Pattern: "/api/graphql", MatchType: MatchContains, Comment: "GraphQL API"},
	{Pattern: "/api/xml", MatchType: MatchContains, Comment: "XML API"},
	{Pattern: "/api/system/readiness", MatchType: MatchContains, Comment: "Generic system API"},
	{Pattern: "/api/system/liveness", MatchType: MatchContains, Comment: "Generic system API"},
	{Pattern: "/api/login_settings", MatchType: MatchContains, Comment: "Login settings API"},

	// Config files (JavaScript configs that shouldn't be exposed)
	{Pattern: "karma.conf.js", MatchType: MatchContains, Comment: "Karma test config"},
	{Pattern: "jsconfig.json", MatchType: MatchContains, Comment: "JS config"},
	{Pattern: "keycloak.json", MatchType: MatchContains, Comment: "Keycloak config"},
}
