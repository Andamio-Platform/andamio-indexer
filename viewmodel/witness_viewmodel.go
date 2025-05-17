package viewmodel

import "errors"

// Witness represents the view model for a Witness API response.
type Witness struct {
	TransactionHash string `json:"transaction_hash"`
	PlutusData      []string `json:"plutus_data"` // Slice of CBOR string representations
	PlutusV1Scripts []string `json:"plutus_v1_scripts"` // Slice of CBOR string representations
	PlutusV2Scripts []string `json:"plutus_v2_scripts"` // Slice of CBOR string representations
	PlutusV3Scripts []string `json:"plutus_v3_scripts"` // Slice of CBOR string representations
	Redeemers       []Redeemer `json:"redeemers"`
}

// IsValid performs validation on the Witness view model.
func (v *Witness) IsValid() error {
	if v.TransactionHash == "" {
		return errors.New("transaction_hash cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}