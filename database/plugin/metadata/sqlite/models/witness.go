package models

type Witness struct {
	ID              uint       `gorm:"primaryKey"`
	TransactionHash []byte     `gorm:"type:blob;index" json:"transaction_hash"` // Hash of the transaction this witness belongs to
	PlutusData      [][]byte   `gorm:"type:blob" json:"plutus_data"` // Removed index tag as it's not effective on blob types
	PlutusV1Scripts [][]byte   `gorm:"type:blob" json:"plutus_v1_scripts"` // Removed index tag as it's not effective on blob types
	PlutusV2Scripts [][]byte   `gorm:"type:blob" json:"plutus_v2_scripts"` // Removed index tag as it's not effective on blob types
	PlutusV3Scripts [][]byte   `gorm:"type:blob" json:"plutus_v3_scripts"` // Removed index tag as it's not effective on blob types
	Redeemers       []Redeemer `gorm:"foreignKey:WitnessID;references:ID" json:"redeemers"`
	// NativeScripts []NativeScript
}

func (Witness) TableName() string {
	return "witnesses"
}
