package google

type GoogleOAuthConfig struct {
	RedirectURL  string `env:"REDIRECT_URL" envDefault:"http://localhost:3000/oauth/google"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}
