package sqlite

import (
	"encoding/hex"
	"fmt"

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
	d.logger.Debug(fmt.Sprintf("[ASSET_DEBUG] SetAsset: Saving Asset with UTxOID: %x, UTxOIDIndex: %d, PolicyId: %x, Name: %x", asset.UTxOID, asset.UTxOIDIndex, asset.PolicyId, asset.Name))
	result := db.Save(asset) // Save will create or update based on primary key
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetTxsByPolicyId retrieves transactions associated with a given policy ID with pagination support.
func (d *MetadataStoreSqlite) GetTxsByPolicyId(txn *gorm.DB, policyId []byte, limit, offset int) ([]models.Transaction, error) {
	d.logger.Debug("GetTxsByPolicyId", "policyId_hex", hex.EncodeToString(policyId))
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Select("transactions.*").
		Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.transaction_hash").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.policy_id = ?", policyId)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Log the generated SQL query
	sqlQuery := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&transactions)
	})
	d.logger.Debug(fmt.Sprintf("GetTxsByPolicyId SQL: %s", sqlQuery))

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
		d.logger.Error("GetTxsByPolicyId: database query failed", "error", result.Error)
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
	query := db.Select("transactions.*").
		Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.transaction_hash").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.name = ?", tokenName)

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
func (d *MetadataStoreSqlite) GetTxsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}

	var transactions []models.Transaction
	d.logger.Debug(fmt.Sprintf("[ASSET_DEBUG] GetTxsByAssetFingerprint: assetFingerprint: %x, limit: %d, offset: %d", assetFingerprint, limit, offset))
	query := db.Select("transactions.*").
		Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.transaction_hash").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.fingerprint = ?", assetFingerprint)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Log the generated SQL query
	sqlQuery := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&transactions)
	})
	d.logger.Debug(fmt.Sprintf("[ASSET_DEBUG] GetTxsByAssetFingerprint SQL: %s", sqlQuery))

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
	d.logger.Debug("GetTxsByPolicyIdAndTokenName", "policyId_hex", hex.EncodeToString(policyId), "tokenName_hex", hex.EncodeToString(tokenName))
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Select("transactions.*").
		Joins("JOIN transaction_outputs ON transactions.transaction_hash = transaction_outputs.transaction_hash").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.policy_id = ? AND assets.name = ?", policyId, tokenName)

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
		d.logger.Error("GetTxsByPolicyIdAndTokenName: database query failed", "error", result.Error)
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

// GetTransactionInputsByAssetFingerprint retrieves TransactionInput models associated with a given asset fingerprint with pagination support.
func (d *MetadataStoreSqlite) GetTransactionInputsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var inputs []models.TransactionInput
	query := db.Table("transaction_inputs").
		Joins("JOIN assets ON transaction_inputs.utxo_id = assets.utxo_id AND transaction_inputs.utxo_index = assets.utxo_index").
		Where("assets.fingerprint = ?", assetFingerprint)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Preload("Asset").Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}

// GetTransactionOutputsByAssetFingerprint retrieves TransactionOutput models associated with a given asset fingerprint with pagination support.
func (d *MetadataStoreSqlite) GetTransactionOutputsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var outputs []models.TransactionOutput
	query := db.Table("transaction_outputs").
		Joins("JOIN assets ON transaction_outputs.utxo_id = assets.utxo_id AND transaction_outputs.utxo_index = assets.utxo_index").
		Where("assets.fingerprint = ?", assetFingerprint)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Preload("Asset").Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}

// CountUniqueAssets counts the number of unique assets based on their fingerprint.
func (d *MetadataStoreSqlite) CountUniqueAssets(txn *gorm.DB) (int64, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var count int64
	result := db.Model(&models.Asset{}).
		Select("COUNT(DISTINCT fingerprint)").
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
