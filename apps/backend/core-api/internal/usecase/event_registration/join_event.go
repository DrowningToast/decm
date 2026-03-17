package event_registration

import (
	"apps/backend/common/customerror"
	"apps/backend/common/hashutils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"

	eventcontract_datagateway "apps/backend/core-api/internal/datagateway/onchain/event_contract"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
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
func (uc *EventRegistrationUsecase) GetJoinEventSignMessage(ctx context.Context, walletAddress common.Address, currentUser auth.JwtClaims, eventContractAddress common.Address, deadlineBlock *uint64) (*string, *common.Hash, error) {
	// Validation
	if deadlineBlock == nil {
		calculatedDeadlineBlock, err := uc.BlockchainClientDg.GetCalculatedDeadlineBlock(ctx)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get calculated deadline block")
		}
		deadlineBlock = &calculatedDeadlineBlock
	}

	// Use wallet address from JWT (what the user will sign with via wallet extension)
	// This must match what JoinEventWithSignature uses for validation
	participantAddress := common.HexToAddress(currentUser.WalletAddress)

	signMessage, err := cyptoutils.GetSignMessage(participantAddress, eventContractAddress, *deadlineBlock)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get sign message")
	}

	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	return &signMessage, &messageHash, nil
}

func (uc *EventRegistrationUsecase) CheckEventAttendeeAndMaxAttendeeCount(ctx context.Context, eventContractInstance eventcontract_datagateway.EventContractDataGateway, currentUser *auth.JwtClaims) (*int, *int, error) {
	// Check attendee count
	attendeeCount, err := eventContractInstance.GetCurrentSeatsCount(ctx)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if attendeeCount == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("attendee count is nil"))
	}
	maxAttendeeCount, err := eventContractInstance.GetMaxSeatsCount(ctx)
	if err != nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if maxAttendeeCount == nil {
		return nil, nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("max attendee count is nil"))
	}

	count := int(*attendeeCount)
	maxCount := int(*maxAttendeeCount)

	return &count, &maxCount, nil
}

type CheckRegistrationEligibilityParams struct {
	EventPassword *string
}

// Performs checks if the user is able to join the event
func (uc *EventRegistrationUsecase) CheckRegistrationEligibility(ctx context.Context, entityEventContract *entity.EventContract, currentUser *auth.JwtClaims, params CheckRegistrationEligibilityParams) (bool, error) {
	if currentUser == nil {
		return false, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}

	// if password is provided, assume it's a password based registration
	if params.EventPassword != nil {
		config, err := uc.EventRegistrationConfigurationDg.GetEventRegistrationConfigPasswordByEventId(ctx, entityEventContract.EventId)
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
	invitation, _, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, entityEventContract.EventId, currentUser.UserId, currentUser.Email, &currentUser.WalletAddress)
	if err != nil {
		return false, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if invitation == nil {
		return false, customerror.NewWithPreset(&customerror.ErrInternalServer, errors.New("invitation not found"))
	}

	return true, nil
}

func (uc *EventRegistrationUsecase) JoinEventWithPin(ctx context.Context, currentUser *auth.JwtClaims, eventId uuid.UUID, eligibilityProof CheckRegistrationEligibilityParams, joinPayload JoinEventPayload, password string) (*entity.EventAttendee, error) {
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

	isEligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check registration eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user is not eligible to join the event"))
	}

	// Get credential to decrypt private key
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if credential == nil || credential.EncryptedPrivateKey == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("encrypted private key not found"))
	}

	// Decrypt private key to get participant address (like fuck_join_event.go)
	privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, password)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("invalid password or failed to decrypt private key"))
	}

	calculatedDeadlineBlock, err := uc.BlockchainClientDg.GetCalculatedDeadlineBlock(ctx)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Use participantAddress from private key (like fuck_join_event.go)
	signMessage, err := cyptoutils.GetSignMessage(*participantAddress, common.HexToAddress(entityEventContract.EventContractAddress), calculatedDeadlineBlock)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Sign directly with private key (like fuck_join_event.go)
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return uc.queueEventJoin(ctx, currentUser, *entityEventContract, joinPayload, signature, signMessage, participantAddress)
}

func (uc *EventRegistrationUsecase) JoinEventWithSignature(ctx context.Context, currentUser *auth.JwtClaims, eventId uuid.UUID, eligibilityProof CheckRegistrationEligibilityParams, joinPayload JoinEventPayload, signature []byte, signMessage string) (*entity.EventAttendee, error) {
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

	// Use wallet address from JWT (what the user signed with via wallet extension)
	// This must match what GetJoinEventSignMessage used to generate the sign message
	participantAddress := common.HexToAddress(currentUser.WalletAddress)

	// validate original sign message
	deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract deadline block from sign message")
	}
	isValid, err := cyptoutils.ValidateSignMessage(signMessage, participantAddress, common.HexToAddress(entityEventContract.EventContractAddress), deadlineBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate sign message")
	}
	if !isValid {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("sign message is not valid"))
	}
	messageHash := cyptoutils.HashEthereumMessage(signMessage)

	// CRITICAL: Make a copy of signature before verification to prevent mutation
	// VerifySignatureByDigest changes v from 27/28 to 0/1, but contracts expect 27/28
	signatureCopy := make([]byte, len(signature))
	copy(signatureCopy, signature)

	// check if the signature matches the sign message or not (using copy)
	isValidHash, err := cyptoutils.VerifySignatureByDigest(participantAddress, messageHash, signatureCopy)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify signature by digest")
	}
	if !isValidHash {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature does not match the sign message"))
	}

	// check if the user is eligible to join the event
	isEligible, err := uc.CheckRegistrationEligibility(ctx, entityEventContract, currentUser, eligibilityProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check registration eligibility")
	}
	if !isEligible {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user is not eligible to join the event"))
	}

	return uc.queueEventJoin(ctx, currentUser, *entityEventContract, joinPayload, signature, signMessage, &participantAddress)
}
