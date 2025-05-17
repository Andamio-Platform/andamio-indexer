package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsByAddress tests retrieving UTxOs by address.
func TestGetUTxOsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	// Use an address from config.json that is expected to have indexed UTxOs.
	// This address is associated with "globalStateRefMS" in config.json.
	testAddress := "addr_test1xp69xmvurx2uesfydnz9ms7huvzafvlejwfna8rer72hlwuke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzqjte2n7"
	// The expected UTxO reference from config.json "globalStateRefMS" mSCTxRef.
	expectedTxHash := "4df3ebc0592b39124c5cc3a1cf680a5d7ac393531dd308e34ee499fbad7257e7"
	expectedIndex := uint32(1) // Index is 1 based on "4df3ebc0592b39124c5cc3a1cf680a5d7ac393531dd308e34ee499fbad7257e7#1"

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos", constants.API_BASE_URL, testAddress)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	var utxos []viewmodel.SimpleUTxO
	if err := json.NewDecoder(resp.Body).Decode(&utxos); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	if len(utxos) == 0 {
		t.Fatalf("Expected to receive UTxOs for address %s, but got an empty list", testAddress)
	}

	// Check if the expected UTxO is present in the list
	foundExpectedUTxO := false
	for _, utxo := range utxos {
		if string(utxo.TransactionHash) == expectedTxHash && utxo.UTxOIDIndex == expectedIndex {
			foundExpectedUTxO = true
			// Optional: Add more specific validation for this UTxO if needed,
			// e.g., checking amount, datum hash, or inline datum if known.
			t.Logf("Found expected UTxO: TxHash=%s, Index=%d, Amount=%s, Address=%s",
				utxo.TransactionHash, utxo.UTxOIDIndex)
			break
		}
	}

	if !foundExpectedUTxO {
		t.Errorf("Expected UTxO with TxHash %s and Index %d not found for address %s",
			expectedTxHash, expectedIndex, testAddress)
	}

	t.Logf("Successfully retrieved and validated UTxOs for address %s. Found %d UTxOs.", testAddress, len(utxos))
}
