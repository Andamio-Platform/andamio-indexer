package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsWitnessesByAddress tests retrieving UTxO witnesses by address.
func TestGetUTxOsWitnessesByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	// Use an address from config.json that is expected to have indexed UTxO witnesses.
	testAddress := "addr_test1xp69xmvurx2uesfydnz9ms7huvzafvlejwfna8rer72hlwuke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzqjte2n7"

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/witnesses", constants.API_BASE_URL, testAddress)
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
		t.Fatalf("Expected to receive transaction witnesses for address %s, but got an empty list", testAddress)
	}

	t.Logf("Successfully retrieved and validated transaction witnesses for address %s. Found %d witnesses.", testAddress, len(witnesses))
}
