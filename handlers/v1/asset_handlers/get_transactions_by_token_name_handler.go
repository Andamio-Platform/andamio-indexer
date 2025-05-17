package asset_handlers

import (
	"encoding/hex"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetTransactionsByTokenNameHandler handles the GET /api/v1/indexer/assets/token/{tokenname}/transactions endpoint.
// @Summary		Get Transactions by Token Name
// @Description	Retrieves transactions associated with a given token name with pagination.
// @ID				getTransactionsByTokenName
// @Tags			Assets
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			tokenname	path		string	true	"The token name to retrieve transactions for (hex-encoded)."
// @Param			limit	query		int		false	"Maximum number of results to return."	default(100)
// @Param			offset	query		int		false	"Number of results to skip."	default(0)
// @Success		200		{array}		viewmodel.Transaction	"Successfully retrieved transactions."
// @Failure		400		{object}	object{error=string}		"Invalid token name or pagination parameters."
// @Failure		404		{object}	object{error=string}		"Token name not found or no transactions found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/indexer/assets/token/{tokenname}/transactions [get]
func GetTransactionsByTokenNameHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenNameStr := c.Params("tokenname")
		tokenName, err := hex.DecodeString(tokenNameStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid token name"})
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

		transactions, err := db.Metadata().GetTxsByTokenName(nil, tokenName, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transactions by token name"})
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