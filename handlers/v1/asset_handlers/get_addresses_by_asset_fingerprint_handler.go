package asset_handlers

import (
	"encoding/hex"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
)

// GetAddressesByAssetFingerprintHandler godoc
// @Summary Get Addresses by Asset Fingerprint
// @Description Retrieve a list of addresses that hold a specific asset fingerprint, with support for pagination.
// @ID getAddressesByAssetFingerprint
// @Tags Assets
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param asset_fingerprint path string true "The asset fingerprint (hex-encoded) to retrieve addresses for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} string "Successfully retrieved addresses."
// @Failure 400 {object} object{error=string} "Invalid asset fingerprint or pagination parameters."
// @Failure 404 {object} object{error=string} "Asset fingerprint not found or no addresses found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /assets/{asset_fingerprint}/addresses [get]
func GetAddressesByAssetFingerprintHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
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

		var assets []models.Asset
		result := db.Metadata().DB().Model(&models.Asset{}).
			Where("fingerprint = ?", assetFingerprint).
			Limit(limit).Offset(offset).
			Find(&assets)

		if result.Error != nil {
			logger.Error("failed to get assets by fingerprint", "asset_fingerprint", assetFingerprintHex, "error", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve addresses"})
		}

		if len(assets) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no assets found for this asset fingerprint"})
		}

		// Extract unique addresses from the transaction outputs associated with the assets
		addressMap := make(map[string]bool)
		var addresses []string
		for _, asset := range assets {
			output, err := db.Metadata().GetTxOutputByUTxO(nil, asset.UTxOID, asset.UTxOIDIndex)
			if err != nil {
				// Log the error but continue processing other assets
				logger.Error("failed to get transaction output for asset", "utxo_id", asset.UTxOID, "utxo_index", asset.UTxOIDIndex, "error", err)
				continue
			}
			if output != nil {
				address := string(output.Address)
				if _, ok := addressMap[address]; !ok {
					addressMap[address] = true
					addresses = append(addresses, address)
				}
			}
		}

		if len(addresses) == 0 {
			// This case might be hit if all outputs for the found assets could not be retrieved
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no addresses found for this asset fingerprint"})
		}

		return c.Status(fiber.StatusOK).JSON(addresses)
	}
}
