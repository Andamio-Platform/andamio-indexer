package viewmodel

import "errors"

// TransactionInput represents the view model for a TransactionInput API response.
type TransactionInput struct {
	TransactionHash []byte `json:"transaction_hash"`
	UTxOID          []byte `json:"utxo_id"`
	UTxOIDIndex     uint32 `json:"utxo_index"`
	Address         []byte `json:"address"`
	Amount          uint64 `json:"amount"`
	Asset           []Asset `json:"asset"`
	Datum           Datum   `json:"datum"`
	Cbor            []byte `json:"cbor"`
}

// IsValid performs validation on the TransactionInput view model.
func (v *TransactionInput) IsValid() error {
	if len(v.TransactionHash) == 0 {
		return errors.New("transaction_hash cannot be empty")
	}
	if len(v.UTxOID) == 0 {
		return errors.New("utxo_id cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}