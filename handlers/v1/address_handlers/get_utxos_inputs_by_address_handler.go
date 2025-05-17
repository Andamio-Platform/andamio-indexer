package address_handlers

import (
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// GetUTxOsInputsByAddressHandler handles the GET /api/v1/indexer/addresses/{address}/utxos/inputs endpoint.
// @Summary		Get UTXO Inputs by Address
// @Description	Retrieves transaction inputs for a given address with pagination.
// @ID				getUTxOsInputsByAddress
// @Tags			Addresses
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			address	path		string	true	"The address to retrieve UTXO inputs for."
// @Param			limit	query		int		false	"Maximum number of results to return."	default(100)
// @Param			offset	query		int		false	"Number of results to skip."	default(0)
// @Success		200		{array}		viewmodel.TransactionInput	"Successfully retrieved UTXO inputs."
// @Failure		400		{object}	object{error=string}		"Invalid address or pagination parameters."
// @Failure		404		{object}	object{error=string}		"Address not found or no UTXO inputs found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/indexer/addresses/{address}/utxos/inputs [get]
func GetUTxOsInputsByAddressHandler(db *database.Database) fiber.Handler {
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

		inputs, err := db.Metadata().GetTxInputsByAddress(nil, address, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get UTXO inputs by address"})
		}

		// Convert database models to view models
		inputViewModels := viewmodel.ConvertTransactionInputsToViewModels(inputs)

		return c.JSON(inputViewModels)
	}
}