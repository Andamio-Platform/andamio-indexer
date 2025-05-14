package models

type Asset struct {
	ID                  uint   `gorm:"primaryKey"`
	UTxOID              []byte `gorm:"type:blob;index:utxo_idx" json:"utxo_id"` // Removed standalone index tag as it's not effective on blob types, composite index is on UTxOIDIndex
	UTxOIDIndex         uint32 `gorm:"index:utxo_idx" json:"utxo_index"` // Foreign key to TransactionInput and TransactionOutput and the reference would be UTxOIDIndex of either TransactionInput or TransactionOutput
	Name                []byte `gorm:"type:blob" json:"name"` // Removed index tag as it's not effective on blob types
	NameHex             []byte `gorm:"type:blob" json:"name_hex"` // Removed index tag as it's not effective on blob types
	PolicyId            []byte `gorm:"type:blob" json:"policy_id"` // Removed index tag as it's not effective on blob types
	Fingerprint         []byte `gorm:"type:blob" json:"fingerprint"` // Removed index tag as it's not effective on blob types
	Amount              uint64 `gorm:"index" json:"amount"`
}

func (Asset) TableName() string {
	return "assets"
}
