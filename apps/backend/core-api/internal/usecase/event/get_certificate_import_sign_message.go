package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/common/utils"
	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
)

type GetCertificateImportSignMessageResponse struct {
	SignMessage             string `json:"sign_message"`
	EventCertificateAddress string `json:"event_certificate_address"`
}

// GetCertificateImportSignMessage deploys the certificate contract if needed,
// computes the receiver hashes, and returns the sign message JSON that the
// host must sign (used by BYOK wallet users before calling ImportCertificateReceivers).
func (uc *EventUsecase) GetCertificateImportSignMessage(
	ctx context.Context,
	eventID uuid.UUID,
	requests []ImportCertificateReceiversRequest,
	currentUser *auth.JwtClaims,
) (*GetCertificateImportSignMessageResponse, error) {
	// 1. Validate requests
	if err := validateReceiverRequests(requests); err != nil {
		return nil, err
	}

	// 2. Authorisation
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}
	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}

	// 2. Event exists + ownership
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}
	if event.OwnerCredentialId != currentUser.UserId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("only the event owner can import certificate receivers"))
	}

	// 3. Get event contract
	eventContract, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if eventContract == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event contract not found"))
	}

	// 4. Deploy certificate contract if not already deployed
	eventCertificateAddressStr := ""
	if eventContract.CertificateContractAddress == nil {
		transactor, err := uc.BlockchainClientDg.GetTransactOpts(ctx)
		if err != nil {
			return nil, err
		}
		eventCertificateAddress, err := uc.deployCertificateContract(
			ctx,
			transactor,
			common.HexToAddress(eventContract.AccessManagerContractAddress),
			common.HexToAddress(eventContract.EventContractAddress),
		)
		if err != nil {
			return nil, err
		}
		eventCertificateAddressStr = eventCertificateAddress.Hex()

		// Persist the newly deployed address so the import step can re-use it
		_, err = uc.EventContractDataGateway.UpdateEventContract(ctx, generated.UpdateEventContractParams{
			EventID:                      eventContract.EventId,
			AccessManagerContractAddress: eventContract.AccessManagerContractAddress,
			EventContractAddress:         eventContract.EventContractAddress,
			TicketContractAddress:        pgmapper.StringPtrToPgText(eventContract.TicketContractAddress),
			CertificateContractAddress:   pgmapper.StringPtrToPgText(&eventCertificateAddressStr),
		})
		if err != nil {
			return nil, err
		}
	} else {
		eventCertificateAddressStr = *eventContract.CertificateContractAddress
	}

	// 5. Compute receiver hashes and build sign message JSON
	receivers := buildReceiverHashes(requests)

	signMessage := SignMessage{
		EventContractAddress: eventCertificateAddressStr,
		Receivers:            receivers,
	}
	signMessageJSONBytes, err := json.Marshal(signMessage)
	if err != nil {
		return nil, err
	}

	return &GetCertificateImportSignMessageResponse{
		SignMessage:             string(signMessageJSONBytes),
		EventCertificateAddress: eventCertificateAddressStr,
	}, nil
}

// buildReceiverHashes computes the keccak256 hex hash for each receiver —
// the same algorithm used in ImportCertificateReceivers.
func buildReceiverHashes(requests []ImportCertificateReceiversRequest) []string {
	hashes := make([]string, 0, len(requests))
	for _, req := range requests {
		firstName := utils.DerefOrEmpty(req.FirstName)
		lastName := utils.DerefOrEmpty(req.LastName)
		academicInstitution := utils.DerefOrEmpty(req.AcademicInstitution)
		certificateTitle := utils.DerefOrEmpty(req.CertificateTitle)
		certificateSubtitle := utils.DerefOrEmpty(req.CertificateSubtitle)

		name := fmt.Sprintf("%s %s", firstName, lastName)
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)
		hash := cyptoutils.HashMessage(csvValue)
		hashes = append(hashes, hexutil.Encode(hash))
	}
	return hashes
}
