package asset_handlers

import (
	"encoding/hex"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetTransactionsByPolicyIdAndTokenNameHandler handles the GET /api/v1/indexer/assets/policy/{policyId}/token/{tokenname}/transactions endpoint.
// @Summary		Get Transactions by Policy ID and Token Name
// @Description	Retrieves transactions associated with a given policy ID and token name with pagination.
// @ID				getTransactionsByPolicyIdAndTokenName
// @Tags			Assets
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			policyId	path		string	true	"The policy ID to retrieve transactions for (hex-encoded)."
// @Param			tokenname	path		string	true	"The token name to retrieve transactions for (hex-encoded)."
// @Param			limit	query		int		false	"Maximum number of results to return."	default(100)
// @Param			offset	query		int		false	"Number of results to skip."	default(0)
// @Success		200		{array}		viewmodel.Transaction	"Successfully retrieved transactions."
// @Failure		400		{object}	object{error=string}		"Invalid policy ID, token name, or pagination parameters."
// @Failure		404		{object}	object{error=string}		"Policy ID and token name combination not found or no transactions found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/assets/policy/{policyId}/token/{tokenname}/transactions [get]
func GetTransactionsByPolicyIdAndTokenNameHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		policyId := c.Params("policyId")

		tokenName := c.Params("tokenname")

		// Debug logging for policyId and tokenName
		db.Logger().Debug("GetTransactionsByPolicyIdAndTokenNameHandler", "policyId_raw", policyId, "tokenName_raw", tokenName)
		db.Logger().Debug("GetTransactionsByPolicyIdAndTokenNameHandler", "policyId_bytes_hex", hex.EncodeToString([]byte(policyId)), "tokenName_bytes_hex", hex.EncodeToString([]byte(tokenName)))

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

		transactions, err := db.GetTxsByPolicyIdAndTokenName([]byte(policyId), []byte(tokenName), limit, offset, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get transactions by policy ID and token name"})
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
