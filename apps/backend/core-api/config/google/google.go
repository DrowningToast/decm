package google

type GoogleOAuthConfig struct {
	RedirectURL  string `env:"REDIRECT_URL" envDefault:"http://localhost:8000/api/v1/onboard/google/callback"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}
