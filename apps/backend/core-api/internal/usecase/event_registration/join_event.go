package event_registration

import (
	"context"

	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	datagateway "apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"

	eventContract "apps/backend/contracts/event"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type JoinEventUserError string

const (
	JoinEventUserErrorEventAttendeeFull          JoinEventUserError = "event_attendee_full"
	JoinEventUserErrorEventAttendeeAlreadyJoined JoinEventUserError = "event_attendee_already_joined"
)

type JoinEventWithPasswordUserError string

type JoinEventPayload struct {
	FirstName           *string `json:"first_name,omitempty"`
	LastName            *string `json:"last_name,omitempty"`
	Email               *string `json:"email,omitempty"`
	PhoneNumber         *string `json:"phone_number,omitempty"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
	AcademicEmail       *string `json:"academic_email,omitempty"`
	Address             *string `json:"address,omitempty"`
	Bio                 *string `json:"bio,omitempty"`
}

// returns raw string, then message hash
func (uc *EventRegistrationUsecase) GetJoinEventSignMessage(ctx context.Context, client *ethclient.Client, walletAddress common.Address, currentUser auth.JwtClaims, eventContractAddress common.Address) (*string, *ethcommon.Hash, error) {
	// Validation
	calculatedDeadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get calculated deadline block")
	}

	signMessage, err := cyptoutils.GetSignMessage(common.HexToAddress(currentUser.WalletAddress), eventContractAddress, calculatedDeadlineBlock)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get sign message")
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	return &signMessage, &messageHash, nil
}

func (uc *EventRegistrationUsecase) CheckEventAttendeeAndMaxAttendeeCount(ctx context.Context, client *ethclient.Client, eventContractInstance *eventContract.Event, currentUser *auth.JwtClaims) (*int, *int, error) {
	// Check attendee count
	attendeeCount, err := eventContractInstance.CurrentSeatsCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if attendeeCount == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("attendee count is nil"))
	}
	maxAttendeeCount, err := eventContractInstance.SeatsCount(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if maxAttendeeCount == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("max attendee count is nil"))
	}

	count := int(attendeeCount.Int64())
	maxCount := int(maxAttendeeCount.Int64())

	return &count, &maxCount, nil
}

type CheckRegistrationEligibilityParams struct {
	EventPassword *string
}

// Performs checks if the user is able to join the event
func (uc *EventRegistrationUsecase) CheckRegistrationEligibility(ctx context.Context, client *ethclient.Client, entityEventContract *entity.EventContract, currentUser *auth.JwtClaims, params CheckRegistrationEligibilityParams) (bool, error) {
	if currentUser == nil {
		return false, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// if password is provided, assume it's a password based registration
	if params.EventPassword != nil {
		config, err := uc.EventRegistrationConfigurationDg.GetEventRegistrationConfigPasswordByEventId(ctx, entityEventContract.EventID)
		if err != nil {
			return false, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		if config == nil || config.RegistrationPassword == nil {
			return false, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("registration password not found"))
		}

		match, err := hashutils.CompareHash(*params.EventPassword, *config.RegistrationPassword)
		if err != nil {
			return false, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if !match {
			return false, customerror.NewWithPreset(&customerror.ErrUnauthorized, errors.New("invalid registration password"))
		}
		return true, nil
	}

	// if not, assume it's an invitation based registration
	invitation, _, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, entityEventContract.EventID, currentUser.UserId, currentUser.Email, &currentUser.WalletAddress)
	if err != nil {
		return false, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if invitation == nil {
		return false, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("invitation not found"))
	}

	return true, nil
}

func (uc *EventRegistrationUsecase) JoinEventWithPin(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, eventId uuid.UUID, eligibilityProof CheckRegistrationEligibilityParams, joinPayload JoinEventPayload, password string) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	entityEventContract, err := uc.EventContractDg.GetEventContractByEventID(ctx, eventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if entityEventContract == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("event contract not found"))
	}

	isEligible, err := uc.CheckRegistrationEligibility(ctx, client, entityEventContract, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check registration eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user is not eligible to join the event"))
	}

	signature, _, err := uc.AuthUsecase.SecuredSignActionForManagedUser(ctx, currentUser, password, common.HexToAddress(currentUser.WalletAddress), common.HexToAddress(entityEventContract.EventContractAddress))
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature)
}

func (uc *EventRegistrationUsecase) JoinEventWithSignature(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, eventId uuid.UUID, eligibilityProof CheckRegistrationEligibilityParams, joinPayload JoinEventPayload, signature []byte, signMessage string) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	entityEventContract, err := uc.EventContractDg.GetEventContractByEventID(ctx, eventId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if entityEventContract == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("event contract not found"))
	}

	// validate original sign message
	isValid, err := cyptoutils.ValidateSignMessage(signMessage, common.HexToAddress(currentUser.WalletAddress), common.HexToAddress(entityEventContract.EventContractAddress), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate sign message")
	}
	if !isValid {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is not valid"))
	}
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	// check if the signature matches the sign message or not
	isValidHash, err := cyptoutils.VerifySignatureByDigest(common.HexToAddress(currentUser.WalletAddress), messageHash, signature)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify signature by digest")
	}
	if !isValidHash {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature does not match the sign message"))
	}

	// check if the user is eligible to join the event
	isEligible, err := uc.CheckRegistrationEligibility(ctx, client, entityEventContract, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check registration eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user is not eligible to join the event"))
	}

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature)
}

func (uc *EventRegistrationUsecase) joinEvent(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, entityEventContract entity.EventContract, joinEventPayload JoinEventPayload, signature []byte) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	eventContractInstance, err := eventContract.NewEvent(common.HexToAddress(entityEventContract.EventContractAddress), client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Validation
	// Check if the user has already joined the event
	participants, err := eventContractInstance.GetParticipants(nil)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	hasJoinedOnChain := false
	for _, participant := range participants {
		if participant.Cmp(common.HexToAddress(currentUser.WalletAddress)) == 0 {
			hasJoinedOnChain = true
			break
		}
	}

	// If the data is not on chain, check max count
	if !hasJoinedOnChain {
		// Check attendee count
		attendeeCount, maxAttendeeCount, err := uc.CheckEventAttendeeAndMaxAttendeeCount(ctx, client, eventContractInstance, currentUser)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if *attendeeCount >= *maxAttendeeCount {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New(string(JoinEventUserErrorEventAttendeeFull)))
		}
	}

	// If the data is not on chain, add participant on blockchain
	if !hasJoinedOnChain {
		transactor, err := cyptoutils.GetKeyedTransactor()
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		signMessage, _, err := uc.GetJoinEventSignMessage(ctx, client, common.HexToAddress(currentUser.WalletAddress), *currentUser, common.HexToAddress(entityEventContract.EventContractAddress))
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if signMessage == nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("sign message not found"))
		}
		messageHash := cyptoutils.HashEthereumMessage(*signMessage)
		// Compare the message hash with the signature
		isValidHash, err := cyptoutils.VerifySignatureByDigest(common.HexToAddress(currentUser.WalletAddress), messageHash, signature)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if !isValidHash {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature is not valid"))
		}

		// Add participant on blockchain
		tx, err := eventContractInstance.AddParticipant(transactor, common.HexToAddress(currentUser.WalletAddress), messageHash.String(), signature)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New(string(JoinEventUserErrorEventAttendeeFull)))
		}
	}

	// If hasJoined, check if the user is also in the event attendee or not
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, entityEventContract.EventID, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if eventAttendee != nil {
		return eventAttendee, nil
	}

	// Save user joined event to database
	eventAttendee, err = uc.EventAttendeeDg.AddParticipant(ctx, datagateway.AddParticipantParameters{
		EventId:         entityEventContract.EventID,
		CredentialId:    currentUser.UserId,
		ContractAddress: entityEventContract.EventContractAddress,
		// TODO: Will work on joinning on request later
		IsParticipantAccepted: true,
		FirstName:             joinEventPayload.FirstName,
		LastName:              joinEventPayload.LastName,
		Email:                 joinEventPayload.Email,
		Bio:                   joinEventPayload.Bio,
		PhoneNumber:           joinEventPayload.PhoneNumber,
		AcademicInstitution:   joinEventPayload.AcademicInstitution,
		AcademicEmail:         joinEventPayload.AcademicEmail,
		Address:               joinEventPayload.Address,
	})
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return eventAttendee, nil
}
