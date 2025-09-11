package oauth

import (
	"context"

	"apps/backend/common/customerror"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type OAuthUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type OAuthService interface {
	Login(session *session.Session) (*string, *customerror.Err)
	Callback(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, *customerror.Err)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*OAuthUser, *customerror.Err)
}

func ParseToken(accessToken string, refreshToken string) (*oauth2.Token, *customerror.Err) {
	return &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
