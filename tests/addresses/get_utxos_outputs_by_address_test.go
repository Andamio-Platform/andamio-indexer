package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/tests"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	// Import the tests package using a relative path
)

// TestGetUTxOsOutputsByAddress tests retrieving UTxO outputs by address.
func TestGetUTxOsOutputsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/outputs", tests.API_BASE_URL, tests.TEST_ADDRESS)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	// Decode the JSON response into a slice of TransactionOutput
	var outputs []viewmodel.TransactionOutput
	if err := json.NewDecoder(resp.Body).Decode(&outputs); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	if len(outputs) == 0 {
		t.Fatalf("Expected to receive transaction outputs for address %s, but got an empty list", tests.TEST_ADDRESS)
	}

	// Check if the expected transaction output is present in the list
	foundExpectedOutput := false
	for _, output := range outputs {
		if output.UTxOID == tests.EXPECTED_OUTPUT_TX_HASH && output.UTxOIDIndex == tests.EXPECTED_OUTPUT_INDEX {
			foundExpectedOutput = true
			// Optional: Add more specific validation for this output if needed,
			// e.g., checking amount, datum hash, or inline datum if known.
			t.Logf("Found expected transaction output: TxHash=%s, Index=%d, Amount=%d, Address=%s, DatumHash=%s, DatumCbor=%s",
				output.UTxOID, output.UTxOIDIndex, output.Amount, output.Address, output.Datum.DatumHash, output.Datum.DatumCbor)
			break
		}
	}

	if !foundExpectedOutput {
		t.Errorf("Expected transaction output with TxHash %s and Index %d not found for address %s",
			tests.EXPECTED_OUTPUT_TX_HASH, tests.EXPECTED_OUTPUT_INDEX, tests.TEST_ADDRESS)
	}

	t.Logf("Successfully retrieved and validated transaction outputs for address %s. Found %d outputs.", tests.TEST_ADDRESS, len(outputs))
}
