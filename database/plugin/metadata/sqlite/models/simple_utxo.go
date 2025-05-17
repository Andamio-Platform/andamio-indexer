package models

type SimpleUTxO struct {
	ID              uint   `gorm:"primaryKey"`
	TransactionHash []byte `gorm:"type:blob" json:"transaction_hash"` // Removed index tag as it's not effective on blob types
	UTxOID          []byte `gorm:"type:blob" json:"utxo_id"` // Removed index tag as it's not effective on blob types, composite index is on UTxOIDIndex
	UTxOIDIndex     uint32 `gorm:"index:simple_utxo_idx" json:"utxo_index"`
}

func (SimpleUTxO) TableName() string {
	return "simple_utxos"
}
