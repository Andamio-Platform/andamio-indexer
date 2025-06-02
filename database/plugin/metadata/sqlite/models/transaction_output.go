package models

type TransactionOutput struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	TransactionHash []byte  `gorm:"index;type:blob" json:"transaction_hash"`
	UTxOID          []byte  `gorm:"type:blob;column:utxo_id" json:"utxo_id"`
	UTxOIDIndex     uint32  `gorm:"index:tx_output_utxo_idx;column:utxo_index" json:"utxo_index"`
	Address         []byte  `gorm:"type:blob" json:"address"`
	Amount          uint64  `gorm:"index" json:"amount"`
	Asset           []Asset `gorm:"foreignKey:UTxOID,UTxOIDIndex;references:UTxOID,UTxOIDIndex" json:"asset"`
	Datum           Datum   `gorm:"foreignKey:UTxOID,UTxOIDIndex;references:UTxOID,UTxOIDIndex" json:"datum"`
	Cbor            []byte  `gorm:"type:blob" json:"cbor"`
}

func (TransactionOutput) TableName() string {
	return "transaction_outputs"
}
