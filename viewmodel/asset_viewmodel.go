package viewmodel

import "errors"

// Asset represents the view model for an Asset API response.
type Asset struct {
	UTxOID      string `json:"utxo_id"`
	UTxOIDIndex uint32 `json:"utxo_index"`
	Name        string `json:"name"`
	NameHex     string `json:"name_hex"`
	PolicyId    string `json:"policy_id"`
	Fingerprint string `json:"fingerprint"`
	Amount      uint64 `json:"amount"`
}

// IsValid performs validation on the Asset view model.
func (v *Asset) IsValid() error {
	if v.UTxOID == "" {
		return errors.New("utxo_id cannot be empty")
	}
	if v.PolicyId == "" {
		return errors.New("policy_id cannot be empty")
	}
	if v.Fingerprint == "" {
		return errors.New("fingerprint cannot be empty")
	}
	return nil
}