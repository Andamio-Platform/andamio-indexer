package sqlite

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// GetAssets retrieves all Assets associated with a given UTxO by UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetAssets(txn *gorm.DB, utxoID []byte, utxoIndex uint32) ([]models.Asset, error) {
	var assets []models.Asset
	var result *gorm.DB
	if txn != nil {
		result = txn.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).Find(&assets)
	} else {
		result = d.db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).Find(&assets)
	}
	if result.Error != nil {
		return nil, result.Error // Return error if query fails
	}
	return assets, nil
}

// SetAsset stores or updates a single Asset.
func (d *MetadataStoreSqlite) SetAsset(txn *gorm.DB, asset *models.Asset) error {
	var db *gorm.DB
	if txn != nil {
		db = txn
	} else {
		db = d.DB()
	}
	result := db.Save(asset) // Save will create or update based on primary key
	if result.Error != nil {
		return result.Error
	}

	return nil
}
