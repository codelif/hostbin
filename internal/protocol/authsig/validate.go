package authsig

import (
	"fmt"
	"strings"
)

func ValidateSharedSecret(secret string) error {
	if len(strings.TrimSpace(secret)) < 32 {
		return fmt.Errorf("shared secret must be at least 32 bytes")
	}

	return nil
}
