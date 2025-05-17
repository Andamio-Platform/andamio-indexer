package transaction_handlers

import (
	"encoding/hex"
	"log/slog"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetInputRedeemerHandler godoc
// @Summary Get Redeemer for Transaction Input
// @Description Retrieves the redeemer associated with a specific transaction input.
// @ID getInputRedeemer
// @Tags Transactions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tx_hash path string true "The transaction hash (hex-encoded)."
// @Param input_index path int true "The index of the transaction input."
// @Success 200 {object} viewmodel.Redeemer "Successfully retrieved redeemer."
// @Failure 400 {object} object{error=string} "Invalid transaction hash or input index."
// @Failure 404 {object} object{error=string} "Redeemer not found for the specified input."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /transactions/{tx_hash}/inputs/{input_index}/redeemer [get]
func GetInputRedeemerHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
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
				"error": "Failed to get input redeemer",
			})
		}

		if tx == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Transaction not found.",
			})
		}

		// Check if the transaction has a witness and redeemers
		if tx.Witness.ID == 0 || len(tx.Witness.Redeemers) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Redeemer not found for the specified input.",
			})
		}

		// Find the redeemer with the matching index and spending tag (0)
		var foundRedeemer *models.Redeemer
		for i := range tx.Witness.Redeemers {
			redeemer := &tx.Witness.Redeemers[i]
			// Assuming input index corresponds to redeemer index and tag 0 for spending redeemers
			if redeemer.Index == uint(inputIndex) && redeemer.Tag == 0 {
				foundRedeemer = redeemer
				break
			}
		}

		if foundRedeemer == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Redeemer not found for the specified input.",
			})
		}

		// Convert database model to view model
		redeemerViewModel := viewmodel.Redeemer{
			Cbor: string(foundRedeemer.Cbor), // CBOR string representation
		}

		return c.Status(fiber.StatusOK).JSON(redeemerViewModel)
	}
}
