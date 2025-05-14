package sqlite

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetRedeemer stores or updates a Redeemer.
func (d *MetadataStoreSqlite) SetRedeemer(redeemer *models.Redeemer, txn *gorm.DB) error {
	result := txn.Save(redeemer) // Save will create or update based on primary key
	return result.Error
}

// GetRedeemersByWitnessId retrieves a slice of Redeemer structs associated with a given witness ID.
func (d *MetadataStoreSqlite) GetRedeemersByWitnessId(witnessID uint, txn *gorm.DB) ([]models.Redeemer, error) {
	var redeemers []models.Redeemer
	result := txn.Where("witness_id = ?", witnessID).Find(&redeemers)
	if result.Error != nil {
		return nil, result.Error
	}
	return redeemers, nil
}

// GetRedeemersByWitnessIdAndIndexAndTag retrieves a Redeemer struct matching a given witness ID, index, and tag.
func (d *MetadataStoreSqlite) GetRedeemersByWitnessIdAndIndexAndTag(witnessID uint, index uint, tag []byte, txn *gorm.DB) (*models.Redeemer, error) {
	var redeemer models.Redeemer
	result := txn.Where("witness_id = ? AND index = ? AND tag = ?", witnessID, index, tag).First(&redeemer)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil Redeemer and nil error if not found
		}
		return nil, result.Error // Return other errors
	}
	return &redeemer, nil
}

// GetRedeemersByWitnessIdAndTag retrieves a slice of Redeemer structs matching a given witness ID and tag.
func (d *MetadataStoreSqlite) GetRedeemersByWitnessIdAndTag(witnessID uint, tag []byte, txn *gorm.DB) ([]models.Redeemer, error) {
	var redeemers []models.Redeemer
	result := txn.Where("witness_id = ? AND tag = ?", witnessID, tag).Find(&redeemers)
	if result.Error != nil {
		return nil, result.Error
	}
	return redeemers, nil
}
