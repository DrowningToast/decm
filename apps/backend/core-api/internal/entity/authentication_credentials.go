package entity

import (
	"decm-database/go/generated"
	"time"

	"apps/backend/common/pgmapper"

	"github.com/jackc/pgx/v5/pgtype"
)

type SolutionStatus int32

const (
	ManagedSolutionStatus SolutionStatus = 0
	BYOSSolutionStatus    SolutionStatus = 1
)

type AuthenticationCredential struct {
	Id             int32          `json:"id"`
	SolutionStatus SolutionStatus `json:"solution_status"`
	// Hashed of a hashed password using Argon2id
	HashedPassword *string `json:"password"`
	// Encrypted by a hash of password using AES-256-GCM
	EncryptedPrivateKey *string `json:"private_key"`
	PublicKey           string  `json:"public_key"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GoogleConnectorRef *string `json:"google_connector_ref"`
	// Encrypted by a PII Encryption Key using AES-256-GCM
	GithubConnectorRef  *string   `json:"github_connector_ref"`
	IsVerifiedOrganizer bool      `json:"is_verified_organizer"`
	IsVerifiedStudent   bool      `json:"is_verified_student"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (entity *AuthenticationCredential) ToModel() *generated.AuthenticationCredential {
	var isVerifiedStudent int32 = 0
	if entity.IsVerifiedStudent {
		isVerifiedStudent = 1
	}
	var isVerifiedOrganizer int32 = 0
	if entity.IsVerifiedOrganizer {
		isVerifiedOrganizer = 1
	}

	return &generated.AuthenticationCredential{
		ID:             entity.Id,
		SolutionStatus: int32(entity.SolutionStatus),
		Password: pgtype.Text{
			String: *entity.Password,
			Valid:  entity.Password != nil,
		},
		PrivateKey: pgtype.Text{
			String: *entity.PrivateKey,
			Valid:  entity.PrivateKey != nil,
		},
		PublicKey: entity.PublicKey,
		GoogleConnectorRef: pgtype.Text{
			String: *entity.GoogleConnectorRef,
			Valid:  entity.GoogleConnectorRef != nil,
		},
		GithubConnectorRef: pgtype.Text{
			String: *entity.GithubConnectorRef,
			Valid:  entity.GithubConnectorRef != nil,
		},
		IsVerifiedOrganizer: isVerifiedOrganizer,
		IsVerifiedStudent:   isVerifiedStudent,
		CreatedAt: pgtype.Timestamptz{
			Time:  entity.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  entity.UpdatedAt,
			Valid: true,
		},
	}
}

func MapAuthenticationCredentialsToEntities(models []generated.AuthenticationCredential) []AuthenticationCredential {
	entities := make([]AuthenticationCredential, len(models))
	for i, model := range models {
		entities[i] = AuthenticationCredential{
			Id:             model.ID,
			SolutionStatus: SolutionStatus(model.SolutionStatus),
			Password:       pgmapper.PgTextToStringPtr(model.Password),
			PrivateKey:     pgmapper.PgTextToStringPtr(model.PrivateKey),
			PublicKey:      model.PublicKey,
		}
	}
	return entities
}
