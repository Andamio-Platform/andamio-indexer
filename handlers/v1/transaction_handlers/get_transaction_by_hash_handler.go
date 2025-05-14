package transaction_handlers

import (
	"net/http"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2"                 // Use fiber
	fiberLogger "github.com/gofiber/fiber/v2/log" // Use fiber logger
	"gorm.io/gorm"
)

// GetTransactionByHashHandler handles the request to get a transaction by its hash.
func GetTransactionByHashHandler(c *fiber.Ctx, db *database.Database) error { // Use fiber.Ctx and accept db
	txHash := c.Params("tx_hash")

	if db == nil {
		fiberLogger.Error("database not available")                                                        // Use fiber logger
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "database not available"}) // Use fiber JSON
	}

	txn := db.MetadataTxn(false) // Read-only transaction
	defer txn.Discard()          // Use Discard for read-only transaction

	var transaction models.Transaction
	result := txn.Metadata().Where("tx_hash = ?", txHash).First(&transaction)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fiberLogger.Errorf("Transaction not found for hash %s", txHash)                        // Use fiber logger
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "transaction not found"}) // Use fiber JSON
		} else {
			fiberLogger.Errorf("Failed to fetch transaction for hash %s: %v", txHash, result.Error)                                                  // Use fiber logger
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch transaction", "details": result.Error.Error()}) // Use fiber JSON
		}
	}

	// Fetch and include related inputs, outputs, metadata, witnesses, datum, and redeemers
	result = txn.Metadata().
		Preload("Block").
		Preload("Inputs.Utxo").
		Preload("Inputs.Redeemer").
		Preload("Outputs.Utxo").
		Preload("Outputs.Address").
		Preload("Outputs.Datum").
		Preload("Witnesses").
		Where("tx_hash = ?", txHash).
		First(&transaction)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fiberLogger.Errorf("Transaction not found for hash %s with relations", txHash)         // Use fiber logger
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "transaction not found"}) // Use fiber JSON
		} else {
			fiberLogger.Errorf("Failed to fetch transaction with relations for hash %s: %v", txHash, result.Error)                                                  // Use fiber logger
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch transaction with relations", "details": result.Error.Error()}) // Use fiber JSON
		}
	}

	return c.Status(http.StatusOK).JSON(transaction) // Use fiber JSON
}
