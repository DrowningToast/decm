package event_registration

import (
	"context"
	"encoding/hex"
	"strings"

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
func (uc *EventRegistrationUsecase) GetJoinEventSignMessage(ctx context.Context, client *ethclient.Client, walletAddress common.Address, currentUser auth.JwtClaims, eventContractAddress common.Address, deadlineBlock *uint64) (*string, *ethcommon.Hash, error) {
	// Validation
	if deadlineBlock == nil {
		calculatedDeadlineBlock, err := cyptoutils.GetCalculatedDeadlineBlock(client)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get calculated deadline block")
		}
		deadlineBlock = &calculatedDeadlineBlock
	}

	signMessage, err := cyptoutils.GetSignMessage(common.HexToAddress(currentUser.WalletAddress), eventContractAddress, *deadlineBlock)
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

	signature, _, deadlineBlock, err := uc.AuthUsecase.SecuredSignActionForManagedUser(ctx, currentUser, password, common.HexToAddress(currentUser.WalletAddress), common.HexToAddress(entityEventContract.EventContractAddress), nil)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature, deadlineBlock)
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
	deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract deadline block from sign message")
	}
	isValid, err := cyptoutils.ValidateSignMessage(signMessage, common.HexToAddress(currentUser.WalletAddress), common.HexToAddress(entityEventContract.EventContractAddress), deadlineBlock)
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

	return uc.joinEvent(ctx, client, currentUser, *entityEventContract, joinPayload, signature, deadlineBlock)
}

func (uc *EventRegistrationUsecase) joinEvent(ctx context.Context, client *ethclient.Client, currentUser *auth.JwtClaims, entityEventContract entity.EventContract, joinEventPayload JoinEventPayload, signature []byte, deadlineBlock *uint64) (*entity.EventAttendee, error) {
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

		signMessage, messageHashPtr, err := uc.GetJoinEventSignMessage(ctx, client, common.HexToAddress(currentUser.WalletAddress), *currentUser, common.HexToAddress(entityEventContract.EventContractAddress), deadlineBlock)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if signMessage == nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.New("sign message not found"))
		}
		messageHash := *messageHashPtr

		// Compare the message hash with the signature
		isValidHash, err := cyptoutils.VerifySignatureByDigest(common.HexToAddress(currentUser.WalletAddress), messageHash, signature)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if !isValidHash {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("signature is not valid"))
		}

		// Debug logging before calling AddParticipant
		participantAddr := common.HexToAddress(currentUser.WalletAddress)
		println("=== DEBUG AddParticipant ===")
		println("Participant Address:", participantAddr.Hex())
		println("Sign Message (RAW):", *signMessage)
		println("Message Hash:", messageHash.String())
		println("Signature Length:", len(signature))
		println("Contract Address:", entityEventContract.EventContractAddress)
		println("Transactor Address:", transactor.From.Hex())

		// Check current state before adding
		currentCount, err := eventContractInstance.CurrentSeatsCount(&bind.CallOpts{Context: ctx})
		if err == nil {
			println("Current Seats:", currentCount.String())
		}
		maxCount, err := eventContractInstance.SeatsCount(&bind.CallOpts{Context: ctx})
		if err == nil {
			println("Max Seats:", maxCount.String())
		}
		println("========================")

		// Pre-flight check: Try to simulate the call first to get detailed error
		println("")
		println("🔍 Pre-flight check: Simulating contract call...")
		callOpts := &bind.CallOpts{
			Context: ctx,
			From:    transactor.From,
		}

		// Try calling GetParticipants to verify contract is accessible
		participants, err := eventContractInstance.GetParticipants(callOpts)
		if err != nil {
			println("⚠️  Warning: Cannot read contract state:", err.Error())
		} else {
			println("✅ Contract is accessible. Current participants:", len(participants))
			// Check if participant already exists on-chain (double-check)
			for _, p := range participants {
				if p.Hex() == participantAddr.Hex() {
					println("⚠️  WARNING: Participant already exists on-chain!")
					break
				}
			}
		}

		println("🚀 Attempting to submit transaction...")
		println("   Passing RAW sign message to contract (not hash)")
		println("")

		// CRITICAL: Pass the RAW sign message, NOT the hash!
		// The smart contract will hash it itself using toEthSignedMessageHash
		tx, err := eventContractInstance.AddParticipant(transactor, participantAddr, *signMessage, signature)
		if err != nil {
			// Enhanced error logging - transaction was rejected before submission
			println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			println("ERROR: AddParticipant transaction REJECTED")
			println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			println("Full Error:", err.Error())
			println("")

			// Try to extract revert reason from error message
			errStr := err.Error()

			// Common patterns in go-ethereum errors
			if strings.Contains(errStr, "execution reverted:") {
				parts := strings.SplitN(errStr, "execution reverted:", 2)
				if len(parts) == 2 {
					revertReason := strings.TrimSpace(parts[1])
					println("🔴 REVERT REASON:", revertReason)
				}
			} else if strings.Contains(errStr, "execution reverted") {
				println("🔴 Transaction would revert (no specific reason provided)")
			}

			// Check for common issues
			println("")
			println("Possible causes:")
			if strings.Contains(errStr, "insufficient funds") {
				println("  ❌ Transactor wallet has insufficient funds for gas")
			}
			if strings.Contains(errStr, "nonce") {
				println("  ❌ Nonce issue - transaction ordering problem")
			}
			if strings.Contains(errStr, "gas") && strings.Contains(errStr, "intrinsic") {
				println("  ❌ Gas limit too low")
			}
			if strings.Contains(errStr, "reverted") {
				println("  ❌ Smart contract rejected the call during gas estimation")
				println("     This means one of the contract's require/revert statements failed:")
				println("     - Event__SeatsCountReached: Event is full")
				println("     - Event__ParticipantIsAlreadyJoined: Already registered")
				println("     - Event__AddressCannotBeZero: Invalid address")
				println("     - Access control: Transactor doesn't have HOST/ADMIN role")
				println("     - Invalid signature: Signature verification failed")
			}

			println("")
			println("Transaction Details:")
			println("  Participant:", participantAddr.Hex())
			println("  Contract:", entityEventContract.EventContractAddress)
			println("  Transactor:", transactor.From.Hex())
			println("  Message Hash:", messageHash.String())
			println("  Signature:", hex.EncodeToString(signature))
			println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to add participant to blockchain: wallet=%s, contract=%s", participantAddr.Hex(), entityEventContract.EventContractAddress))
		}

		println("Transaction submitted:", tx.Hash().Hex())

		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			println("ERROR: WaitMined failed:", err.Error())
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
		}

		println("Transaction mined. Status:", receipt.Status, "Gas Used:", receipt.GasUsed)

		if receipt.Status != types.ReceiptStatusSuccessful {
			// Try to get revert reason
			revertReason, err := cyptoutils.GetRevertReason(ctx, client, tx, receipt)
			if err != nil {
				println("Failed to get revert reason:", err.Error())
				revertReason = "failed to decode revert reason: " + err.Error()
			}
			println("Revert Reason:", revertReason)

			if len(receipt.Logs) > 0 {
				println("Transaction logs count:", len(receipt.Logs))
				for i, log := range receipt.Logs {
					println("Log", i, "- Topics:", len(log.Topics))
				}
			}
			return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d): %s", tx.Hash().Hex(), receipt.GasUsed, revertReason))
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
