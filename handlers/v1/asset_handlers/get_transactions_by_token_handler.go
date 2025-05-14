package asset_handlers

import (
	"time"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// GetTransactionsByTokenHandler handles the request to get transactions for a specific asset token name.
//
//	@Summary		Get Transactions by Token Name
//	@Description	Get transactions associated with a specific asset token name.
//	@ID			getTransactionsByToken
//	@Tags			Assets
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			tokenname	path		string	true	"Asset Token Name to get transactions for"
//	@Param			limit		query		int		false	"Pagination limit"
//	@Param			offset		query		int		false	"Pagination offset"
//	@Param			startTime	query		string	false	"Start time (RFC3339)"
//	@Param			endTime		query		string	false	"End time (RFC3339)"
//	@Success		200		{object}	[]models.Transaction	"Success response" // TODO: Define a proper response struct
//	@Failure		400		{object}	errors.ServerError	"Bad request" // TODO: Use project's error handling
//	@Failure		500		{object}	errors.ServerError	"Server error" // TODO: Use project's error handling
//	@Router			/assets/token/{tokenname}/transactions [get]
func GetTransactionsByTokenHandler(c *fiber.Ctx, db *database.Database) error {
	tokenname := c.Params("tokenname")

	if tokenname == "" {
		fiberLogger.Error("Token name is required")
		return fiber.NewError(fiber.StatusBadRequest, "Token name is required")
	}

	// Get pagination parameters
	limit := c.QueryInt("limit", 100) // Default limit to 100
	offset := c.QueryInt("offset", 0) // Default offset to 0

	// Get time range parameters
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime *time.Time

	if startTimeStr != "" {
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			fiberLogger.Errorf("Invalid startTime format: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid startTime format, use RFC3339")
		}
		startTime = &t
	}

	if endTimeStr != "" {
		t, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			fiberLogger.Errorf("Invalid endTime format: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "invalid endTime format, use RFC3339")
		}
		endTime = &t
	}

	globalDB := database.GetGlobalDB()
	if globalDB == nil {
		fiberLogger.Error("database not available")
		return fiber.NewError(fiber.StatusInternalServerError, "database not available")
	}

	txn := globalDB.MetadataTxn(false) // Read-only transaction
	defer txn.Discard()

	// Fetch asset IDs for the given token name
	var assetIDs []uint
	result := txn.Metadata().Model(&models.Asset{}).
		Where("asset_name = ?", tokenname).
		Pluck("id", &assetIDs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch asset IDs for token %s: %v", tokenname, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	if len(assetIDs) == 0 {
		return c.JSON([]models.Transaction{}) // No assets found with this token name
	}

	// Fetch utxo IDs associated with these asset IDs
	var utxoIDs []uint
	result = txn.Metadata().Model(&models.UtxoAsset{}).
		Where("asset_id IN (?)", assetIDs).
		Pluck("utxo_id", &utxoIDs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch utxo IDs for assets with token %s: %v", tokenname, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	if len(utxoIDs) == 0 {
		return c.JSON([]models.Transaction{}) // No utxos found for these assets
	}

	// Fetch transaction IDs associated with these utxo IDs (either as input or output)
	var inputTxIDs []uint
	result = txn.Metadata().Model(&models.TransactionInput{}).
		Joins("JOIN transactions ON transaction_inputs.transaction_id = transactions.id").
		Joins("JOIN blocks ON transactions.block_id = blocks.id").
		Where("transaction_inputs.utxo_id IN (?)", utxoIDs)

	if startTime != nil {
		result = result.Where("blocks.block_time >= ?", *startTime)
	}
	if endTime != nil {
		result = result.Where("blocks.block_time <= ?", *endTime)
	}
	result = result.Pluck("transaction_id", &inputTxIDs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch input transaction IDs for utxos: %v", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	var outputTxIDs []uint
	result = txn.Metadata().Model(&models.TransactionOutput{}).
		Joins("JOIN transactions ON transaction_outputs.transaction_id = transactions.id").
		Joins("JOIN blocks ON transactions.block_id = blocks.id").
		Where("transaction_outputs.utxo_id IN (?)", utxoIDs)

	if startTime != nil {
		result = result.Where("blocks.block_time >= ?", *startTime)
	}
	if endTime != nil {
		result = result.Where("blocks.block_time <= ?", *endTime)
	}
	result = result.Pluck("transaction_id", &outputTxIDs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch output transaction IDs for utxos: %v", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	// Combine and get unique transaction IDs
	uniqueTxIDsMap := make(map[uint]bool)
	for _, id := range inputTxIDs {
		uniqueTxIDsMap[id] = true
	}
	for _, id := range outputTxIDs {
		uniqueTxIDsMap[id] = true
	}

	var uniqueTxIDs []uint
	for id := range uniqueTxIDsMap {
		uniqueTxIDs = append(uniqueTxIDs, id)
	}

	// Apply pagination to unique transaction IDs
	start := offset
	end := offset + limit
	if start > len(uniqueTxIDs) {
		start = len(uniqueTxIDs)
	}
	if end > len(uniqueTxIDs) {
		end = len(uniqueTxIDs)
	}
	paginatedTxIDs := uniqueTxIDs[start:end]

	// Fetch the full transactions for the paginated IDs, preloading relationships
	var transactions []models.Transaction
	result = txn.Metadata().
		Preload("Block").
		Preload("Inputs.Utxo").
		Preload("Inputs.Redeemer").
		Preload("Outputs.Utxo").
		Preload("Outputs.Address").
		Preload("Outputs.Datum").
		Preload("Witnesses").
		Where("id IN (?)", paginatedTxIDs).
		Find(&transactions)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch paginated transactions for token %s: %v", tokenname, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	return c.JSON(transactions)
}