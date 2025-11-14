package event

import (
	"context"
	"decm-database/go/generated"
	"encoding/json"
	"fmt"

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
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	AcademicInstitution string `json:"academic_institution"`
	CertificateTitle    string `json:"certificate_title"`
	CertificateSubtitle string `json:"certificate_subtitle"`
	HostPin             string `json:"host_pin"`
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

	// 2. Check if event exists
	event, err := uc.EventDataGateway.GetEventById(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, fmt.Errorf("event not found"))
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

	privateKey, _, err := cyptoutils.DecryptPrivateKey(*credential.EncryptedPrivateKey, requests[0].HostPin)
	if err != nil {
		return nil, err
	}

	eventCertificateAddressStr := ""

	if eventContract.CertificateContractAddress.String == "" {
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
			common.HexToAddress(eventContract.EventContractAddress),
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
		eventCertificateAddressStr = eventContract.CertificateContractAddress.String
	}

	// 5. Update eventContract.certificate_contract_address
	_, err = uc.EventContractDataGateway.UpdateEventContract(ctx, generated.UpdateEventContractParams{
		EventID:                      eventID,
		AccessManagerContractAddress: eventContract.AccessManagerContractAddress,
		EventContractAddress:         eventContract.EventContractAddress,
		TicketContractAddress:        eventContract.TicketContractAddress,
		CertificateContractAddress:   pgmapper.StringPtrToPgText(&eventCertificateAddressStr),
	})
	if err != nil {
		return nil, err
	}

	// 6. Save certificate data to event_certificates
	certificates := make([]*entity.EventCertificate, 0, len(requests))
	var certificateIDs []uuid.UUID

	for _, req := range requests {
		// Combine first and last name
		name := fmt.Sprintf("%s %s", req.FirstName, req.LastName)

		// Create CSV value
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, req.AcademicInstitution, req.CertificateTitle, req.CertificateSubtitle)

		// Hash the CSV value
		hash := cyptoutils.HashMessage(csvValue)
		encodedHash := hexutil.Encode(hash)

		// Create certificate
		certificate, err := uc.EventCertificateDataGateway.CreateEventCertificate(ctx, eventdatagateway.CreateEventCertificateParameters{
			EventID:                 eventID,
			ReceiverCredentialID:    nil, // Will be set when receiver claims certificate
			ReceiverEmail:           nil, // Will be set when receiver claims certificate
			Name:                    &name,
			AcademicInstitution:     &req.AcademicInstitution,
			CertificateTitle:        &req.CertificateTitle,
			CertificateSubtitle:     &req.CertificateSubtitle,
			EventContractAddress:    eventContract.EventContractAddress,
			EventCertificateAddress: &eventCertificateAddressStr,
			CertificateTokenID:      nil, // Will be set when minted,
			Digest:                  &encodedHash,
		})
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certificate)
		certificateIDs = append(certificateIDs, certificate.Id)
	}

	// 7. Create sign_message with hashes
	receivers := make([]string, 0, len(requests))
	for _, req := range requests {
		// Combine first and last name
		name := fmt.Sprintf("%s %s", req.FirstName, req.LastName)

		// Create CSV value
		csvValue := fmt.Sprintf("%s,%s,%s,%s", name, req.AcademicInstitution, req.CertificateTitle, req.CertificateSubtitle)

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
	signMessageDigest := cyptoutils.HashMessage(signMessageJSON)
	signature, err := cyptoutils.Sign(signMessageDigest[:], privateKey)
	if err != nil {
		return nil, err
	}

	// 9. Create event_certificate_signatures for each certificate
	for _, certificateID := range certificateIDs {
		// Get event issuers for this event
		eventIssuers, err := uc.EventIssuerDataGateway.GetEventIssuersByEventID(ctx, eventID)
		if err != nil {
			return nil, err
		}

		// Create signature for each issuer
		for _, issuer := range eventIssuers {
			encodedSignMessageDigestStr := hexutil.Encode(signMessageDigest[:])
			encodedHostSignature := hexutil.Encode(signature)

			_, err := uc.EventCertificateSignatureDataGateway.CreateEventCertificateSignature(ctx, eventdatagateway.CreateEventCertificateSignatureParameters{
				EventCertificateID: certificateID,
				IssuerCredentialID: issuer.IssuerCredentialID,
				IssuerSignature:    nil, // Will be set when issuer signs
				HostSignature:      encodedHostSignature,
				SignMessage:        &signMessageJSON,
				SignMessageDigest:  &encodedSignMessageDigestStr,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return &ImportCertificateReceiversResponse{
		EventID:                 eventID,
		EventCertificateAddress: eventCertificateAddressStr,
		Certificates:            certificates,
	}, nil
}
