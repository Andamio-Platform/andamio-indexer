package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gofiber/fiber/v2"
)

// GetAssetsCountHandler godoc
// @Summary Get Total Unique Assets Count
// @Description Retrieves the total number of unique assets from the database.
// @ID getAssetsCount
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{count=int64} "Successfully retrieved unique assets count."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/assets/count [get]
func GetAssetsCountHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		count, err := db.Metadata().CountUniqueAssets(nil)
		if err != nil {
			logger.Error("Error getting unique asset count", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get unique asset count",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
		})
	}
}