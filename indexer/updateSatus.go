package indexer

// UpdateSlotTimestamp updates the commit timestamp in the database with the current slot number.
// func UpdateSlotTimestamp(status chainsync.ChainSyncStatus) {
// 	// Start a database transaction
// 	globalDB := database.GetGlobalDB()
// 	txn := globalDB.Transaction(true) // true for writable transaction
// 	defer txn.Discard()               // Ensure transaction is discarded if not committed

// 	// Update the commit timestamp with the current slot number
// 	// Update metadata commit timestamp
// 	if err := globalDB.Metadata().SetCommitTimestamp(txn.Metadata(), int64(status.SlotNumber)); err != nil {
// 		fmt.Errorf("failed to update metadata commit timestamp: %w", err)
// 	}

// 	// Update blob commit timestamp
// 	if err := globalDB.Blob().SetCommitTimestamp(txn.Blob(), int64(status.SlotNumber)); err != nil {
// 		fmt.Errorf("failed to update blob commit timestamp: %w", err)
// 	}

// 	// Commit the transaction
// 	if err := txn.Commit(); err != nil {
// 		fmt.Errorf("failed to commit database transaction: %w", err)
// 	}

// 	slog.Debug(fmt.Sprintf("updated commit timestamp to slot %d", status.SlotNumber))

// 	return
// }
