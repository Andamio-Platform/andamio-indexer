package viewmodel

import "errors"

// Datum represents the view model for a Datum API response.
type Datum struct {
	UTxOID              string `json:"utxo_id"`
	UTxOIDIndex         uint32 `json:"utxo_index"`
	DatumHash           string `json:"datum_hash"`
	DatumCbor           string `json:"datum_cbor"`
}

// IsValid performs validation on the Datum view model.
func (v *Datum) IsValid() error {
	if v.DatumHash == "" {
		return errors.New("datum_hash cannot be empty")
	}
	if v.DatumCbor == "" {
		return errors.New("datum_cbor cannot be empty")
	}
	return nil
}