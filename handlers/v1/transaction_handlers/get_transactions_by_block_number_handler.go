package transaction_handlers

import (
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetTransactionsByBlockNumberHandler handles the GET /api/v1/indexer/transactions/{block_number} endpoint.
// @Summary		Get Transactions by Block Number
// @Description	Retrieves transactions for a given block number with pagination.
// @ID				getTransactionsByBlockNumber
// @Tags			Transactions
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			block_number	path		int		true	"The block number to retrieve transactions for."
// @Param			limit	query		int		false	"Maximum number of results to return."	default(100)
// @Param			offset	query		int		false	"Number of results to skip."	default(0)
// @Success		200		{array}		viewmodel.Transaction	"Successfully retrieved transactions."
// @Failure		400		{object}	object{error=string}		"Invalid block number or pagination parameters."
// @Failure		404		{object}	object{error=string}		"Block number not found or no transactions found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/indexer/transactions/{block_number} [get]
func GetTransactionsByBlockNumberHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		blockNumberStr := c.Params("block_number")
		blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid block number"})
		}

		// Get pagination parameters
		limitStr := c.Query("limit", "100")
		offsetStr := c.Query("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter"})
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid offset parameter"})
		}

		transactions, err := db.Metadata().GetTxsByBlockNumber(nil, blockNumber, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transactions by block number"})
		}

		// Convert database models to view models
		transactionViewModels := []viewmodel.Transaction{}
		for _, tx := range transactions {
			transactionViewModels = append(transactionViewModels, viewmodel.Transaction{
				TransactionHash: string(tx.TransactionHash),
				BlockNumber:     tx.BlockNumber,
				SlotNumber:      tx.SlotNumber,
				Inputs:          viewmodel.ConvertTransactionInputsToViewModels(tx.Inputs),
				Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(tx.Outputs),
				Fee:             tx.Fee,
				TTL:             tx.TTL,
				BlockHash:       string(tx.BlockHash),
				Metadata:        string(tx.Metadata), // CBOR string representation
				ReferenceInputs: viewmodel.ConvertSimpleUTxOModelsToViewModels(tx.ReferenceInputs),
				Withdrawals:     tx.Withdrawals,
				Certificates:    viewmodel.ConvertByteSliceSliceToStringSlice(tx.Certificates), // Convert [][]byte to []string
				Witness:         viewmodel.ConvertWitnessModelToViewModel(tx.Witness),
			})
		}

		return c.JSON(transactionViewModels)
	}
}
