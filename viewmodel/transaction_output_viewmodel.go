package viewmodel

import "errors"

// TransactionOutput represents the view model for a TransactionOutput API response.
type TransactionOutput struct {
	UTxOID      []byte `json:"utxo_id"`
	UTxOIDIndex uint32 `json:"utxo_index"`
	Address     []byte `json:"address"`
	Amount      uint64 `json:"amount"`
	Asset       []Asset `json:"asset"`
	Datum       Datum   `json:"datum"`
	Cbor        []byte `json:"cbor"`
}

// IsValid performs validation on the TransactionOutput view model.
func (v *TransactionOutput) IsValid() error {
	if len(v.UTxOID) == 0 {
		return errors.New("utxo_id cannot be empty")
	}
	if len(v.Address) == 0 {
		return errors.New("address cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}