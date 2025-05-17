package witness_handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetWitnessByIdHandler godoc
// @Summary Get Witness by ID
// @Description Retrieve details of a specific witness by its ID.
// @ID getWitnessById
// @Tags Witnesses
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param witness_id path string true "The witness ID to retrieve."
// @Success 200 {object} viewmodel.Witness "Successfully retrieved witness."
// @Failure 400 {object} object{error=string} "Invalid witness ID."
// @Failure 404 {object} object{error=string} "Witness not found."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /witnesses/{witness_id} [get]
func GetWitnessByIdHandler(db *database.Database, log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		witnessIDParam := c.Params("witness_id")
		if witnessIDParam == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing witness_id path parameter"})
		}

		_, err := strconv.ParseUint(witnessIDParam, 10, 64) // Parse to validate, but don't store as it's unused in the placeholder
		if err != nil {
			log.Error("failed to parse witness ID", "error", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid witness ID"})
		}

		witnessID, err := strconv.ParseUint(witnessIDParam, 10, 64)
		if err != nil {
			log.Error("failed to parse witness ID", "error", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid witness ID"})
		}

		witness, err := db.Metadata().GetWitnessByID(nil, uint(witnessID))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "witness not found"})
			}
			log.Error("failed to get witness by ID", "witness_id", witnessID, "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		if witness == nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "witness not found"})
		}

		// Convert database model to view model
		witnessViewModel := viewmodel.Witness{
			TransactionHash: string(witness.TransactionHash),
			PlutusData:      viewmodel.ConvertByteSliceSliceToStringSlice(witness.PlutusData),      // Convert [][]byte to []string
			PlutusV1Scripts: viewmodel.ConvertByteSliceSliceToStringSlice(witness.PlutusV1Scripts), // Convert [][]byte to []string
			PlutusV2Scripts: viewmodel.ConvertByteSliceSliceToStringSlice(witness.PlutusV2Scripts), // Convert [][]byte to []string
			PlutusV3Scripts: viewmodel.ConvertByteSliceSliceToStringSlice(witness.PlutusV3Scripts), // Convert [][]byte to []string
			Redeemers:       viewmodel.ConvertRedeemersToViewModels(witness.Redeemers),
		}

		return c.Status(http.StatusOK).JSON(witnessViewModel)
	}
}