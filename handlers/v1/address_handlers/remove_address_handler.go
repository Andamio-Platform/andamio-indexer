package address_handlers

import (
	"log/slog"

	database "github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// RemoveAddressHandler handles the request to remove an address from the tracking list.
//
//	@Summary		Remove Address
//	@Description	Removes an address from the tracking list.
//	@ID				removeAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			addressRequest	body		viewmodel.AddressRequest	true	"The address object containing the address string to be removed."
//	@Success		200		{object}	object{message=string}		"Successfully removed address."
//	@Failure		400		{object}	object{error=string}		"Invalid request body or missing address."
//	@Failure		500		{object}	object{error=string}		"Internal server error."
//	@Router			/addresses [delete]
func RemoveAddressHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		addressRequest := new(viewmodel.AddressRequest)
		if err := c.BodyParser(addressRequest); err != nil {
			fiberLogger.Errorf("failed to parse request body: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		if addressRequest.Address == "" {
			fiberLogger.Error("address is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Address is required"})
		}

		// Use a database transaction
		txn := db.Transaction(true)
		defer txn.Discard() // Ensure rollback on error

		err := db.Metadata().RemoveAddress(txn.Metadata(), addressRequest.Address)
		if err != nil {
			txn.Rollback()
			fiberLogger.Errorf("failed to remove address from database: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove address"})
		}

		if err := txn.Commit(); err != nil {
			fiberLogger.Errorf("failed to commit transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove address"})
		}

		fiberLogger.Infof("address removed successfully: %s", addressRequest.Address)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Address removed successfully"})
	}
}
