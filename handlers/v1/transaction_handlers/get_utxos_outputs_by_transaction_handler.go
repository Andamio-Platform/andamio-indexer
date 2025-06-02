package transaction_handlers

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetUTxOsOutputsByTransactionHandler handles the GET /api/v1/indexer/transactions/{tx_hash}/utxos/outputs endpoint.
// @Summary		Get UTXO Outputs by Transaction Hash
// @Description	Retrieves transaction outputs for a given transaction hash.
// @ID				getUTxOsOutputsByTransaction
// @Tags			Transactions
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			tx_hash	path		string	true	"The transaction hash to retrieve UTXO outputs for."
// @Success		200		{array}		viewmodel.TransactionOutput	"Successfully retrieved UTXO outputs."
// @Failure		400		{object}	object{error=string}		"Invalid transaction hash."
// @Failure		404		{object}	object{error=string}		"Transaction not found or no UTXO outputs found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/transactions/{tx_hash}/utxos/outputs [get]
func GetUTxOsOutputsByTransactionHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHashStr := c.Params("tx_hash")
		txHash, err := hex.DecodeString(txHashStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction hash"})
		}

		transaction, err := db.Metadata().GetTxByTxHash(nil, txHash)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transaction by hash"})
		}

		if transaction == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
		}

		// Convert database models to view models
		outputViewModels := viewmodel.ConvertTransactionOutputsToViewModels(transaction.Outputs)

		return c.JSON(outputViewModels)
	}
}
