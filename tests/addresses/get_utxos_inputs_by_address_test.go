package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetUTxOsInputsByAddress tests retrieving UTxO inputs by address.
func TestGetUTxOsInputsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	// Use an address from config.json that is expected to have indexed UTxO inputs.
	// This address is associated with "globalStateRefMS" in config.json, and we assume
	// the UTxO referenced there has been spent.
	testAddress := "addr_test1xp69xmvurx2uesfydnz9ms7huvzafvlejwfna8rer72hlwuke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzqjte2n7"
	// The expected input corresponds to the spent UTxO reference from config.json "globalStateRefMS" mSCTxRef.
	expectedInputTxHash := "4df3ebc0592b39124c5cc3a1cf680a5d7ac393531dd308e34ee499fbad7257e7"
	expectedInputIndex := uint32(1) // Index is 1 based on "4df3ebc0592b39124c5cc3a1cf680a5d7ac393531dd308e34ee499fbad7257e7#1"

	endpoint := fmt.Sprintf("%s/addresses/%s/utxos/inputs", constants.API_BASE_URL, testAddress)
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
		t.Fatalf("Expected to receive transaction inputs for address %s, but got an empty list", testAddress)
	}

	// Check if the expected transaction input is present in the list
	foundExpectedInput := false
	for _, input := range inputs {
		if input.UTxOID == expectedInputTxHash && input.UTxOIDIndex == expectedInputIndex {
			foundExpectedInput = true
			// Optional: Add more specific validation for this input if needed.
			t.Logf("Found expected transaction input: TxHash=%s, Index=%d, Address=%s",
				input.UTxOID, input.UTxOIDIndex, input.Address)
			break
		}
	}

	if !foundExpectedInput {
		t.Errorf("Expected transaction input with TxHash %s and Index %d not found for address %s",
			expectedInputTxHash, expectedInputIndex, testAddress)
	}

	t.Logf("Successfully retrieved and validated transaction inputs for address %s. Found %d inputs.", testAddress, len(inputs))
}
