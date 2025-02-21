package lock

import (
	"crypto/rand"
	"encoding/hex"
)

// =============================================================================
// Helper Functions
// =============================================================================

// generateRandomToken returns a securely generated random token (hex-encoded).
// This helps identify the owner of the lock.
func generateRandomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
