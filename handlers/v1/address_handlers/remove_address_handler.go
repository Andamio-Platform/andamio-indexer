package address_handlers

import (
	"encoding/json"
	"io"

	"github.com/Andamio-Platform/andamio-indexer/errors"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// RemoveAddressHandler handles the request to remove an address from the tracking list.
//
// This endpoint allows clients to submit an address to be removed from the tracking list.
//
//	@Summary		Remove Address
//	@Description	Remove an address from the tracking list.
//	@ID				removeAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//
// header
//
//	@Param			addressRequest	body		viewmodel.AddressRequest	true	"Address to remove"
//
// body
//
//	@Success		200			{object}	viewmodel.TxResponse				"Success response"
//
// in auth
//
//	@Failure		408			{object}	errors.ServerError					"Request Timeout"
//	@Failure		401			{object}	errors.ServerError					"Access Denied"
//	@Failure		498			{object}	errors.ServerError					"Invalid Token"
//
// in request
//
//	@Failure		400			{function}	errors.GeneralError					"Bad request or empty body"
//	@Failure		400			{function}	errors.FieldError					"Validation failed for required fields"
//
// to request
//
//	@Router			/addresses [delete]
func RemoveAddressHandler(c *fiber.Ctx) error {
	var addressRequest viewmodel.AddressRequest

	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetIndent("", "    ")

	c.Response().Header.Set("Content-Type", "application/json")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")
	c.Response().Header.Set("Access-Control-Allow-Methods", "DELETE")

	err := c.BodyParser(&addressRequest)

	if err != nil {
		fiberLogger.Error(err)
		return errors.BadRequestErrorHandler(c, err)
	}

	switch {
	//The EOF check if body is empty
	case err == io.EOF:
		c.Status(fiber.StatusBadRequest)
		err := enc.Encode(errors.GeneralError("Nothing Was Sent"))
		if err != nil {
			fiberLogger.Error(err)
			return errors.ServerErrorHandler(c, err)
		}
		return err
	}

	// Check if Necessary fields are empty or not,
	// if it is empty it'll send an error message
	if addressRequest.IsValid() != nil {
		c.Status(fiber.StatusBadRequest)
		err := enc.Encode(errors.FieldError("Address"))
		if err != nil {
			fiberLogger.Error(err)
			return errors.ServerErrorHandler(c, err)
		}
		return err
	}

	// TODO: Implement the logic to remove the address from the tracking list.
	// This may involve modifying the configuration or using a database.
	// For now, just return a success message.
	// Remove address from tracking list here

	response := map[string]string{"message": "Address removed successfully"}
	c.Status(fiber.StatusOK)
	return enc.Encode(response)
}
