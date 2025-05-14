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
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/dgraph-io/badger/v4"
)

type Transaction struct {
	ID              uint                       `gorm:"primaryKey"`
	BlockHash       string                     `gorm:"index" json:"block_hash"`
	BlockNumber     uint64                     `gorm:"index" json:"block_number"`
	SlotNumber      uint64                     `gorm:"index" json:"slot_number"`
	TransactionHash []byte                     `gorm:"index" json:"transaction_hash"`
	Inputs          []models.TransactionInput  `gorm:"-" json:"inputs"`
	Outputs         []models.TransactionOutput `gorm:"-" json:"outputs"`
	ReferenceInputs []models.SimpleUTxO        `gorm:"-" json:"reference_inputs"`
	Metadata        []byte                     `gorm:"type:blob" json:"metadata"`
	Fee             uint64                     `gorm:"index" json:"fee"`
	TTL             uint64                     `gorm:"index" json:"ttl"`
	Withdrawals     map[byte]uint64            `gorm:"index" json:"withdrawals"`
	Witness         models.Witness             `gorm:"-" json:"witness"`
	Certificate     []byte                     `gorm:"type:blob" json:"certificate"`
	TransactionCBOR []byte                     `gorm:"-" json:"transaction_cbor"`
}

// TableName overrides the table name
func (Transaction) TableName() string {
	return "transactions"
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
func (d *Database) NewTx(tx *models.Transaction, cbor []byte, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}

	// Store CBOR in blob DB
	key := TxBlobKey(tx.TransactionHash)
	err := txn.Blob().Set(key, cbor)
	if err != nil {
		return err
	}

	// Store metadata in metadata DB
	return d.metadata.SetTx(tx, txn.Metadata())
}

// GetTxByHash retrieves a transaction's metadata and CBOR by its hash
func (d *Database) GetTxByHash(txHash []byte, txn *Txn) (*models.Transaction, []byte, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}

	// Get metadata from metadata DB
	tx, err := d.metadata.GetTxByHash(txHash, txn.Metadata())
	if err != nil {
		return nil, nil, err
	}
	if tx == nil {
		return nil, nil, errors.New("transaction not found") // Or return a specific not found error
	}

	// Get CBOR from blob DB
	key := TxBlobKey(txHash)
	item, err := txn.Blob().Get(key)
	if err != nil {
		if errors.Is(err, errors.New("Key not found")) { // Check for specific blob store not found error
			return tx, nil, nil // Return metadata even if CBOR is not found
		}
		return tx, nil, err
	}
	cbor, err := item.ValueCopy(nil)
	if err != nil {
		return tx, nil, err
	}

	return tx, cbor, nil
}

// GetTxsByBlockNumber retrieves transaction metadata by block number with pagination support
func (d *Database) GetTxsByBlockNumber(blockNumber uint64, limit, offset int, txn *Txn) ([]models.Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	return d.metadata.GetTxsByBlockNumber(blockNumber, limit, offset, txn.Metadata())
}

// GetTxsByInputAddress retrieves transaction metadata by input address with pagination support
func (d *Database) GetTxsByInputAddress(address string, limit, offset int, txn *Txn) ([]models.Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	return d.metadata.GetTxsByInputAddress(address, limit, offset, txn.Metadata())
}

// GetTxsByOutputAddress retrieves transaction metadata by output address with pagination support
func (d *Database) GetTxsByOutputAddress(address string, limit, offset int, txn *Txn) ([]models.Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	return d.metadata.GetTxsByOutputAddress(address, limit, offset, txn.Metadata())
}

// GetTxsByAnyAddress retrieves transaction metadata by any address (input or output) with pagination support
func (d *Database) GetTxsByAnyAddress(address string, limit, offset int, txn *Txn) ([]models.Transaction, error) {
	if txn == nil {
		txn = d.Transaction(false)
		defer txn.Commit() //nolint:errcheck
	}
	return d.metadata.GetTxsByAnyAddress(address, limit, offset, txn.Metadata())
}

// DeleteTxByHash deletes a transaction's metadata and CBOR by its hash
func (d *Database) DeleteTxByHash(txHash []byte, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}

	// Delete CBOR from blob DB
	key := TxBlobKey(txHash)
	err := txn.Blob().Delete(key)
	if err != nil {
		// Log error but continue to delete metadata
		d.Logger().Error("failed to delete transaction CBOR from blob store", "txHash", txHash, "error", err)
	}

	// Delete metadata from metadata DB
	return d.metadata.DeleteTxByHash(txHash, txn.Metadata())
}

// DeleteTxsByBlockNumber deletes transaction metadata and CBOR by block number
func (d *Database) DeleteTxsByBlockNumber(blockNumber uint64, txn *Txn) error {
	if txn == nil {
		txn = d.Transaction(true)
		defer txn.Commit() //nolint:errcheck
	}

	// Get transaction hashes for the block from metadata DB
	txs, err := d.metadata.GetTxsByBlockNumber(blockNumber, 0, -1, txn.Metadata()) // Get all txs for the block
	if err != nil {
		return err
	}

	// Delete CBOR for each transaction in blob DB
	for _, tx := range txs {
		key := TxBlobKey(tx.TransactionHash)
		err := txn.Blob().Delete(key)
		if err != nil {
			// Log error but continue
			d.Logger().Error("failed to delete transaction CBOR from blob store during block deletion", "txHash", tx.TransactionHash, "blockNumber", blockNumber, "error", err)
		}
	}

	// Delete metadata from metadata DB
	return d.metadata.DeleteTxsByBlockNumber(blockNumber, txn.Metadata())
}
