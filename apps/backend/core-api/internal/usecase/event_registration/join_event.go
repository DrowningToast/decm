package event_registration

import (
	"context"
	"time"

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
// CRITICAL: walletAddress parameter should be the address derived from the user's private key,
// NOT from the JWT claims, to ensure cryptographic consistency
func (uc *EventRegistrationUsecase) GetJoinEventSignMessage(ctx context.Context, client *ethclient.Client, walletAddress common.Address, currentUser auth.JwtClaims, eventContractAddress common.Address, deadlineBlock *uint64) (*string, *ethcommon.Hash, error) {
	// Validation
	if deadlineBlock == nil {
		calculatedDeadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get calculated deadline block")
		}
		deadlineBlock = &calculatedDeadlineBlock
	}

	// Use the provided walletAddress parameter (which should be derived from private key)
	// instead of currentUser.WalletAddress (which comes from JWT claims)
	signMessage, err := cyptoutils.GetSignMessage(walletAddress, eventContractAddress, *deadlineBlock)
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

	// CRITICAL FIX: Get the actual wallet address derived from the private key
	// instead of using currentUser.WalletAddress from JWT claims
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if credential == nil || credential.EncryptedPrivateKey == nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("encrypted private key not found"))
	}

	// Decrypt to get the address that cryptographically matches the private key
	privateKey, participantAddress, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, password)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, errors.New("invalid password or failed to decrypt private key"))
	}

	deadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Use the derived address from the private key, NOT the JWT address
	signMessage, err := cyptoutils.GetSignMessage(*participantAddress, common.HexToAddress(entityEventContract.EventContractAddress), deadlineBlock)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// Sign the message directly with the decrypted private key
	messageHash := cyptoutils.HashEthereumMessage(signMessage)
	signature, err := cyptoutils.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature, signMessage, participantAddress)
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

	// CRITICAL FIX: Get the actual wallet address derived from the private key
	// The signature verification must use the address that actually signed the message
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if credential == nil || credential.WalletAddress == "" {
		return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("wallet address not found"))
	}

	// Use the wallet address from the credential (which is derived from the private key)
	participantAddress := common.HexToAddress(credential.WalletAddress)

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
	// check if the signature matches the sign message or not
	isValidHash, err := cyptoutils.VerifySignatureByDigest(participantAddress, messageHash, signature)
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

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature, signMessage, &participantAddress)
}

func (uc *EventRegistrationUsecase) joinEvent(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, entityEventContract entity.EventContract, joinEventPayload JoinEventPayload, signature []byte, signMessage string, participantAddress *common.Address) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if participantAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant address is required"))
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
		if participant.Cmp(*participantAddress) == 0 {
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

		// Use the original signMessage that was used to create the signature
		messageHash := cyptoutils.HashEthereumMessage(signMessage)

		// Compare the message hash with the signature
		// CRITICAL: Use the actual participant address derived from the private key
		isValidHash, err := cyptoutils.VerifySignatureByDigest(*participantAddress, messageHash, signature)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if !isValidHash {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature is not valid"))
		}

		// Check current state before adding
		currentCountPtr, err := eventContractInstance.CurrentSeatsCount(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}

		maxCountPtr, err := eventContractInstance.SeatsCount(&bind.CallOpts{Context: ctx})
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if currentCountPtr == nil || maxCountPtr == nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("current count or max count is nil"))
		}

		currentCount := int(currentCountPtr.Int64())
		maxCount := int(maxCountPtr.Int64())
		if currentCount >= maxCount {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New(string(JoinEventUserErrorEventAttendeeFull)))
		}

		// Pre-flight check: Try to simulate the call first to get detailed error
		callOpts := &bind.CallOpts{
			Context: ctx,
			From:    transactor.From,
		}

		// Try calling GetParticipants to verify contract is accessible
		_, err = eventContractInstance.GetParticipants(callOpts)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrap(err, "failed to get participants"))
		}

		// CRITICAL: Pass the RAW sign message, NOT the hash!
		// The smart contract will hash it itself using toEthSignedMessageHash
		tx, err := eventContractInstance.AddParticipant(transactor, *participantAddress, signMessage, signature)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to add participant to blockchain: wallet=%s, contract=%s", participantAddress.Hex(), entityEventContract.EventContractAddress))
		}

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d): %s", tx.Hash().Hex(), receipt.GasUsed))
		}
	}

	// If hasJoined, check if the user is also in the event attendee or not
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, entityEventContract.EventID, currentUser.UserId)
	if err != nil {
		var customError *customerror.Err
		if errors.As(err, &customError) {
			if *customError.Code != customerror.ErrNotFound.Code {
				return nil, errors.Wrap(err, "failed to get event attendee by event id and credential id")
			}
		} else {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
	}
	hasJoinedOnDatabase := eventAttendee != nil
	if hasJoinedOnChain && hasJoinedOnDatabase {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user has already joined the event"))
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

	// Mark invitation as accepted if it exists (for invitation-based registration)
	// Only do this silently - if invitation doesn't exist, that's fine (password-based registration)
	invitation, _, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, entityEventContract.EventID, currentUser.UserId, currentUser.Email, &currentUser.WalletAddress)
	var customError *customerror.Err
	// If the invitation is not found, that's fine
	if errors.As(err, &customError) {
		if *customError.Code != customerror.ErrNotFound.Code {
			return nil, errors.Wrap(err, "failed to get event registration invitation by event id and credential")
		}
	}
	// If the invitation is found, mark it as accepted with current timestamp
	if err == nil && invitation != nil {
		now := time.Now()
		_, err = uc.EventRegistrationInvitationDg.UpdateEventRegistrationInvitationAcceptedStatus(ctx, invitation.Id, &now)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update event registration invitation accepted status")
		}
	}

	return eventAttendee, nil
}
