package viewmodel

import "errors"

// SimpleUTxO represents the view model for a SimpleUTxO API response.
type SimpleUTxO struct {
	TransactionHash string `json:"transaction_hash"`
	UTxOID          string `json:"utxo_id"`
	UTxOIDIndex     uint32 `json:"utxo_index"`
}

// IsValid performs validation on the SimpleUTxO view model.
func (v *SimpleUTxO) IsValid() error {
	if v.TransactionHash == "" {
		return errors.New("transaction_hash cannot be empty")
	}
	if v.UTxOID == "" {
		return errors.New("utxo_id cannot be empty")
	}
	return nil
}