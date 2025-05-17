package viewmodel

import "errors"

// TransactionOutput represents the view model for a TransactionOutput API response.
type TransactionOutput struct {
	UTxOID      string `json:"utxo_id"`
	UTxOIDIndex uint32 `json:"utxo_index"`
	Address     string `json:"address"`
	Amount      uint64 `json:"amount"`
	Asset       []Asset `json:"asset"`
	Datum       Datum   `json:"datum"`
	Cbor        string `json:"cbor"` // CBOR string representation
}

// IsValid performs validation on the TransactionOutput view model.
func (v *TransactionOutput) IsValid() error {
	if v.UTxOID == "" {
		return errors.New("utxo_id cannot be empty")
	}
	if v.Address == "" {
		return errors.New("address cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}