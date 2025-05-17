package transaction_handlers

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetTransactionByHashHandler handles the request to get a transaction by its hash.
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
func GetTransactionByHashHandler(db *database.Database) fiber.Handler { // Use fiber.Ctx and accept db
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

		// Convert database model to view model
		transactionViewModel := viewmodel.Transaction{
			TransactionHash: string(transaction.TransactionHash),
			BlockNumber:     transaction.BlockNumber,
			SlotNumber:      transaction.SlotNumber,
			Inputs:          viewmodel.ConvertTransactionInputsToViewModels(transaction.Inputs),
			Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(transaction.Outputs),
			Fee:             transaction.Fee,
			TTL:             transaction.TTL,
			BlockHash:       string(transaction.BlockHash),
			Metadata:        string(transaction.Metadata), // CBOR string representation
			ReferenceInputs: viewmodel.ConvertSimpleUTxOModelsToViewModels(transaction.ReferenceInputs),
			Withdrawals:     transaction.Withdrawals,
			Certificates:    viewmodel.ConvertByteSliceSliceToStringSlice(transaction.Certificates), // Convert [][]byte to []string
			Witness:         viewmodel.ConvertWitnessModelToViewModel(transaction.Witness),
		}

		return c.Status(fiber.StatusOK).JSON(transactionViewModel) // Use fiber JSON
	}
}
