package hashutils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Configuration struct {
	HashRaw    []byte
	Salt       []byte
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
	KeyLength  uint32
}

func GenerateCryptographicSalt(saltSize uint32) ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}
	return salt, nil
}

// Utility: Hash password with default Argon2id configuration
func HashPassword(password string) (string, error) {
	config := &Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	salt, err := GenerateCryptographicSalt(16)
	if err != nil {
		return "", fmt.Errorf("password hashing failed: %w", err)
	}
	config.Salt = salt

	// Execute Argon2id hashing algorithm
	config.HashRaw = argon2.IDKey(
		[]byte(password),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Threads,
		config.KeyLength,
	)

	// Generate standardized hash format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.MemoryCost,
		config.TimeCost,
		config.Threads,
		base64.RawStdEncoding.EncodeToString(config.Salt),
		base64.RawStdEncoding.EncodeToString(config.HashRaw),
	)

	return encodedHash, nil
}

func HashPasswordWithCustomSalt(password string, salt []byte) (string, error) {
	config := &Argon2Configuration{
		TimeCost:   2,
		MemoryCost: 64 * 1024,
		Threads:    4,
		KeyLength:  32,
	}

	config.Salt = salt

	config.HashRaw = argon2.IDKey(
		[]byte(password),
		config.Salt,
		config.TimeCost,
		config.MemoryCost,
		config.Threads,
		config.KeyLength,
	)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.MemoryCost,
		config.TimeCost,
		config.Threads,
		base64.RawStdEncoding.EncodeToString(config.Salt),
		base64.RawStdEncoding.EncodeToString(config.HashRaw),
	)
	return encodedHash, nil
}

func CompareHash(password, hashedPassword string) (bool, error) {
	var version int
	var memory, time uint32
	var threads uint8
	var salt []byte

	_, err := fmt.Sscanf(hashedPassword, "$argon2id$v=%d$m=%d,t=%d,p=%d", &version, &memory, &time, &threads)
	if err != nil {
		return false, fmt.Errorf("invalid hash format: %w", err)
	}

	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// Hash SHA256
func HashSHA256(data string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
