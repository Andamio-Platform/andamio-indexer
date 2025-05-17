package witness_handlers

import (
	"log/slog"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetTransactionsByWitnessIdHandler godoc
// @Summary Get Transactions by Witness ID
// @Description Retrieve a list of transactions that include a specific witness ID, with support for pagination.
// @ID getTransactionsByWitnessId
// @Tags Witnesses
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param witness_id path string true "The witness ID to retrieve transactions for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} viewmodel.Transaction "Successfully retrieved transactions."
// @Failure 400 {object} object{error=string} "Invalid witness ID or pagination parameters."
// @Failure 404 {object} object{error=string} "Witness ID not found or no transactions found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /witnesses/{witness_id}/transactions [get]
func GetTransactionsByWitnessIdHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		witnessIDStr := c.Params("witness_id")
		if witnessIDStr == "" {
			logger.Error("witness_id path parameter is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "witness_id path parameter is missing"})
		}

		witnessID, err := strconv.ParseUint(witnessIDStr, 10, 64)
		if err != nil {
			logger.Error("invalid witness_id format", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid witness_id format"})
		}

		witness, err := db.Metadata().GetWitnessByID(nil, uint(witnessID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "witness not found"})
			}
			logger.Error("failed to get witness by ID", "witness_id", witnessID, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transactions"})
		}

		if witness == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "witness not found"})
		}

		tx, err := db.Metadata().GetTxByTxHash(nil, witness.TransactionHash)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction for witness not found"})
			}
			logger.Error("failed to get transaction by hash for witness", "witness_id", witnessID, "tx_hash", witness.TransactionHash, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transactions"})
		}

		if tx == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction for witness not found"})
		}

		// Convert database model to view model and return in a slice
		transactionViewModels := []viewmodel.Transaction{
			{
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
			},
		}

		return c.Status(fiber.StatusOK).JSON(transactionViewModels)
	}
}
