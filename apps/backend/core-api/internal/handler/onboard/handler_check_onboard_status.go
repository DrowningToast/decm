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
	Method *RegistrationMethod `json:"method,omitempty" validate:"omitempty,oneof=google wallet"`

	AccessToken *string `json:"access_token,omitempty"`
	ExpiresIn   *int    `json:"expires_in,omitempty"`

	MessageSignature *string `json:"message_signature,omitempty"`
}

type CheckOnboardStatusResponse struct {
	AuthenticationCredentialId *string `json:"authentication_credential_id,omitempty"`
	ProfileId                  *string `json:"profile_id,omitempty"`
}

// @Summary Check onboard status
// @Description Check onboard status
// @ID check-onboard-status
// @Tags Onboard
// @Param method body onboard.CheckOnboardStatusRequest.Method false "Method"
// @Param access_token body onboard.CheckOnboardStatusRequest.AccessToken false "Access token"
// @Param expires_in body onboard.CheckOnboardStatusRequest.ExpiresIn false "Expires in"
// @Param message_signature body onboard.CheckOnboardStatusRequest.MessageSignature false "Message signature"
// @Accept json
// @Produce json
// @Success 200 {object} CheckOnboardStatusResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Router /api/v1/onboard/check-onboard-status [post]
func (h Handler) CheckOnboardStatus(ctx *fiber.Ctx) error {
	currentUser, err := h.AuthService.GetUserContext(ctx)
	if err != nil {
		return err
	}
	// If user is authenticated, check onboard status with JWT
	if currentUser != nil {
		jwt, profileId, err := h.OnboardUc.CheckOnboardStatusWithJwt(ctx.UserContext(), *currentUser)
		if err != nil {
			return errors.Wrap(err, "failed to check onboard status with jwt")
		}
		authenticationCredentialIdStr := jwt.UserId.String()
		if profileId == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: &authenticationCredentialIdStr,
				ProfileId:                  nil,
			})
		}
		profileIdStr := profileId.String()
		return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
			AuthenticationCredentialId: &authenticationCredentialIdStr,
			ProfileId:                  &profileIdStr,
		})
	}

	requestBody := CheckOnboardStatusRequest{}
	if err := requestBody.Parse(ctx); err != nil {
		return err
	}
	if err := requestBody.IsValid(); err != nil {
		return err
	}

	if requestBody.Method == nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("method is required when not authenticated"))
	}

	switch *requestBody.Method {
	// If user is not authenticated, check onboard status with Google connector ref
	case RegistrationMethodGoogle:
		if requestBody.AccessToken == nil || requestBody.ExpiresIn == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("access token and expires in are required"))
		}
		jwt, profileId, err := h.OnboardUc.CheckOnboardStatusWithGoogleConnectorRef(ctx.UserContext(), *requestBody.AccessToken, *requestBody.ExpiresIn)
		if err != nil {
			return errors.Wrap(err, "failed to check onboard status with google connector ref")
		}
		if jwt == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: nil,
				ProfileId:                  nil,
			})
		}
		jwtToken, err := h.AuthService.CreateToken(*jwt)
		if err != nil {
			return errors.Wrap(err, "failed to create jwt token")
		}
		h.AuthService.SetJwtCookie(ctx, jwtToken)
		authenticationCredentialIdStr := jwt.UserId.String()
		if profileId == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: &authenticationCredentialIdStr,
				ProfileId:                  nil,
			})
		}
		profileIdStr := profileId.String()
		return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
			AuthenticationCredentialId: &authenticationCredentialIdStr,
			ProfileId:                  &profileIdStr,
		})
	// If user is not authenticated, check onboard status with wallet address
	case RegistrationMethodWallet:
		messageSignature := requestBody.MessageSignature
		if messageSignature == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is required"))
		}
		jwt, profileId, err := h.OnboardUc.CheckOnboardStatusWithWalletAddress(ctx.UserContext(), *messageSignature)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("message signature is required"))
		}
		if jwt == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: nil,
				ProfileId:                  nil,
			})
		}
		jwtToken, err := h.AuthService.CreateToken(*jwt)
		if err != nil {
			return errors.Wrap(err, "failed to create jwt token")
		}
		h.AuthService.SetJwtCookie(ctx, jwtToken)
		authenticationCredentialIdStr := jwt.UserId.String()
		if profileId == nil {
			return ctx.Status(fiber.StatusOK).JSON(CheckOnboardStatusResponse{
				AuthenticationCredentialId: &authenticationCredentialIdStr,
				ProfileId:                  nil,
			})
		}
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
