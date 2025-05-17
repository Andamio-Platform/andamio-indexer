package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/viewmodel"
)

// TestGetTransactionsByAddress tests retrieving transactions by address.
func TestGetTransactionsByAddress(t *testing.T) {
	// Ensure the indexer is running and connected to the testnet.
	// The setupTestEnvironment function in api_test.go handles starting the indexer and waiting for readiness.

	// Use an address from config.json that is expected to have indexed transactions.
	// This address is associated with "globalStateS" in config.json.
	testAddress := "addr_test1xr7xs02kjwr7v3frqrx4exearkd5nmx5ashhzsj5l3nja7yke8x9mpjf7aerjt3n3nfd5tnzkfhlprp09mpf4sdy8dzq6ptcdp"

	endpoint := fmt.Sprintf("%s/addresses/%s/transactions", constants.API_BASE_URL, testAddress)
	resp, err := http.Get(endpoint)
	if err != nil {
		t.Fatalf("Failed to send GET request to %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	// Validate the HTTP status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d for %s, but got %d", http.StatusOK, endpoint, resp.StatusCode)
	}

	// Decode the JSON response into a slice of TransactionViewModel
	var transactions []viewmodel.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		t.Fatalf("Failed to decode JSON response from %s: %v", endpoint, err)
	}

	// Validate the response data
	if len(transactions) == 0 {
		t.Fatalf("Expected to receive transactions for address %s, but got an empty list", testAddress)
	}

	// Validate that each retrieved transaction is associated with the test address
	for _, tx := range transactions {
		isAssociated := false
		// Check inputs
		for _, input := range tx.Inputs {
			if string(input.Address) == testAddress {
				isAssociated = true
				break
			}
		}
		if isAssociated {
			continue
		}
		// Check outputs
		for _, output := range tx.Outputs {
			if string(output.Address) == testAddress {
				isAssociated = true
				break
			}
		}

		if !isAssociated {
			t.Errorf("Retrieved transaction %s is not associated with address %s", tx.TransactionHash, testAddress)
		}

		// Optional: Add more detailed validation of transaction content if needed.
	}

	t.Logf("Successfully retrieved and validated transactions for address %s. Found %d transactions.", testAddress, len(transactions))
}
