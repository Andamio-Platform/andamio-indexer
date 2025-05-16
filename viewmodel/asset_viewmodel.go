package viewmodel

import "errors"

// Asset represents the view model for an Asset API response.
type Asset struct {
	UTxOID      []byte `json:"utxo_id"`
	UTxOIDIndex uint32 `json:"utxo_index"`
	Name        []byte `json:"name"`
	NameHex     []byte `json:"name_hex"`
	PolicyId    []byte `json:"policy_id"`
	Fingerprint []byte `json:"fingerprint"`
	Amount      uint64 `json:"amount"`
}

// IsValid performs validation on the Asset view model.
func (v *Asset) IsValid() error {
	if len(v.UTxOID) == 0 {
		return errors.New("utxo_id cannot be empty")
	}
	if len(v.PolicyId) == 0 {
		return errors.New("policy_id cannot be empty")
	}
	if len(v.Fingerprint) == 0 {
		return errors.New("fingerprint cannot be empty")
	}
	return nil
}