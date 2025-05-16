package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models" // Corrected import path based on go.mod
	"gorm.io/gorm"
)

// GetDatum retrieves a Datum by UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetDatum(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.Datum, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var datum models.Datum
	result := db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).First(&datum)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil Datum and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &datum, nil
}

// GetDatumByHash retrieves a Datum by its hash.
func (d *MetadataStoreSqlite) GetDatumByHash(txn *gorm.DB, datumHash []byte) (*models.Datum, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var datum models.Datum
	result := db.Where("datum_hash = ?", datumHash).First(&datum)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil Datum and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &datum, nil
}


// SetDatum stores or updates a Datum.
func (d *MetadataStoreSqlite) SetDatum(txn *gorm.DB, datum *models.Datum) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if datum == nil {
		return errors.New("datum cannot be nil")
	}
	if len(datum.DatumHash) == 0 {
		return errors.New("datum hash cannot be empty")
	}
	if len(datum.DatumCbor) == 0 {
		return errors.New("datum cbor cannot be empty")
	}

	result := db.Save(datum) // Save will create or update based on primary key
	return result.Error
}
