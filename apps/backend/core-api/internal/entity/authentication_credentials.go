package entity

import (
	"time"

	"github.com/google/uuid"
)

type SolutionStatus int32

const (
	SolutionStatusManaged SolutionStatus = 0
	SolutionStatusBYOK    SolutionStatus = 1
)

// User represented by an authentication credential
// @description User represented by an authentication credential
type AuthenticationCredential struct {
	Id             uuid.UUID      `json:"id"`
	SolutionStatus SolutionStatus `json:"solution_status"`
	// Hashed of a hashed password using Argon2id
	HashedPassword *string `json:"password"`
	// Encrypted by a hash of password using AES-256-GCM
	EncryptedPrivateKey *string `json:"private_key"`
	WalletAddress       string  `json:"wallet_address"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GoogleConnectorRef *string `json:"google_connector_ref"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GithubConnectorRef  *string   `json:"github_connector_ref"`
	IsVerifiedOrganizer bool      `json:"is_verified_organizer"`
	IsVerifiedStudent   bool      `json:"is_verified_student"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
} // @name AuthenticationCredential
