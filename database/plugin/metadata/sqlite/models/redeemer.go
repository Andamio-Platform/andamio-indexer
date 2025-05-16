package models

type Redeemer struct {
	ID        uint   `gorm:"primaryKey"`
	WitnessID uint   `gorm:"index"`
	Index     uint   `gorm:"index" json:"index"`
	Tag       uint8  `gorm:"index" json:"tag"`
	Cbor      []byte `gorm:"type:blob" json:"cbor"`
}

func (Redeemer) TableName() string {
	return "redeemers"
}
