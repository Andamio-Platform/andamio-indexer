package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/tests"
)

// Define a struct to match the expected response structure for the latest block endpoint.
// This should align with the actual API response or a corresponding viewmodel if available.
type LatestBlockResponse struct {
	BlockNumber uint64 `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	Slot        uint64 `json:"slot"`
	// Add other fields if the API returns more information
}

// TestGetLatestBlock tests retrieving the latest indexed block information.
func TestGetLatestBlock(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet and has indexed some blocks.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.
	// We assume that by the time this test runs, the indexer has indexed at least one block.

	endpoint := tests.API_BASE_URL + "/metrics/latest-block"
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	// Decode the JSON response
	var latestBlock LatestBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&latestBlock); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	// Basic validation: Check if block number, hash, and slot are non-zero or in a valid format.
	// More specific validation might require knowing expected values from a test dataset.
	if latestBlock.BlockNumber == 0 {
		t.Errorf("Expected non-zero block number, but got %d", latestBlock.BlockNumber)
	}
	if latestBlock.BlockHash == "" {
		t.Errorf("Expected non-empty block hash, but got empty string")
	}
	if latestBlock.Slot == 0 {
		t.Errorf("Expected non-zero slot, but got %d", latestBlock.Slot)
	}

	t.Logf("Successfully retrieved latest block: BlockNumber=%d, BlockHash=%s, Slot=%d",
		latestBlock.BlockNumber, latestBlock.BlockHash, latestBlock.Slot)
}
