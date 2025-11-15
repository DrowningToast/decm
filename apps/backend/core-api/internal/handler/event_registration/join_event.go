package event_registration

import (
	"encoding/hex"

	customerror "apps/backend/common/customerror"
	"apps/backend/common/validatorutils"
	"apps/backend/core-api/internal/usecase/cyptoutils"

	event_registration_uc "apps/backend/core-api/internal/usecase/event_registration"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JoinEventRequest struct {
	Password  *string `json:"password,omitempty"`
	Signature *string `json:"signature,omitempty"`
}

type GetJoinEventSignMessageResponse struct {
	SignMessage string `json:"sign_message"`
}

// @Summary Get join event sign message
// @Description Get join event sign message
// @ID get-join-event-sign-message
// @Tags Event Registration
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Success 200 {object} GetJoinEventSignMessageResponse
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/join/{event_id}/sign-message [get]
func (h *Handler) GetJoinEventSignMessage(ctx *fiber.Ctx) error {
	eventIdStr := ctx.Params("event_id")
	if eventIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("event_id is required"))
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}
	eventContract, err := h.EventUc.GetEventContractByEventID(ctx.UserContext(), eventId)
	if err != nil {
		return errors.Wrap(err, "failed to get event contract by event id")
	}
	if eventContract == nil {
		return customerror.Parse(&customerror.ErrNotFound, errors.New("event contract not found"))
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return errors.Wrap(err, "failed to get ethereum client")
	}

	signMessage, _, err := h.EventRegistrationUc.GetJoinEventSignMessage(ctx.UserContext(), client, common.HexToAddress(currentUser.WalletAddress), *currentUser, common.HexToAddress(eventContract.EventContractAddress))
	if err != nil {
		return errors.Wrap(err, "failed to get join event sign message")
	}

	return ctx.Status(fiber.StatusOK).JSON(GetJoinEventSignMessageResponse{
		SignMessage: *signMessage,
	})
}

// @name JoinEventBody
// @description Body for joining an event
// @Accept json
// @Produce json
// @Param password body string false "Password"
// @Param signature body string false "Signature"
// @Param sign_message body string false "Sign message"
// @Param registration_data body event_registration_uc.JoinEventPayload true "Registration data"
// @Success 200 {object} entity.EventAttendee
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/join/{event_id} [post]
type JoinEventBody struct {
	Password    *string `json:"password,omitempty"`
	Signature   *string `json:"signature,omitempty"`
	SignMessage *string `json:"sign_message,omitempty"`

	RegistrationData event_registration_uc.JoinEventPayload `json:"registration_data" validate:"required"`
}

func (r *JoinEventBody) IsValid() error {
	return validatorutils.ValidateStruct(r)
}

func (r *JoinEventBody) Parse(ctx *fiber.Ctx) error {
	return ctx.BodyParser(r)
}

// @Summary Join event
// @Description Join an event
// @ID join-event
// @Tags Event Registration
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body JoinEventRequest true "Join event request"
// @Failure 400 {object} customerror.ErrResponse
// @Failure 401 {object} customerror.ErrResponse
// @Failure 404 {object} customerror.ErrResponse
// @Failure 500 {object} customerror.ErrResponse
// @Router /api/v1/event-registration/join/{event_id} [post]
func (h *Handler) JoinEvent(ctx *fiber.Ctx) error {
	eventIdStr := ctx.Params("event_id")
	if eventIdStr == "" {
		return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("event_id is required"))
	}
	eventId, err := uuid.Parse(eventIdStr)
	if err != nil {
		return customerror.Parse(&customerror.ErrInvalidArgument, err)
	}

	var req JoinEventBody
	if err := req.Parse(ctx); err != nil {
		return err
	}
	if err := req.IsValid(); err != nil {
		return err
	}

	currentUser, err := h.AuthenticationService.GetUserContext(ctx)
	if err != nil {
		return err
	}

	client, err := cyptoutils.GetEthereumClient()
	if err != nil {
		return errors.Wrap(err, "failed to get ethereum client")
	}

	proof := event_registration_uc.CheckRegistrationEligibilityParams{
		Password: req.Password,
	}

	if req.Signature != nil {
		signature, err := hex.DecodeString(*req.Signature)
		if err != nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, err)
		}
		if req.SignMessage == nil {
			return customerror.Parse(&customerror.ErrInvalidArgument, errors.New("original sign message is required"))
		}
		eventAttendee, err := h.EventRegistrationUc.JoinEventWithSignature(ctx.UserContext(), client, currentUser, eventId, proof, req.RegistrationData, signature, *req.SignMessage)
		if err != nil {
			return errors.Wrap(err, "failed to join event with signature")
		}
		return ctx.Status(fiber.StatusOK).JSON(eventAttendee)
	} else {
		eventAttendee, err := h.EventRegistrationUc.JoinEventWithPin(ctx.UserContext(), client, currentUser, eventId, proof, req.RegistrationData, *req.Password)
		if err != nil {
			return errors.Wrap(err, "failed to join event with pin")
		}
		return ctx.Status(fiber.StatusOK).JSON(eventAttendee)
	}
}
