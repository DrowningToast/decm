package eventconfig

import (
	"apps/backend/common/customerror"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ToggleCertificatePublishedRequest struct {
	IsPublished bool `json:"is_published"`
}

// ToggleCertificatePublished godoc
// @Summary Toggle certificate published status
// @Description Toggle the published status of event certificate configuration. When setting to published, inbox messages will be created for all certificate receivers. Only accessible by verified organizers.
// @ID toggle-certificate-published
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body ToggleCertificatePublishedRequest true "Toggle published request"
// @Success 200 {object} EventCertificateConfigResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 403 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/events/{event_id}/config/certificate/published [patch]
func (h *Handler) ToggleCertificatePublished(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("event_id"))
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req ToggleCertificatePublishedRequest
	if err := ctx.BodyParser(&req); err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	// Get current user from JWT claims
	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return customerror.Parse(&customerror.ErrUnauthorized, err)
	}

	// Toggle the published status and create inbox messages if publishing
	result, err := h.EventConfigUc.ToggleCertificatePublished(
		ctx.UserContext(),
		eventID,
		req.IsPublished,
		currentUser.UserId,
	)
	if err != nil {
		return errors.Wrap(err, "failed to toggle certificate published status")
	}

	// Get the full config with presigned URL
	config, err := h.EventConfigUc.GetEventCertificateConfigByEventID(ctx.UserContext(), eventID)
	if err != nil {
		return errors.Wrap(err, "failed to get event certificate config")
	}

	// Update the is_published field from the result
	config.IsPublished = result.IsPublished

	return ctx.Status(http.StatusOK).JSON(map[string]interface{}{
		"id":                             config.ID,
		"event_id":                       config.EventID,
		"base_certificate_storage_key":   config.BaseCertificateStorageKey,
		"base_certificate_presigned_url": config.BaseCertificatePresignedURL,
		"event_name_pos_x":               config.EventNamePosX,
		"event_name_pos_y":               config.EventNamePosY,
		"name_pos_x":                     config.NamePosX,
		"name_pos_y":                     config.NamePosY,
		"academic_institution_pos_x":     config.AcademicInstitutionPosX,
		"academic_institution_pos_y":     config.AcademicInstitutionPosY,
		"is_published":                   config.IsPublished,
		"created_at":                     config.CreatedAt,
		"updated_at":                     config.UpdatedAt,
		"mint_readiness":                 config.MintReadiness,
		"inbox_messages_created":         result.InboxMessagesCreated,
	})
}
