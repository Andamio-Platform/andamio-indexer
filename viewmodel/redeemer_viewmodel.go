package viewmodel

import "errors"

// Redeemer represents the view model for a Redeemer API response.
type Redeemer struct {
	WitnessID uint   `json:"witness_id"`
	Index     uint   `json:"index"`
	Tag       uint8  `json:"tag"`
	Cbor      []byte `json:"cbor"`
}

// IsValid performs validation on the Redeemer view model.
func (v *Redeemer) IsValid() error {
	if len(v.Cbor) == 0 {
		return errors.New("cbor cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}