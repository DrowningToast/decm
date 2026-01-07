package entity

import (
	"apps/backend/common"
	"time"

	"github.com/google/uuid"
)

// User represented by an authentication credential
// @description User represented by an authentication credential
type AuthenticationCredential struct {
	Id             uuid.UUID             `json:"id"`
	SolutionStatus common.SolutionStatus `json:"solution_status"`
	// Hashed of a hashed password using Argon2id
	HashedPassword *string `json:"password,omitempty"`
	// Encrypted by a hash of password using AES-256-GCM
	EncryptedPrivateKey *string `json:"private_key,omitempty"`
	WalletAddress       string  `json:"wallet_address"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GoogleConnectorRef *string `json:"google_connector_ref,omitempty"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GithubConnectorRef  *string   `json:"github_connector_ref,omitempty"`
	IsVerifiedOrganizer bool      `json:"is_verified_organizer"`
	IsVerifiedStudent   bool      `json:"is_verified_student"`
	IsVerifiedIssuer    bool      `json:"is_verified_issuer"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
} // @name AuthenticationCredential
