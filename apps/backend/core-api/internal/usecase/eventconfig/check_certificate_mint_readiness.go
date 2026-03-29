package eventconfig

import (
	"apps/backend/common/utils"
	"context"
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CertificateMintReadinessResponse represents the response for certificate mint readiness check
type CertificateMintReadinessResponse struct {
	IsReady                    bool     `json:"is_ready"`
	HasCertificateConfig       bool     `json:"has_certificate_config"`
	AllIssuersHaveSigned       bool     `json:"all_issuers_have_signed"`
	SignedIssuersCount         int64    `json:"signed_issuers_count"`
	TotalIssuersCount          int64    `json:"total_issuers_count"`
	HasCertificateContract     bool     `json:"has_certificate_contract"`
	CertificateContractAddress *string  `json:"certificate_contract_address,omitempty"`
	MissingRequirements        []string `json:"missing_requirements,omitempty"`
}

// CheckCertificateMintReadiness checks if an event certificate is ready to be minted
// Returns detailed information about the readiness status and what requirements are missing
func (uc *EventConfigUsecase) CheckCertificateMintReadiness(ctx context.Context, eventID uuid.UUID) (*CertificateMintReadinessResponse, error) {
	response := &CertificateMintReadinessResponse{
		IsReady:                false,
		HasCertificateConfig:   false,
		AllIssuersHaveSigned:   false,
		SignedIssuersCount:     0,
		TotalIssuersCount:      0,
		HasCertificateContract: false,
		MissingRequirements:    []string{},
	}

	// 1. Check if certificate config exists
	certificateConfig, err := uc.EventCertificateDg.GetEventCertificateConfigByEventID(ctx, eventID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "No row found in event_certificate_configs table", "event_id", eventID)
			response.MissingRequirements = append(response.MissingRequirements, "Certificate configuration is not set up")
		} else {
			return nil, errors.Wrap(err, "failed to check certificate config")
		}
	} else if certificateConfig != nil {
		response.HasCertificateConfig = true
	}

	// 2. Check if ALL issuers have signed
	allIssuersHaveSigned, err := uc.EventIssuerDg.AllIssuersHaveSigned(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check if all issuers have signed")
	}
	response.AllIssuersHaveSigned = allIssuersHaveSigned

	// Get signed issuers count
	signedCount, err := uc.EventIssuerDg.GetSignedIssuersCount(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signed issuers count")
	}
	response.SignedIssuersCount = signedCount

	// Get total issuers count
	totalCount, err := uc.EventIssuerDg.GetTotalIssuersCount(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total issuers count")
	}
	response.TotalIssuersCount = totalCount

	// Add specific message about issuer signing status
	if totalCount == 0 {
		response.MissingRequirements = append(response.MissingRequirements, "No issuers have been assigned to this event")
	} else if !allIssuersHaveSigned {
		response.MissingRequirements = append(response.MissingRequirements,
			fmt.Sprintf("Not all issuers have signed (%d/%d signed)", signedCount, totalCount))
	}

	// 3. Check if certificate contract is deployed
	eventContract, err := uc.EventContractDg.GetEventContractByEventID(ctx, eventID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.InfoContext(ctx, "No row found in event_contracts table", "event_id", eventID)
			response.MissingRequirements = append(response.MissingRequirements, "Event contracts are not deployed")
		} else {
			return nil, errors.Wrap(err, "failed to check event contract")
		}
	} else if eventContract != nil && utils.DerefOrEmpty(eventContract.CertificateContractAddress) != "" {
		response.HasCertificateContract = true
		response.CertificateContractAddress = eventContract.CertificateContractAddress
	} else {
		response.MissingRequirements = append(response.MissingRequirements, "Certificate contract address is not set")
	}

	// Determine if certificate is ready to mint (ALL requirements must be met)
	response.IsReady = response.HasCertificateConfig &&
		response.AllIssuersHaveSigned &&
		response.HasCertificateContract

	return response, nil
}
