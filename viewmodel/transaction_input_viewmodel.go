package viewmodel

import "errors"

// TransactionInput represents the view model for a TransactionInput API response.
type TransactionInput struct {
	TransactionHash string `json:"transaction_hash"`
	UTxOID          string `json:"utxo_id"`
	UTxOIDIndex     uint32 `json:"utxo_index"`
	Address         string `json:"address"`
	Amount          uint64 `json:"amount"`
	Asset           []Asset `json:"asset"`
	Datum           Datum   `json:"datum"`
	Cbor            string `json:"cbor"` // CBOR string representation
}

// IsValid performs validation on the TransactionInput view model.
func (v *TransactionInput) IsValid() error {
	if v.TransactionHash == "" {
		return errors.New("transaction_hash cannot be empty")
	}
	if v.UTxOID == "" {
		return errors.New("utxo_id cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}