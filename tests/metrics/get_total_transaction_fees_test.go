package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/router"
	"github.com/gofiber/fiber/v2"
)

// TestGetTotalTransactionFees tests retrieving the total sum of all transaction fees.
func TestGetTotalTransactionFees(t *testing.T) {
	t.Skip("Implement total transaction fees retrieval and validation")

	// TODO: Initialize a test database with some transaction data
	// For now, we'll use a mock database or a simple in-memory setup if possible
	// Or, ensure the indexer is running and connected to the testnet.

	// Setup Fiber app
	app := fiber.New()
	mockDB := &database.Database{}      // You'll need a proper mock or test DB here
	router.RouterInit(app, mockDB, nil) // Pass nil for logger in tests

	// Create a new HTTP request to the endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/v1/metrics/total_transaction_fees", nil)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := app.Test(req, -1) // -1 for no timeout
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Parse the response body
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	// Assert the presence of the key and its format
	assert.Contains(t, result, "total_transaction_fees")
	// You might want to add more specific validation for the format (e.g., regex for decimal with 6 precision)
	// For example: assert.Regexp(t, `^\d+\.\d{6}$`, result["total_transaction_fees"])
}
