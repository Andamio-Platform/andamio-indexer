package transaction_handlers

import (
	"encoding/hex"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// GetTransactionsBySlotRangeHandler godoc
// @Summary Get Transactions by Slot Range
// @Description Retrieves transactions within a specified slot number range, with support for pagination.
// @ID getTransactionsBySlotRange
// @Tags Transactions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param start_slot query uint64 true "The start slot number of the range (inclusive)."
// @Param end_slot query uint64 true "The end slot number of the range (inclusive)."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} viewmodel.Transaction "Successfully retrieved transactions."
// @Failure 400 {object} object{error=string} "Invalid slot number or pagination parameters."
// @Failure 404 {object} object{error=string} "No transactions found within the specified slot range."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /transactions/by-slot-range [get]
func GetTransactionsBySlotRangeHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startSlotStr := c.Query("start_slot")
		endSlotStr := c.Query("end_slot")
		limitStr := c.Query("limit", "100")
		offsetStr := c.Query("offset", "0")

		startSlot, err := strconv.ParseUint(startSlotStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start_slot format. Must be a valid unsigned integer.",
			})
		}

		endSlot, err := strconv.ParseUint(endSlotStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end_slot format. Must be a valid unsigned integer.",
			})
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid limit parameter.",
			})
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid offset parameter.",
			})
		}

		transactions, err := db.GetTxsBySlotRange(startSlot, endSlot, limit, offset, nil)
		if err != nil {
			logger.Error("Error getting transactions by slot range", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		if len(transactions) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No transactions found within the specified slot range.",
			})
		}

		// Convert database models to view models
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

		return c.Status(fiber.StatusOK).JSON(transactionViewModels)
	}
}
