package viewmodel

import "errors"

// SimpleUTxO represents the view model for a SimpleUTxO API response.
type SimpleUTxO struct {
	TransactionHash []byte `json:"transaction_hash"`
	UTxOID          []byte `json:"utxo_id"`
	UTxOIDIndex     uint32 `json:"utxo_index"`
}

// IsValid performs validation on the SimpleUTxO view model.
func (v *SimpleUTxO) IsValid() error {
	if len(v.TransactionHash) == 0 {
		return errors.New("transaction_hash cannot be empty")
	}
	if len(v.UTxOID) == 0 {
		return errors.New("utxo_id cannot be empty")
	}
	return nil
}