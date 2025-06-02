package models

type SimpleUTxO struct {
	ID              uint   `gorm:"primaryKey"`
	TransactionHash []byte `gorm:"type:blob" json:"transaction_hash"`
	UTxOID          []byte `gorm:"type:blob;column:utxo_id" json:"utxo_id"`
	UTxOIDIndex     uint32 `gorm:"index:simple_utxo_idx;column:utxo_index" json:"utxo_index"`
}

func (SimpleUTxO) TableName() string {
	return "simple_utxos"
}
