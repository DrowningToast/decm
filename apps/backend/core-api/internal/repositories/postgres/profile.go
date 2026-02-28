package postgres

import (
	"apps/backend/common"
	"apps/backend/common/pgerrutils"
	"apps/backend/common/pgmapper"
	offchain_datagateway "apps/backend/core-api/internal/datagateway/offchain"
	"apps/backend/core-api/internal/entity"
	"context"
	"decm-database/go/generated"

	customerror "apps/backend/common/customerror"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var _ offchain_datagateway.ProfileDataGateway = (*Repository)(nil)

func (r *Repository) GetProfileById(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	query, err := r.queries.GetProfileByID(ctx, id)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	profilePictureUrl, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstName, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastName, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bio, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	address, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmail, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstName,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastName,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       email,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumber,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmail,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, error) {
	query, err := r.queries.GetProfileByAuthCredentialID(ctx, authenticationCredentialId)
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	profilePictureUrl, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstName, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastName, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	email, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bio, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	address, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmail, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstName,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastName,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       email,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumber,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmail,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetProfileByEmail(ctx context.Context, email string) (*entity.Profile, error) {
	// Encrypt the email for searching (database stores encrypted values)
	encryptedEmail, err := pgmapper.EncryptPII(email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.GetProfileByEmail(ctx, pgtype.Text{
		String: encryptedEmail,
		Valid:  true,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields
	profilePictureUrl, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstName, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastName, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDecrypted, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bio, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	address, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitution, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrl,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstName,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastName,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       emailDecrypted,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bio,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumber,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     address,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitution,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmailDec,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetProfileAndCredentialWithCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID) (*entity.Profile, *entity.AuthenticationCredential, error) {
	result, err := r.queries.GetProfileAndCredentialWithCredentialId(ctx, authenticationCredentialId)
	if err != nil {
		return nil, nil, pgerrutils.ParsePgError(err)
	}

	decryptedProfilePictureUrl, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedFirstName, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileFirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedLastName, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileLastName, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedBio, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileBio, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedPhoneNumber, err := pgmapper.DecryptPgTextToStringPtr(result.ProfilePhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedAddress, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileAddress, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedAcademicInstitution, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileAcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedAcademicEmail, err := pgmapper.DecryptPgTextToStringPtr(result.ProfileAcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}

	profile := &entity.Profile{
		Id:                          result.ProfileID,
		AuthenticationCredentialId:  result.ProfileAuthenticationCredentialID,
		IsProfilePicturePublic:      result.ProfileIsProfilePicturePublic == 1,
		ProfilePictureUrl:           decryptedProfilePictureUrl,
		IsFirstNamePublic:           result.ProfileIsFirstNamePublic == 1,
		FirstName:                   decryptedFirstName,
		IsLastNamePublic:            result.ProfileIsLastNamePublic == 1,
		LastName:                    decryptedLastName,
		IsEmailPublic:               result.ProfileIsEmailPublic == 1,
		Email:                       decryptedEmail,
		IsBioPublic:                 result.ProfileIsBioPublic == 1,
		Bio:                         decryptedBio,
		IsPhoneNumberPublic:         result.ProfileIsPhoneNumberPublic == 1,
		PhoneNumber:                 decryptedPhoneNumber,
		IsAddressPublic:             result.ProfileIsAddressPublic == 1,
		Address:                     decryptedAddress,
		IsAcademicInstitutionPublic: result.ProfileIsAcademicInstitutionPublic == 1,
		AcademicInstitution:         decryptedAcademicInstitution,
		IsAcademicEmailPublic:       result.ProfileIsAcademicEmailPublic == 1,
		AcademicEmail:               decryptedAcademicEmail,
		CreatedAt:                   result.ProfileCreatedAt.Time,
		UpdatedAt:                   result.ProfileUpdatedAt.Time,
	}

	decryptedGoogleConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(result.GoogleConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedGithubConnectorRef, err := pgmapper.DecryptPgTextToStringPtr(result.GithubConnectorRef, r.piiEncryptionKey)
	if err != nil {
		return nil, nil, err
	}

	authenticationCredential := &entity.AuthenticationCredential{
		Id:                  result.AuthenticationCredentialID,
		SolutionStatus:      common.SolutionStatus(result.SolutionStatus),
		HashedPassword:      pgmapper.PgTextToStringPtr(result.HashedPassword),
		EncryptedPrivateKey: nil,
		WalletAddress:       result.WalletAddress,
		GoogleConnectorRef:  decryptedGoogleConnectorRef,
		GithubConnectorRef:  decryptedGithubConnectorRef,
		IsVerifiedOrganizer: result.IsVerifiedOrganizer == 1,
		IsVerifiedStudent:   result.IsVerifiedStudent == 1,
		IsVerifiedIssuer:    result.IsVerifiedIssuer == 1,
		CreatedAt:           result.AuthenticationCredentialCreatedAt.Time,
		UpdatedAt:           result.AuthenticationCredentialUpdatedAt.Time,
	}

	return profile, authenticationCredential, nil
}

func (r *Repository) CreateProfile(ctx context.Context, profile entity.Profile) (*entity.Profile, error) {
	// Encrypt PII fields
	profilePictureUrlEnc, err := pgmapper.EncryptStringPtrToPgText(profile.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberEnc, err := pgmapper.EncryptStringPtrToPgText(profile.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.CreateProfile(ctx, generated.CreateProfileParams{
		AuthenticationCredentialID:  profile.AuthenticationCredentialId,
		IsProfilePicturePublic:      pgmapper.BoolToInt32(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           profilePictureUrlEnc,
		IsFirstNamePublic:           pgmapper.BoolToInt32(profile.IsFirstNamePublic),
		FirstName:                   firstNameEnc,
		IsLastNamePublic:            pgmapper.BoolToInt32(profile.IsLastNamePublic),
		LastName:                    lastNameEnc,
		IsEmailPublic:               pgmapper.BoolToInt32(profile.IsEmailPublic),
		Email:                       emailEnc,
		IsBioPublic:                 pgmapper.BoolToInt32(profile.IsBioPublic),
		Bio:                         bioEnc,
		IsPhoneNumberPublic:         pgmapper.BoolToInt32(profile.IsPhoneNumberPublic),
		PhoneNumber:                 phoneNumberEnc,
		IsAddressPublic:             pgmapper.BoolToInt32(profile.IsAddressPublic),
		Address:                     addressEnc,
		IsAcademicInstitutionPublic: pgmapper.BoolToInt32(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         academicInstitutionEnc,
		IsAcademicEmailPublic:       pgmapper.BoolToInt32(profile.IsAcademicEmailPublic),
		AcademicEmail:               academicEmailEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	profilePictureUrlDec, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioDec, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressDec, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrlDec,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstNameDec,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastNameDec,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       emailDec,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bioDec,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumberDec,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     addressDec,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitutionDec,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmailDec,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, profile offchain_datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	// Encrypt PII fields
	profilePictureUrlEnc, err := pgmapper.EncryptStringPtrToPgText(profile.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberEnc, err := pgmapper.EncryptStringPtrToPgText(profile.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.UpdateProfile(ctx, generated.UpdateProfileParams{
		ID:                          id,
		IsProfilePicturePublic:      pgmapper.BoolPtrToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           profilePictureUrlEnc,
		IsFirstNamePublic:           pgmapper.BoolPtrToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   firstNameEnc,
		IsLastNamePublic:            pgmapper.BoolPtrToPgInt4(profile.IsLastNamePublic),
		LastName:                    lastNameEnc,
		IsEmailPublic:               pgmapper.BoolPtrToPgInt4(profile.IsEmailPublic),
		Email:                       emailEnc,
		IsBioPublic:                 pgmapper.BoolPtrToPgInt4(profile.IsBioPublic),
		Bio:                         bioEnc,
		IsPhoneNumberPublic:         pgmapper.BoolPtrToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 phoneNumberEnc,
		IsAddressPublic:             pgmapper.BoolPtrToPgInt4(profile.IsAddressPublic),
		Address:                     addressEnc,
		IsAcademicInstitutionPublic: pgmapper.BoolPtrToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         academicInstitutionEnc,
		IsAcademicEmailPublic:       pgmapper.BoolPtrToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               academicEmailEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	profilePictureUrlDec, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioDec, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressDec, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrlDec,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstNameDec,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastNameDec,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       emailDec,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bioDec,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumberDec,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     addressDec,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitutionDec,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmailDec,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) UpdateProfileByAuthenticationCredentialId(ctx context.Context, authenticationCredentialId uuid.UUID, profile offchain_datagateway.UpdateProfileParameters) (*entity.Profile, error) {
	// Encrypt PII fields
	profilePictureUrlEnc, err := pgmapper.EncryptStringPtrToPgText(profile.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameEnc, err := pgmapper.EncryptStringPtrToPgText(profile.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberEnc, err := pgmapper.EncryptStringPtrToPgText(profile.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressEnc, err := pgmapper.EncryptStringPtrToPgText(profile.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailEnc, err := pgmapper.EncryptStringPtrToPgText(profile.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	query, err := r.queries.UpdateProfileByAuthenticationCredentialId(ctx, generated.UpdateProfileByAuthenticationCredentialIdParams{
		AuthenticationCredentialID:  authenticationCredentialId,
		IsProfilePicturePublic:      pgmapper.BoolPtrToPgInt4(profile.IsProfilePicturePublic),
		ProfilePictureUrl:           profilePictureUrlEnc,
		IsFirstNamePublic:           pgmapper.BoolPtrToPgInt4(profile.IsFirstNamePublic),
		FirstName:                   firstNameEnc,
		IsLastNamePublic:            pgmapper.BoolPtrToPgInt4(profile.IsLastNamePublic),
		LastName:                    lastNameEnc,
		IsEmailPublic:               pgmapper.BoolPtrToPgInt4(profile.IsEmailPublic),
		Email:                       emailEnc,
		IsBioPublic:                 pgmapper.BoolPtrToPgInt4(profile.IsBioPublic),
		Bio:                         bioEnc,
		IsPhoneNumberPublic:         pgmapper.BoolPtrToPgInt4(profile.IsPhoneNumberPublic),
		PhoneNumber:                 phoneNumberEnc,
		IsAddressPublic:             pgmapper.BoolPtrToPgInt4(profile.IsAddressPublic),
		Address:                     addressEnc,
		IsAcademicInstitutionPublic: pgmapper.BoolPtrToPgInt4(profile.IsAcademicInstitutionPublic),
		AcademicInstitution:         academicInstitutionEnc,
		IsAcademicEmailPublic:       pgmapper.BoolPtrToPgInt4(profile.IsAcademicEmailPublic),
		AcademicEmail:               academicEmailEnc,
	})
	if err != nil {
		return nil, pgerrutils.ParsePgError(err)
	}

	// Decrypt PII fields for return
	profilePictureUrlDec, err := pgmapper.DecryptPgTextToStringPtr(query.ProfilePictureUrl, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	firstNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.FirstName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	lastNameDec, err := pgmapper.DecryptPgTextToStringPtr(query.LastName, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	emailDec, err := pgmapper.DecryptPgTextToStringPtr(query.Email, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	bioDec, err := pgmapper.DecryptPgTextToStringPtr(query.Bio, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	phoneNumberDec, err := pgmapper.DecryptPgTextToStringPtr(query.PhoneNumber, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	addressDec, err := pgmapper.DecryptPgTextToStringPtr(query.Address, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicInstitutionDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicInstitution, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}
	academicEmailDec, err := pgmapper.DecryptPgTextToStringPtr(query.AcademicEmail, r.piiEncryptionKey)
	if err != nil {
		return nil, err
	}

	return &entity.Profile{
		Id:                          query.ID,
		AuthenticationCredentialId:  query.AuthenticationCredentialID,
		IsProfilePicturePublic:      query.IsProfilePicturePublic == 1,
		ProfilePictureUrl:           profilePictureUrlDec,
		IsFirstNamePublic:           query.IsFirstNamePublic == 1,
		FirstName:                   firstNameDec,
		IsLastNamePublic:            query.IsLastNamePublic == 1,
		LastName:                    lastNameDec,
		IsEmailPublic:               query.IsEmailPublic == 1,
		Email:                       emailDec,
		IsBioPublic:                 query.IsBioPublic == 1,
		Bio:                         bioDec,
		IsPhoneNumberPublic:         query.IsPhoneNumberPublic == 1,
		PhoneNumber:                 phoneNumberDec,
		IsAddressPublic:             query.IsAddressPublic == 1,
		Address:                     addressDec,
		IsAcademicInstitutionPublic: query.IsAcademicInstitutionPublic == 1,
		AcademicInstitution:         academicInstitutionDec,
		IsAcademicEmailPublic:       query.IsAcademicEmailPublic == 1,
		AcademicEmail:               academicEmailDec,
		CreatedAt:                   query.CreatedAt.Time,
		UpdatedAt:                   query.UpdatedAt.Time,
	}, nil
}

func (r *Repository) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteProfile(ctx, id)
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}
	return nil
}
