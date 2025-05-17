package viewmodel

import "errors"

// Transaction represents the view model for a Transaction API response.
type Transaction struct {
	BlockHash       string              `json:"block_hash"`
	BlockNumber     uint64              `json:"block_number"`
	SlotNumber      uint64              `json:"slot_number"`
	TransactionHash string              `json:"transaction_hash"`
	Inputs          []TransactionInput  `json:"inputs"`
	Outputs         []TransactionOutput `json:"outputs"`
	ReferenceInputs []SimpleUTxO        `json:"reference_inputs"`
	Metadata        string              `json:"metadata"` // CBOR string representation
	Fee             uint64              `json:"fee"`
	TTL             uint64              `json:"ttl"`
	Withdrawals     map[string]uint64   `json:"withdrawals"`
	Witness         Witness             `json:"witness"`
	Certificates    []string            `json:"certificates"` // Slice of CBOR string representations
}

// IsValid performs validation on the Transaction view model.
func (v *Transaction) IsValid() error {
	if v.BlockHash == "" {
		return errors.New("block_hash cannot be empty")
	}
	if v.TransactionHash == "" {
		return errors.New("transaction_hash cannot be empty")
	}
	// Add more specific validation for other fields if needed
	return nil
}