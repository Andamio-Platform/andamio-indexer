package transaction_handlers

import (
	"encoding/hex"
	"log/slog"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetInputDatumHandler godoc
// @Summary Get Datum for Transaction Input
// @Description Retrieves the datum associated with a specific transaction input.
// @ID getInputDatum
// @Tags Transactions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tx_hash path string true "The transaction hash (hex-encoded)."
// @Param input_index path int true "The index of the transaction input."
// @Success 200 {object} viewmodel.Datum "Successfully retrieved datum."
// @Failure 400 {object} object{error=string} "Invalid transaction hash or input index."
// @Failure 404 {object} object{error=string} "Datum not found for the specified input."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /transactions/{tx_hash}/inputs/{input_index}/datum [get]
func GetInputDatumHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHashHex := c.Params("tx_hash")
		inputIndexStr := c.Params("input_index")

		txHash, err := hex.DecodeString(txHashHex)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid transaction hash format. Must be hex-encoded.",
			})
		}

		inputIndex, err := strconv.Atoi(inputIndexStr)
		if err != nil || inputIndex < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input index.",
			})
		}

		tx, err := db.Metadata().GetTxByTxHash(nil, txHash)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Transaction not found.",
				})
			}
			logger.Error("Error getting transaction by hash", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get input datum",
			})
		}

		if tx == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found.",
			})
		}

		if inputIndex < 0 || inputIndex >= len(tx.Inputs) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Input index out of bounds for this transaction.",
			})
		}

		input := tx.Inputs[inputIndex]

		// Get the datum using the UTxO ID and index from the input
		datum, err := db.Metadata().GetDatum(nil, input.UTxOID, input.UTxOIDIndex)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Datum not found for the specified input.",
				})
			}
			logger.Error("Error getting datum by UTxO", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get input datum",
			})
		}

		if datum == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Datum not found for the specified input.",
			})
		}

		// Convert database model to view model
		datumViewModel := viewmodel.Datum{
			DatumHash: string(datum.DatumHash),
			DatumCbor: string(datum.DatumCbor), // Assuming DatumCbor should be a string representation of CBOR
		}

		return c.Status(fiber.StatusOK).JSON(datumViewModel)
	}
}
