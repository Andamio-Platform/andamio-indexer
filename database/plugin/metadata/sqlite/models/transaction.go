package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/types"
)

// Assuming byteSliceJsonHex is intended to be stored as a byte slice
// and the JSON hex encoding is for API representation.
// We will store it as []byte in the database.

// WithdrawalsMap is a custom type to handle map[string]uint64 for GORM
type WithdrawalsMap map[string]uint64

// Value implements the driver.Valuer interface for WithdrawalsMap.
func (w WithdrawalsMap) Value() (driver.Value, error) {
	if w == nil {
		return nil, nil
	}
	return json.Marshal(w)
}

// Scan implements the sql.Scanner interface for WithdrawalsMap.
func (w *WithdrawalsMap) Scan(value interface{}) error {
	if value == nil {
		*w = make(WithdrawalsMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, w)
}

type Transaction struct {
	ID              uint                `gorm:"primaryKey"`
	BlockHash       []byte              `gorm:"index" json:"block_hash"`
	BlockNumber     uint64              `gorm:"index" json:"block_number"`
	SlotNumber      uint64              `gorm:"index" json:"slot_number"`
	TransactionHash []byte              `gorm:"index" json:"transaction_hash"`
	Inputs          []TransactionInput  `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"inputs"`
	Outputs         []TransactionOutput `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"outputs"`
	ReferenceInputs []SimpleUTxO        `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"reference_inputs"`
	Metadata        []byte              `gorm:"type:blob" json:"metadata"`
	Fee             uint64              `gorm:"index" json:"fee"`
	TTL             uint64              `gorm:"index" json:"ttl"`
	Withdrawals     WithdrawalsMap      `gorm:"type:blob" json:"withdrawals"`
	Witness         Witness             `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"witness"`
	Certificates    types.ByteSliceSlice `gorm:"type:blob" json:"certificate"`
}

// TableName overrides the table name
func (Transaction) TableName() string {
	return "transactions"
}

// GetID returns the ID of the Transaction.
func (t *Transaction) GetID() uint {
	return t.ID
}

// GetBlockHash returns the BlockHash of the Transaction.
func (t *Transaction) GetBlockHash() []byte {
	return t.BlockHash
}

// SetBlockHash sets the BlockHash of the Transaction.
func (t *Transaction) SetBlockHash(blockHash []byte) {
	t.BlockHash = blockHash
}

// GetBlockNumber returns the BlockNumber of the Transaction.
func (t *Transaction) GetBlockNumber() uint64 {
	return t.BlockNumber
}

// SetBlockNumber sets the BlockNumber of the Transaction.
func (t *Transaction) SetBlockNumber(blockNumber uint64) {
	t.BlockNumber = blockNumber
}

// GetSlotNumber returns the SlotNumber of the Transaction.
func (t *Transaction) GetSlotNumber() uint64 {
	return t.SlotNumber
}

// SetSlotNumber sets the SlotNumber of the Transaction.
func (t *Transaction) SetSlotNumber(slotNumber uint64) {
	t.SlotNumber = slotNumber
}

// GetTransactionHash returns the TransactionHash of the Transaction.
func (t *Transaction) GetTransactionHash() []byte {
	return t.TransactionHash
}

// GetInputs returns the Inputs of the Transaction.
func (t *Transaction) GetInputs() []TransactionInput {
	return t.Inputs
}

// SetInputs sets the Inputs of the Transaction.
func (t *Transaction) SetInputs(inputs []TransactionInput) {
	t.Inputs = inputs
}

// GetOutputs returns the Outputs of the Transaction.
func (t *Transaction) GetOutputs() []TransactionOutput {
	return t.Outputs
}

// SetOutputs sets the Outputs of the Transaction.
func (t *Transaction) SetOutputs(outputs []TransactionOutput) {
	t.Outputs = outputs
}

// GetReferenceInputs returns the ReferenceInputs of the Transaction.
func (t *Transaction) GetReferenceInputs() []SimpleUTxO {
	return t.ReferenceInputs
}

// SetReferenceInputs sets the ReferenceInputs of the Transaction.
func (t *Transaction) SetReferenceInputs(refInputs []SimpleUTxO) {
	t.ReferenceInputs = refInputs
}

// GetMetadata returns the Metadata of the Transaction.
func (t *Transaction) GetMetadata() []byte {
	return t.Metadata
}

// SetMetadata sets the Metadata of the Transaction.
func (t *Transaction) SetMetadata(metadata []byte) {
	t.Metadata = metadata
}

// GetFee returns the Fee of the Transaction.
func (t *Transaction) GetFee() uint64 {
	return t.Fee
}

// SetFee sets the Fee of the Transaction.
func (t *Transaction) SetFee(fee uint64) {
	t.Fee = fee
}

// GetTTL returns the TTL of the Transaction.
func (t *Transaction) GetTTL() uint64 {
	return t.TTL
}

// SetTTL sets the TTL of the Transaction.
func (t *Transaction) SetTTL(ttl uint64) {
	t.TTL = ttl
}

// GetWithdrawals returns the Withdrawals of the Transaction.
func (t *Transaction) GetWithdrawals() WithdrawalsMap {
	return t.Withdrawals
}

// SetWithdrawals sets the Withdrawals of the Transaction.
func (t *Transaction) SetWithdrawals(withdrawals WithdrawalsMap) {
	t.Withdrawals = withdrawals
}

// GetWitness returns the Witness of the Transaction.
func (t *Transaction) GetWitness() Witness {
	return t.Witness
}

// SetWitness sets the Witness of the Transaction.
func (t *Transaction) SetWitness(witness Witness) {
	t.Witness = witness
}

// GetCertificate returns the Certificate of the Transaction.
func (t *Transaction) GetCertificate() types.ByteSliceSlice {
	return t.Certificates
}

// SetCertificate sets the Certificate of the Transaction.
func (t *Transaction) SetCertificate(certificates types.ByteSliceSlice) {
	t.Certificates = certificates
}
