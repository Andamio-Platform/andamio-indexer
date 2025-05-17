package transaction_handlers

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetUTxOsInputsByTransactionHandler handles the GET /api/v1/indexer/transactions/{tx_hash}/utxos/inputs endpoint.
// @Summary Get UTXO Inputs by Transaction Hash
// @Description Retrieves transaction inputs for a given transaction hash.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param tx_hash path string true "Transaction hash to retrieve UTXO inputs for"
// @Success 200 {array} viewmodel.TransactionInput "Success response"
// @Failure 400 {object} errors.ServerError "Bad request"
// @Failure 404 {object} errors.ServerError "Transaction not found"
// @Failure 500 {object} errors.ServerError "Server error"
// @Router /indexer/transactions/{tx_hash}/utxos/inputs [get]
func GetUTxOsInputsByTransactionHandler(db *database.Database) fiber.Handler {
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
		inputViewModels := viewmodel.ConvertTransactionInputsToViewModels(transaction.Inputs)

		return c.JSON(inputViewModels)
	}
}
