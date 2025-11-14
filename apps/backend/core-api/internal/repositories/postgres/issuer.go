package postgres

import (
	"context"
	"decm-database/go/generated"

	"apps/backend/common"
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/datagateway"
	"apps/backend/core-api/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ datagateway.IssuerDataGateway = (*Repository)(nil)

func (r *Repository) ListVerifiedIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	query, err := r.queries.ListVerifiedIssuerProfiles(ctx, generated.ListVerifiedIssuerProfilesParams{
		LimitCount:  int32(limitCount),
		OffsetCount: int32(offsetCount),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}
	profiles := make([]entity.Profile, len(query))
	for i, profile := range query {
		// Decrypt PII fields for return
		profilePictureUrlDec, err := pgmapper.DecryptPgTextToStringPtr(profile.ProfilePictureUrl, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(profile.FirstName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(profile.LastName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		emailDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Email, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		bioDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Bio, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(profile.PhoneNumber, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		addressDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Address, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(profile.AcademicInstitution, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(profile.AcademicEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		profiles[i] = entity.Profile{
			Id:                          profile.ID,
			AuthenticationCredentialId:  profile.AuthenticationCredentialID,
			IsProfilePicturePublic:      profile.IsProfilePicturePublic == 1,
			ProfilePictureUrl:           profilePictureUrlDec,
			IsFirstNamePublic:           profile.IsFirstNamePublic == 1,
			FirstName:                   firstNameDec,
			IsLastNamePublic:            profile.IsLastNamePublic == 1,
			LastName:                    lastNameDec,
			IsEmailPublic:               profile.IsEmailPublic == 1,
			Email:                       emailDec,
			IsBioPublic:                 profile.IsBioPublic == 1,
			Bio:                         bioDec,
			IsPhoneNumberPublic:         profile.IsPhoneNumberPublic == 1,
			PhoneNumber:                 phoneNumberDec,
			IsAddressPublic:             profile.IsAddressPublic == 1,
			Address:                     addressDec,
			IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic == 1,
			AcademicInstitution:         academicInstitutionDec,
			IsAcademicEmailPublic:       profile.IsAcademicEmailPublic == 1,
			AcademicEmail:               academicEmailDec,
			CreatedAt:                   profile.CreatedAt.Time,
			UpdatedAt:                   profile.UpdatedAt.Time,
		}
	}

	return profiles, nil
}

func (r *Repository) ListIssuerProfiles(ctx context.Context, limitCount int, offsetCount int) ([]entity.Profile, error) {
	query, err := r.queries.ListIssuerProfiles(ctx, generated.ListIssuerProfilesParams{
		LimitCount:  int32(limitCount),
		OffsetCount: int32(offsetCount),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	profiles := make([]entity.Profile, len(query))
	for i, profile := range query {
		// Decrypt PII fields for return using pgmapper
		profilePictureUrlDec, err := pgmapper.DecryptPgTextToStringPtr(profile.ProfilePictureUrl, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(profile.FirstName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(profile.LastName, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		emailDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Email, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		bioDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Bio, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(profile.PhoneNumber, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		addressDec, err := pgmapper.DecryptPgTextToStringPtr(profile.Address, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(profile.AcademicInstitution, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(profile.AcademicEmail, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		profiles[i] = entity.Profile{
			Id:                          profile.ID,
			AuthenticationCredentialId:  profile.AuthenticationCredentialID,
			IsProfilePicturePublic:      profile.IsProfilePicturePublic == 1,
			ProfilePictureUrl:           profilePictureUrlDec,
			IsFirstNamePublic:           profile.IsFirstNamePublic == 1,
			FirstName:                   firstNameDec,
			IsLastNamePublic:            profile.IsLastNamePublic == 1,
			LastName:                    lastNameDec,
			IsEmailPublic:               profile.IsEmailPublic == 1,
			Email:                       emailDec,
			IsBioPublic:                 profile.IsBioPublic == 1,
			Bio:                         bioDec,
			IsPhoneNumberPublic:         profile.IsPhoneNumberPublic == 1,
			PhoneNumber:                 phoneNumberDec,
			IsAddressPublic:             profile.IsAddressPublic == 1,
			Address:                     addressDec,
			IsAcademicInstitutionPublic: profile.IsAcademicInstitutionPublic == 1,
			AcademicInstitution:         academicInstitutionDec,
			IsAcademicEmailPublic:       profile.IsAcademicEmailPublic == 1,
			AcademicEmail:               academicEmailDec,
			CreatedAt:                   profile.CreatedAt.Time,
			UpdatedAt:                   profile.UpdatedAt.Time,
		}
	}

	return profiles, nil
}

func (r *Repository) SearchIssuerCredentialsByWalletAddress(ctx context.Context, searchQuery string, limitCount int, offsetCount int) ([]entity.AuthenticationCredential, error) {
	query, err := r.queries.ListIssuerAuthenticationCredentialsByWalletAddress(ctx, generated.ListIssuerAuthenticationCredentialsByWalletAddressParams{
		SearchQuery: pgtype.Text{String: searchQuery, Valid: true},
		LimitCount:  int32(limitCount),
		OffsetCount: int32(offsetCount),
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	credentials := make([]entity.AuthenticationCredential, len(query))
	for i, cred := range query {
		// Decrypt PII fields using pgmapper
		googleConnectorRefDec, err := pgmapper.DecryptPgTextToStringPtr(cred.GoogleConnectorRef, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}
		githubConnectorRefDec, err := pgmapper.DecryptPgTextToStringPtr(cred.GithubConnectorRef, r.piiEncryptionKey)
		if err != nil {
			return nil, err
		}

		credentials[i] = entity.AuthenticationCredential{
			Id:                  cred.ID,
			SolutionStatus:      common.SolutionStatus(cred.SolutionStatus),
			WalletAddress:       cred.WalletAddress,
			GoogleConnectorRef:  googleConnectorRefDec,
			GithubConnectorRef:  githubConnectorRefDec,
			IsVerifiedOrganizer: cred.IsVerifiedOrganizer == 1,
			IsVerifiedStudent:   cred.IsVerifiedStudent == 1,
			IsVerifiedIssuer:    cred.IsVerifiedIssuer == 1,
			CreatedAt:           cred.CreatedAt.Time,
			UpdatedAt:           cred.UpdatedAt.Time,
		}
	}

	return credentials, nil
}

func (r *Repository) GetEventsByIssuerCredentialID(ctx context.Context, issuerCredentialID string, limitCount int32, offsetCount int32) ([]generated.GetEventIssuersByCredentialIDRow, error) {
	query, err := r.queries.GetEventIssuersByCredentialID(ctx, generated.GetEventIssuersByCredentialIDParams{
		IssuerCredentialID: uuid.MustParse(issuerCredentialID),
		LimitCount:         limitCount,
		OffsetCount:        offsetCount,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	return query, nil
}
