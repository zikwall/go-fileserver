package lib

import (
	"crypto/rand"
	"fmt"
)

// todo use more powerful protection mechanism
func GenerateToken() (string, error) {
	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
