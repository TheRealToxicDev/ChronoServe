package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Password hashing parameters
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLength   = 16
)

// HashPassword creates an Argon2id hash of a password
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash the password with Argon2id
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Format the hash parameters
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Return formatted hash string
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory, argonTime, argonThreads, encodedSalt, encodedHash), nil
}

// VerifyPassword checks if a password matches a hashed value
func VerifyPassword(password string, encodedHash string) (bool, error) {
	// Extract parameters from encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Verify hash type and version
	if parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, fmt.Errorf("unsupported hash algorithm or version")
	}

	// Extract parameters
	var memory uint32
	var time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, fmt.Errorf("invalid hash parameters: %v", err)
	}

	// Decode salt and hash values
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("invalid salt: %v", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("invalid hash: %v", err)
	}

	// Compute hash with the same parameters
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(expectedHash)))

	// Compare hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// GenerateRandomString creates a cryptographically secure random string
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateSecureKey generates a secure key for JWT signing
func GenerateSecureKey(length int) (string, error) {
	if length < 32 {
		length = 32 // Minimum recommended length for security
	}

	return GenerateRandomString(length)
}

// IsDefaultPassword checks if a password hash matches any of the known default passwords
func IsDefaultPassword(hash string) bool {
	// Define known default password hashes
	defaultHashes := []string{
		"$argon2id$v=19$m=65536,t=1,p=4$...",
		"$argon2id$v=19$m=65536,t=1,p=4$...",
	}

	// Check if the hash matches any default
	for _, defaultHash := range defaultHashes {
		if hash == defaultHash {
			return true
		}
	}

	return false
}
