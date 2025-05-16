package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetSimpleUTxO stores or updates a SimpleUTxO struct.
func (d *MetadataStoreSqlite) SetSimpleUTxO(txn *gorm.DB, utxo *models.SimpleUTxO) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if utxo == nil {
		return errors.New("simple utxo cannot be nil")
	}
	// Basic validation
	if len(utxo.TransactionHash) == 0 {
		return errors.New("transaction hash cannot be empty")
	}
	if len(utxo.UTxOID) == 0 {
		return errors.New("utxo id cannot be empty")
	}
	// UTxOIDIndex can be 0, so no validation needed

	result := db.Save(utxo) // Save will create or update based on primary key
	return result.Error
}

// GetSimpleUTxOByUTxO retrieves a SimpleUTxO struct by its UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetSimpleUTxOByUTxO(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.SimpleUTxO, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var utxo models.SimpleUTxO
	result := db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&utxo)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil SimpleUTxO and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &utxo, nil
}

// GetSimpleUTxOByID retrieves SimpleUTxOs struct by its ID.
// NOTE: This function retrieves by UTxOID, not the primary key ID.
func (d *MetadataStoreSqlite) GetSimpleUTxOByID(txn *gorm.DB, utxoID []byte) ([]models.SimpleUTxO, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var utxos []models.SimpleUTxO
	result := db.Where("utxo_id = ?", utxoID).Find(&utxos)
	if result.Error != nil {
		return nil, result.Error
	}
	return utxos, nil
}

// GetSimpleUTxOByPrimaryKey retrieves a SimpleUTxO struct by its primary key ID.
func (d *MetadataStoreSqlite) GetSimpleUTxOByPrimaryKey(txn *gorm.DB, id uint) (*models.SimpleUTxO, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var utxo models.SimpleUTxO
	result := db.First(&utxo, id) // Find by primary key
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil SimpleUTxO and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &utxo, nil
}


// GetSimpleUTxOsByTransactionHash retrieves all SimpleUTxO structs for a given transaction hash.
func (d *MetadataStoreSqlite) GetSimpleUTxOsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.SimpleUTxO, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var utxos []models.SimpleUTxO
	result := db.Where("transaction_hash = ?", transactionHash).Find(&utxos)
	if result.Error != nil {
		return nil, result.Error
	}
	return utxos, nil
}
