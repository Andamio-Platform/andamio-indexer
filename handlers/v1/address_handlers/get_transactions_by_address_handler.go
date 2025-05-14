package address_handlers

import (
	"time"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// GetTransactionsByAddressHandler handles the request to get transactions for a specific address.
//
//	@Summary		Get Transactions by Address
//	@Description	Get transactions associated with a specific address.
//	@ID			getTransactionsByAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string	true	"Address to get transactions for"
//	@Success		200		{object}	[]models.Transaction	"Success response" // TODO: Define a proper response struct
//	@Failure		400		{object}	errors.ServerError	"Bad request" // TODO: Use project's error handling
//	@Failure		404		{object}	errors.ServerError	"Address not found" // TODO: Use project's error handling
//	@Failure		500		{object}	errors.ServerError	"Server error" // TODO: Use project's error handling
//	@Router			/addresses/{address}/transactions [get]
func GetTransactionsByAddressHandler(c *fiber.Ctx, db *database.Database) error {
	address := c.Params("address")

	if address == "" {
		fiberLogger.Error("Address is required")
		return fiber.NewError(fiber.StatusBadRequest, "Address is required")
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

	var addrModel models.Address
	result := txn.Metadata().Where("address = ?", address).First(&addrModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fiberLogger.Errorf("Address not found: %s", address)
			return fiber.NewError(fiber.StatusNotFound, "address not found")
		} else {
			fiberLogger.Errorf("Failed to fetch address %s: %v", address, result.Error)
			return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch address")
		}
	}

	// Fetch transaction IDs where the address is in the outputs
	var outputTxIDs []uint
	result = txn.Metadata().Model(&models.TransactionOutput{}).
		Joins("JOIN transactions ON transaction_outputs.transaction_id = transactions.id").
		Joins("JOIN blocks ON transactions.block_id = blocks.id").
		Where("transaction_outputs.address_id = ?", addrModel.ID)

	if startTime != nil {
		result = result.Where("blocks.block_time >= ?", *startTime)
	}
	if endTime != nil {
		result = result.Where("blocks.block_time <= ?", *endTime)
	}
	result = result.Pluck("transaction_id", &outputTxIDs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch transaction output IDs for address %s: %v", address, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	// Fetch transaction IDs where the address is in the inputs
	var inputTxIDs []uint
	result = txn.Metadata().Model(&models.TransactionInput{}).
		Joins("JOIN utxos ON transaction_inputs.utxo_id = utxos.id").
		Joins("JOIN addresses ON utxos.payment_key = addresses.address OR utxos.staking_key = addresses.address"). // Join with addresses to match address string
		Joins("JOIN transactions ON transaction_inputs.transaction_id = transactions.id").
		Joins("JOIN blocks ON transactions.block_id = blocks.id").
		Where("addresses.id = ?", addrModel.ID) // Filter by the fetched address ID

	if startTime != nil {
		result = result.Where("blocks.block_time >= ?", *startTime)
	}
	if endTime != nil {
		result = result.Where("blocks.block_time <= ?", *endTime)
	}
	result = result.Pluck("transaction_id", &inputTxIDs)


	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch transaction input IDs for address %s: %v", address, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	// Combine and get unique transaction IDs
	uniqueTxIDsMap := make(map[uint]bool)
	for _, id := range outputTxIDs {
		uniqueTxIDsMap[id] = true
	}
	for _, id := range inputTxIDs {
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
		fiberLogger.Errorf("Failed to fetch paginated transactions for address %s: %v", address, result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch transactions")
	}

	return c.JSON(transactions)
}