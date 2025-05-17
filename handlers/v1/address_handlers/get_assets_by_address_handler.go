package address_handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
)

// GetAssetsByAddressHandler godoc
// @Summary Get Assets by Address
// @Description Retrieve a list of assets held at a specific address, with support for pagination.
// @ID getAssetsByAddress
// @Tags Addresses
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param address path string true "The address to retrieve assets for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} viewmodel.Asset "Successfully retrieved assets."
// @Failure 400 {object} object{error=string} "Invalid address or pagination parameters."
// @Failure 404 {object} object{error=string} "Address not found or no assets found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /addresses/{address}/assets [get]
func GetAssetsByAddressHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		address := c.Params("address")
		if address == "" {
			logger.Error("address path parameter is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "address path parameter is missing"})
		}

		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		if limit < 0 || offset < 0 {
			logger.Error("invalid pagination parameters", "limit", limit, "offset", offset)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pagination parameters"})
		}

		outputs, err := db.Metadata().GetTxOutputsByAddress(nil, address, limit, offset)
		if err != nil {
			logger.Error("failed to get transaction outputs by address", "address", address, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve assets"})
		}

		if len(outputs) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no assets found for this address"})
		}

		var assets []models.Asset
		for _, output := range outputs {
			outputAssets, err := db.Metadata().GetAssets(nil, output.UTxOID, output.UTxOIDIndex)
			if err != nil {
				logger.Error("failed to get assets for output", "utxo_id", output.UTxOID, "utxo_index", output.UTxOIDIndex, "error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve assets"})
			}
			assets = append(assets, outputAssets...)
		}

		if len(assets) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no assets found for this address"})
		}

		return c.Status(fiber.StatusOK).JSON(assets)
	}
}
