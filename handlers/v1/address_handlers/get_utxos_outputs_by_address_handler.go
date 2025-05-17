package address_handlers

import (
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetUTxOsOutputsByAddressHandler handles the GET /api/v1/indexer/addresses/{address}/utxos/outputs endpoint.
// @Summary		Get UTXO Outputs by Address
// @Description Retrieves transaction outputs for a given address with pagination.
// @Tags 		Addresses
// @Accept 		json
// @Produce json
// @Param address path string true "Address to retrieve UTXO outputs for"
// @Param limit query int false "Maximum number of results to return" default(100)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {array} viewmodel.TransactionOutput "Success response"
// @Failure 400 {object} errors.ServerError "Bad request"
// @Failure 500 {object} errors.ServerError "Server error"
// @Router /indexer/addresses/{address}/utxos/outputs [get]
func GetUTxOsOutputsByAddressHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		address := c.Params("address")

		// Get pagination parameters
		limitStr := c.Query("limit", "100")
		offsetStr := c.Query("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit parameter"})
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid offset parameter"})
		}

		outputs, err := db.Metadata().GetTxOutputsByAddress(nil, address, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get UTXO outputs by address"})
		}

		// Convert database models to view models
		outputViewModels := viewmodel.ConvertTransactionOutputsToViewModels(outputs)

		return c.JSON(outputViewModels)
	}
}
