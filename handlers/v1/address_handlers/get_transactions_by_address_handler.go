package address_handlers

import (
	"encoding/hex"
	"strconv"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

// GetTransactionsByAddressHandler handles the request to get transactions for a specific address.
//
//	@Summary		Get Transactions by Address
//	@Description	Retrieves transactions associated with a specific address with pagination.
//	@ID				getTransactionsByAddress
//	@Tags			Addresses
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			address	path		string	true	"The address to retrieve transactions for."
//	@Param			limit	query		int		false	"Maximum number of results to return."	default(100)
//	@Param			offset	query		int		false	"Number of results to skip."	default(0)
//	@Success		200		{array}		viewmodel.Transaction	"Successfully retrieved transactions."
//	@Failure		400		{object}	object{error=string}		"Invalid address or pagination parameters."
//	@Failure		404		{object}	object{error=string}		"Address not found or no transactions found."
//	@Failure		500		{object}	object{error=string}		"Internal server error."
//	@Router			/addresses/{address}/transactions [get]
func GetTransactionsByAddressHandler(db *database.Database) fiber.Handler {
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
		transactions, err := db.GetTxsByAnyAddress(address, limit, offset, nil)
		if err != nil {
			fiberLogger.Errorf("Failed to get transactions for address %s: %v", address, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get transactions"})
		}

		if len(transactions) == 0 {
			fiberLogger.Error("No transactions found for the address")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No transactions found for the address"})
		}

		// Convert database models to view models
		transactionViewModels := []viewmodel.Transaction{}
		for _, tx := range transactions {
			transactionViewModels = append(transactionViewModels, viewmodel.Transaction{
				TransactionHash: hex.EncodeToString(tx.TransactionHash),
				BlockNumber:     tx.BlockNumber,
				SlotNumber:      tx.SlotNumber,
				Inputs:          viewmodel.ConvertTransactionInputsToViewModels(tx.Inputs),
				Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(tx.Outputs),
				Fee:             tx.Fee,
				TTL:             tx.TTL,
				BlockHash:       string(tx.BlockHash),
				Metadata:        hex.EncodeToString(tx.Metadata),
				ReferenceInputs: viewmodel.ConvertSimpleUTxOModelsToViewModels(tx.ReferenceInputs),
				Withdrawals:     tx.Withdrawals,
				Certificates:    viewmodel.ConvertByteSliceSliceToStringSlice(tx.Certificates),
				Witness:         viewmodel.ConvertWitnessModelToViewModel(tx.Witness),
				TransactionCBOR: hex.EncodeToString(tx.TransactionCBOR),
			})
		}

		return c.JSON(transactionViewModels)
	}
}
