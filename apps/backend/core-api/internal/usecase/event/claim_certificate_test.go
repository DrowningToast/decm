package event_test

import (
	"apps/backend/common/pgmapper"
	"apps/backend/core-api/internal/entity"
	"apps/backend/core-api/internal/usecase/cyptoutils"
	"apps/backend/services/auth"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	event_usecase "apps/backend/core-api/internal/usecase/event"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ============================================
// TEST HELPERS
// ============================================

func setupTestUsecase(_ *testing.T) *event_usecase.EventUsecase {
	// For unit tests without blockchain/database, we just test logic
	uc := &event_usecase.EventUsecase{}
	return uc
}

func createTestPrivateKey(t *testing.T) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return privateKey, address
}

func createTestCertificate(credentialId uuid.UUID, _ string, tokenId *string, revoked bool) *entity.EventCertificate {
	eventId := uuid.New()
	certId := uuid.New()
	contractAddr := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"
	name := "John Doe"

	cert := &entity.EventCertificate{
		Id:                      certId,
		EventId:                 eventId,
		ReceiverCredentialId:    &credentialId,
		ReceiverEmail:           nil,
		Name:                    &name,
		AcademicInstitution:     strPtr("Chulalongkorn University"),
		CertificateTitle:        strPtr("Certificate of Completion"),
		CertificateSubtitle:     strPtr("Blockchain Development"),
		EventCertificateAddress: &contractAddr,
		CertificateTokenId:      tokenId,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	if revoked {
		now := time.Now()
		cert.RevokedAt = &now
	}

	return cert
}

func createTestAttendee(eventId, credentialId uuid.UUID, walletAddress string) *entity.EventAttendee {
	email := "john.doe@example.com"
	firstName := "John"
	lastName := "Doe"
	bio := "Software engineer"
	phone := "+66812345678"
	address := "123 Main St"
	institution := "Chulalongkorn University"
	academicEmail := "john.doe@student.chula.ac.th"
	contractAddress := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"

	return &entity.EventAttendee{
		Id:                   uuid.New(),
		EventId:              eventId,
		AttendeeCredentialId: credentialId,
		WalletAddress:        walletAddress,
		ContractAddress:      contractAddress,
		IsAttendeeAccepted:   true,
		Email:                &email,
		FirstName:            &firstName,
		LastName:             &lastName,
		Bio:                  &bio,
		PhoneNumber:          &phone,
		Address:              &address,
		AcademicInstitution:  &institution,
		AcademicEmail:        &academicEmail,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

func createTestJwtClaims(credentialId uuid.UUID, walletAddress string) *auth.JwtClaims {
	email := "john.doe@example.com"
	return &auth.JwtClaims{
		UserId:        credentialId,
		WalletAddress: walletAddress,
		Email:         &email,
	}
}

func strPtr(s string) *string {
	return &s
}

// ============================================
// ELIGIBILITY TESTS
// ============================================

func TestCheckClaimEligibility_Success_CredentialIdMatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	credentialId := uuid.New()
	certificate := createTestCertificate(credentialId, "0xAddress", nil, false)
	currentUser := createTestJwtClaims(credentialId, "0xAddress")

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)
}

func TestCheckClaimEligibility_Success_EmailMatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	email := "john.doe@example.com"
	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           &email,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: "0xAddress",
		Email:         &email,
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)
}

func TestCheckClaimEligibility_Success_OpenCertificate(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           nil,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: "0xAddress",
		Email:         strPtr("any@email.com"),
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)
}

func TestCheckClaimEligibility_Fail_NotAuthenticated(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certificate := createTestCertificate(uuid.New(), "0xAddress", nil, false)

	err := uc.CheckClaimEligibility(ctx, certificate, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestCheckClaimEligibility_Fail_CertificateRevoked(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	credentialId := uuid.New()
	certificate := createTestCertificate(credentialId, "0xAddress", nil, true)
	currentUser := createTestJwtClaims(credentialId, "0xAddress")

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "revoked")
}

func TestCheckClaimEligibility_Fail_NotPublished(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	credentialId := uuid.New()
	certificate := createTestCertificate(credentialId, "0xAddress", nil, false)
	certificate.EventCertificateAddress = nil // Not published
	currentUser := createTestJwtClaims(credentialId, "0xAddress")

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not published")
}

func TestCheckClaimEligibility_Fail_AlreadyClaimed(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	credentialId := uuid.New()
	tokenId := "123"
	certificate := createTestCertificate(credentialId, "0xAddress", &tokenId, false)
	currentUser := createTestJwtClaims(credentialId, "0xAddress")

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "certificate_already_claimed")
}

func TestCheckClaimEligibility_Fail_CredentialIdMismatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certificateOwner := uuid.New()
	certificate := createTestCertificate(certificateOwner, "0xAddress", nil, false)

	differentUser := uuid.New()
	currentUser := createTestJwtClaims(differentUser, "0xAddress")

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not_eligible")
}

func TestCheckClaimEligibility_Success_WalletAddressMatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	walletAddress := "0x742d35cc6634c0532925a3b844bc9e7595f0beb1"
	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           nil,
		ReceiverWalletAddress:   &walletAddress,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: walletAddress,
		Email:         strPtr("byok@example.com"),
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)
}

func TestCheckClaimEligibility_Success_WalletAddressMatchCaseInsensitive(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	// Certificate stores checksummed address; user JWT has lowercase
	certWalletAddress := "0x742d35CC6634C0532925a3b844Bc9e7595f0bEb1"
	userWalletAddress := "0x742d35cc6634c0532925a3b844bc9e7595f0beb1"

	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           nil,
		ReceiverWalletAddress:   &certWalletAddress,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: userWalletAddress,
		Email:         strPtr("byok@example.com"),
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)
}

func TestCheckClaimEligibility_Fail_WalletAddressMismatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certWalletAddress := "0x742d35cc6634c0532925a3b844bc9e7595f0beb1"
	differentWalletAddress := "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           nil,
		ReceiverWalletAddress:   &certWalletAddress,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: differentWalletAddress,
		Email:         strPtr("byok@example.com"),
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not_eligible")
}

func TestCheckClaimEligibility_Fail_WalletAddressEmpty(t *testing.T) {
	// User whose JWT has no wallet address cannot claim a wallet-gated certificate
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certWalletAddress := "0x742d35cc6634c0532925a3b844bc9e7595f0beb1"
	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           nil,
		ReceiverWalletAddress:   &certWalletAddress,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: "", // no wallet address in JWT
		Email:         strPtr("byok@example.com"),
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not_eligible")
}

func TestCheckClaimEligibility_Fail_EmailMismatch(t *testing.T) {
	uc := setupTestUsecase(t)
	ctx := context.Background()

	certificateEmail := "certificate@example.com"
	certificate := &entity.EventCertificate{
		Id:                      uuid.New(),
		EventId:                 uuid.New(),
		ReceiverCredentialId:    nil,
		ReceiverEmail:           &certificateEmail,
		EventCertificateAddress: strPtr("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1"),
		CertificateTokenId:      nil,
		RevokedAt:               nil,
		CreatedAt:               time.Now(),
	}

	differentEmail := "different@example.com"
	currentUser := &auth.JwtClaims{
		UserId:        uuid.New(),
		WalletAddress: "0xAddress",
		Email:         &differentEmail,
	}

	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not_eligible")
}

// ============================================
// SIGNATURE VALIDATION TESTS
// ============================================

func TestClaimCertificateWithSignature_ValidSignature(t *testing.T) {
	// This test validates signature verification logic
	privateKey, address := createTestPrivateKey(t)

	signMessage := "I want to claim certificate 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1 with deadline 12345678"

	// Sign the message
	hash := cyptoutils.HashEthereumMessage(signMessage)
	hashBytes := hash[:]
	signature, err := crypto.Sign(hashBytes, privateKey)
	assert.NoError(t, err)
	assert.Equal(t, 65, len(signature))

	// Verify signature
	isValid, err := cyptoutils.VerifySignatureByDigest(address, hash, signature)
	assert.NoError(t, err)
	assert.True(t, isValid)

	// Recover public key
	recoveredPubKey, err := cyptoutils.RecoverPublicKeyFromSignature(hash, signature)
	assert.NoError(t, err)
	assert.NotNil(t, recoveredPubKey)

	// Verify recovered address matches
	recoveredAddress := cyptoutils.PublicKeyToAddress(recoveredPubKey)
	assert.Equal(t, address, recoveredAddress)
}

func TestClaimCertificateWithSignature_InvalidSignature(t *testing.T) {
	_, address := createTestPrivateKey(t)

	signMessage := "I want to claim certificate 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1 with deadline 12345678"
	hash := cyptoutils.HashEthereumMessage(signMessage)

	// Create random invalid signature
	invalidSignature := make([]byte, 65)
	for i := range invalidSignature {
		invalidSignature[i] = byte(i % 256)
	}

	// Verify should fail (returns error for invalid signature)
	isValid, err := cyptoutils.VerifySignatureByDigest(address, hash, invalidSignature)
	if err != nil {
		// Invalid signature causes error, which is expected
		assert.Error(t, err)
	} else {
		// If no error, verification should be false
		assert.False(t, isValid)
	}
}

// ============================================
// DATA ENCRYPTION TESTS
// ============================================

func TestDualEncryption_AttendeeData(t *testing.T) {
	// Test attendee profile data encryption (JSON format)
	privateKey, _ := createTestPrivateKey(t)
	publicKey := privateKey.PublicKey

	attendeeJSON := `{"first_name":"John","last_name":"Doe","email":"john@example.com","bio":"Engineer","phone_number":"+66812345678","address":"123 Main St","academic_institution":"Chula","academic_email":"john@student.chula.ac.th"}`

	// Test ECIES encryption (user-facing)
	encryptedData, err := cyptoutils.EncryptWithPublicKeyBytes(attendeeJSON, &publicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData)
	assert.NotEqual(t, attendeeJSON, encryptedData)

	// Test AES-GCM encryption (backend-facing)
	piiKey := "test-encryption-key-32-bytes!!"
	backendEncrypted, err := pgmapper.EncryptPII(attendeeJSON, piiKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, backendEncrypted)
	assert.NotEqual(t, attendeeJSON, backendEncrypted)

	// Test decryption
	decrypted, err := pgmapper.DecryptPII(backendEncrypted, piiKey)
	assert.NoError(t, err)
	assert.Equal(t, attendeeJSON, decrypted)
}

func TestDualEncryption_CertificatePII(t *testing.T) {
	// Test certificate PII encryption (CSV format)
	privateKey, _ := createTestPrivateKey(t)
	publicKey := privateKey.PublicKey

	certificateCSV := "John Doe,Chulalongkorn University,Certificate of Completion,Blockchain Development"

	// Test ECIES encryption
	encryptedProof, err := cyptoutils.EncryptWithPublicKeyBytes(certificateCSV, &publicKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedProof)

	// Test AES-GCM encryption
	piiKey := "test-encryption-key-32-bytes!!"
	backendEncryptedProof, err := pgmapper.EncryptPII(certificateCSV, piiKey)
	assert.NoError(t, err)
	assert.NotEmpty(t, backendEncryptedProof)

	// Test decryption
	decrypted, err := pgmapper.DecryptPII(backendEncryptedProof, piiKey)
	assert.NoError(t, err)
	assert.Equal(t, certificateCSV, decrypted)
}

// ============================================
// PASSWORD VALIDATION TESTS
// ============================================

func TestPasswordEncryptionDecryption(t *testing.T) {
	// Test that password encryption/decryption works correctly
	// Note: EncryptPrivateKey/DecryptPrivateKey are in common/cyptoutils
	// This test validates the password flow logic

	t.Skip("Skipping password encryption test - EncryptPrivateKey not implemented yet")
}

// ============================================
// HASH VALIDATION TESTS
// ============================================

func TestCertificateDataHash(t *testing.T) {
	certificateCSV := "John Doe,Chulalongkorn University,Certificate of Completion,Blockchain Development"

	// Generate hash
	hash := cyptoutils.HashMessage(certificateCSV)
	assert.NotEmpty(t, hash)
	assert.Equal(t, 32, len(hash)) // SHA256 produces 32 bytes

	// Verify hash is deterministic
	hash2 := cyptoutils.HashMessage(certificateCSV)
	assert.Equal(t, hash, hash2)

	// Verify different data produces different hash
	differentCSV := "Jane Doe,MIT,Certificate of Achievement,AI Development"
	hash3 := cyptoutils.HashMessage(differentCSV)
	assert.NotEqual(t, hash, hash3)

	// Convert to hex string with 0x prefix
	hashStr := hex.EncodeToString(hash)
	assert.NotEmpty(t, hashStr)
	assert.Equal(t, 64, len(hashStr)) // 32 bytes = 64 hex characters
}

// ============================================
// INTEGRATION-STYLE TESTS (NO BLOCKCHAIN)
// ============================================

func TestClaimWorkflow_CompleteDataFlow(t *testing.T) {
	// This test validates the complete data flow without blockchain
	ctx := context.Background()
	uc := setupTestUsecase(t)

	// Setup test data
	privateKey, address := createTestPrivateKey(t)

	credentialId := uuid.New()
	certificate := createTestCertificate(credentialId, address.Hex(), nil, false)
	attendee := createTestAttendee(certificate.EventId, credentialId, address.Hex())
	currentUser := createTestJwtClaims(credentialId, address.Hex())

	// Step 1: Validate eligibility
	err := uc.CheckClaimEligibility(ctx, certificate, currentUser)
	assert.NoError(t, err)

	// Step 2: Build attendee data (DATA section) - JSON format
	type AttendeeProfileData struct {
		FirstName           *string `json:"first_name"`
		LastName            *string `json:"last_name"`
		Email               *string `json:"email"`
		Bio                 *string `json:"bio"`
		PhoneNumber         *string `json:"phone_number"`
		Address             *string `json:"address"`
		AcademicInstitution *string `json:"academic_institution"`
		AcademicEmail       *string `json:"academic_email"`
	}

	attendeeData := AttendeeProfileData{
		FirstName:           attendee.FirstName,
		LastName:            attendee.LastName,
		Email:               attendee.Email,
		Bio:                 attendee.Bio,
		PhoneNumber:         attendee.PhoneNumber,
		Address:             attendee.Address,
		AcademicInstitution: attendee.AcademicInstitution,
		AcademicEmail:       attendee.AcademicEmail,
	}

	attendeeJSON, err := json.Marshal(attendeeData)
	assert.NoError(t, err)
	assert.NotEmpty(t, attendeeJSON)

	// Step 3: Build certificate PII (PROOF section) - CSV format
	certificateCSV := *certificate.Name + "," + *certificate.AcademicInstitution + "," +
		*certificate.CertificateTitle + "," + *certificate.CertificateSubtitle
	assert.NotEmpty(t, certificateCSV)

	// Step 4: Test dual encryption
	piiKey := "test-encryption-key-32-bytes!!"

	// Encrypt attendee data (DATA section)
	userEncryptedData, err := cyptoutils.EncryptWithPublicKeyBytes(string(attendeeJSON), &privateKey.PublicKey)
	assert.NoError(t, err)
	backendEncryptedData, err := pgmapper.EncryptPII(string(attendeeJSON), piiKey)
	assert.NoError(t, err)

	// Encrypt certificate PII (PROOF section)
	userEncryptedProof, err := cyptoutils.EncryptWithPublicKeyBytes(certificateCSV, &privateKey.PublicKey)
	assert.NoError(t, err)
	backendEncryptedProof, err := pgmapper.EncryptPII(certificateCSV, piiKey)
	assert.NoError(t, err)

	// Step 5: Generate hash
	hash := cyptoutils.HashMessage(certificateCSV)
	assert.NotEmpty(t, hash)

	// Verify all encrypted data is valid
	assert.NotEmpty(t, userEncryptedData)
	assert.NotEmpty(t, backendEncryptedData)
	assert.NotEmpty(t, userEncryptedProof)
	assert.NotEmpty(t, backendEncryptedProof)
}

// ============================================
// SIGN MESSAGE GENERATION TESTS
// ============================================

func TestGetClaimCertificateSignMessage(t *testing.T) {
	// Test sign message format
	contractAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1")
	deadlineBlock := big.NewInt(12345678)

	signMessage := "I want to claim certificate " + contractAddress.Hex() + " with deadline " + deadlineBlock.String()

	assert.NotEmpty(t, signMessage)
	// Note: Hex() returns checksummed address, which may differ in casing
	assert.Contains(t, signMessage, contractAddress.Hex())
	assert.Contains(t, signMessage, "12345678")
	assert.Contains(t, signMessage, "I want to claim certificate")
}

func TestExtractDeadlineBlockFromSignMessage(t *testing.T) {
	// The sign message format is: "signerAddress,contractAddress,deadlineBlock"
	// Example: "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,12345678"

	t.Run("should extract deadline block from valid sign message", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,12345678"

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.NoError(t, err)
		assert.NotNil(t, deadlineBlock)
		assert.Equal(t, uint64(12345678), *deadlineBlock)
	})

	t.Run("should fail when sign message has wrong number of parts", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,12345678" // Only 2 parts

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.Error(t, err)
		assert.Nil(t, deadlineBlock)
		assert.Contains(t, err.Error(), "invalid sign message")
	})

	t.Run("should fail when deadline block is not a valid number", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,invalid"

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.Error(t, err)
		assert.Nil(t, deadlineBlock)
	})

	t.Run("should handle large deadline block numbers", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,18446744073709551615" // Max uint64

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.NoError(t, err)
		assert.NotNil(t, deadlineBlock)
		assert.Equal(t, uint64(18446744073709551615), *deadlineBlock)
	})

	t.Run("should fail when deadline block is negative", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,-1"

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.Error(t, err)
		assert.Nil(t, deadlineBlock)
	})

	t.Run("should extract zero as valid deadline block", func(t *testing.T) {
		signMessage := "0x1234567890123456789012345678901234567890,0xAbCdEf1234567890123456789012345678901234,0"

		deadlineBlock, err := cyptoutils.ExtractDeadlineBlockFromSignMessage(signMessage)

		assert.NoError(t, err)
		assert.NotNil(t, deadlineBlock)
		assert.Equal(t, uint64(0), *deadlineBlock)
	})
}
