package eventHandlers

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"      // Import the database package
	"github.com/Andamio-Platform/andamio-indexer/indexer/cache" // Import the cache package
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
)

// AddToTransactionBatch adds a transaction event to the cache
func AddToTransactionBatch(db *database.Database, eventTx input_chainsync.TransactionEvent, eventCtx input_chainsync.TransactionContext) {
	txCache := cache.GetTransactionCache()
	if txCache != nil {
		txCache.Add(eventTx, eventCtx)
		slog.Info("Added transaction to batch cache.", "txHash", eventTx.Transaction.Hash(), "currentBatchSize", txCache.Len(), "batchLimit", txCache.Limit())
		// Check if cache limit is reached and process the batch
		if txCache.Len() >= txCache.Limit() {
			slog.Info("Transaction batch limit reached, processing batch.")
			go ProcessTransactionBatch(db) // Process in a goroutine to avoid blocking, passing the database instance
		}
	} else {
		slog.Error("Transaction cache not initialized when trying to add transaction.")
	}
}

// ProcessTransactionBatch processes the cached transactions
func ProcessTransactionBatch(db *database.Database) {
	txCache := cache.GetTransactionCache()
	if txCache == nil {
		slog.Error("Transaction cache not initialized")
		return
	}

	slog.Info("Getting transactions from cache for batch processing.")
	// Get all transactions from the cache and clear it
	transactionsToProcess := txCache.GetAll()
	slog.Info("Retrieved transactions from cache.", "count", len(transactionsToProcess))

	if len(transactionsToProcess) == 0 {
		slog.Info("No transactions to process in the batch")
		return
	}

	slog.Info("Processing transaction batch", "count", len(transactionsToProcess))

	// Start a single transaction for the batch
	txn := db.Transaction(true)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic occurred during batch processing, rolling back transaction.", "panic", r)
			if err := txn.Rollback(); err != nil {
				slog.Error("Failed to rollback transaction after panic.", "error", err)
			}
			panic(r) // Re-panic after rollback
		}
	}()

	var batchErr error
	for _, item := range transactionsToProcess {
		slog.Debug("Processing individual transaction in batch.", "txHash", item.Event.Transaction.Hash())
		// Pass the transaction to the TxEvent function
		err := TxEvent(db.Logger(), item.Event, item.Context, txn)
		if err != nil {
			slog.Error("Error processing transaction in batch", "txHash", string(item.Event.Transaction.Hash().Bytes()), "error", err)
			batchErr = err // Store the first error encountered
			// Continue processing other transactions in the batch to potentially log more errors,
			// but the transaction will be rolled back due to batchErr being set.
		} else {
			slog.Debug("Finished processing individual transaction in batch.", "txHash", item.Event.Transaction.Hash())
		}
	}

	// Commit or rollback the transaction based on whether an error occurred
	if batchErr != nil {
		slog.Error("Rolling back transaction due to batch processing error.")
		if err := txn.Rollback(); err != nil {
			slog.Error("Failed to rollback transaction.", "error", err)
			return // Return the original batchErr
		}
		return // Return the original batchErr
	} else {
		slog.Info("Committing transaction batch.")
		if err := txn.Commit(); err != nil {
			slog.Error("Failed to commit transaction.", "error", err)
			return // Return the commit error
		}
		slog.Info("Finished processing transaction batch.")
	}
}
