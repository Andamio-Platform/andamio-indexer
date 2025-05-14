package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTransactionOutput stores or updates a TransactionOutput struct.
func (d *MetadataStoreSqlite) SetTransactionOutput(output *models.TransactionOutput, txn *gorm.DB) error {
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

	result := txn.Save(output) // Save will create or update based on primary key
	return result.Error
}

// GetTransactionOutput retrieves a TransactionOutput struct by its UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetTransactionOutput(utxoID []byte, utxoIndex uint32, txn *gorm.DB) (*models.TransactionOutput, error) {
	var output models.TransactionOutput
	result := txn.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&output)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionOutput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &output, nil
}

// GetTransactionOutputByID retrieves a TransactionOutput struct by its primary key ID.
func (d *MetadataStoreSqlite) GetTransactionOutputByID(id uint, txn *gorm.DB) (*models.TransactionOutput, error) {
	var output models.TransactionOutput
	result := txn.First(&output, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil TransactionOutput and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &output, nil
}

// GetTransactionOutputsByTransactionHash retrieves all TransactionOutput structs for a given transaction hash.
func (d *MetadataStoreSqlite) GetTransactionOutputsByTransactionHash(transactionHash []byte, txn *gorm.DB) ([]models.TransactionOutput, error) {
	var outputs []models.TransactionOutput
	result := txn.Where("transaction_hash = ?", transactionHash).Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}

// GetTransactionOutputsByAddress retrieves all TransactionOutput structs for a given address.
func (d *MetadataStoreSqlite) GetTransactionOutputsByAddress(address []byte, txn *gorm.DB) ([]models.TransactionOutput, error) {
	var outputs []models.TransactionOutput
	result := txn.Where("address = ?", address).Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}
