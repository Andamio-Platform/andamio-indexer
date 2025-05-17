package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTx inserts or updates a transaction record
// SetTx inserts or updates a transaction record and its nested data
func (d *MetadataStoreSqlite) SetTx(txn *gorm.DB, tx *models.Transaction) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if tx == nil {
		return errors.New("transaction cannot be nil")
	}
	if len(tx.TransactionHash) == 0 {
		return errors.New("transaction hash cannot be empty")
	}
	// Basic validation for nested data could be added here if needed,
	// but the individual setter functions should handle validation for their respective types.

	// Save main transaction metadata
	result := db.Save(tx)
	if result.Error != nil {
		return result.Error
	}

	// Save nested data
	if err := d.setInputs(txn, tx.Inputs, tx.TransactionHash); err != nil {
		return err
	}
	if err := d.setOutputs(txn, tx.Outputs); err != nil {
		return err
	}
	if err := d.setReferenceInputs(txn, tx.ReferenceInputs, tx.TransactionHash); err != nil {
		return err
	}
	// Check if Witness is present and should be saved (e.g., by checking a required field)
	if len(tx.Witness.TransactionHash) > 0 { // Assuming TransactionHash is a required field for Witness
		if err := d.setWitness(txn, tx.Witness, tx.TransactionHash); err != nil {
			return err
		}
	}

	return nil
}

// setInputs saves a slice of transaction inputs to the database
func (d *MetadataStoreSqlite) setInputs(txn *gorm.DB, inputs []models.TransactionInput, txHash []byte) error {
	db := txn
	if db == nil {
		db = d.db
	}
	for _, input := range inputs {
		input.TransactionHash = txHash // Set foreign key
		result := db.Save(&input)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// setOutputs saves a slice of transaction outputs to the database
func (d *MetadataStoreSqlite) setOutputs(txn *gorm.DB, outputs []models.TransactionOutput) error {
	db := txn
	if db == nil {
		db = d.db
	}
	for _, output := range outputs {
		// Save the TransactionOutput record
		result := db.Save(&output)
		if result.Error != nil {
			return result.Error
		}
		// Save nested assets within the output
		for _, asset := range output.Asset {
			// The relationship between Asset and TransactionOutput is via UTxOID and UTxOIDIndex
			// Ensure these fields are set correctly on the asset before saving
			asset.UTxOID = output.UTxOID
			asset.UTxOIDIndex = output.UTxOIDIndex
			if err := d.SetAsset(txn, &asset); err != nil { // Assuming SetAsset takes a single asset and handles its saving
				return err
			}
		}
	}
	return nil
}

// setReferenceInputs saves a slice of reference inputs to the database
func (d *MetadataStoreSqlite) setReferenceInputs(txn *gorm.DB, refInputs []models.SimpleUTxO, txHash []byte) error {
	db := txn
	if db == nil {
		db = d.db
	}
	for _, refInput := range refInputs {
		refInput.TransactionHash = txHash // Set foreign key
		result := db.Save(&refInput)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// setWitness saves a transaction witness and its nested data to the database
func (d *MetadataStoreSqlite) setWitness(txn *gorm.DB, witness models.Witness, txHash []byte) error {
	db := txn
	if db == nil {
		db = d.db
	}
	witness.TransactionHash = txHash // Set foreign key to Transaction

	// Save the Witness record
	result := db.Save(&witness)
	if result.Error != nil {
		return result.Error
	}

	// Save nested Redeemers
	for _, redeemer := range witness.Redeemers {
		redeemer.TransactionHash = witness.TransactionHash // Set foreign key to Witness's transaction hash
		result := db.Save(&redeemer)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// GetTxByTxHash retrieves a single transaction by its hash
func (d *MetadataStoreSqlite) GetTxByTxHash(txn *gorm.DB, txHash []byte) (*models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transaction models.Transaction
	result := db.Where("transaction_hash = ?", txHash).
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		First(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil transaction and nil error if not found
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// GetTxByID retrieves a single transaction by its primary key ID.
func (d *MetadataStoreSqlite) GetTxByID(txn *gorm.DB, id uint) (*models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transaction models.Transaction
	result := db.
		Preload("Inputs").
		Preload("Inputs.Asset").
		Preload("Inputs.Datum").
		Preload("Outputs").
		Preload("Outputs.Asset").
		Preload("Outputs.Datum").
		Preload("ReferenceInputs").
		Preload("Witness").
		Preload("Witness.Redeemers").
		First(&transaction, id) // Find by primary key
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil transaction and nil error if not found
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// GetTxsByBlockNumber retrieves all transactions for a given block number with pagination support
func (d *MetadataStoreSqlite) GetTxsByBlockNumber(txn *gorm.DB, blockNumber uint64, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	query := db.Where("block_number = ?", blockNumber)

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

// GetTxsByInputAddress retrieves transactions where the given address appears in the inputs with pagination support.
// This function joins the transactions and transaction_inputs tables on the transaction hash to find transactions
// that have an input from the given address.
func (d *MetadataStoreSqlite) GetTxsByInputAddress(txn *gorm.DB, address string, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction

	query := db.Joins("JOIN transaction_inputs ON transactions.transaction_hash = transaction_inputs.transaction_hash").
		Where("transaction_inputs.address = ?", []byte(address)). // Filter by the address in the input
		Distinct("transactions.transaction_hash")                 // Ensure unique transactions

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

// GetTxsByOutputAddress retrieves transactions where the given address appears in the outputs with pagination support.
// This function finds transactions that contain an output sent to a specified address.
// OutputID in transaction_outputs is the hash of the transaction that created the output.
func (d *MetadataStoreSqlite) GetTxsByOutputAddress(txn *gorm.DB, address string, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	var outputTxHashes [][]byte

	// Find the OutputID (transaction hash) of transactions that have an output to the given address
	result := db.Model(&models.TransactionOutput{}).
		Select("DISTINCT output_id").
		Where("address = ?", []byte(address)). // Filter by the address in the output
		Find(&outputTxHashes)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(outputTxHashes) == 0 {
		return transactions, nil // No transactions found with outputs to this address
	}

	// Retrieve the Transaction records using these transaction hashes with pagination
	query := db.Where("transaction_hash IN (?)", outputTxHashes)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result = query.
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

// GetTxsByAnyAddress retrieves transactions where the given address appears in either inputs or outputs with pagination support.
// This function queries both the transaction_inputs and transaction_outputs tables
// to find transactions associated with the given address.
func (d *MetadataStoreSqlite) GetTxsByAnyAddress(txn *gorm.DB, address string, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	uniqueTxHashes := make(map[string]bool)

	// Get transaction hashes from inputs
	var inputTxHashes [][]byte
	result := db.Model(&models.TransactionInput{}).
		Select("DISTINCT transaction_inputs.transaction_hash").
		Where("transaction_inputs.address = ?", []byte(address)).
		Find(&inputTxHashes)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, hash := range inputTxHashes {
		uniqueTxHashes[string(hash)] = true
	}

	// Get transaction hashes from outputs
	var outputTxHashes [][]byte
	result = db.Model(&models.TransactionOutput{}).
		Select("DISTINCT output_id").
		Where("address = ?", []byte(address)).
		Find(&outputTxHashes)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, hash := range outputTxHashes {
		uniqueTxHashes[string(hash)] = true
	}

	// Collect unique hashes into a slice
	var finalTxHashes [][]byte
	for hashStr := range uniqueTxHashes {
		finalTxHashes = append(finalTxHashes, []byte(hashStr))
	}

	if len(finalTxHashes) == 0 {
		return transactions, nil // No transactions found for this address
	}

	// Retrieve the Transaction records with pagination
	query := db.Where("transaction_hash IN (?)", finalTxHashes)

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result = query.
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

// SetTxs inserts or updates a batch of transaction records.
func (d *MetadataStoreSqlite) SetTxs(txn *gorm.DB, txs []*models.Transaction) error {
	db := txn
	if db == nil {
		db = d.db
	}
	if len(txs) == 0 {
		return nil // Nothing to save
	}
	for _, tx := range txs {
		if tx == nil {
			return errors.New("transaction in batch cannot be nil")
		}
		if len(tx.TransactionHash) == 0 {
			return errors.New("transaction hash in batch cannot be empty")
		}
	}

	// Using Save in a loop handles both insertion and updating based on primary key/unique index.
	for _, tx := range txs {
		result := db.Save(tx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// GetTxs retrieves a list of transactions with pagination
func (d *MetadataStoreSqlite) GetTxs(txn *gorm.DB, limit, offset int) ([]models.Transaction, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var transactions []models.Transaction
	result := db.Limit(limit).Offset(offset).
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

// CountTxs gets the total count of transaction records
func (d *MetadataStoreSqlite) CountTxs(txn *gorm.DB) (int64, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var count int64
	result := db.Model(&models.Transaction{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// DeleteTxByHash deletes a single transaction by its hash
func (d *MetadataStoreSqlite) DeleteTxByHash(txn *gorm.DB, txHash []byte) error {
	db := txn
	if db == nil {
		db = d.db
	}
	result := db.Where("transaction_hash = ?", txHash).Delete(&models.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteTxsByBlockNumber deletes all transactions for a given block number
func (d *MetadataStoreSqlite) DeleteTxsByBlockNumber(txn *gorm.DB, blockNumber uint64) error {
	db := txn
	if db == nil {
		db = d.db
	}
	result := db.Where("block_number = ?", blockNumber).Delete(&models.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
