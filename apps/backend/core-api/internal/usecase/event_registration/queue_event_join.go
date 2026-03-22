package event_registration

import (
	"apps/backend/common/customerror"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	event_datagateway "apps/backend/core-api/internal/datagateway/offchain/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (uc *EventRegistrationUsecase) queueEventJoin(ctx context.Context, currentUser *auth.JwtClaims, entityEventContract entity.EventContract, joinEventPayload JoinEventPayload, signature []byte, signMessage string, participantAddress *common.Address) (*entity.EventAttendee, error) {
	if currentUser == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthenticated, errors.New("user is not authenticated"))
	}
	if participantAddress == nil {
		return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("participant address is required"))
	}

	eventContractInstance, err := uc.EventContractFactoryDg.GetContract(common.HexToAddress(entityEventContract.EventContractAddress))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event contract by address")
	}

	// Validation
	// Check if the user has already joined the event
	participants, err := eventContractInstance.GetParticipants(ctx)
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

	// If the data is not on chain, check max count.
	// Use DB attendee count as the source of truth since the on-chain
	// AddParticipant write is now deferred to the background worker.
	if !hasJoinedOnChain {
		dbAttendees, err := uc.EventAttendeeDg.ListEventAttendeesByEventID(ctx, entityEventContract.EventId)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		_, maxAttendeeCount, err := uc.CheckEventAttendeeAndMaxAttendeeCount(ctx, eventContractInstance, currentUser)
		if err != nil {
			return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		}
		if len(dbAttendees) >= *maxAttendeeCount {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New(string(JoinEventUserErrorEventAttendeeFull)))
		}
	}

	// If the data is not on chain, add participant on blockchain
	if !hasJoinedOnChain {
		// transactor, err := cyptoutils.GetKeyedTransactor(ctx, client)
		// if err != nil {
		// 	return nil, customerror.Parse(&customerror.ErrInternalServer, err)
		// }

		// // messageHash := cyptoutils.HashEthereumMessage(signMessage)
		// participantAddr := *participantAddress

		// // Log transaction details before submission
		// slog.InfoContext(ctx, "Adding participant to blockchain contract",
		// 	slog.String("participant_address", participantAddr.Hex()),
		// 	slog.String("contract_address", entityEventContract.EventContractAddress),
		// 	slog.String("transactor_address", transactor.From.Hex()),
		// 	slog.String("event_id", entityEventContract.EventID.String()),
		// 	slog.Int("signature_length", len(signature)),
		// )

		if len(signature) == 65 {
			slog.DebugContext(ctx, "Signature details",
				slog.Uint64("recovery_id", uint64(signature[64])),
			)
		}

		// DEPRECATED: Will work on blockchain joinning on request later
		// CRITICAL: Pass the RAW sign message, NOT the hash!
		// The smart contract will hash it itself using toEthSignedMessageHash
		// tx, err := eventContractInstance.AddParticipant(transactor, *participantAddress, signMessage, signature)
		// if err != nil {
		// 	// Extract and log revert reason if available
		// 	errStr := err.Error()
		// 	var revertReason string
		// 	if strings.Contains(errStr, "execution reverted:") {
		// 		parts := strings.SplitN(errStr, "execution reverted:", 2)
		// 		if len(parts) == 2 {
		// 			revertReason = strings.TrimSpace(parts[1])
		// 		}
		// 	}

		// 	slog.ErrorContext(ctx, "Failed to add participant to blockchain",
		// 		slog.String("error", err.Error()),
		// 		slog.String("revert_reason", revertReason),
		// 		slog.String("participant_address", participantAddr.Hex()),
		// 		slog.String("contract_address", entityEventContract.EventContractAddress),
		// 		slog.String("transactor_address", transactor.From.Hex()),
		// 		slog.String("message_hash", messageHash.String()),
		// 	)

		// 	return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "failed to add participant to blockchain: wallet=%s, contract=%s", participantAddr.Hex(), entityEventContract.EventContractAddress))
		// }

		// slog.InfoContext(ctx, "Transaction submitted to blockchain",
		// 	slog.String("tx_hash", tx.Hash().Hex()),
		// 	slog.String("participant_address", participantAddr.Hex()),
		// )

		// receipt, err := bind.WaitMined(ctx, client, tx)
		// if err != nil {
		// 	slog.ErrorContext(ctx, "Transaction mining failed",
		// 		slog.String("error", err.Error()),
		// 		slog.String("tx_hash", tx.Hash().Hex()),
		// 	)
		// 	return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Wrapf(err, "transaction mining failed: tx=%s", tx.Hash().Hex()))
		// }

		// slog.InfoContext(ctx, "Transaction mined successfully",
		// 	slog.String("tx_hash", tx.Hash().Hex()),
		// 	slog.Uint64("status", receipt.Status),
		// 	slog.Uint64("gas_used", receipt.GasUsed),
		// )

		// if receipt.Status != types.ReceiptStatusSuccessful {
		// 	// Try to get revert reason
		// 	revertReason, err := cyptoutils.GetRevertReason(ctx, client, tx, receipt)
		// 	if err != nil {
		// 		revertReason = "failed to decode revert reason: " + err.Error()
		// 	}

		// 	slog.ErrorContext(ctx, "Transaction reverted on blockchain",
		// 		slog.String("tx_hash", tx.Hash().Hex()),
		// 		slog.Uint64("gas_used", receipt.GasUsed),
		// 		slog.String("revert_reason", revertReason),
		// 		slog.Int("logs_count", len(receipt.Logs)),
		// 	)

		// 	return nil, customerror.Parse(&customerror.ErrInternalServer, errors.Errorf("transaction reverted (tx=%s, gas=%d): %s", tx.Hash().Hex(), receipt.GasUsed, revertReason))
		// }
	}

	// If hasJoined, check if the user is also in the event attendee or not
	eventAttendee, err := uc.EventAttendeeDg.GetEventAttendeeByEventIdAndCredentialId(ctx, entityEventContract.EventId, currentUser.UserId)
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
	// Use DB as the source of truth for the already-joined check.
	// hasJoinedOnChain is always false while the blockchain write is deferred
	// to the background worker, so checking both conditions would make this
	// guard permanently unreachable.
	if hasJoinedOnDatabase {
		if eventAttendee.UserSignatureID == nil {
			// Old version of EventAttendee, has no stored signature.
			// already joined
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user has already joined the event"))
		}
		userSignature, err := uc.UserSignatureDg.GetUserSignatureByID(ctx, *eventAttendee.UserSignatureID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get user signature by id")
		}
		if userSignature == nil {
			return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("user signature not found"))
		}
		// Already broadcasted, returns error
		if userSignature.BroadcastedAt != nil {
			return nil, customerror.Parse(&customerror.ErrInvalidArgument, errors.New("user has already joined the event"))
		}
		// Check whether the queued signature has expired.
		isExpired := userSignature.MarkAsExpiredAt != nil ||
			(userSignature.EstimatedDeadline != nil && time.Now().After(*userSignature.EstimatedDeadline))
		if !isExpired {
			// Not broadcasted yet and not expired — already queued, nothing to do.
			return nil, customerror.Parse(&customerror.ErrAlreadyPerformed, errors.New("join event already queued and pending broadcast"))
		}
		// Signature expired — remove the stale attendee record and signature so
		// the user can re-join with a fresh signature below.
		if deleteErr := uc.EventAttendeeDg.DeleteEventAttendeeById(ctx, eventAttendee.Id); deleteErr != nil {
			return nil, errors.Wrap(deleteErr, "failed to delete expired event attendee for re-join")
		}
		if deleteErr := uc.UserSignatureDg.DeleteUserSignature(ctx, userSignature.Id); deleteErr != nil {
			slog.ErrorContext(ctx, "Failed to delete expired user signature during re-join", slog.String("error", deleteErr.Error()))
		}
		// Fall through — the code below will create a fresh attendee + signature.
	}

	// Check the existence of invitation first
	invitation, _, err := uc.EventRegistrationInvitationDg.GetEventRegistrationInvitationByEventIDAndCredential(ctx, entityEventContract.EventId, currentUser.UserId, currentUser.Email, &currentUser.WalletAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event registration invitation by event id and credential")
	}
	if invitation == nil {
		return nil, customerror.Parse(&customerror.ErrForbidden, errors.New("invitation not found"))
	}

	// Extract deadline block from the sign message for storage
	deadlineBlockUint64, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract deadline block from sign message")
	}
	var deadlineBlock *int32
	var estimatedDeadline *time.Time
	if deadlineBlockUint64 != nil {
		db := int32(*deadlineBlockUint64)
		deadlineBlock = &db
		// Estimate deadline using blockchain client
		estimatedDeadline, err = uc.BlockchainClientDg.EstimateDeadlineTime(ctx, *deadlineBlockUint64)
		if err != nil {
			// Log but don't fail - estimation is not critical
			slog.WarnContext(ctx, "failed to estimate deadline time", "error", err)
		}
	}

	userSignature, err := uc.UserSignatureDg.CreateUserSignature(ctx, offchain_datagateway.CreateUserSignatureParameters{
		AuthenticationCredentialId: currentUser.UserId,
		SignMessage:                signMessage,
		Signature:                  hex.EncodeToString(signature),
		DeadlineBlock:              deadlineBlock,
		EstimatedDeadline:          estimatedDeadline,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user signature")
	}

	// Save user joined event to database
	eventAttendee, err = uc.EventAttendeeDg.AddParticipant(ctx, event_datagateway.AddParticipantParameters{
		EventId:               entityEventContract.EventId,
		CredentialId:          currentUser.UserId,
		ContractAddress:       entityEventContract.EventContractAddress,
		IsParticipantAccepted: true,
		FirstName:             joinEventPayload.FirstName,
		LastName:              joinEventPayload.LastName,
		Email:                 joinEventPayload.Email,
		Bio:                   joinEventPayload.Bio,
		PhoneNumber:           joinEventPayload.PhoneNumber,
		AcademicInstitution:   joinEventPayload.AcademicInstitution,
		AcademicEmail:         joinEventPayload.AcademicEmail,
		Address:               joinEventPayload.Address,
		UserSignatureId:       &userSignature.Id,
	})
	if err != nil {
		// rollback user signature
		if deleteErr := uc.UserSignatureDg.DeleteUserSignature(ctx, userSignature.Id); deleteErr != nil {
			slog.ErrorContext(ctx, "Failed to rollback user signature after AddParticipant failure", slog.String("error", deleteErr.Error()))
		}
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	// mark invitation as accepted
	now := time.Now()
	_, updateErr := uc.EventRegistrationInvitationDg.UpdateEventRegistrationInvitationAcceptedStatus(ctx, invitation.Id, &now)
	if updateErr != nil {
		// attempt to rollback both attendee and user signature
		errAttendee := uc.EventAttendeeDg.DeleteEventAttendeeById(ctx, eventAttendee.Id)
		errSignature := uc.UserSignatureDg.DeleteUserSignature(ctx, userSignature.Id)
		if errAttendee != nil || errSignature != nil {
			return nil, errors.Wrapf(updateErr, "rollback also failed: attendee=%v, signature=%v", errAttendee, errSignature)
		}
		return nil, errors.Wrap(updateErr, "rolled back attendee and signature after failing to mark invitation accepted")
	}

	return eventAttendee, nil
}
