package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gofiber/fiber/v2"
)

// GetAddressesCountHandler godoc
// @Summary Get Total Unique Addresses Count
// @Description Retrieves the total number of unique addresses from all transactions in the database.
// @ID getAddressesCount
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{count=int64} "Successfully retrieved unique addresses count."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/addresses/count [get]
func GetAddressesCountHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := db.GetUniqueAddressesCount(nil)
		if err != nil {
			logger.Error("Error getting unique addresses count", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get unique addresses count",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
		})
	}
}