package profile

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProfileParameters_IsValid(t *testing.T) {
	t.Run("should pass with valid required fields", func(t *testing.T) {
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail when authentication_credential_id is missing", func(t *testing.T) {
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.Nil,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AuthenticationCredentialId")
	})

	t.Run("should pass with valid first name", func(t *testing.T) {
		firstName := "John"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			FirstName:                  &firstName,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short first name", func(t *testing.T) {
		firstName := "Jo" // Less than 3 characters
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			FirstName:                  &firstName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FirstName")
	})

	t.Run("should fail with too long first name", func(t *testing.T) {
		firstName := "ThisIsAnExtremelyLongFirstNameThatExceedsTheMaximumAllowedLength"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			FirstName:                  &firstName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FirstName")
	})

	t.Run("should pass with valid last name", func(t *testing.T) {
		lastName := "Doe"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			LastName:                   &lastName,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short last name", func(t *testing.T) {
		lastName := "Do" // Less than 3 characters
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			LastName:                   &lastName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LastName")
	})

	t.Run("should fail with too long last name", func(t *testing.T) {
		lastName := "ThisIsAnExtremelyLongLastNameThatExceedsTheMaximumAllowedLength"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			LastName:                   &lastName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LastName")
	})

	t.Run("should pass with valid email", func(t *testing.T) {
		email := "test@example.com"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Email:                      &email,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		email := "invalid-email"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Email:                      &email,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("should fail with too long email", func(t *testing.T) {
		email := "thisissuchalongemailaddressthatitexceedsthemaximumallowedlength@example.com"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Email:                      &email,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("should pass with valid bio", func(t *testing.T) {
		bio := "This is a valid bio with enough characters."
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Bio:                        &bio,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short bio", func(t *testing.T) {
		bio := "Short" // Less than 10 characters
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Bio:                        &bio,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bio")
	})

	t.Run("should fail with too long bio", func(t *testing.T) {
		bio := "This is an extremely long bio that exceeds the maximum allowed length of 255 characters. " +
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Bio:                        &bio,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bio")
	})

	t.Run("should pass with valid E.164 format phone number", func(t *testing.T) {
		phoneNumber := "+66812345678" // Valid E.164 format with Thailand country code
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			PhoneNumber:                &phoneNumber,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should pass with valid US E.164 format phone number", func(t *testing.T) {
		phoneNumber := "+15551234567" // Valid E.164 format with US country code
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			PhoneNumber:                &phoneNumber,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with phone number without country code", func(t *testing.T) {
		phoneNumber := "0812345678" // Missing + prefix and country code
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			PhoneNumber:                &phoneNumber,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PhoneNumber")
	})

	t.Run("should fail with invalid phone number format", func(t *testing.T) {
		phoneNumber := "invalid-phone" // Invalid format
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			PhoneNumber:                &phoneNumber,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PhoneNumber")
	})

	t.Run("should pass with valid address", func(t *testing.T) {
		address := "123 Main Street, City"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Address:                    &address,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short address", func(t *testing.T) {
		address := "Short" // Less than 10 characters
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Address:                    &address,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Address")
	})

	t.Run("should fail with too long address", func(t *testing.T) {
		address := "This is an extremely long address that exceeds the maximum allowed length of 255 characters. " +
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris."
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			Address:                    &address,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Address")
	})

	t.Run("should pass with valid academic institution", func(t *testing.T) {
		academicInstitution := "MIT"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			AcademicInstitution:        &academicInstitution,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short academic institution", func(t *testing.T) {
		academicInstitution := "AB" // Less than 3 characters
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			AcademicInstitution:        &academicInstitution,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AcademicInstitution")
	})

	t.Run("should pass with valid academic email", func(t *testing.T) {
		academicEmail := "student@university.edu"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			AcademicEmail:              &academicEmail,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid academic email format", func(t *testing.T) {
		academicEmail := "invalid-email"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			AcademicEmail:              &academicEmail,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AcademicEmail")
	})

	t.Run("should pass with valid profile picture URL", func(t *testing.T) {
		profilePictureUrl := "https://example.com/profile.jpg"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			ProfilePictureUrl:          &profilePictureUrl,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid profile picture URL", func(t *testing.T) {
		profilePictureUrl := "not-a-valid-url"
		params := CreateProfileParameters{
			AuthenticationCredentialId: uuid.New(),
			ProfilePictureUrl:          &profilePictureUrl,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ProfilePictureUrl")
	})

	t.Run("should pass with all valid fields", func(t *testing.T) {
		firstName := "John"
		lastName := "Doe"
		email := "john.doe@example.com"
		bio := "This is a valid bio with enough characters for the test."
		phoneNumber := "+66812345678" // Valid E.164 format
		address := "123 Main Street, City, Country"
		academicInstitution := "MIT"
		academicEmail := "john.doe@mit.edu"
		profilePictureUrl := "https://example.com/profile.jpg"

		params := CreateProfileParameters{
			AuthenticationCredentialId:  uuid.New(),
			IsProfilePicturePublic:      true,
			ProfilePictureUrl:           &profilePictureUrl,
			IsFirstNamePublic:           true,
			FirstName:                   &firstName,
			IsLastNamePublic:            true,
			LastName:                    &lastName,
			IsEmailPublic:               true,
			Email:                       &email,
			IsBioPublic:                 true,
			Bio:                         &bio,
			IsPhoneNumberPublic:         true,
			PhoneNumber:                 &phoneNumber,
			IsAddressPublic:             true,
			Address:                     &address,
			IsAcademicInstitutionPublic: true,
			AcademicInstitution:         &academicInstitution,
			IsAcademicEmailPublic:       true,
			AcademicEmail:               &academicEmail,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})
}

func TestUpdateProfileParameters_IsValid(t *testing.T) {
	t.Run("should pass with no fields (all optional)", func(t *testing.T) {
		params := UpdateProfileParameters{}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should pass with valid first name", func(t *testing.T) {
		firstName := "John"
		params := UpdateProfileParameters{
			FirstName: &firstName,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short first name", func(t *testing.T) {
		firstName := "Jo"
		params := UpdateProfileParameters{
			FirstName: &firstName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FirstName")
	})

	t.Run("should fail with too long first name", func(t *testing.T) {
		firstName := "ThisIsAnExtremelyLongFirstNameThatExceedsTheMaximumAllowedLength"
		params := UpdateProfileParameters{
			FirstName: &firstName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FirstName")
	})

	t.Run("should pass with valid last name", func(t *testing.T) {
		lastName := "Doe"
		params := UpdateProfileParameters{
			LastName: &lastName,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short last name", func(t *testing.T) {
		lastName := "Do"
		params := UpdateProfileParameters{
			LastName: &lastName,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "LastName")
	})

	t.Run("should pass with valid email", func(t *testing.T) {
		email := "test@example.com"
		params := UpdateProfileParameters{
			Email: &email,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		email := "invalid-email"
		params := UpdateProfileParameters{
			Email: &email,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})

	t.Run("should pass with valid bio", func(t *testing.T) {
		bio := "This is a valid bio with enough characters."
		params := UpdateProfileParameters{
			Bio: &bio,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short bio", func(t *testing.T) {
		bio := "Short"
		params := UpdateProfileParameters{
			Bio: &bio,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Bio")
	})

	t.Run("should pass with valid E.164 format phone number", func(t *testing.T) {
		phoneNumber := "+66812345678" // Valid E.164 format with Thailand country code
		params := UpdateProfileParameters{
			PhoneNumber: &phoneNumber,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should pass with valid US E.164 format phone number", func(t *testing.T) {
		phoneNumber := "+15551234567" // Valid E.164 format with US country code
		params := UpdateProfileParameters{
			PhoneNumber: &phoneNumber,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with phone number without country code", func(t *testing.T) {
		phoneNumber := "0812345678" // Missing + prefix and country code
		params := UpdateProfileParameters{
			PhoneNumber: &phoneNumber,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PhoneNumber")
	})

	t.Run("should fail with invalid phone number format", func(t *testing.T) {
		phoneNumber := "invalid-phone" // Invalid format
		params := UpdateProfileParameters{
			PhoneNumber: &phoneNumber,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PhoneNumber")
	})

	t.Run("should pass with valid address", func(t *testing.T) {
		address := "123 Main Street, City"
		params := UpdateProfileParameters{
			Address: &address,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short address", func(t *testing.T) {
		address := "Short"
		params := UpdateProfileParameters{
			Address: &address,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Address")
	})

	t.Run("should pass with valid academic institution", func(t *testing.T) {
		academicInstitution := "MIT"
		params := UpdateProfileParameters{
			AcademicInstitution: &academicInstitution,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with too short academic institution", func(t *testing.T) {
		academicInstitution := "AB"
		params := UpdateProfileParameters{
			AcademicInstitution: &academicInstitution,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AcademicInstitution")
	})

	t.Run("should pass with valid academic email", func(t *testing.T) {
		academicEmail := "student@university.edu"
		params := UpdateProfileParameters{
			AcademicEmail: &academicEmail,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid academic email format", func(t *testing.T) {
		academicEmail := "invalid-email"
		params := UpdateProfileParameters{
			AcademicEmail: &academicEmail,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "AcademicEmail")
	})

	t.Run("should pass with valid profile picture URL", func(t *testing.T) {
		profilePictureUrl := "https://example.com/profile.jpg"
		params := UpdateProfileParameters{
			ProfilePictureUrl: &profilePictureUrl,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail with invalid profile picture URL", func(t *testing.T) {
		profilePictureUrl := "not-a-valid-url"
		params := UpdateProfileParameters{
			ProfilePictureUrl: &profilePictureUrl,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ProfilePictureUrl")
	})

	t.Run("should pass with multiple valid fields", func(t *testing.T) {
		firstName := "John"
		email := "john.doe@example.com"
		bio := "Updated bio with enough characters for the test."
		isProfilePicturePublic := true

		params := UpdateProfileParameters{
			FirstName:              &firstName,
			Email:                  &email,
			Bio:                    &bio,
			IsProfilePicturePublic: &isProfilePicturePublic,
		}

		err := params.IsValid()
		require.NoError(t, err)
	})

	t.Run("should fail when one field is invalid among multiple fields", func(t *testing.T) {
		firstName := "John"
		email := "invalid-email" // Invalid
		bio := "Valid bio with enough characters for the test."

		params := UpdateProfileParameters{
			FirstName: &firstName,
			Email:     &email,
			Bio:       &bio,
		}

		err := params.IsValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Email")
	})
}
