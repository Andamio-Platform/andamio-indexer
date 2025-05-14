package sqlite

import (
	"errors"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// SetTx inserts or updates a transaction record
// SetTx inserts or updates a transaction record and its nested data
func (d *MetadataStoreSqlite) SetTx(tx *models.Transaction, txn *gorm.DB) error {
	if tx == nil {
		return errors.New("transaction cannot be nil")
	}
	if len(tx.TransactionHash) == 0 {
		return errors.New("transaction hash cannot be empty")
	}
	// Basic validation for nested data could be added here if needed,
	// but the individual setter functions should handle validation for their respective types.

	// Save main transaction metadata
	result := txn.Save(tx)
	if result.Error != nil {
		return result.Error
	}

	// Save nested data
	if err := d.setInputs(tx.Inputs, tx.TransactionHash, txn); err != nil {
		return err
	}
	if err := d.setOutputs(tx.Outputs, tx.TransactionHash, txn); err != nil {
		return err
	}
	if err := d.setReferenceInputs(tx.ReferenceInputs, tx.TransactionHash, txn); err != nil {
		return err
	}
	// Check if Witness is present and should be saved (e.g., by checking a required field)
	if len(tx.Witness.TransactionHash) > 0 { // Assuming TransactionHash is a required field for Witness
		if err := d.setWitness(tx.Witness, tx.TransactionHash, txn); err != nil {
			return err
		}
	}


	return nil
}

// setInputs saves a slice of transaction inputs to the database
func (d *MetadataStoreSqlite) setInputs(inputs []models.TransactionInput, txHash []byte, txn *gorm.DB) error {
	for _, input := range inputs {
		input.TransactionHash = txHash // Set foreign key
		result := txn.Save(&input)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// setOutputs saves a slice of transaction outputs to the database
func (d *MetadataStoreSqlite) setOutputs(outputs []models.TransactionOutput, txHash []byte, txn *gorm.DB) error {
	for _, output := range outputs {
		// Save the TransactionOutput record
		result := txn.Save(&output)
		if result.Error != nil {
			return result.Error
		}
		// Save nested assets within the output
		for _, asset := range output.Asset {
			// The relationship between Asset and TransactionOutput is via UTxOID and UTxOIDIndex
			// Ensure these fields are set correctly on the asset before saving
			asset.UTxOID = output.UTxOID
			asset.UTxOIDIndex = output.UTxOIDIndex
			if err := d.SetAsset(&asset); err != nil { // Assuming SetAsset takes a single asset and handles its saving
				return err
			}
		}
	}
	return nil
}

// setReferenceInputs saves a slice of reference inputs to the database
func (d *MetadataStoreSqlite) setReferenceInputs(refInputs []models.SimpleUTxO, txHash []byte, txn *gorm.DB) error {
	for _, refInput := range refInputs {
		refInput.TransactionHash = txHash // Set foreign key
		result := txn.Save(&refInput)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// setWitness saves a transaction witness and its nested data to the database
func (d *MetadataStoreSqlite) setWitness(witness models.Witness, txHash []byte, txn *gorm.DB) error {
	witness.TransactionHash = txHash // Set foreign key to Transaction

	// Save the Witness record
	result := txn.Save(&witness)
	if result.Error != nil {
		return result.Error
	}

	// Save nested Redeemers
	for _, redeemer := range witness.Redeemers {
		redeemer.WitnessID = witness.ID // Set foreign key to Witness
		result := txn.Save(&redeemer)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// GetTxByHash retrieves a single transaction by its hash
func (d *MetadataStoreSqlite) GetTxByHash(txHash []byte, txn *gorm.DB) (*models.Transaction, error) {
	var transaction models.Transaction
	result := txn.Where("transaction_hash = ?", txHash).
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
func (d *MetadataStoreSqlite) GetTxByID(id uint, txn *gorm.DB) (*models.Transaction, error) {
	var transaction models.Transaction
	result := txn.
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
func (d *MetadataStoreSqlite) GetTxsByBlockNumber(blockNumber uint64, limit, offset int, txn *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := txn.Where("block_number = ?", blockNumber)

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
func (d *MetadataStoreSqlite) GetTxsByInputAddress(address string, limit, offset int, txn *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := txn.Joins("JOIN transaction_inputs ON transactions.transaction_hash = transaction_inputs.transaction_hash").
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
func (d *MetadataStoreSqlite) GetTxsByOutputAddress(address string, limit, offset int, txn *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	var outputTxHashes [][]byte

	// Find the OutputID (transaction hash) of transactions that have an output to the given address
	result := txn.Model(&models.TransactionOutput{}).
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
	query := txn.Where("transaction_hash IN (?)", outputTxHashes)

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
func (d *MetadataStoreSqlite) GetTxsByAnyAddress(address string, limit, offset int, txn *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	uniqueTxHashes := make(map[string]bool)

	// Get transaction hashes from inputs
	var inputTxHashes [][]byte
	result := txn.Model(&models.TransactionInput{}).
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
	result = txn.Model(&models.TransactionOutput{}).
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
	query := txn.Where("transaction_hash IN (?)", finalTxHashes)

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
func (d *MetadataStoreSqlite) SetTxs(txs []*models.Transaction, txn *gorm.DB) error {
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
		result := txn.Save(tx)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// GetTxs retrieves a list of transactions with pagination
func (d *MetadataStoreSqlite) GetTxs(limit, offset int, txn *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := txn.Limit(limit).Offset(offset).
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
	var count int64
	result := txn.Model(&models.Transaction{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// DeleteTxByHash deletes a single transaction by its hash
func (d *MetadataStoreSqlite) DeleteTxByHash(txHash []byte, txn *gorm.DB) error {
	result := txn.Where("transaction_hash = ?", txHash).Delete(&models.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteTxsByBlockNumber deletes all transactions for a given block number
func (d *MetadataStoreSqlite) DeleteTxsByBlockNumber(blockNumber uint64, txn *gorm.DB) error {
	result := txn.Where("block_number = ?", blockNumber).Delete(&models.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
