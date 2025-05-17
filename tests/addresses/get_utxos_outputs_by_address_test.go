package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
	// Import the tests package using a relative path
)

// TestGetUTxOsOutputsByAddress tests retrieving UTxO outputs by address.
func TestGetUTxOsOutputsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	// Use an address from config.json that is expected to have indexed UTxO outputs.
	// This address is associated with "globalStateS" in config.json.
	testAddress := "addr_test1xr7xs02kjwr7v3frqrx4exearkd5nmx5ashhzsj5l3nja7yke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzq6ptcdp"
	// The expected UTxO output reference from config.json "globalStateS" sCTxRef.
	expectedOutputTxHash := "8a3a9c393bec05d40b73ed459a10a5c9c7a11f197c88d1aaca48080a2e48e7c5"
	expectedOutputIndex := uint32(0) // Index is 0 based on "8a3a9c393bec05d40b73ed459a10a5c9c7a11f197c88d1aaca48080a2e48e7c5#0"

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/outputs", constants.API_BASE_URL, testAddress)
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
		t.Fatalf("Expected to receive transaction outputs for address %s, but got an empty list", testAddress)
	}

	// Check if the expected transaction output is present in the list
	foundExpectedOutput := false
	for _, output := range outputs {
		if output.UTxOID == expectedOutputTxHash && output.UTxOIDIndex == expectedOutputIndex {
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
			expectedOutputTxHash, expectedOutputIndex, testAddress)
	}

	t.Logf("Successfully retrieved and validated transaction outputs for address %s. Found %d outputs.", testAddress, len(outputs))
}
