package event

import (
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"
	"strings"

	"apps/backend/common/customerror"
	"apps/backend/common/pgmapper"
	eventdatagateway "apps/backend/core-api/internal/datagateway/event"
	"apps/backend/core-api/internal/entity"
	"apps/backend/services/auth"

	cyptoutils "apps/backend/core-api/internal/usecase/cyptoutils"

	eventCertificateContract "apps/backend/contracts/certificate"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
)

type ImportCertificateReceiversRequest struct {
	Email               string  `json:"email"`
	FirstName           *string `json:"first_name"`
	LastName            *string `json:"last_name"`
	AcademicInstitution *string `json:"academic_institution"`
	CertificateTitle    *string `json:"certificate_title"`
	CertificateSubtitle *string `json:"certificate_subtitle"`
	HostPin             string  `json:"host_pin"`
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

func (uc *EventUsecase) ImportCertificateReceivers(ctx context.Context, eventID uuid.UUID, requests []ImportCertificateReceiversRequest, currentUser *auth.JwtClaims) (*ImportCertificateReceiversResponse, error) {
	// 1. Check if current user is authorized
	credential, err := uc.AuthenticationCredentialDg.GetAuthenticationCredentialByIdWithEncryptedPrivateKey(ctx, currentUser.UserId)
	if err != nil {
		return nil, err
	}

	if !credential.IsVerifiedOrganizer {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user is not a verified organizer"))
	}

	if credential.EncryptedPrivateKey == nil {
		return nil, customerror.Parse(&customerror.ErrUnauthorized, fmt.Errorf("user does not have an encrypted private key"))
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

	privateKey, _, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, requests[0].HostPin)
	if err != nil {
		return nil, err
	}

	eventCertificateAddressStr := ""

	if eventContract.CertificateContractAddress == nil {
		// 4. Deploy event certificate contract
		client, err := cyptoutils.GetEthereumClient()
		if err != nil {
			return nil, err
		}

		auth, err := cyptoutils.GetKeyedTransactor()
		if err != nil {
			return nil, err
		}

		eventCertificateAddress, tx, _, err := eventCertificateContract.DeployEventCertificate(
			auth,
			client,
			common.HexToAddress(eventContract.AccessManagerContractAddress),
			common.HexToAddress(eventContract.EventContractAddress),
		)
		if err != nil {
			return nil, err
		}

		// Wait for transaction to be mined
		_, err = bind.WaitMined(ctx, client, tx)
		if err != nil {
			return nil, err
		}

		eventCertificateAddressStr = eventCertificateAddress.Hex()
	} else {
		eventCertificateAddressStr = *eventContract.CertificateContractAddress
	}

	// 5. Update eventContract.certificate_contract_address
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

	// 6. Save certificate data to event_certificates
	certificates := make([]*entity.EventCertificate, 0, len(requests))

	for _, req := range requests {
		// Safely dereference pointer fields (use empty string if nil)
		firstName := ""
		if req.FirstName != nil {
			firstName = *req.FirstName
		}
		lastName := ""
		if req.LastName != nil {
			lastName = *req.LastName
		}
		academicInstitution := ""
		if req.AcademicInstitution != nil {
			academicInstitution = *req.AcademicInstitution
		}
		certificateTitle := ""
		if req.CertificateTitle != nil {
			certificateTitle = *req.CertificateTitle
		}
		certificateSubtitle := ""
		if req.CertificateSubtitle != nil {
			certificateSubtitle = *req.CertificateSubtitle
		}

		// Combine first and last name
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Create CSV value for hash (used for blockchain verification)
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)

		// Create certificate - pass original pointers to preserve nil vs empty distinction
		certificate, err := uc.EventCertificateDataGateway.CreateEventCertificate(ctx, eventdatagateway.CreateEventCertificateParameters{
			EventID:                 eventID,
			ReceiverCredentialID:    nil,        // Will be set when receiver claims certificate
			ReceiverEmail:           &req.Email, // Set email from import for matching
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
	receivers := make([]string, 0, len(requests))
	for _, req := range requests {
		// Safely dereference pointer fields (use empty string if nil)
		firstName := ""
		if req.FirstName != nil {
			firstName = *req.FirstName
		}
		lastName := ""
		if req.LastName != nil {
			lastName = *req.LastName
		}
		academicInstitution := ""
		if req.AcademicInstitution != nil {
			academicInstitution = *req.AcademicInstitution
		}
		certificateTitle := ""
		if req.CertificateTitle != nil {
			certificateTitle = *req.CertificateTitle
		}
		certificateSubtitle := ""
		if req.CertificateSubtitle != nil {
			certificateSubtitle = *req.CertificateSubtitle
		}

		// Combine first and last name
		name := fmt.Sprintf("%s %s", firstName, lastName)

		// Create CSV value
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, academicInstitution, certificateTitle, certificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)
		receivers = append(receivers, encodedHash)
	}

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

	// 8. Host signs the message
	// Use HashEthereumMessage to match contract's recoverSigner which applies Ethereum prefix
	// The contract expects the original message string and will apply the prefix itself
	signMessageDigest := cyptoutils.HashEthereumMessage(signMessageJSON)
	signature, err := cyptoutils.Sign(signMessageDigest.Bytes(), privateKey)
	if err != nil {
		return nil, err
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
		encodedHostSignature := hexutil.Encode(signature)

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
