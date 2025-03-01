package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

var defaultPasswordConfig = PasswordConfig{
	time:    1,
	memory:  64 * 1024,
	threads: 4,
	keyLen:  32,
}

// HashPassword creates a secure hash using Argon2id
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	cfg := defaultPasswordConfig
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		cfg.time,
		cfg.memory,
		cfg.threads,
		cfg.keyLen,
	)

	// Format: $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		cfg.memory,
		cfg.time,
		cfg.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

// VerifyPassword checks if a password matches its hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the hash string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var cfg PasswordConfig
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &cfg.memory, &cfg.time, &cfg.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	cfg.keyLen = uint32(len(decodedHash))

	// Compute hash from the password and compare
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		cfg.time,
		cfg.memory,
		cfg.threads,
		cfg.keyLen,
	)

	return subtle.ConstantTimeCompare(hash, decodedHash) == 1, nil
}
