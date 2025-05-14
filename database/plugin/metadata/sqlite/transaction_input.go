package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTransactionInput stores or updates a TransactionInput struct.
func (d *MetadataStoreSqlite) SetTransactionInput(input *models.TransactionInput, txn *gorm.DB) error {
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

	result := txn.Save(input) // Save will create or update based on primary key
	return result.Error
}

// GetTransactionInput retrieves a TransactionInput struct by its UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetTxInputByUTxO(utxoID []byte, utxoIndex uint32, txn *gorm.DB) (*models.TransactionInput, error) {
	var input models.TransactionInput
	result := txn.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&input)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionInput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &input, nil
}

// GetTransactionInputByID retrieves a TransactionInput struct by its primary key ID.
func (d *MetadataStoreSqlite) GetTransactionInputByID(id uint, txn *gorm.DB) (*models.TransactionInput, error) {
	var input models.TransactionInput
	result := txn.First(&input, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionInput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &input, nil
}

// GetTransactionInputsByTransactionHash retrieves all TransactionInput structs for a given transaction hash.
func (d *MetadataStoreSqlite) GetTransactionInputsByTransactionHash(transactionHash []byte, txn *gorm.DB) ([]models.TransactionInput, error) {
	var inputs []models.TransactionInput
	result := txn.Where("transaction_hash = ?", transactionHash).Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}

// GetTransactionInputsByAddress retrieves all TransactionInput structs for a given address.
func (d *MetadataStoreSqlite) GetTransactionInputsByAddress(address []byte, txn *gorm.DB) ([]models.TransactionInput, error) {
	var inputs []models.TransactionInput
	result := txn.Where("address = ?", address).Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}