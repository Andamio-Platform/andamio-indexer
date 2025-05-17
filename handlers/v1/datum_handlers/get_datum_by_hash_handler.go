package datum_handlers

import (
	"encoding/hex"
	"log/slog"
	"net/http"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetDatumByHashHandler godoc
// @Summary Get Datum by Hash
// @Description Retrieve details of a specific datum by its hash.
// @ID getDatumByHash
// @Tags Datums
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param datum_hash path string true "The datum hash (hex-encoded) to retrieve."
// @Success 200 {object} viewmodel.Datum "Successfully retrieved datum."
// @Failure 400 {object} object{error=string} "Invalid datum hash."
// @Failure 404 {object} object{error=string} "Datum not found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /datums/{datum_hash} [get]
func GetDatumByHashHandler(db *database.Database, log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		datumHashHex := c.Params("datum_hash")
		if datumHashHex == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing datum_hash path parameter"})
		}

		_, err := hex.DecodeString(datumHashHex) // Decode to validate, but don't store as it's unused in the placeholder
		if err != nil {
			log.Error("failed to decode datum hash", "error", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid datum hash"})
		}

		datumHash, err := hex.DecodeString(datumHashHex)
		if err != nil {
			log.Error("failed to decode datum hash", "error", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid datum hash"})
		}

		datum, err := db.Metadata().GetDatumByHash(nil, datumHash)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "datum not found"})
			}
			log.Error("failed to get datum by hash", "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		// Convert database model to view model
		datumViewModel := viewmodel.Datum{
			DatumHash: string(datum.DatumHash),
			DatumCbor: string(datum.DatumCbor), // Assuming DatumCbor should be a string representation of CBOR
		}

		return c.Status(http.StatusOK).JSON(datumViewModel)
	}
}
