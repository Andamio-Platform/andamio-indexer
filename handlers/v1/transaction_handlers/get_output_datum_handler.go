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

// GetOutputDatumHandler godoc
// @Summary Get Datum for Transaction Output
// @Description Retrieves the datum associated with a specific transaction output.
// @ID getOutputDatum
// @Tags Transactions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tx_hash path string true "The transaction hash (hex-encoded)."
// @Param output_index path int true "The index of the transaction output."
// @Success 200 {object} viewmodel.Datum "Successfully retrieved datum."
// @Failure 400 {object} object{error=string} "Invalid transaction hash or output index."
// @Failure 404 {object} object{error=string} "Datum not found for the specified output."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /transactions/{tx_hash}/outputs/{output_index}/datum [get]
func GetOutputDatumHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHashHex := c.Params("tx_hash")
		outputIndexStr := c.Params("output_index")

		txHash, err := hex.DecodeString(txHashHex)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid transaction hash format. Must be hex-encoded.",
			})
		}

		outputIndex, err := strconv.Atoi(outputIndexStr)
		if err != nil || outputIndex < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid output index.",
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
				"error": "Failed to get output datum",
			})
		}

		if tx == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found.",
			})
		}

		if outputIndex < 0 || outputIndex >= len(tx.Outputs) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Output index out of bounds for this transaction.",
			})
		}

		output := tx.Outputs[outputIndex]

		// Get the datum using the UTxO ID and index from the output
		datum, err := db.Metadata().GetDatum(nil, output.UTxOID, output.UTxOIDIndex)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Datum not found for the specified output.",
				})
			}
			logger.Error("Error getting datum by UTxO", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get output datum",
			})
		}

		if datum == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Datum not found for the specified output.",
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
