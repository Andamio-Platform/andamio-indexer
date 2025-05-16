package viewmodel

import "errors"

// Datum represents the view model for a Datum API response.
type Datum struct {
	UTxOID              []byte `json:"utxo_id"`
	UTxOIDIndex         uint32 `json:"utxo_index"`
	DatumHash           []byte `json:"datum_hash"`
	DatumCbor           []byte `json:"datum_cbor"`
}

// IsValid performs validation on the Datum view model.
func (v *Datum) IsValid() error {
	if len(v.DatumHash) == 0 {
		return errors.New("datum_hash cannot be empty")
	}
	if len(v.DatumCbor) == 0 {
		return errors.New("datum_cbor cannot be empty")
	}
	return nil
}