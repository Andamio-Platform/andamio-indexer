package metrics_handlers

import (
	"fmt"
	"net/http"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gofiber/fiber/v2"
)

// GetTotalTransactionFeesHandler godoc
// @Summary Get Total Transaction Fees
// @Description Retrieves the total sum of all transaction fees across the entire dataset.
// @ID getTotalTransactionFees
// @Tags Metrics
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} object{total_transaction_fees=string} "Successfully retrieved total transaction fees."
// @Failure 500 {object} object{error=string} "Internal server error."
// @Router /metrics/total_transaction_fees [get]
func GetTotalTransactionFeesHandler(db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		totalFees, err := db.GetTotalTransactionFees()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Convert uint64 to float64 and format to 6 decimal places
		// Assuming the fee is in lovelace (1 ADA = 1,000,000 lovelace)
		totalFeesDecimal := float64(totalFees) / 1_000_000.0
		formattedFees := fmt.Sprintf("%.6f", totalFeesDecimal)

		return c.Status(http.StatusOK).JSON(fiber.Map{"total_transaction_fees": formattedFees})
	}
}
