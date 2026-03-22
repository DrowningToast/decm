package event

import "apps/backend/core-api/internal/entity"

// ClaimCertificateStatus and its values are defined in the entity layer.
// These aliases keep existing references in this package working.
type ClaimCertificateStatus = entity.ClaimCertificateStatus

var (
	ClaimCertificateStatusPending = entity.ClaimCertificateStatusPending
	ClaimCertificateStatusClaimed = entity.ClaimCertificateStatusClaimed
	ClaimCertificateStatusExpired = entity.ClaimCertificateStatusExpired
	ClaimCertificateStatusAborted = entity.ClaimCertificateStatusAborted
)

// GetClaimCertificateStatus is an alias for entity.GetClaimCertificateStatus.
var GetClaimCertificateStatus = entity.GetClaimCertificateStatus
