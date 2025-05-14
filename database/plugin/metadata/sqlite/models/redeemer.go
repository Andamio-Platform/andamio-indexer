package models

type Redeemer struct {
	ID          uint   `gorm:"primaryKey"`
	WitnessID   uint   `gorm:"index"`
	Index       uint   `gorm:"index" json:"index"`
	Tag         []byte `gorm:"type:blob" json:"tag"`
	ScriptHash  []byte `gorm:"type:blob" json:"script_hash"`
	Cbor        []byte `gorm:"type:blob" json:"cbor"`
}

func (Redeemer) TableName() string {
	return "redeemers"
}
