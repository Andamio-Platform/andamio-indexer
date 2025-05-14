package models

type Datum struct {
	ID                  uint   `gorm:"primaryKey"`
	UTxOID              []byte `gorm:"type:blob;index:utxo_idx" json:"utxo_id"` // Removed standalone index tag as it's not effective on blob types, composite index is on UTxOIDIndex
	UTxOIDIndex         uint32 `gorm:"index:utxo_idx" json:"utxo_index"` // Foreign key to TransactionInput and TransactionOutput and the reference would be UTxOIDIndex of either TransactionInput or TransactionOutput
	DatumHash           []byte `gorm:"type:blob;not null;unique" json:"datum_hash"` // Removed index tag as it's not effective on blob types
	DatumCbor           []byte `gorm:"type:blob;not null" json:"datum_cbor"` // Removed index tag as it's not effective on blob types and fixed typo (double semicolon)
}

func (Datum) TableName() string {
	return "datum"
}
