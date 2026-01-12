package oauth

import (
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"context"

	"github.com/cockroachdb/errors"

	customerror "apps/backend/common/customerror"

	oauth_services "apps/backend/services/oauth"

	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type OAuthUsecase struct {
	googleOAuthService *oauth_services.GoogleOAuthService
}

func NewOAuthUsecase(googleOAuthService *oauth_services.GoogleOAuthService, authenticationCredentialDg offchain_datagateway.AuthenticationCredentialDataGateway) *OAuthUsecase {
	return &OAuthUsecase{
		googleOAuthService: googleOAuthService,
	}
}

func (u *OAuthUsecase) VerifyGoogleOAuthCode(ctx context.Context, session *session.Session, code string, state string) (*oauth2.Token, error) {
	token, err := u.googleOAuthService.Callback(ctx, session, code, state)
	if err != nil {
		var customErr *customerror.Err
		if errors.As(err, &customErr) {
			return nil, customErr.Extend("failed to verify google oauth code")
		}
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return token, nil
}
