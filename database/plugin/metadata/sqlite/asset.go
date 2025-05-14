package sqlite

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
)

// GetAssets retrieves all Assets associated with a given UTxO by UTxOID and UTxOIDIndex.
func (d *MetadataStoreSqlite) GetAssets(utxoID []byte, utxoIndex uint32) ([]models.Asset, error) {
	var assets []models.Asset
	result := d.db.Where("utxo_id = ? AND utxo_index = ?", utxoID, utxoIndex).Find(&assets)
	if result.Error != nil {
		return nil, result.Error // Return error if query fails
	}
	return assets, nil
}

// SetAsset stores or updates a single Asset.
func (d *MetadataStoreSqlite) SetAsset(asset *models.Asset) error {
	result := d.db.Save(asset) // Save will create or update based on primary key
	return result.Error
}
