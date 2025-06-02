package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/tests"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsWitnessesByAddress tests retrieving UTxO witnesses by address.
func TestGetUTxOsWitnessesByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/witnesses", tests.API_BASE_URL, tests.TEST_ADDRESS)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	// Decode the JSON response into a slice of TransactionWitnessViewModel
	var witnesses []viewmodel.Witness
	if err := json.NewDecoder(resp.Body).Decode(&witnesses); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	if len(witnesses) == 0 {
		t.Fatalf("Expected to receive transaction witnesses for address %s, but got an empty list", tests.TEST_ADDRESS)
	}

	t.Logf("Successfully retrieved and validated transaction witnesses for address %s. Found %d witnesses.", tests.TEST_ADDRESS, len(witnesses))
}
