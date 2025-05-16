package viewmodel

import "errors"

// Witness represents the view model for a Witness API response.
type Witness struct {
	TransactionHash []byte `json:"transaction_hash"`
	PlutusData      [][]byte `json:"plutus_data"`
	PlutusV1Scripts [][]byte `json:"plutus_v1_scripts"`
	PlutusV2Scripts [][]byte `json:"plutus_v2_scripts"`
	PlutusV3Scripts [][]byte `json:"plutus_v3_scripts"`
	Redeemers       []Redeemer `json:"redeemers"`
}

// IsValid performs validation on the Witness view model.
func (v *Witness) IsValid() error {
	if len(v.TransactionHash) == 0 {
		return errors.New("transaction_hash cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}