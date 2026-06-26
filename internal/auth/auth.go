// Package auth provides passwordless (magic-link) authentication and
// cookie-based session management for the PCC-TimeBot web API.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/tpoulsen/pcc-timebot/shared/database"
)

const (
	// LoginTokenTTL is how long a magic-link token remains valid.
	LoginTokenTTL = 15 * time.Minute
	// SessionTTL is the rolling session window. Each authenticated request
	// resets the expiry to now + SessionTTL, so active users stay logged in.
	SessionTTL = 15 * 24 * time.Hour
	// SessionCookieName is the name of the session cookie.
	SessionCookieName = "pcc_session"
)

// generateToken returns a cryptographically random URL-safe token string.
func generateToken(numBytes int) (string, error) {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken returns the hex-encoded SHA-256 hash of a token. Only the hash is
// persisted, so a database leak does not expose usable tokens.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreateLoginToken generates a magic-link token for the given employee, stores
// only its hash, and returns the raw token to embed in the emailed link.
func CreateLoginToken(employeeID int) (string, error) {
	token, err := generateToken(32)
	if err != nil {
		return "", err
	}
	if err := database.CreateLoginToken(employeeID, hashToken(token), time.Now().Add(LoginTokenTTL)); err != nil {
		return "", err
	}
	return token, nil
}

// ConsumeLoginToken validates and consumes a magic-link token, returning the
// authenticated employee ID. The token is single-use.
func ConsumeLoginToken(token string) (int, error) {
	if token == "" {
		return 0, fmt.Errorf("missing token")
	}
	return database.ConsumeLoginToken(hashToken(token))
}

// CreateSession issues a new session for an employee and returns the opaque
// session id to store in the client's cookie.
func CreateSession(employeeID int) (string, time.Time, error) {
	sessionID, err := generateToken(32)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Now().Add(SessionTTL)
	if err := database.CreateSession(sessionID, employeeID, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return sessionID, expiresAt, nil
}
