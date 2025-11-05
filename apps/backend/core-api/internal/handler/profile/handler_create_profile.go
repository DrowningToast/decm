package profile

import (
	"errors"
	"fmt"

	"apps/backend/common/customerror"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/profile"

	"github.com/gofiber/fiber/v2"
)

type CreateProfileRequest struct {
	profile.CreateProfileParameters
}

type CreateProfileResponse struct {
	entity.Profile
}

func (r *CreateProfileRequest) Parse(c *fiber.Ctx) error {
	if err := c.BodyParser(r); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	return nil
}

// @Summary Create a profile
// @Description Create a profile
// @ID create-profile
// @Tags Profile
// @Accept json
// @Produce json
// @Param profile body CreateProfileRequest true "Profile"
// @Success 200 {object} CreateProfileResponse
// @Failure 400 {object} customerror.Err
// @Failure 403 {object} customerror.Err
// @Failure 404 {object} customerror.Err
// @Failure 500 {object} customerror.Err
// @Router /api/v1/profile [post]
func (h *Handler) CreateProfile(c *fiber.Ctx) error {
	var createProfileRequest CreateProfileRequest
	err := createProfileRequest.Parse(c)
	if err != nil {
		return err
	}
	if err := createProfileRequest.CreateProfileParameters.IsValid(); err != nil {
		return err
	}
	isCurrentUser, err := h.AuthenticationService.IsCurrentUser(c, createProfileRequest.CreateProfileParameters.AuthenticationCredentialId)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if !isCurrentUser {
		return customerror.Parse(&customerror.ErrUnauthorized, errors.New("mismatch authentication credential id"))
	}

	fmt.Println(`createProfileRequest.CreateProfileParameters`, createProfileRequest.CreateProfileParameters)
	fmt.Println(`createProfileRequest.CreateProfileParameters.AuthenticationCredentialId`, createProfileRequest.CreateProfileParameters.AuthenticationCredentialId)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsProfilePicturePublic`, createProfileRequest.CreateProfileParameters.IsProfilePicturePublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.ProfilePictureUrl`, createProfileRequest.CreateProfileParameters.ProfilePictureUrl)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsFirstNamePublic`, createProfileRequest.CreateProfileParameters.IsFirstNamePublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.FirstName`, createProfileRequest.CreateProfileParameters.FirstName)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsLastNamePublic`, createProfileRequest.CreateProfileParameters.IsLastNamePublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.LastName`, createProfileRequest.CreateProfileParameters.LastName)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsEmailPublic`, createProfileRequest.CreateProfileParameters.IsEmailPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.Email`, createProfileRequest.CreateProfileParameters.Email)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsBioPublic`, createProfileRequest.CreateProfileParameters.IsBioPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.Bio`, createProfileRequest.CreateProfileParameters.Bio)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsPhoneNumberPublic`, createProfileRequest.CreateProfileParameters.IsPhoneNumberPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.PhoneNumber`, createProfileRequest.CreateProfileParameters.PhoneNumber)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsAddressPublic`, createProfileRequest.CreateProfileParameters.IsAddressPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.Address`, createProfileRequest.CreateProfileParameters.Address)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsAcademicInstitutionPublic`, createProfileRequest.CreateProfileParameters.IsAcademicInstitutionPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.AcademicInstitution`, createProfileRequest.CreateProfileParameters.AcademicInstitution)
	fmt.Println(`createProfileRequest.CreateProfileParameters.IsAcademicEmailPublic`, createProfileRequest.CreateProfileParameters.IsAcademicEmailPublic)
	fmt.Println(`createProfileRequest.CreateProfileParameters.AcademicEmail`, createProfileRequest.CreateProfileParameters.AcademicEmail)
	profile, err := h.ProfileUc.CreateProfile(c.Context(), createProfileRequest.CreateProfileParameters)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(CreateProfileResponse{
		Profile: *profile,
	})
}
