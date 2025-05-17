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

// GetTxsByPolicyId retrieves transactions associated with a given policy ID with pagination support.
func (d *MetadataStoreSqlite) GetTxsByPolicyId(txn *gorm.DB, policyId []byte, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.output_id").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.policy_id = ?", policyId).
		Distinct("transactions.transaction_hash")

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

// GetTxsByTokenName retrieves transactions associated with a given token name with pagination support.
func (d *MetadataStoreSqlite) GetTxsByTokenName(txn *gorm.DB, tokenName []byte, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.output_id").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.asset_name = ?", tokenName).
		Distinct("transactions.transaction_hash")

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

// GetTxsByAssetFingerprint retrieves transactions associated with a given asset fingerprint with pagination support.
func (d *MetadataStoreSqlite) GetTxsByAssetFingerprint(txn *gorm.DB, assetFingerprint string, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.output_id").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.asset_fingerprint = ?", assetFingerprint).
		Distinct("transactions.transaction_hash")

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

// GetTxsByPolicyIdAndTokenName retrieves transactions associated with a given policy ID and token name with pagination support.
func (d *MetadataStoreSqlite) GetTxsByPolicyIdAndTokenName(txn *gorm.DB, policyId []byte, tokenName []byte, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.output_id").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.policy_id = ? AND assets.asset_name = ?", policyId, tokenName).
		Distinct("transactions.transaction_hash")

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

// GetUTxOsByAssetFingerprint retrieves SimpleUTxO view models associated with a given asset fingerprint with pagination support.
func (d *MetadataStoreSqlite) GetUTxOsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.SimpleUTxO, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var utxos []models.SimpleUTxO
	query := db.Table("assets").
		Select("transaction_outputs.transaction_hash, transaction_outputs.utxo_id, transaction_outputs.utxo_index, transaction_outputs.address, transaction_outputs.amount, transaction_outputs.cbor").
		Joins("JOIN transaction_outputs ON assets.utxo_id = transaction_outputs.utxo_id AND assets.utxo_index = transaction_outputs.utxo_index").
		Where("assets.fingerprint = ?", assetFingerprint)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Find(&utxos)
	if result.Error != nil {
		return nil, result.Error
	}

	return utxos, nil
}
