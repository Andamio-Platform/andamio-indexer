package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetLatestBlockHandler godoc
// @Summary Get Latest Indexed Block
// @Description Retrieves the block number and slot number of the most recently indexed block.
// @ID getLatestBlock
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{block_number=uint64,slot_number=uint64} "Successfully retrieved latest block information."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/latest-block [get]
func GetLatestBlockHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var latestTx models.Transaction
		result := db.Metadata().DB().Model(&models.Transaction{}).
			Order("block_number DESC, slot_number DESC").
			Limit(1).
			First(&latestTx)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// No transactions found, return 0 for block and slot
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"block_number": 0,
					"slot_number":  0,
				})
			}
			logger.Error("Error getting latest block", "error", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get latest block information",
			})
		}

		blockNumber := latestTx.BlockNumber
		slotNumber := latestTx.SlotNumber

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"block_number": blockNumber,
			"slot_number":  slotNumber,
		})
	}
}