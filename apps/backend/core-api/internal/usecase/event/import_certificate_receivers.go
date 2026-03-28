package event

import (
	"apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	"apps/backend/common/utils"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"
	"strings"

	eventdatagateway "apps/backend/core-api/internal/datagateway/offchain/event"

	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"

	eventCertificateContract "apps/backend/contracts/certificate"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

type ImportCertificateReceiversRequest struct {
	Email               *string `json:"email,omitempty"`
	WalletAddress       *string `json:"wallet_address,omitempty"`
	FirstName           *string `json:"first_name"`
	LastName            *string `json:"last_name"`
	AcademicInstitution *string `json:"academic_institution"`
	CertificateTitle    *string `json:"certificate_title"`
	CertificateSubtitle *string `json:"certificate_subtitle"`
}

// ImportCertificateReceiversOptions carries the host authentication method.
// Exactly one of HostPin (PIN-based) or (HostSignMessage + HostSignature) (wallet-based) must be set.
type ImportCertificateReceiversOptions struct {
	HostPin         *string // non-nil for PIN-based users
	HostSignMessage *string // non-nil for wallet-based (BYOK) users
	HostSignature   *string // non-nil for wallet-based (BYOK) users
}

type ImportCertificateReceiversResponse struct {
	EventID                 uuid.UUID                  `json:"event_id"`
	EventCertificateAddress string                     `json:"event_certificate_address"`
	Certificates            []*entity.EventCertificate `json:"certificates"`
}

type SignMessage struct {
	EventContractAddress string   `json:"eventContractAddress"`
	Receivers            []string `json:"receivers"`
}

// validateReceiverRequests enforces that every receiver has exactly one identifier
// (email XOR wallet address) and that wallet addresses are valid hex strings.
func validateReceiverRequests(requests []ImportCertificateReceiversRequest) error {
	for i, receiver := range requests {
		hasEmail := receiver.Email != nil && *receiver.Email != ""
		hasWallet := receiver.WalletAddress != nil && *receiver.WalletAddress != ""
		if hasEmail == hasWallet {
			return customerror.Parse(
				&customerror.ErrInvalidArgument,
				fmt.Errorf("receiver at index %d: must have exactly one of email or wallet_address, not both and not neither", i),
			)
		}
		if hasWallet && !common.IsHexAddress(*receiver.WalletAddress) {
			return customerror.Parse(
				&customerror.ErrInvalidArgument,
				fmt.Errorf("receiver at index %d: invalid wallet address format", i),
			)
		}
	}
	return nil
}

func (uc *EventUsecase) ImportCertificateReceivers(ctx context.Context, eventID uuid.UUID, requests []ImportCertificateReceiversRequest, opts ImportCertificateReceiversOptions, currentUser *auth.JwtClaims) (*ImportCertificateReceiversResponse, error) {
	// 0. Validate requests
	if err := validateReceiverRequests(requests); err != nil {
		return nil, err
	}

	// 1. Check if current user is authorized
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}

	// Determine auth path: PIN-based (non-BYOK) or wallet-based (BYOK)
	isByok := credential.EncryptedPrivateKey == nil
	if !isByok && utils.DerefOrEmpty(opts.HostPin) == "" {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("host_pin is required for non-BYOK users"))
	}
	if isByok && (opts.HostSignMessage == nil || opts.HostSignature == nil) {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("host_sign_message and host_signature are required for BYOK users"))
	}

	// 2. Check if event exists
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
	}

	// 2.5. Verify that current user is the event owner (only event owner can import certificates)
	if event.OwnerCredentialId != currentUser.UserId {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("only the event owner can import certificate receivers"))
	}

	// 3. Get eventContract from eventContracts table using eventID
	eventContract, err := uc.EventContractDataGateway.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if eventContract == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event contract not found"))
	}

	// Reset all event issuers' signing status
	err = uc.EventIssuerDataGateway.ResetAllEventIssuersSigningStatus(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Delete all existing certificates and their signatures for this event
	// This revokes the old receiver list before importing new ones.
	// This also removes old sign messages (stored in certificate signatures) which contained the old receiver list.
	// New sign messages will be generated below with the new receiver list.
	existingCertificates, err := uc.EventCertificateDataGateway.GetEventCertificatesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Get certificate config (same for all certificates in the event)
	certificateConfig, err := uc.EventCertificateConfigDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Delete all certificate signatures first (foreign key constraint)
	// Since signatures are now linked to config, we can get them directly
	signatures, err := uc.EventCertificateSignatureDataGateway.GetEventCertificateSignaturesByEventCertificateConfigID(ctx, certificateConfig.ID)
	if err != nil {
		return nil, err
	}

	for _, signature := range signatures {
		err = uc.EventCertificateSignatureDataGateway.DeleteEventCertificateSignature(ctx, signature.Id)
		if err != nil {
			return nil, err
		}
	}

	// Delete all existing certificates
	for _, certificate := range existingCertificates {
		err = uc.EventCertificateDataGateway.DeleteEventCertificate(ctx, certificate.Id)
		if err != nil {
			return nil, err
		}
	}

	eventCertificateAddressStr := ""

	if eventContract.CertificateContractAddress == nil {
		// 4. Deploy event certificate contract
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

		// 5. Persist newly deployed certificate contract address
		_, err = uc.EventContractDataGateway.UpdateEventContract(ctx, generated.UpdateEventContractParams{
			EventID:                      eventID,
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

	// 6. Save certificate data to event_certificates
	certificates := make([]*entity.EventCertificate, 0, len(requests))
	receivers := make([]string, 0, len(requests))

	for _, req := range requests {
		// Safely dereference pointer fields (use empty string if nil)
		firstName := utils.DerefOrEmpty(req.FirstName)
		lastName := utils.DerefOrEmpty(req.LastName)
		academicInstitution := utils.DerefOrEmpty(req.AcademicInstitution)
		certificateTitle := utils.DerefOrEmpty(req.CertificateTitle)
		certificateSubtitle := utils.DerefOrEmpty(req.CertificateSubtitle)

		// Combine first and last name
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Create CSV value for hash (used for blockchain verification)
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)
		receivers = append(receivers, encodedHash)

		// Create certificate - pass original pointers to preserve nil vs empty distinction
		certificate, err := uc.EventCertificateDataGateway.CreateEventCertificate(ctx, eventdatagateway.CreateEventCertificateParameters{
			EventID:                 eventID,
			ReceiverCredentialID:    nil,               // Will be set when receiver claims certificate
			ReceiverEmail:           req.Email,         // Set email from import for matching (may be nil for wallet-only)
			ReceiverWalletAddress:   req.WalletAddress, // Set wallet address for wallet-only receivers
			Name:                    stringPtrIfNotEmpty(name),
			AcademicInstitution:     req.AcademicInstitution,
			CertificateTitle:        req.CertificateTitle,
			CertificateSubtitle:     req.CertificateSubtitle,
			EventContractAddress:    eventContract.EventContractAddress,
			EventCertificateAddress: &eventCertificateAddressStr,
			CertificateTokenID:      nil, // Will be set when minted,
			Digest:                  &encodedHash,
		})
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
	}

	// 7. Create NEW sign_message with NEW receiver hashes
	// Note: The sign message contains the list of ALL receivers, so it must be regenerated
	// when the receiver list changes. Old sign messages were deleted with old certificate signatures above.
	signMessage := SignMessage{
		EventContractAddress: eventCertificateAddressStr,
		Receivers:            receivers,
	}

	// Convert sign message to JSON
	signMessageJSONBytes, err := json.Marshal(signMessage)
	if err != nil {
		return nil, err
	}
	signMessageJSON := string(signMessageJSONBytes)

	// 8. Produce host signature — two paths depending on BYOK status
	var encodedHostSignature string
	signMessageDigest := cyptoutils.HashEthereumMessage(signMessageJSON)

	if isByok {
		// Wallet path: frontend signed the message; verify and use the provided signature.
		valid, verifyErr := cyptoutils.VerifySignedMessageByAddress(
			common.HexToAddress(currentUser.WalletAddress),
			signMessageJSON,
			*opts.HostSignature,
		)
		if verifyErr != nil {
			uc.logger.WarnContext(ctx, "security event: failed signature verification",
				"event_id", eventID,
				"wallet_address", currentUser.WalletAddress,
				"error", verifyErr.Error(),
				"context", "ImportCertificateReceivers",
			)
			return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("signature verification failed: %w", verifyErr))
		}
		if !valid {
			uc.logger.WarnContext(ctx, "security event: failed signature verification",
				"event_id", eventID,
				"wallet_address", currentUser.WalletAddress,
				"context", "ImportCertificateReceivers",
			)
			return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("host signature does not match wallet address"))
		}
		encodedHostSignature = *opts.HostSignature
	} else {
		// PIN path: decrypt private key server-side and sign.
		privateKey, _, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, *opts.HostPin)
		if err != nil {
			return nil, err
		}
		rawSignature, err := cyptoutils.Sign(signMessageDigest.Bytes(), privateKey)
		if err != nil {
			return nil, err
		}
		encodedHostSignature = hexutil.Encode(rawSignature)
	}

	// 9. Create event_certificate_signatures for the certificate config (one set per config, not per certificate)
	// Get event issuers for this event
	eventIssuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Create signature for each issuer (linked to config, not individual certificates)
	for _, issuer := range eventIssuers {
		// Store the hash digest as hex string (for compatibility, though contract uses original message)
		encodedSignMessageDigestStr := hexutil.Encode(signMessageDigest[:])

		_, err := uc.EventCertificateSignatureDataGateway.CreateEventCertificateSignature(ctx, eventdatagateway.CreateEventCertificateSignatureParameters{
			EventCertificateConfigID: certificateConfig.ID,
			IssuerCredentialID:       issuer.IssuerCredentialID,
			IssuerSignature:          nil, // Will be set when issuer signs
			HostSignature:            encodedHostSignature,
			SignMessage:              &signMessageJSON,
			SignMessageDigest:        &encodedSignMessageDigestStr,
		})
		if err != nil {
			return nil, err
		}
	}

	return &ImportCertificateReceiversResponse{
		EventID:                 eventID,
		EventCertificateAddress: eventCertificateAddressStr,
		Certificates:            certificates,
	}, nil
}

// stringPtrIfNotEmpty returns a pointer to the string if it's not empty, otherwise nil
func stringPtrIfNotEmpty(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (uc *EventUsecase) deployCertificateContractImpl(ctx context.Context, transactor *bind.TransactOpts, accessManagerAddr, eventAddr common.Address) (common.Address, error) {
	eventCertificateAddress, tx, _, err := eventCertificateContract.DeployEventCertificate(
		transactor,
		uc.ethClient,
		accessManagerAddr,
		eventAddr,
	)
	if err != nil {
		return common.Address{}, err
	}

	receipt, err := bind.WaitMined(ctx, uc.ethClient, tx)
	if err != nil {
		return common.Address{}, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("transaction mining failed: tx=%s: %w", tx.Hash().Hex(), err))
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, customerror.Parse(&customerror.ErrInternalServer, fmt.Errorf("transaction reverted (tx=%s)", tx.Hash().Hex()))
	}

	return eventCertificateAddress, nil
}
