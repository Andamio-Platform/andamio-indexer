package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/tests"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsInputsByAddress tests retrieving UTxO inputs by address.
func TestGetUTxOsInputsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/inputs", tests.API_BASE_URL, tests.TEST_ADDRESS)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	// Decode the JSON response into a slice of TransactionInput
	var inputs []viewmodel.TransactionInput
	if err := json.NewDecoder(resp.Body).Decode(&inputs); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	if len(inputs) == 0 {
		t.Fatalf("Expected to receive transaction inputs for address %s, but got an empty list", tests.TEST_ADDRESS)
	}

	// Check if the expected transaction input is present in the list
	foundExpectedInput := false
	for _, input := range inputs {
		if input.UTxOID == tests.EXPECTED_OUTPUT_TX_HASH && input.UTxOIDIndex == tests.EXPECTED_OUTPUT_INDEX {
			foundExpectedInput = true
			// Optional: Add more specific validation for this input if needed.
			t.Logf("Found expected transaction input: TxHash=%s, Index=%d, Address=%s",
				input.UTxOID, input.UTxOIDIndex, input.Address)
			break
		}
	}

	if !foundExpectedInput {
		t.Errorf("Expected transaction input with TxHash %s and Index %d not found for address %s",
			tests.EXPECTED_OUTPUT_TX_HASH, tests.EXPECTED_OUTPUT_INDEX, tests.TEST_ADDRESS)
	}

	t.Logf("Successfully retrieved and validated transaction inputs for address %s. Found %d inputs.", tests.TEST_ADDRESS, len(inputs))
}
