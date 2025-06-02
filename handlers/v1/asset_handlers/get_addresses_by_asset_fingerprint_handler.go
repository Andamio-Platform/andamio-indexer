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
// @Param asset_fingerprint path string true "The asset fingerprint to retrieve addresses for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {array} string "Successfully retrieved addresses."
// @Failure 400 {object} object{error=string} "Invalid asset fingerprint or pagination parameters."
// @Failure 404 {object} object{error=string} "Asset fingerprint not found or no addresses found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /assets/fingerprint/{asset_fingerprint}/addresses [get]
func GetAddressesByAssetFingerprintHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assetFingerprint := c.Params("asset_fingerprint")
		logger.Info("Received asset_fingerprint", "asset_fingerprint", assetFingerprint)
		if assetFingerprint == "" {
			logger.Error("asset_fingerprint path parameter is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "asset_fingerprint path parameter is missing"})
		}

		limit := c.QueryInt("limit", 100)
		offset := c.QueryInt("offset", 0)

		if limit < 0 || offset < 0 {
			logger.Error("invalid pagination parameters", "limit", limit, "offset", offset)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid pagination parameters"})
		}

		var assets []models.Asset
		fingerprintBytes := []byte(assetFingerprint)
		logger.Info("Querying with fingerprint bytes", "fingerprint_bytes", fingerprintBytes)
		result := db.Metadata().DB().Model(&models.Asset{}).
			Where("fingerprint = ?", fingerprintBytes).
			Limit(limit).Offset(offset).
			Find(&assets)

		logger.Info("Database query result", "rows_affected", result.RowsAffected, "error", result.Error)
		if result.Error != nil {
			logger.Error("failed to get assets by fingerprint", "asset_fingerprint", assetFingerprint, "error", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve addresses"})
		}

		if len(assets) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no assets found for this asset fingerprint"})
		}

		// Extract unique addresses from the transaction outputs associated with the assets
		addressMap := make(map[string]bool)
		var addresses []string
		for _, asset := range assets {
			logger.Debug("Attempting to retrieve transaction output for asset", "asset_utxo_id", hex.EncodeToString(asset.UTxOID), "asset_utxo_index", asset.UTxOIDIndex)
			var address string
			output, err := db.Metadata().GetTxOutputByUTxO(nil, asset.UTxOID, asset.UTxOIDIndex)
			if err != nil {
				logger.Error("failed to get transaction output for asset", "utxo_id", hex.EncodeToString(asset.UTxOID), "utxo_index", asset.UTxOIDIndex, "error", err)
				continue
			}
			if output != nil {
				logger.Debug("Successfully retrieved transaction output", "output_id", output.ID, "output_address", string(output.Address))
				address = string(output.Address)
			} else {
				logger.Debug("Transaction output not found for asset, checking transaction inputs", "asset_utxo_id", hex.EncodeToString(asset.UTxOID), "asset_utxo_index", asset.UTxOIDIndex)
				input, err := db.Metadata().GetTxInputByUTxO(nil, asset.UTxOID, asset.UTxOIDIndex)
				if err != nil {
					logger.Error("failed to get transaction input for asset", "utxo_id", hex.EncodeToString(asset.UTxOID), "utxo_index", asset.UTxOIDIndex, "error", err)
					continue
				}
				if input != nil {
					logger.Debug("Successfully retrieved transaction input", "input_id", input.ID, "input_address", string(input.Address))
					address = string(input.Address)
				} else {
					logger.Debug("Transaction input not found for asset", "asset_utxo_id", hex.EncodeToString(asset.UTxOID), "asset_utxo_index", asset.UTxOIDIndex)
					continue
				}
			}

			if address != "" {
				logger.Debug("Found address for asset", "address", address)
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
