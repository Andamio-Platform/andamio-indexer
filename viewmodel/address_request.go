package viewmodel

import (
	"fmt"
	"regexp"
)

type AddressRequest struct {
	Address string `json:"address"`
}

func (ar *AddressRequest) IsValid() error {
	// Regular expression for a valid address (example: Cardano address)
	addressRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if !addressRegex.MatchString(ar.Address) {
		return fmt.Errorf("invalid address format")
	}

	return nil
}
