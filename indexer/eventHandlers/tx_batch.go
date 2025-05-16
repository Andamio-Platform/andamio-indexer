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
		// Check if cache limit is reached and process the batch
		if txCache.Len() >= txCache.Limit() {
			go ProcessTransactionBatch() // Process in a goroutine to avoid blocking
		}
	}
}

// ProcessTransactionBatch processes the cached transactions
func ProcessTransactionBatch() {
	txCache := cache.GetTransactionCache()
	if txCache == nil {
		slog.Error("Transaction cache not initialized")
		return
	}

	// Get all transactions from the cache and clear it
	transactionsToProcess := txCache.GetAll()

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
			err := TxEvent(txEvent, txContext) // Call the TxEvent function to process the transaction
			if err != nil {
				slog.Error("Error processing transaction", "txHash", string(txEvent.Transaction.Hash().Bytes()), "error", err)
			}
		}(item.Event, item.Context)
	}

	wg.Wait() // Wait for all goroutines to complete

	slog.Info("Finished processing transaction batch")
}