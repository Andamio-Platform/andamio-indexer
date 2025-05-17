package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database" // Corrected import path
	"github.com/gofiber/fiber/v2"
)

// GetTransactionsCountHandler godoc
// @Summary Get Total Transaction Count
// @Description Retrieves the total number of indexed transactions.
// @ID getTransactionsCount
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{count=int} "Successfully retrieved transaction count."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/transactions/count [get]
func GetTransactionsCountHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := db.Metadata().CountTxs(nil)
		if err != nil {
			logger.Error("Error getting transaction count", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get transaction count",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
		})
	}
}