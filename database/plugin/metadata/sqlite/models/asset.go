package models

type Asset struct {
	ID          uint   `gorm:"primaryKey"`
	UTxOID      []byte `gorm:"type:blob;index:asset_utxo_idx;column:utxo_id" json:"utxo_id"`
	UTxOIDIndex uint32 `gorm:"index:asset_utxo_idx;column:utxo_index" json:"utxo_index"`
	Name        []byte `gorm:"type:blob" json:"name"`
	NameHex     []byte `gorm:"type:blob" json:"name_hex"`
	PolicyId    []byte `gorm:"type:blob" json:"policy_id"`
	Fingerprint []byte `gorm:"type:blob" json:"fingerprint"`
	Amount      uint64 `gorm:"index" json:"amount"`
}

func (Asset) TableName() string {
	return "assets"
}
