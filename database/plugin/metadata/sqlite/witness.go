package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetWitness stores or updates a Witness struct.
func (d *MetadataStoreSqlite) SetWitness(txn *gorm.DB, witness *models.Witness) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if witness == nil {
		return errors.New("witness cannot be nil")
	}
	if len(witness.TransactionHash) == 0 {
		return errors.New("witness transaction hash cannot be empty")
	}
	// Basic validation for nested Redeemers could be added here if needed,
	// but the SetRedeemer function should handle individual redeemer validation.

	result := db.Save(witness) // Save will create or update based on primary key
	return result.Error
}

// GetWitnessByTransactionHash retrieves a Witness struct by its transaction hash.
func (d *MetadataStoreSqlite) GetWitnessByTransactionHash(txn *gorm.DB, transactionHash []byte) (*models.Witness, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var witness models.Witness
	result := db.Where("transaction_hash = ?", transactionHash).First(&witness)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil Witness and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &witness, nil
}

// GetWitnessByID retrieves a Witness struct by its primary key ID.
func (d *MetadataStoreSqlite) GetWitnessByID(txn *gorm.DB, id uint) (*models.Witness, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var witness models.Witness
	result := db.First(&witness, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil Witness and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &witness, nil
}