package transaction_handlers

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetTransactionByTxHashHandler handles the request to get a transaction by its hash.
//
//	@Summary		Get Transaction by Hash
//	@Description	Retrieves a transaction by its hash.
//	@ID				getTransactionByHash
//	@Tags			Transactions
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			tx_hash	path		string	true	"The transaction hash to retrieve."
//	@Success		200		{object}	viewmodel.Transaction	"Successfully retrieved transaction."
//	@Failure		400		{object}	object{error=string}		"Invalid transaction hash."
//	@Failure		404		{object}	object{error=string}		"Transaction not found."
//	@Failure		500		{object}	object{error=string}		"Internal server error."
//	@Router			/transactions/{tx_hash} [get]
func GetTransactionByTxHashHandler(db *database.Database) fiber.Handler { // Use fiber.Ctx and accept db
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

		tx, err := db.GetTxByTxHash(txHash, nil)
		if err != nil {
			fiberLogger.Errorf("Failed to get transaction by hash %s: %v", txHashStr, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transaction by hash"})
		}

		if tx == nil {
			fiberLogger.Errorf("Transaction not found for hash %s", txHashStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
		}

		// Convert database model to view model
		transactionViewModel := viewmodel.Transaction{
			TransactionHash: hex.EncodeToString(tx.TransactionHash),
			BlockNumber:     tx.BlockNumber,
			SlotNumber:      tx.SlotNumber,
			Inputs:          viewmodel.ConvertTransactionInputsToViewModels(tx.Inputs),
			Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(tx.Outputs),
			Fee:             tx.Fee,
			TTL:             tx.TTL,
			BlockHash:       string(tx.BlockHash),
			Metadata:        hex.EncodeToString(tx.Metadata),
			ReferenceInputs: viewmodel.ConvertSimpleUTxOModelsToViewModels(tx.ReferenceInputs),
			Withdrawals:     tx.Withdrawals,
			Certificates:    viewmodel.ConvertByteSliceSliceToStringSlice(tx.Certificates),
			Witness:         viewmodel.ConvertWitnessModelToViewModel(tx.Witness),
			TransactionCBOR: hex.EncodeToString(tx.TransactionCBOR),
		}

		return c.Status(fiber.StatusOK).JSON(transactionViewModel) // Use fiber JSON
	}
}
