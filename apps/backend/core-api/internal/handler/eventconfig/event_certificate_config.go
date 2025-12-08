package eventconfig

import (
	"github.com/google/uuid"
)

// MintReadinessInfo contains mint readiness information for OpenAPI docs
type MintReadinessInfo struct {
	IsReady                    bool     `json:"is_ready"`
	HasCertificateConfig       bool     `json:"has_certificate_config"`
	AllIssuersHaveSigned       bool     `json:"all_issuers_have_signed"`
	SignedIssuersCount         int64    `json:"signed_issuers_count"`
	TotalIssuersCount          int64    `json:"total_issuers_count"`
	HasCertificateContract     bool     `json:"has_certificate_contract"`
	CertificateContractAddress *string  `json:"certificate_contract_address,omitempty"`
	MissingRequirements        []string `json:"missing_requirements,omitempty"`
}

type EventCertificateConfigResponse struct {
	ID                          uuid.UUID          `json:"id"`
	EventID                     uuid.UUID          `json:"event_id"`
	BaseCertificateStorageKey   string             `json:"base_certificate_storage_key"`
	BaseCertificatePresignedURL string             `json:"base_certificate_presigned_url"`
	EventNamePosX               float64            `json:"event_name_pos_x"`
	EventNamePosY               float64            `json:"event_name_pos_y"`
	NamePosX                    float64            `json:"name_pos_x"`
	NamePosY                    float64            `json:"name_pos_y"`
	AcademicInstitutionPosX     *float64           `json:"academic_institution_pos_x,omitempty"`
	AcademicInstitutionPosY     *float64           `json:"academic_institution_pos_y,omitempty"`
	CertificateTitlePosX        *float64           `json:"certificate_title_pos_x,omitempty"`
	CertificateTitlePosY        *float64           `json:"certificate_title_pos_y,omitempty"`
	CertificateSubtitlePosX     *float64           `json:"certificate_subtitle_pos_x,omitempty"`
	CertificateSubtitlePosY     *float64           `json:"certificate_subtitle_pos_y,omitempty"`
	IsPublished                 bool               `json:"is_published"`
	CreatedAt                   string             `json:"created_at"`
	UpdatedAt                   string             `json:"updated_at"`
	MintReadiness               *MintReadinessInfo `json:"mint_readiness,omitempty"`
}
