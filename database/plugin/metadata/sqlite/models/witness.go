package models

import (
	"github.com/Andamio-Platform/andamio-indexer/database/types"
)

type Witness struct {
	ID              uint                 `gorm:"primaryKey"`
	TransactionHash []byte               `gorm:"type:blob;index" json:"transaction_hash"`
	PlutusData      [][]byte             `gorm:"type:blob" json:"plutus_data"`
	PlutusV1Scripts types.ByteSliceSlice `gorm:"type:blob" json:"plutus_v1_scripts"`
	PlutusV2Scripts types.ByteSliceSlice `gorm:"type:blob" json:"plutus_v2_scripts"`
	PlutusV3Scripts types.ByteSliceSlice `gorm:"type:blob" json:"plutus_v3_scripts"`
	Redeemers       []Redeemer           `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"redeemers"`
	// NativeScripts []NativeScript
}

func (Witness) TableName() string {
	return "witnesses"
}
