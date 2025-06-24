package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies if the provided password matches the hashed password
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomToken generates a random token for email verification, password reset, etc.
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateEmailVerificationToken generates a token for email verification
func GenerateEmailVerificationToken() (string, error) {
	return GenerateRandomToken(32) // 64 character hex string
}

// GeneratePasswordResetToken generates a token for password reset
func GeneratePasswordResetToken() (string, error) {
	return GenerateRandomToken(32) // 64 character hex string
}
