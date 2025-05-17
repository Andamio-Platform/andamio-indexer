package transaction_handlers

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetUTxOsByTransactionHandler handles the request to get UTxOs (inputs and outputs) for a specific transaction.
//
//	@Summary		Get UTxOs by Transaction Hash
//	@Description	Retrieves the unspent transaction outputs (UTxOs) and inputs associated with a specific transaction hash.
//	@ID				getUTxOsByTransaction
//	@Tags			Transactions
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			tx_hash	path		string	true	"The transaction hash (hex-encoded) to retrieve UTxOs for."
//	@Success		200		{object}	viewmodel.TransactionUTxOs	"Successfully retrieved UTxOs."
//	@Failure		400		{object}	object{error=string}		"Invalid transaction hash."
//	@Failure		404		{object}	object{error=string}		"Transaction not found or no UTxOs found for the given hash."
//	@Failure		500		{object}	object{error=string}		"Internal server error."
//	@Router			/transactions/{tx_hash}/utxos [get]
func GetUTxOsByTransactionHandler(db *database.Database) fiber.Handler { // Use fiber.Ctx and accept db
	return func(c *fiber.Ctx) error {
		txHashStr := c.Params("tx_hash")
		txHash, err := hex.DecodeString(txHashStr)
		if err != nil {
			fiberLogger.Errorf("Invalid transaction hash: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction hash"})
		}

		if db == nil {
			fiberLogger.Error("database not available")                                                         // Use fiber logger
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database not available"}) // Use fiber JSON
		}

		transaction, err := db.Metadata().GetTxByTxHash(nil, txHash)
		if err != nil {
			fiberLogger.Errorf("Failed to get transaction by hash %s: %v", txHashStr, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transaction by hash"})
		}

		if transaction == nil {
			fiberLogger.Errorf("Transaction not found for hash %s", txHashStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
		}

		// Convert database models to view models
		transactionUTxOs := viewmodel.TransactionUTxOs{
			Inputs:  viewmodel.ConvertTransactionInputsToViewModels(transaction.Inputs),
			Outputs: viewmodel.ConvertTransactionOutputsToViewModels(transaction.Outputs),
		}

		return c.Status(fiber.StatusOK).JSON(transactionUTxOs) // Use fiber JSON
	}
}
