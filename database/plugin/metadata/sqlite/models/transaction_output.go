package models

type TransactionOutput struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UTxOID      []byte  `gorm:"type:blob" json:"utxo_id"` // Removed index tag as it's not effective on blob types, composite index is on UTxOIDIndex
	UTxOIDIndex uint32  `gorm:"index:tx_output_utxo_idx" json:"utxo_index"`
	Address     []byte  `gorm:"type:blob" json:"address"` // Removed index tag as it's not effective on blob types
	Amount      uint64  `gorm:"index" json:"amount"`
	Asset       []Asset `gorm:"foreignKey:UTxOID;references:UTxOID;foreignKey:UTxOIDIndex;references:UTxOIDIndex" json:"asset"` // Removed foreign key tag, should be on Asset struct
	Datum       Datum   `gorm:"foreignKey:UTxOID;references:UTxOID;foreignKey:UTxOIDIndex;references:UTxOIDIndex" json:"datum"` // Removed foreign key tag, should be on Datum struct
	Cbor        []byte  `gorm:"type:blob" json:"cbor"`                                                                          // Removed index tag as it's not effective on blob types
}

func (TransactionOutput) TableName() string {
	return "transaction_outputs"
}
