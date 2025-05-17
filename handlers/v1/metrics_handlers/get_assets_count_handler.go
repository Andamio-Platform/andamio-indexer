package metrics_handlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/gofiber/fiber/v2"
)

// GetAssetsCountHandler godoc
// @Summary Get Total Indexed Assets Count
// @Description Retrieves the total number of indexed assets.
// @ID getAssetsCount
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{count=int} "Successfully retrieved assets count."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/assets/count [get]
func GetAssetsCountHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var count int64
		err := db.Metadata().DB().Model(&models.Asset{}).Count(&count).Error

		if err != nil {
			logger.Error("Error getting asset count", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get asset count",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"count": count,
		})
	}
}