package tests

import (
	// Import encoding/json
	"fmt" // Import io
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Andamio-Platform/andamio-indexer/constants"
	// Assuming viewmodel package exists and contains TransactionViewModel
	// Add other necessary imports (e.g., for HTTP client, database interaction)
)

var (
	indexerCmd *exec.Cmd
	testDBPath string
)

func TestMain(m *testing.M) {
	// Setup logging for tests
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo, // Adjust log level as needed for test verbosity
	})))

	// Setup: Start indexer process, wait for readiness
	setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Teardown: Stop indexer process, clean up database
	teardownTestEnvironment()

	os.Exit(code)
}

func setupTestEnvironment() {
	slog.Info("Setting up test environment...")

	// Define the path to the built indexer executable
	// Assuming the indexer binary is built in the root directory
	indexerBinaryPath := "../andamio-indexer" // Adjust path as needed

	// Define the path to the test configuration file
	// Create a test-specific config file that uses a temporary database.
	testConfigPath := "../config/test_config.json" // Adjust path as needed

	// Create a temporary database directory
	testDBPath = "./testdb_" + time.Now().Format("20060102150405")
	if err := os.MkdirAll(testDBPath, 0755); err != nil {
		slog.Error("Failed to create test database directory", "error", err)
		os.Exit(1)
	}
	// TODO: Ensure the test config points to this testDBPath - This requires modifying test_config.json or passing the path as an argument if the indexer supports it. For now, assuming test_config.json is set up to use a relative path or will be modified separately.

	// Start the indexer process
	slog.Info("Starting indexer process...")
	// Pass the test database path as an environment variable or command-line argument if the indexer supports it.
	// For now, assuming the config file handles the database path.
	indexerCmd = exec.Command(indexerBinaryPath, "--config", testConfigPath)
	// Redirect stdout/stderr for visibility during testing
	indexerCmd.Stdout = os.Stdout
	indexerCmd.Stderr = os.Stderr

	if err := indexerCmd.Start(); err != nil {
		slog.Error("Failed to start indexer process", "error", err)
		teardownTestEnvironment() // Clean up before exiting
		os.Exit(1)
	}
	slog.Info("Indexer process started.")

	// Wait for the indexer API to be ready
	slog.Info("Waiting for indexer API to become ready...")
	// Use a known, lightweight endpoint for the readiness check.
	readinessEndpoint := constants.API_BASE_URL + "/metrics/latest-block" // Adjust if a dedicated health check exists
	if err := waitForAPIReady(readinessEndpoint); err != nil {
		slog.Error("Indexer API did not become ready", "error", err)
		teardownTestEnvironment() // Clean up before exiting
		os.Exit(1)
	}
	slog.Info("Indexer API is ready.")

	slog.Info("Test environment setup complete.")
}

func teardownTestEnvironment() {
	slog.Info("Tearing down test environment...")

	// Stop the indexer process
	if indexerCmd != nil && indexerCmd.Process != nil {
		slog.Info("Stopping indexer process...")
		// Attempt graceful shutdown first (e.g., SIGINT)
		if err := indexerCmd.Process.Signal(os.Interrupt); err != nil {
			slog.Error("Failed to send interrupt signal to indexer process", "error", err)
			// If graceful shutdown fails, force kill after a timeout
			time.Sleep(5 * time.Second) // Give it some time to shut down
			if err := indexerCmd.Process.Kill(); err != nil {
				slog.Error("Failed to kill indexer process", "error", err)
			}
		}
		// Wait for the process to exit
		indexerCmd.Wait()
		slog.Info("Indexer process stopped.")
	}

	// Clean up the temporary database directory
	if testDBPath != "" {
		slog.Info("Cleaning up test database directory...")
		if err := os.RemoveAll(testDBPath); err != nil {
			slog.Error("Error cleaning up test database directory", "error", err)
		}
		slog.Info("Test database directory cleaned up.")
	}

	slog.Info("Test environment teardown complete.")
}

// waitForAPIReady polls a given URL until it returns a 200 OK status or a timeout occurs.
func waitForAPIReady(url string) error {
	timeout := 60 * time.Second // Adjust timeout as needed
	interval := 1 * time.Second // Adjust polling interval
	startTime := time.Now()

	for time.Since(startTime) < timeout {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil // API is ready
		}
		if resp != nil {
			resp.Body.Close()
		}
		slog.Info("Waiting for API...", "url", url, "elapsed", time.Since(startTime).Round(time.Second))
		time.Sleep(interval)
	}

	return fmt.Errorf("timeout waiting for API to become ready at %s", url)
}

// TestAPIReachability is a basic test to ensure the indexer API is reachable.
func TestAPIReachability(t *testing.T) {
	// This test is now covered by the waitForAPIReady in setupTestEnvironment,
	// but can be kept as a simple check if needed.
	t.Skip("Covered by setupTestEnvironment readiness check")
}
