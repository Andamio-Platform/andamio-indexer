package datum_handlers

import (
	"encoding/hex"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// GetTransactionsByDatumHashHandler godoc
// @Summary		Get Transactions by Datum Hash
// @Description	Retrieves a list of transactions that reference a specific datum hash.
// @ID				getTransactionsByDatumHash
// @Tags			Datums
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			datum_hash	path		string	true	"The datum hash (hex-encoded) to retrieve transactions for."
// @Success		200		{array}		viewmodel.Transaction	"Successfully retrieved transactions."
// @Failure		400		{object}	object{error=string}		"Invalid datum hash."
// @Failure		404		{object}	object{error=string}		"Datum hash not found or no transactions found."
// @Failure		500		{object}	object{error=string}		"Internal server error."
// @Router			/datums/{datum_hash}/transactions [get]
func GetTransactionsByDatumHashHandler(db *database.Database, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		datumHashHex := c.Params("datum_hash")
		if datumHashHex == "" {
			logger.Error("datum_hash path parameter is missing")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "datum_hash path parameter is missing"})
		}

		datumHash, err := hex.DecodeString(datumHashHex)
		if err != nil {
			logger.Error("invalid datum_hash hex encoding", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid datum_hash hex encoding"})
		}

		datum, err := db.Metadata().GetDatumByHash(nil, datumHash)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "datum not found"})
			}
			logger.Error("failed to get datum by hash", "datum_hash", datumHashHex, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transactions"})
		}

		if datum == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "datum not found"})
		}

		output, err := db.Metadata().GetTxOutputByUTxO(nil, datum.UTxOID, datum.UTxOIDIndex)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction output for datum not found"})
			}
			logger.Error("failed to get transaction output for datum", "datum_hash", datumHashHex, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transactions"})
		}

		if output == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction output for datum not found"})
		}

		tx, err := db.Metadata().GetTxByTxHash(nil, output.UTxOID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction for datum not found"})
			}
			logger.Error("failed to get transaction by hash", "tx_hash", output.UTxOID, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve transactions"})
		}

		if tx == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction for datum not found"})
		}

		// Convert database model to view model and return in a slice
		transactionViewModels := []viewmodel.Transaction{
			{
				TransactionHash: string(tx.TransactionHash),
				BlockNumber:     tx.BlockNumber,
				SlotNumber:      tx.SlotNumber,
				Inputs:          viewmodel.ConvertTransactionInputsToViewModels(tx.Inputs),
				Outputs:         viewmodel.ConvertTransactionOutputsToViewModels(tx.Outputs),
				Fee:             tx.Fee,
				TTL:             tx.TTL,
				BlockHash:       string(tx.BlockHash),
				Metadata:        string(tx.Metadata),
				ReferenceInputs: viewmodel.ConvertSimpleUTxOModelsToViewModels(tx.ReferenceInputs),
				Withdrawals:     tx.Withdrawals,
				Certificates:    viewmodel.ConvertByteSliceSliceToStringSlice(tx.Certificates),
				Witness:         viewmodel.ConvertWitnessModelToViewModel(tx.Witness),
			},
		}

		return c.Status(fiber.StatusOK).JSON(transactionViewModels)
	}
}
