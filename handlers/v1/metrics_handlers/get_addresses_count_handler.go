package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gofiber/fiber/v2"
)

// GetAddressesCountHandler godoc
// @Summary Get Total Monitored Addresses Count
// @Description Retrieves the total number of monitored addresses.
// @ID getAddressesCount
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{count=int} "Successfully retrieved addresses count."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/addresses/count [get]
func GetAddressesCountHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		addresses, err := db.Metadata().GetAllAddresses(nil)
		if err != nil {
			logger.Error("Error getting addresses count", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get addresses count",
			})
		}

		count := len(addresses)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
		})
	}
}