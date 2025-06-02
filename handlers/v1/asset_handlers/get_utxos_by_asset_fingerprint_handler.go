package asset_handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// GetUTxOsByAssetFingerprintHandler godoc
// @Summary Get UTxOs by Asset Fingerprint
// @Description Retrieve a list of UTxOs containing a specific asset fingerprint, with support for pagination.
// @ID getUTxOsByAssetFingerprint
// @Tags Assets
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param asset_fingerprint path string true "The asset fingerprint to retrieve UTxOs for."
// @Param limit query int false "Maximum number of results to return." default(100)
// @Param offset query int false "Number of results to skip." default(0)
// @Success 200 {object} viewmodel.TransactionUTxOs "Successfully retrieved UTxOs."
// @Failure 400 {object} object{error=string} "Invalid asset fingerprint or pagination parameters."
// @Failure 404 {object} object{error=string} "Asset fingerprint not found or no UTxOs found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /assets/fingerprint/{asset_fingerprint}/utxos [get]
func GetUTxOsByAssetFingerprintHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		assetFingerprint := c.Params("asset_fingerprint")
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

		inputs, err := db.Metadata().GetTransactionInputsByAssetFingerprint(nil, []byte(assetFingerprint), limit, offset)
		if err != nil {
			logger.Error("failed to get transaction inputs by asset fingerprint", "asset_fingerprint", assetFingerprint, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transaction inputs"})
		}

		outputs, err := db.Metadata().GetTransactionOutputsByAssetFingerprint(nil, []byte(assetFingerprint), limit, offset)
		if err != nil {
			logger.Error("failed to get transaction outputs by asset fingerprint", "asset_fingerprint", assetFingerprint, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transaction outputs"})
		}

		if len(inputs) == 0 && len(outputs) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no UTxOs (inputs or outputs) found for the given asset fingerprint"})
		}

		// Convert models to view models
		inputViewModels := viewmodel.ConvertTransactionInputsToViewModels(inputs)
		outputViewModels := viewmodel.ConvertTransactionOutputsToViewModels(outputs)

		transactionUTxOs := viewmodel.TransactionUTxOs{
			Inputs:  inputViewModels,
			Outputs: outputViewModels,
		}

		return c.Status(fiber.StatusOK).JSON(transactionUTxOs)
	}
}
