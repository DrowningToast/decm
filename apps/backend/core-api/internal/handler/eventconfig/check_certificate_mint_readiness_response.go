package eventconfig

// CertificateMintReadinessResponse represents the response for certificate mint readiness check
type CertificateMintReadinessResponse struct {
	IsReady                    bool     `json:"is_ready" example:"true"`
	HasCertificateConfig       bool     `json:"has_certificate_config" example:"true"`
	AllIssuersHaveSigned       bool     `json:"all_issuers_have_signed" example:"true"`
	SignedIssuersCount         int64    `json:"signed_issuers_count" example:"2"`
	TotalIssuersCount          int64    `json:"total_issuers_count" example:"2"`
	HasCertificateContract     bool     `json:"has_certificate_contract" example:"true"`
	CertificateContractAddress *string  `json:"certificate_contract_address,omitempty" example:"0x1234567890abcdef"`
	MissingRequirements        []string `json:"missing_requirements,omitempty"`
}














