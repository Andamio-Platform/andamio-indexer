package transaction_handlers

import (
	"net/http" // Import http for status codes

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2" // Use fiber
	fiberLogger "github.com/gofiber/fiber/v2/log" // Use fiber logger
	"gorm.io/gorm"
)

// GetUTxOsByTransactionHandler handles the request to get UTxOs created by a specific transaction.
//
//	@Summary		Get UTxOs by Transaction Hash
//	@Description	Get UTxOs created as outputs by a specific transaction.
//	@ID			getUTxOsByTransaction
//	@Tags			Transactions
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			tx_hash	path		string	true	"Transaction hash to get UTxOs for"
//	@Success		200		{object}	[]models.TransactionOutput	"Success response" // TODO: Define a proper response struct
//	@Failure		400		{object}	errors.ServerError	"Bad request" // TODO: Use project's error handling
//	@Failure		404		{object}	errors.ServerError	"Transaction not found" // TODO: Use project's error handling
//	@Failure		500		{object}	errors.ServerError	"Server error" // TODO: Use project's error handling
//	@Router			/transactions/{tx_hash}/utxos [get]
func GetUTxOsByTransactionHandler(c *fiber.Ctx, db *database.Database) error { // Use fiber.Ctx and accept db
	txHash := c.Params("tx_hash")

	if txHash == "" {
		fiberLogger.Error("Transaction hash is required") // Use fiber logger
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Transaction hash is required"}) // Use fiber JSON
	}

	if db == nil {
		fiberLogger.Error("database not available") // Use fiber logger
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database not available"}) // Use fiber JSON
	}

	txn := db.MetadataTxn(false) // Read-only transaction
	defer txn.Discard() // Use Discard for read-only transaction

	var transaction models.Transaction
	result := txn.Metadata().Where("tx_hash = ?", txHash).First(&transaction)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fiberLogger.Errorf("Transaction not found for hash %s", txHash) // Use fiber logger
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction not found"}) // Use fiber JSON
		} else {
			fiberLogger.Errorf("Failed to fetch transaction for hash %s: %v", txHash, result.Error) // Use fiber logger
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch transaction", "details": result.Error.Error()}) // Use fiber JSON
		}
	}

	var utxos []models.TransactionOutput
	result = txn.Metadata().
		Preload("Utxo.UtxoAssets.Asset"). // Preload Utxo and nested UtxoAssets.Asset
		Where("transaction_id = ?", transaction.ID).
		Find(&utxos)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		fiberLogger.Errorf("Failed to fetch UTxOs for transaction hash %s: %v", txHash, result.Error) // Use fiber logger
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch UTxOs"}) // Use fiber JSON
	}

	return c.Status(fiber.StatusOK).JSON(utxos) // Use fiber JSON
}