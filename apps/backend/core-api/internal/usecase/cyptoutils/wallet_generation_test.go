package cyptoutils

import (
	"apps/backend/common/customerror"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
)

func TestGenerateMnemonic(t *testing.T) {
	tests := []struct {
		name              string
		wordsCount        *int
		expectError       bool
		expectedErrorCode *customerror.ErrCode
		expectedWordCount int
	}{
		{
			name:              "should_generate_12_word_mnemonic_by_default",
			wordsCount:        nil,
			expectError:       false,
			expectedWordCount: 12,
		},
		{
			name: "should_generate_12_word_mnemonic_when_specified",
			wordsCount: func() *int {
				count := 12
				return &count
			}(),
			expectError:       false,
			expectedWordCount: 12,
		},
		{
			name: "should_generate_15_word_mnemonic",
			wordsCount: func() *int {
				count := 15
				return &count
			}(),
			expectError:       false,
			expectedWordCount: 15,
		},
		{
			name: "should_generate_18_word_mnemonic",
			wordsCount: func() *int {
				count := 18
				return &count
			}(),
			expectError:       false,
			expectedWordCount: 18,
		},
		{
			name: "should_generate_21_word_mnemonic",
			wordsCount: func() *int {
				count := 21
				return &count
			}(),
			expectError:       false,
			expectedWordCount: 21,
		},
		{
			name: "should_generate_24_word_mnemonic",
			wordsCount: func() *int {
				count := 24
				return &count
			}(),
			expectError:       false,
			expectedWordCount: 24,
		},
		{
			name: "should_fail_with_invalid_word_count_10",
			wordsCount: func() *int {
				count := 10
				return &count
			}(),
			expectError:       true,
			expectedErrorCode: &customerror.ErrInvalidArgument.Code,
		},
		{
			name: "should_fail_with_invalid_word_count_13",
			wordsCount: func() *int {
				count := 13
				return &count
			}(),
			expectError:       true,
			expectedErrorCode: &customerror.ErrInvalidArgument.Code,
		},
		{
			name: "should_fail_with_invalid_word_count_30",
			wordsCount: func() *int {
				count := 30
				return &count
			}(),
			expectError:       true,
			expectedErrorCode: &customerror.ErrInvalidArgument.Code,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnemonic, err := GenerateMnemonic(tt.wordsCount)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, mnemonic)
				if tt.expectedErrorCode != nil {
					var customError *customerror.Err
					assert.ErrorAs(t, err, &customError)
					assert.Equal(t, *tt.expectedErrorCode, *customError.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mnemonic)

				// Verify mnemonic is valid BIP39
				isValid := bip39.IsMnemonicValid(*mnemonic)
				assert.True(t, isValid, "Generated mnemonic should be valid BIP39")

				// Count words
				words := len(strings.Fields(*mnemonic))
				assert.Equal(t, tt.expectedWordCount, words)
			}
		})
	}
}

func TestGenerateMnemonic_Uniqueness(t *testing.T) {
	t.Run("should_generate_unique_mnemonics_each_time", func(t *testing.T) {
		mnemonic1, err1 := GenerateMnemonic(nil)
		assert.NoError(t, err1)
		assert.NotNil(t, mnemonic1)

		mnemonic2, err2 := GenerateMnemonic(nil)
		assert.NoError(t, err2)
		assert.NotNil(t, mnemonic2)

		// Mnemonics should be different
		assert.NotEqual(t, *mnemonic1, *mnemonic2)

		// Both should be valid
		assert.True(t, bip39.IsMnemonicValid(*mnemonic1))
		assert.True(t, bip39.IsMnemonicValid(*mnemonic2))
	})
}

func TestGenerateSeed(t *testing.T) {
	tests := []struct {
		name               string
		wordsCount         *int
		expectError        bool
		expectedSeedLength int
	}{
		{
			name:               "should_generate_seed_with_12_words",
			wordsCount:         nil,
			expectError:        false,
			expectedSeedLength: 64, // BIP39 seeds are always 64 bytes
		},
		{
			name: "should_generate_seed_with_24_words",
			wordsCount: func() *int {
				count := 24
				return &count
			}(),
			expectError:        false,
			expectedSeedLength: 64,
		},
		{
			name: "should_fail_with_invalid_word_count",
			wordsCount: func() *int {
				count := 10
				return &count
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed, err := GenerateSeed(tt.wordsCount)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, seed)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, seed)
				assert.Equal(t, tt.expectedSeedLength, len(seed))
			}
		})
	}
}

func TestGenerateSeedFromMnemonic(t *testing.T) {
	tests := []struct {
		name               string
		mnemonic           string
		expectError        bool
		expectedSeedLength int
	}{
		{
			name:               "should_generate_seed_from_valid_mnemonic",
			mnemonic:           "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			expectError:        false,
			expectedSeedLength: 64,
		},
		{
			name:               "should_generate_deterministic_seed_from_same_mnemonic",
			mnemonic:           "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
			expectError:        false,
			expectedSeedLength: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed, err := GenerateSeedFromMnemonic(&tt.mnemonic)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, seed)
				assert.Equal(t, tt.expectedSeedLength, len(seed))
			}
		})
	}
}

func TestGenerateSeedFromMnemonic_Deterministic(t *testing.T) {
	t.Run("should_generate_same_seed_from_same_mnemonic", func(t *testing.T) {
		mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

		seed1, err1 := GenerateSeedFromMnemonic(&mnemonic)
		assert.NoError(t, err1)

		seed2, err2 := GenerateSeedFromMnemonic(&mnemonic)
		assert.NoError(t, err2)

		// Seeds should be identical for same mnemonic
		assert.Equal(t, seed1, seed2)
	})
}

func TestGeneratePrivateKeyFromSeed(t *testing.T) {
	tests := []struct {
		name        string
		setupSeed   func() []byte
		expectError bool
	}{
		{
			name: "should_generate_private_key_from_valid_seed",
			setupSeed: func() []byte {
				mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
				seed := bip39.NewSeed(mnemonic, "")
				return seed
			},
			expectError: false,
		},
		{
			name: "should_generate_private_key_from_another_valid_seed",
			setupSeed: func() []byte {
				mnemonic := "legal winner thank year wave sausage worth useful legal winner thank yellow"
				seed := bip39.NewSeed(mnemonic, "")
				return seed
			},
			expectError: false,
		},
		{
			name: "should_handle_short_seed",
			setupSeed: func() []byte {
				// Short seed - will still work but derived from limited entropy
				return []byte("short")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed := tt.setupSeed()
			privateKey, err := GeneratePrivateKeyFromSeed(seed)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, privateKey)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, privateKey)
				assert.NotNil(t, privateKey.D)
				assert.NotNil(t, privateKey.PublicKey.X)
				assert.NotNil(t, privateKey.PublicKey.Y)
			}
		})
	}
}

func TestGeneratePrivateKeyFromSeed_Deterministic(t *testing.T) {
	t.Run("should_generate_same_private_key_from_same_seed", func(t *testing.T) {
		mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		seed := bip39.NewSeed(mnemonic, "")

		privateKey1, err1 := GeneratePrivateKeyFromSeed(seed)
		assert.NoError(t, err1)

		privateKey2, err2 := GeneratePrivateKeyFromSeed(seed)
		assert.NoError(t, err2)

		// Private keys should be identical for same seed
		assert.Equal(t, privateKey1.D.Bytes(), privateKey2.D.Bytes())
		assert.Equal(t, privateKey1.PublicKey.X.Bytes(), privateKey2.PublicKey.X.Bytes())
		assert.Equal(t, privateKey1.PublicKey.Y.Bytes(), privateKey2.PublicKey.Y.Bytes())
	})

	t.Run("should_generate_different_private_keys_from_different_seeds", func(t *testing.T) {
		mnemonic1 := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
		seed1 := bip39.NewSeed(mnemonic1, "")

		mnemonic2 := "legal winner thank year wave sausage worth useful legal winner thank yellow"
		seed2 := bip39.NewSeed(mnemonic2, "")

		privateKey1, err1 := GeneratePrivateKeyFromSeed(seed1)
		assert.NoError(t, err1)

		privateKey2, err2 := GeneratePrivateKeyFromSeed(seed2)
		assert.NoError(t, err2)

		// Private keys should be different for different seeds
		assert.NotEqual(t, privateKey1.D.Bytes(), privateKey2.D.Bytes())
	})
}

func TestFullWalletGenerationFlow(t *testing.T) {
	t.Run("should_complete_full_wallet_generation_flow", func(t *testing.T) {
		// Step 1: Generate mnemonic
		mnemonic, err := GenerateMnemonic(nil)
		assert.NoError(t, err)
		assert.NotNil(t, mnemonic)
		assert.True(t, bip39.IsMnemonicValid(*mnemonic))

		// Step 2: Generate seed from mnemonic
		seed, err := GenerateSeedFromMnemonic(mnemonic)
		assert.NoError(t, err)
		assert.NotNil(t, seed)
		assert.Equal(t, 64, len(seed))

		// Step 3: Generate private key from seed
		privateKey, err := GeneratePrivateKeyFromSeed(seed)
		assert.NoError(t, err)
		assert.NotNil(t, privateKey)
		assert.NotNil(t, privateKey.D)
		assert.NotNil(t, privateKey.PublicKey.X)
		assert.NotNil(t, privateKey.PublicKey.Y)

		// Step 4: Verify private key can derive public key and address
		publicKeyBytes := privateKey.PublicKey.X.Bytes()
		assert.NotEmpty(t, publicKeyBytes)
	})
}
