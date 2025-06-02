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
// @Router			/assets/token/{tokenname}/transactions [get]
func GetTransactionsByTokenNameHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenName := c.Params("tokenname")

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

		transactions, err := db.GetTxsByTokenName([]byte(tokenName), limit, offset, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transactions by token name"})
		}

		transactionViewModels := []viewmodel.Transaction{}
		for _, tx := range transactions {
			transactionViewModels = append(transactionViewModels, viewmodel.Transaction{
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
			})
		}

		return c.JSON(transactionViewModels)
	}
}
