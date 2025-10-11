package oauth

import (
	"context"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type OAuthUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type OAuthService interface {
	Login(session *session.Session) (*string, error)
	Callback(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUser, error)
}

func ParseToken(accessToken string) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: accessToken,
	}, nil
}
