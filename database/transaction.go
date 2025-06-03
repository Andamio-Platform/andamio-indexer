// Copyright 2025 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/Andamio-Platform/andamio-indexer/database/types"
	"github.com/dgraph-io/badger/v4"
)

type Transaction struct {
	ID              uint                       `gorm:"primaryKey"`
	BlockHash       []byte                     `gorm:"index" json:"block_hash"`
	BlockNumber     uint64                     `gorm:"index" json:"block_number"`
	SlotNumber      uint64                     `gorm:"index" json:"slot_number"`
	TransactionHash []byte                     `gorm:"index" json:"transaction_hash"`
	Inputs          []models.TransactionInput  `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"inputs"`
	Outputs         []models.TransactionOutput `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"outputs"`
	ReferenceInputs []models.SimpleUTxO        `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"reference_inputs"`
	Metadata        []byte                     `gorm:"type:blob" json:"metadata"`
	Fee             uint64                     `gorm:"index" json:"fee"`
	TTL             uint64                     `gorm:"index" json:"ttl"`
	Withdrawals     map[string]uint64          `gorm:"index" json:"withdrawals"`
	Witness         models.Witness             `gorm:"foreignKey:TransactionHash;references:TransactionHash" json:"witness"`
	Certificates    types.ByteSliceSlice       `gorm:"type:blob" json:"certificate"`
	TransactionCBOR []byte                     `gorm:"type:blob" json:"transaction_cbor"`
}

// TableName overrides the table name
func (Transaction) TableName() string {
	return "transactions"
}

// GetID returns the ID of the Transaction.
func (tx *Transaction) GetID() uint {
	return tx.ID
}

// GetBlockHash returns the BlockHash of the Transaction.
func (tx *Transaction) GetBlockHash() []byte {
	return tx.BlockHash
}

// SetBlockHash sets the BlockHash of the Transaction.
func (tx *Transaction) SetBlockHash(blockHash []byte) {
	tx.BlockHash = blockHash
}

// GetBlockNumber returns the BlockNumber of the Transaction.
func (tx *Transaction) GetBlockNumber() uint64 {
	return tx.BlockNumber
}

// SetBlockNumber sets the BlockNumber of the Transaction.
func (tx *Transaction) SetBlockNumber(blockNumber uint64) {
	tx.BlockNumber = blockNumber
}

// GetSlotNumber returns the SlotNumber of the Transaction.
func (tx *Transaction) GetSlotNumber() uint64 {
	return tx.SlotNumber
}

// SetSlotNumber sets the SlotNumber of the Transaction.
func (tx *Transaction) SetSlotNumber(slotNumber uint64) {
	tx.SlotNumber = slotNumber
}

// GetTransactionHash returns the TransactionHash of the Transaction.
func (tx *Transaction) GetTransactionHash() []byte {
	return tx.TransactionHash
}

// GetInputs returns the Inputs of the Transaction.
func (tx *Transaction) GetInputs() []models.TransactionInput {
	return tx.Inputs
}

// SetInputs sets the Inputs of the Transaction.
func (tx *Transaction) SetInputs(inputs []models.TransactionInput) {
	tx.Inputs = inputs
}

// GetOutputs returns the Outputs of the Transaction.
func (tx *Transaction) GetOutputs() []models.TransactionOutput {
	return tx.Outputs
}

// SetOutputs sets the Outputs of the Transaction.
func (tx *Transaction) SetOutputs(outputs []models.TransactionOutput) {
	tx.Outputs = outputs
}

// GetReferenceInputs returns the ReferenceInputs of the Transaction.
func (tx *Transaction) GetReferenceInputs() []models.SimpleUTxO {
	return tx.ReferenceInputs
}

// SetReferenceInputs sets the ReferenceInputs of the Transaction.
func (tx *Transaction) SetReferenceInputs(refInputs []models.SimpleUTxO) {
	tx.ReferenceInputs = refInputs
}

// GetMetadata returns the Metadata of the Transaction.
func (tx *Transaction) GetMetadata() []byte {
	return tx.Metadata
}

// SetMetadata sets the Metadata of the Transaction.
func (tx *Transaction) SetMetadata(metadata []byte) {
	tx.Metadata = metadata
}

// GetFee returns the Fee of the Transaction.
func (tx *Transaction) GetFee() uint64 {
	return tx.Fee
}

// SetFee sets the Fee of the Transaction.
func (tx *Transaction) SetFee(fee uint64) {
	tx.Fee = fee
}

// GetTTL returns the TTL of the Transaction.
func (tx *Transaction) GetTTL() uint64 {
	return tx.TTL
}

// SetTTL sets the TTL of the Transaction.
func (tx *Transaction) SetTTL(ttl uint64) {
	tx.TTL = ttl
}

// GetWithdrawals returns the Withdrawals of the Transaction.
func (tx *Transaction) GetWithdrawals() map[string]uint64 {
	return tx.Withdrawals
}

// SetWithdrawals sets the Withdrawals of the Transaction.
func (tx *Transaction) SetWithdrawals(withdrawals map[string]uint64) {
	tx.Withdrawals = withdrawals
}

// GetWitness returns the Witness of the Transaction.
func (tx *Transaction) GetWitness() models.Witness {
	return tx.Witness
}

// SetWitness sets the Witness of the Transaction.
func (tx *Transaction) SetWitness(witness models.Witness) {
	tx.Witness = witness
}

// GetCertificate returns the Certificate of the Transaction.
func (tx *Transaction) GetCertificate() types.ByteSliceSlice {
	return tx.Certificates
}

// SetCertificate sets the Certificate of the Transaction.
func (tx *Transaction) SetCertificate(certificates types.ByteSliceSlice) {
	tx.Certificates = certificates
}

// GetTransactionCBOR returns the TransactionCBOR of the Transaction.
func (tx *Transaction) GetTransactionCBOR() []byte {
	return tx.TransactionCBOR
}

// SetTransactionCBOR sets the TransactionCBOR of the Transaction.
func (tx *Transaction) SetTransactionCBOR(cbor []byte) {
	tx.TransactionCBOR = cbor
}

// TxBlobKey generates the blob store key for a transaction's CBOR
func TxBlobKey(txHash []byte) []byte {
	key := []byte("t")
	key = append(key, txHash...)
	return key
}

func (tx *Transaction) loadCbor(txn *Txn) error {
	key := TxBlobKey(tx.TransactionHash)
	item, err := txn.Blob().Get(key)
	if err != nil {
		return err
	}
	tx.TransactionCBOR, err = item.ValueCopy(nil)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// NewTx stores a transaction's metadata in the metadata store and its CBOR in the blob store
func (d *Database) NewTx(blockHash []byte, blockNumber uint64, slotNumber uint64, transactionHash []byte, inputs []models.TransactionInput, outputs []models.TransactionOutput, referenceInputs []models.SimpleUTxO, metadata []byte, fee uint64, ttl uint64, withdrawals map[string]uint64, witness models.Witness, certificates [][]byte, transactionCBOR []byte, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}

	// Store CBOR in blob DB
	key := TxBlobKey(transactionHash)
	err := txn.Blob().Set(key, transactionCBOR)
	if err != nil {
		return err
	}

	tempTx := models.Transaction{
		BlockHash:       blockHash,
		BlockNumber:     blockNumber,
		SlotNumber:      slotNumber,
		TransactionHash: transactionHash,
		Inputs:          inputs,
		Outputs:         outputs,
		ReferenceInputs: referenceInputs,
		Metadata:        metadata,
		Fee:             fee,
		TTL:             ttl,
		Withdrawals:     withdrawals,
		Witness:         witness,
		Certificates:    types.ByteSliceSlice(certificates),
	}

	// Store metadata in metadata DB
	return d.metadata.SetTx(txn.Metadata(), &tempTx)
}

// GetTxByTxHash retrieves a transaction's metadata and CBOR by its hash
func (d *Database) GetTxByTxHash(txHash []byte, txn *Txn) (*Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}

	// Get metadata from metadata DB
	modelTx, err := d.metadata.GetTxByTxHash(txn.Metadata(), txHash)
	if err != nil {
		return nil, err
	}
	if modelTx == nil {
		return nil, errors.New("transaction not found") // Or return a specific not found error
	}

	dbTx := &Transaction{
		ID:              modelTx.ID,
		BlockHash:       modelTx.BlockHash,
		BlockNumber:     modelTx.BlockNumber,
		SlotNumber:      modelTx.SlotNumber,
		TransactionHash: modelTx.TransactionHash,
		Inputs:          modelTx.Inputs,
		Outputs:         modelTx.Outputs,
		ReferenceInputs: modelTx.ReferenceInputs,
		Metadata:        modelTx.Metadata,
		Fee:             modelTx.Fee,
		TTL:             modelTx.TTL,
		Withdrawals:     modelTx.Withdrawals,
		Witness:         modelTx.Witness,
		Certificates:    modelTx.Certificates,
	}

	// Get CBOR from blob DB
	cborKey := TxBlobKey(txHash)
	cborItem, err := txn.Blob().Get(cborKey)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			d.Logger().Debug("Transaction CBOR not found in blob store", "txHash", hex.EncodeToString(txHash))
		} else {
			d.Logger().Warn("Failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(txHash), "error", err)
		}
	} else {
		cborBytes, err := cborItem.ValueCopy(nil)
		if err != nil {
			d.Logger().Warn("Failed to copy transaction CBOR value", "txHash", hex.EncodeToString(txHash), "error", err)
		} else {
			dbTx.TransactionCBOR = cborBytes
		}
	}

	return dbTx, nil
}

// GetTxsByBlockNumber retrieves transaction metadata by block number with pagination support
func (d *Database) GetTxsByBlockNumber(blockNumber uint64, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByBlockNumber(txn.Metadata(), blockNumber, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert models.Transaction to database.Transaction and load CBOR
	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
			// TransactionCBOR will be loaded below
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			// Log error but continue if CBOR not found
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsBySlotRange retrieves transaction metadata and CBOR by slot range with pagination support
func (d *Database) GetTxsBySlotRange(startSlot, endSlot uint64, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}

	modelsTxs, err := d.metadata.GetTxsBySlotRange(txn.Metadata(), startSlot, endSlot, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				d.Logger().Debug("Transaction CBOR not found in blob store", "txHash", hex.EncodeToString(modelTx.TransactionHash))
			} else {
				d.Logger().Warn("Failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			}
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("Failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByInputAddress retrieves transaction inputs by address with pagination support
func (d *Database) GetTxsByInputAddress(address string, limit, offset int, txn *Txn) ([]models.TransactionInput, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	return d.metadata.GetTxInputsByAddress(txn.Metadata(), address, limit, offset)
}

// GetTxsByOutputAddress retrieves transaction metadata by output address with pagination support
func (d *Database) GetTxsByOutputAddress(address string, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByOutputAddress(txn.Metadata(), address, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByAnyAddress retrieves transaction metadata by any address (input or output) with pagination support
func (d *Database) GetTxsByAnyAddress(address string, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByAnyAddress(txn.Metadata(), address, limit, offset)
	if err != nil {
		return nil, err
	}
	d.Logger().Debug("GetTxsByAnyAddress: retrieved transactions from metadata", "count", len(modelsTxs))
	if len(modelsTxs) == 0 {
		d.Logger().Debug("GetTxsByAnyAddress: no transactions found for address", "address", address)
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		d.Logger().Debug("GetTxsByAnyAddress: processing transaction for CBOR retrieval", "index", i, "txHash", hex.EncodeToString(modelTx.TransactionHash))
		// Log inputs and their assets
		for j, input := range modelTx.Inputs {
			d.Logger().Debug(fmt.Sprintf("GetTxsByAnyAddress:   Input[%d]: UTxOID=%x, UTxOIDIndex=%d", j, input.UTxOID, input.UTxOIDIndex))
			for k, asset := range input.Asset {
				d.Logger().Debug(fmt.Sprintf("GetTxsByAnyAddress:     Input Asset[%d]: UTxOID=%x, UTxOIDIndex=%d, PolicyId=%x, Name=%x", k, asset.UTxOID, asset.UTxOIDIndex, asset.PolicyId, asset.Name))
			}
		}
		// Log outputs and their assets
		for j, output := range modelTx.Outputs {
			d.Logger().Debug(fmt.Sprintf("GetTxsByAnyAddress:   Output[%d]: UTxOID=%x, UTxOIDIndex=%d", j, output.UTxOID, output.UTxOIDIndex))
			for k, asset := range output.Asset {
				d.Logger().Debug(fmt.Sprintf("GetTxsByAnyAddress:     Output Asset[%d]: UTxOID=%x, UTxOIDIndex=%d, PolicyId=%x, Name=%x", k, asset.UTxOID, asset.UTxOIDIndex, asset.PolicyId, asset.Name))
			}
		}

		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}
		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByAssetFingerprint retrieves transaction metadata by asset fingerprint with pagination support
func (d *Database) GetTxsByAssetFingerprint(assetFingerprint []byte, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByAssetFingerprint(txn.Metadata(), assetFingerprint, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByPolicyIdAndTokenName retrieves transaction metadata by policy ID and token name with pagination support
func (d *Database) GetTxsByPolicyIdAndTokenName(policyID []byte, tokenName []byte, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByPolicyIdAndTokenName(txn.Metadata(), policyID, tokenName, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByPolicyId retrieves transaction metadata by policy ID with pagination support
func (d *Database) GetTxsByPolicyId(policyID []byte, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByPolicyId(txn.Metadata(), policyID, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// GetTxsByTokenName retrieves transaction metadata by token name with pagination support
func (d *Database) GetTxsByTokenName(tokenName []byte, limit, offset int, txn *Txn) ([]Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	modelsTxs, err := d.metadata.GetTxsByTokenName(txn.Metadata(), tokenName, limit, offset)
	if err != nil {
		return nil, err
	}

	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			ID:              modelTx.ID,
			BlockHash:       modelTx.BlockHash,
			BlockNumber:     modelTx.BlockNumber,
			SlotNumber:      modelTx.SlotNumber,
			TransactionHash: modelTx.TransactionHash,
			Inputs:          modelTx.Inputs,
			Outputs:         modelTx.Outputs,
			ReferenceInputs: modelTx.ReferenceInputs,
			Metadata:        modelTx.Metadata,
			Fee:             modelTx.Fee,
			TTL:             modelTx.TTL,
			Withdrawals:     modelTx.Withdrawals,
			Witness:         modelTx.Witness,
			Certificates:    modelTx.Certificates,
		}

		// Load CBOR for each transaction
		cborKey := TxBlobKey(modelTx.TransactionHash)
		cborItem, err := txn.Blob().Get(cborKey)
		if err != nil {
			d.Logger().Warn("failed to load transaction CBOR for transaction", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
		} else {
			cborBytes, err := cborItem.ValueCopy(nil)
			if err != nil {
				d.Logger().Warn("failed to copy transaction CBOR value", "txHash", hex.EncodeToString(modelTx.TransactionHash), "error", err)
			} else {
				dbTxs[i].TransactionCBOR = cborBytes
			}
		}
	}
	return dbTxs, nil
}

// DeleteTxByHash deletes a single transaction by its hash
func (d *Database) DeleteTxByHash(txHash []byte, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}
	// Delete CBOR from blob DB
	key := TxBlobKey(txHash)
	if err := txn.Blob().Delete(key); err != nil {
		return err
	}
	// Delete metadata from metadata DB
	return d.metadata.DeleteTxByHash(txn.Metadata(), txHash)
}

// DeleteTxsByBlockNumber deletes transaction metadata and CBOR by block number
func (d *Database) DeleteTxsByBlockNumber(blockNumber uint64, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}

	// Get transaction hashes for the block from metadata DB
	modelsTxs, err := d.metadata.GetTxsByBlockNumber(txn.Metadata(), blockNumber, 0, -1) // Get all txs for the block
	if err != nil {
		return err
	}

	// Convert models.Transaction to database.Transaction for CBOR deletion
	dbTxs := make([]Transaction, len(modelsTxs))
	for i, modelTx := range modelsTxs {
		dbTxs[i] = Transaction{
			TransactionHash: modelTx.TransactionHash, // Only hash is needed for deletion
		}
	}

	// Delete CBOR for each transaction in blob DB
	for _, tx := range dbTxs {
		key := TxBlobKey(tx.TransactionHash)
		err := txn.Blob().Delete(key)
		if err != nil {
			// Log error but continue
			d.Logger().Error("failed to delete transaction CBOR from blob store during block deletion", "txHash", tx.TransactionHash, "blockNumber", blockNumber, "error", err)
		}
	}

	// Delete metadata from metadata DB
	return d.metadata.DeleteTxsByBlockNumber(txn.Metadata(), blockNumber)
}

// GetUniqueAddressesCount retrieves the total count of unique addresses from all transactions.
