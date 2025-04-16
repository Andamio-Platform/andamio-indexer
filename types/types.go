package types

import "gorm.io/gorm"

type AndamioEvent struct {
	TxHash      string `json:"tx_hash" gorm:"index:idx_utxo"`
	TxID        int    `json:"tx_id" gorm:"index:idx_utxo"`
	Alias       string `json:"alias"`
	Datum       string `json:"datum"`
	ThisName    string `json:"this_name"`
	NextName    string `json:"next_name"`
	Blockhash   string `json:"blockhash"`
	BlockNumber uint64 `json:"block_number"`
	SlotNumber  uint64 `json:"slot_number"`
}

type DbStatus struct {
	gorm.Model
	SlotNumber uint64 `json:"slot_number"`
	BlockHash  string `json:"block_hash"`
}

type CardanoDatum struct {
	CBOR string `json:"cbor"`
	JSON struct {
		Constructor int `json:"constructor"`
		Fields      []struct {
			List  []interface{} `json:"list"`
			Bytes string        `json:"bytes"`
		} `json:"fields"`
	} `json:"json"`
}
