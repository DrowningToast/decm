package onboard

import (
	customerror "apps/backend/common/customerror"
	validatorutils "apps/backend/common/validatorutils"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type RegistrationMethod string

const (
	RegistrationMethodGoogle RegistrationMethod = "google"
	RegistrationMethodWallet RegistrationMethod = "wallet"
)

type CheckOnboardStatusRequest struct {
	Method RegistrationMethod `json:"method" validate:"required,oneof=google wallet"`

	AccessToken *string `json:"access_token"`
	ExpiresIn   *int    `json:"expires_in"`

	SignMessage *string `json:"sign_message"`
}

type CheckOnboardStatusResponse struct {
	AuthenticationCredentialId *string `json:"authentication_credential_id"`
	ProfileId                  *string `json:"profile_id"`
}

// @Summary Check onboard status
// @Description Check onboard status
// @ID check-onboard-status
// @Tags Onboard
// @Param method body onboard.CheckOnboardStatusRequest.Method true "Method"
// @Param access_token body onboard.CheckOnboardStatusRequest.AccessToken false "Access token"
// @Param expires_in body onboard.CheckOnboardStatusRequest.ExpiresIn false "Expires in"
// @Param sign_message body onboard.CheckOnboardStatusRequest.SignMessage false "Sign message"
// @Accept json
// @Produce json
// @Success 200 {object} CheckOnboardStatusResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Router /api/v1/onboard/check-onboard-status [post]
func (h Handler) CheckOnboardStatus(ctx *fiber.Ctx) error {
	requestBody := CheckOnboardStatusRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	switch requestBody.Method {
	case RegistrationMethodGoogle:
		accessToken := requestBody.AccessToken
		if accessToken == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("access token is required"))
		}
		expiresIn := requestBody.ExpiresIn
		if expiresIn == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("expires in is required"))
		}
		authenticationCredentialId, profileId, err := h.OnboardUc.CheckOnboardStatusWithGoogleConnectorRef(ctx.UserContext(), *accessToken, *expiresIn)
		if err != nil {
			return errors.Wrap(err, "failed to check onboard status with google connector ref")
		}
		if authenticationCredentialId == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: nil,
				ProfileId:                  nil,
			})
		}
		authenticationCredentialIdStr := authenticationCredentialId.String()
		profileIdStr := profileId.String()
		return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
			AuthenticationCredentialId: &authenticationCredentialIdStr,
			ProfileId:                  &profileIdStr,
		})
	case RegistrationMethodWallet:
		signMessage := requestBody.SignMessage
		if signMessage == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is required"))
		}
		authenticationCredentialId, profileId, err := h.OnboardUc.CheckOnboardStatusWithWalletAddress(ctx.UserContext(), *signMessage)
		if err != nil {
			return errors.Wrap(err, "failed to check onboard status with wallet address")
		}
		if authenticationCredentialId == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: nil,
				ProfileId:                  nil,
			})
		}
		authenticationCredentialIdStr := authenticationCredentialId.String()
		profileIdStr := profileId.String()
		return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
			AuthenticationCredentialId: &authenticationCredentialIdStr,
			ProfileId:                  &profileIdStr,
		})
	default:
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("method is invalid"))
	}
}

func (r *CheckOnboardStatusRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

func (r *CheckOnboardStatusRequest) IsValid() error {
	if err := validatorutils.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}
