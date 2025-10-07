package google

type GoogleOAuthConfig struct {
	RedirectURL  string `env:"REDIRECT_URL" envDefault:"http://localhost:8080/api/v1/auth/verify-google-oauth"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	SuccessURL   string `env:"SUCCESS_URL" envDefault:"http://localhost:3000/oauth/google"`
}
