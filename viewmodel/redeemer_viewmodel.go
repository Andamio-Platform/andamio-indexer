package viewmodel

import "errors"

// Redeemer represents the view model for a Redeemer API response.
type Redeemer struct {
	TransactionHash string `json:"transaction_hash"`
	Index           uint   `json:"index"`
	Tag             uint8  `json:"tag"`
	Cbor            string `json:"cbor"` // CBOR string representation
}

// IsValid performs validation on the Redeemer view model.
func (v *Redeemer) IsValid() error {
	if v.Cbor == "" {
		return errors.New("cbor cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}
