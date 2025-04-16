package address_handlers

import (
	"encoding/json"

	database "github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetUTxOsByAddressHandler handles the request to get UTxOs for a specific address.
//
//	@Summary		Get UTxOs by Address
//	@Description	Get UTxOs associated with a specific address.
//	@ID			getUTxOsByAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string	true	"Address to get UTxOs for"
//	@Success		200		{object}	[]database.UTxO	"Success response"
//	@Failure		400		{object}	errors.ServerError	"Bad request"
//	@Failure		404		{object}	errors.ServerError	"Address not found"
//	@Failure		500		{object}	errors.ServerError	"Server error"
//	@Router			/addresses/{address}/utxos [get]
func GetUTxOsByAddressHandler(c *fiber.Ctx, db *database.Database) error {
	address := c.Params("address")

	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")

	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("Access-Control-Allow-Methods", "GET")

	if address == "" {
		fiberLogger.Error("Address is required")
		return fiber.NewError(fiber.StatusBadRequest, "Address is required")
	}

	// Assuming you have a database instance available in the Fiber context or a global variable
	// and a function in the database package to fetch UTxOs by address.
	txn := db.Transaction(false)
	defer txn.Discard()
	utxos, err := txn.GetUTxOsByAddress(address)
	if err != nil {
		fiberLogger.Error(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(utxos) == 0 {
		fiberLogger.Error("No UTxOs found for the address")
		return fiber.NewError(fiber.StatusNotFound, "No UTxOs found for the address")
	}

	c.Status(fiber.StatusOK)
	return enc.Encode(utxos)
}
