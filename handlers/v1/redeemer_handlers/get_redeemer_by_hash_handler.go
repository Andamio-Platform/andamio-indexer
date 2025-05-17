package redeemer_handlers

import (
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// GetRedeemersByTxHashHandler godoc
// @Summary Get Redeemers by Transaction Hash
// @Description Retrieve a list of redeemers associated with a specific transaction hash.
// @ID getRedeemersByTxHash
// @Tags Transactions
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param tx_hash path string true "The transaction hash (hex-encoded) to retrieve redeemers for."
// @Success 200 {array} viewmodel.Redeemer "Successfully retrieved redeemers."
// @Failure 400 {object} object{error=string} "Invalid transaction hash."
// @Failure 404 {object} object{error=string} "Transaction not found or no redeemers found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /transactions/{tx_hash}/redeemers [get]
func GetRedeemersByTxHashHandler(db *database.Database, log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txHashHex := c.Params("tx_hash")
		if txHashHex == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing tx_hash path parameter"})
		}

		txHash, err := hex.DecodeString(txHashHex)
		if err != nil {
			log.Error("failed to decode transaction hash", "error", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid transaction hash"})
		}

		witness, err := db.Metadata().GetWitnessByTransactionHash(nil, txHash)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "redeemers not found for this transaction hash"})
			}
			log.Error("failed to get witness by transaction hash", "tx_hash", txHashHex, "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if witness == nil || len(witness.Redeemers) == 0 {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "redeemers not found for this transaction hash"})
		}

		// Convert database models to view models
		redeemerViewModels := viewmodel.ConvertRedeemersToViewModels(witness.Redeemers)

		return c.Status(http.StatusOK).JSON(redeemerViewModels)
	}
}
