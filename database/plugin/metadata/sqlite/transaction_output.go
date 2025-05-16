package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTransactionOutput stores or updates a TransactionOutput struct.
func (d *MetadataStoreSqlite) SetTransactionOutput(txn *gorm.DB, output *models.TransactionOutput) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if output == nil {
		return errors.New("transaction output cannot be nil")
	}
	// Basic validation
	if len(output.UTxOID) == 0 {
		return errors.New("utxo id cannot be empty")
	}
	if len(output.Address) == 0 {
		return errors.New("address cannot be empty")
	}
	// UTxOIDIndex can be 0, so no validation needed
	// Amount can be 0, so no validation needed

	result := db.Save(output) // Save will create or update based on primary key
	return result.Error
}

// GetTransactionOutput retrieves a TransactionOutput struct by its UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetTxOutputByUTxO(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var output models.TransactionOutput
	result := db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&output)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionOutput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &output, nil
}

// GetTransactionOutputByID retrieves a TransactionOutput struct by its primary key ID.
func (d *MetadataStoreSqlite) GetTxOutputByID(txn *gorm.DB, id uint) (*models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var output models.TransactionOutput
	result := db.First(&output, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionOutput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &output, nil
}

// GetTransactionOutputsByTransactionHash retrieves all TransactionOutput structs for a given transaction hash.
func (d *MetadataStoreSqlite) GetTransactionOutputsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var outputs []models.TransactionOutput
	result := db.Where("transaction_hash = ?", transactionHash).Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}

// GetTransactionOutputsByAddress retrieves all TransactionOutput structs for a given address.
func (d *MetadataStoreSqlite) GetTransactionOutputsByAddress(txn *gorm.DB, address []byte) ([]models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var outputs []models.TransactionOutput
	result := db.Where("address = ?", address).Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}
