package asset_handlers

import (
	"encoding/hex"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/Andamio-Platform/andamio-indexer/database"
)

// GetUTxOsByAssetFingerprintHandler godoc
// @Summary Get UTxOs by Asset Fingerprint
// @Description Retrieve a list of UTxOs containing a specific asset fingerprint, with support for pagination.
// @ID getUTxOsByAssetFingerprint
// @Tags Assets
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param asset_fingerprint path string true "The asset fingerprint (hex-encoded) to retrieve UTxOs for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} viewmodel.SimpleUTxO "Successfully retrieved UTxOs."
// @Failure 400 {object} object{error=string} "Invalid asset fingerprint or pagination parameters."
// @Failure 404 {object} object{error=string} "Asset fingerprint not found or no UTxOs found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /assets/{asset_fingerprint}/utxos [get]
func GetUTxOsByAssetFingerprintHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assetFingerprintHex := c.Params("asset_fingerprint")
		if assetFingerprintHex == "" {
			logger.Error("asset_fingerprint path parameter is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "asset_fingerprint path parameter is missing"})
		}

		assetFingerprint, err := hex.DecodeString(assetFingerprintHex)
		if err != nil {
			logger.Error("invalid asset_fingerprint hex encoding", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid asset_fingerprint hex encoding"})
		}

		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		if limit < 0 || offset < 0 {
			logger.Error("invalid pagination parameters", "limit", limit, "offset", offset)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pagination parameters"})
		}

		utxos, err := db.Metadata().GetUTxOsByAssetFingerprint(nil, assetFingerprint, limit, offset)
		if err != nil {
			logger.Error("failed to get UTxOs by asset fingerprint", "asset_fingerprint", assetFingerprintHex, "error", err)
			// Handle specific errors like not found
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve UTxOs"})
		}

		if len(utxos) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no UTxOs found for the given asset fingerprint"})
		}

		return c.Status(fiber.StatusOK).JSON(utxos)
	}
}
