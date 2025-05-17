package eventHandlers

import (
	"log/slog"
	"sync"

	"github.com/Andamio-Platform/andamio-indexer/indexer/cache" // Import the cache package
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
)

// AddToTransactionBatch adds a transaction event to the cache
func AddToTransactionBatch(eventTx input_chainsync.TransactionEvent, eventCtx input_chainsync.TransactionContext) {
	txCache := cache.GetTransactionCache()
	if txCache != nil {
		txCache.Add(eventTx, eventCtx)
		slog.Info("Added transaction to batch cache.", "txHash", eventTx.Transaction.Hash(), "currentBatchSize", txCache.Len(), "batchLimit", txCache.Limit())
		// Check if cache limit is reached and process the batch
		if txCache.Len() >= txCache.Limit() {
			slog.Info("Transaction batch limit reached, processing batch.")
			go ProcessTransactionBatch() // Process in a goroutine to avoid blocking
		}
	} else {
		slog.Error("Transaction cache not initialized when trying to add transaction.")
	}
}

// ProcessTransactionBatch processes the cached transactions
func ProcessTransactionBatch() {
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

	var wg sync.WaitGroup
	for _, item := range transactionsToProcess {
		wg.Add(1)
		go func(txEvent input_chainsync.TransactionEvent, txContext input_chainsync.TransactionContext) {
			defer wg.Done()
			slog.Debug("Processing individual transaction in batch.", "txHash", string(txEvent.Transaction.Hash().Bytes()))
			err := TxEvent(txEvent, txContext) // Call the TxEvent function to process the transaction
			if err != nil {
				slog.Error("Error processing transaction in batch", "txHash", string(txEvent.Transaction.Hash().Bytes()), "error", err)
			} else {
				slog.Debug("Finished processing individual transaction in batch.", "txHash", string(txEvent.Transaction.Hash().Bytes()))
			}
		}(item.Event, item.Context)
	}

	wg.Wait() // Wait for all goroutines to complete

	slog.Info("Finished processing transaction batch")
}
