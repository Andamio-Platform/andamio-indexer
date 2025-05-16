package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTransactionInput stores or updates a TransactionInput struct.
func (d *MetadataStoreSqlite) SetTransactionInput(txn *gorm.DB, input *models.TransactionInput) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if input == nil {
		return errors.New("transaction input cannot be nil")
	}
	// Basic validation
	if len(input.TransactionHash) == 0 {
		return errors.New("transaction hash cannot be empty")
	}
	if len(input.UTxOID) == 0 {
		return errors.New("utxo id cannot be empty")
	}
	if len(input.Address) == 0 {
		return errors.New("address cannot be empty")
	}
	// UTxOIDIndex can be 0, so no validation needed
	// Amount can be 0, so no validation needed

	result := db.Save(input) // Save will create or update based on primary key
	return result.Error
}

// GetTransactionInput retrieves a TransactionInput struct by its UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetTxInputByUTxO(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var input models.TransactionInput
	result := db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&input)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionInput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &input, nil
}

// GetTransactionInputByID retrieves a TransactionInput struct by its primary key ID.
func (d *MetadataStoreSqlite) GetTxInputByID(txn *gorm.DB, id uint) (*models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var input models.TransactionInput
	result := db.First(&input, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionInput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &input, nil
}

// GetTransactionInputsByTransactionHash retrieves all TransactionInput structs for a given transaction hash.
func (d *MetadataStoreSqlite) GetTransactionInputsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var inputs []models.TransactionInput
	result := db.Where("transaction_hash = ?", transactionHash).Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}

// GetTransactionInputsByAddress retrieves all TransactionInput structs for a given address.
func (d *MetadataStoreSqlite) GetTransactionInputsByAddress(txn *gorm.DB, address []byte) ([]models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var inputs []models.TransactionInput
	result := db.Where("address = ?", address).Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}