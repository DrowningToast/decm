package entity

import "time"

type ClaimCertificateStatus string

var (
	ClaimCertificateStatusPending ClaimCertificateStatus = "PENDING"
	ClaimCertificateStatusClaimed ClaimCertificateStatus = "CLAIMED"
	ClaimCertificateStatusExpired ClaimCertificateStatus = "EXPIRED"
	ClaimCertificateStatusAborted ClaimCertificateStatus = "ABORTED"
)

// GetClaimCertificateStatus determines the claim status of a certificate from its fields.
// Priority: CLAIMED > ABORTED > EXPIRED > PENDING
func GetClaimCertificateStatus(certificate *EventCertificate) ClaimCertificateStatus {
	if certificate.CertificateTokenId != nil || certificate.BroadcastedAt != nil {
		return ClaimCertificateStatusClaimed
	}
	if certificate.AbortedAt != nil {
		return ClaimCertificateStatusAborted
	}
	if certificate.EstimatedDeadline != nil && certificate.EstimatedDeadline.Before(time.Now()) {
		return ClaimCertificateStatusExpired
	}
	return ClaimCertificateStatusPending
}
