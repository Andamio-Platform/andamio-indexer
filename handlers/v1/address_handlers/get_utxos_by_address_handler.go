package address_handlers

import (
	"strconv"

	database "github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetUTxOsByAddressHandler handles the request to get UTxOs for a specific address.
//
//	@Summary		Get UTxOs by Address
//	@Description	Retrieves UTxOs associated with a specific address with pagination.
//	@ID				getUTxOsByAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string	true	"The address to retrieve UTxOs for."
//	@Param			limit	query		int		false	"Maximum number of results to return."	default(100)
//	@Param			offset	query		int		false	"Number of results to skip."	default(0)
//	@Success		200		{array}		viewmodel.Transaction	"Successfully retrieved UTxOs."
//	@Failure		400		{object}	object{error=string}		"Invalid address or pagination parameters."
//	@Failure		404		{object}	object{error=string}		"Address not found or no UTxOs found."
//	@Failure		500		{object}	object{error=string}		"Internal server error."
//	@Router			/addresses/{address}/utxos [get]
func GetUTxOsByAddressHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		address := c.Params("address")

		if address == "" {
			fiberLogger.Error("Address is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Address is required"})
		}

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

		// Use the database function with pagination
		utxos, err := db.Metadata().GetTxsByAnyAddress(nil, address, limit, offset)
		if err != nil {
			fiberLogger.Errorf("Failed to get UTxOs for address %s: %v", address, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get UTxOs"})
		}

		if len(utxos) == 0 {
			fiberLogger.Error("No UTxOs found for the address")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No UTxOs found for the address"})
		}

		// Convert database models to view models
		utxoViewModels := []viewmodel.Transaction{}
		for _, utxo := range utxos {
			utxoViewModels = append(utxoViewModels, viewmodel.Transaction{
				TransactionHash: string(utxo.TransactionHash),
				BlockNumber:     utxo.BlockNumber,
				SlotNumber:      utxo.SlotNumber,
				Inputs:          viewmodel.ConvertTransactionInputsToViewModels(utxo.Inputs),
				Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(utxo.Outputs),
				Fee:             utxo.Fee,
				TTL:             utxo.TTL,
				// Add other fields as needed
			})
		}

		return c.JSON(utxoViewModels)
	}
}
