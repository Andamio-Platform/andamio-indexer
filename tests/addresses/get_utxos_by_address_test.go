package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/tests"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsByAddress tests retrieving UTxOs by address.
func TestGetUTxOsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos", tests.API_BASE_URL, tests.TEST_ADDRESS)
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
		t.Fatalf("Expected to receive UTxOs for address %s, but got an empty list", tests.TEST_ADDRESS)
	}

	// Check if the expected UTxO is present in the list
	foundExpectedUTxO := false
	for _, utxo := range utxos {
		if string(utxo.TransactionHash) == tests.EXPECTED_OUTPUT_TX_HASH && utxo.UTxOIDIndex == tests.EXPECTED_OUTPUT_INDEX {
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
			tests.EXPECTED_OUTPUT_TX_HASH, tests.EXPECTED_OUTPUT_INDEX, tests.TEST_ADDRESS)
	}

	t.Logf("Successfully retrieved and validated UTxOs for address %s. Found %d UTxOs.", tests.TEST_ADDRESS, len(utxos))
}
