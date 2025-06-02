package models

type Datum struct {
	ID          uint   `gorm:"primaryKey"`
	UTxOID      []byte `gorm:"type:blob;index:datum_utxo_idx;column:utxo_id" json:"utxo_id"`
	UTxOIDIndex uint32 `gorm:"index:datum_utxo_idx;column:utxo_index" json:"utxo_index"`
	DatumHash   []byte `gorm:"type:blob;not null;unique" json:"datum_hash"`
	DatumCbor   []byte `gorm:"type:blob;not null" json:"datum_cbor"`
}

func (Datum) TableName() string {
	return "datum"
}
