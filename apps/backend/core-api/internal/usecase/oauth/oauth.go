package oauth

import (
	"context"

	customerror "apps/backend/common/customerror"
	"apps/backend/core-api/internal/datagateway"
	oauth_services "apps/backend/services/oauth"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type OAuthUsecase struct {
	googleOAuthService *oauth_services.GoogleOAuthService
}

func NewOAuthUsecase(googleOAuthService *oauth_services.GoogleOAuthService, authenticationCredentialDg datagateway.AuthenticationCredentialDataGateway) *OAuthUsecase {
	return &OAuthUsecase{
		googleOAuthService: googleOAuthService,
	}
}

func (u *OAuthUsecase) VerifyGoogleOAuthCode(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, *customerror.Err) {
	token, err := u.googleOAuthService.Callback(ctx, session, code, state)
	if err != nil {
		return nil, err.Extend("failed to verify google oauth code")
	}

	return token, nil
}
