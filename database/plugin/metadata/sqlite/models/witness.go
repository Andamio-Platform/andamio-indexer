package models

type Witness struct {
	ID              uint       `gorm:"primaryKey"`
	TransactionHash []byte     `gorm:"type:blob;index" json:"transaction_hash"`
	PlutusData      [][]byte   `gorm:"type:blob" json:"plutus_data"`
	PlutusV1Scripts [][]byte   `gorm:"type:blob" json:"plutus_v1_scripts"`
	PlutusV2Scripts [][]byte   `gorm:"type:blob" json:"plutus_v2_scripts"`
	PlutusV3Scripts [][]byte   `gorm:"type:blob" json:"plutus_v3_scripts"`
	Redeemers       []Redeemer `gorm:"foreignKey:WitnessID;references:ID" json:"redeemers"`
	// NativeScripts []NativeScript
}

func (Witness) TableName() string {
	return "witnesses"
}
