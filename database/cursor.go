package database

import (
	"encoding/json"
	"fmt"

	badgerdb "github.com/dgraph-io/badger/v4" // Import badgerdb with an alias
)

// CursorState holds the last processed slot number and block hash.
// This struct is used to track the indexer's progress.
type CursorState struct {
	SlotNumber uint64 `json:"slot_number"`
	BlockHash  []byte `json:"block_hash"` // Assuming block hash is a byte slice
}

// TODO: Clarify the exact data types for SlotNumber and BlockHash.
// uint64 and []byte are common, but confirm if different types are needed.

// CursorStore is responsible for managing the cursor state using the main Database instance.
type CursorStore struct {
	db *Database // Reference to the main Database struct
	cursorKey []byte
}

// NewCursorStore creates a new instance of CursorStore.
// It takes a Database instance as an argument.
func NewCursorStore(db *Database) *CursorStore {
	return &CursorStore{
		db:        db,
		cursorKey: []byte("indexer_cursor"), // Unique key for the cursor state
	}
}

// UpdateCursor stores the current slot number and block hash in BadgerDB using a database transaction.
// It serializes the CursorState struct before storing.
func (cs *CursorStore) UpdateCursor(slotNumber uint64, blockHash []byte) error {
	state := CursorState{
		SlotNumber: slotNumber,
		BlockHash:  blockHash,
	}

	// Serialize the cursor state to a byte slice.
	// TODO: Consider alternative serialization methods (e.g., Protocol Buffers, MessagePack)
	// for potentially better performance or smaller storage size, especially if the
	// cursor state grows in complexity. JSON is simple for now.
	serializedState, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to serialize cursor state: %w", err)
	}

	// Use a read-write blob-only transaction from the main database.
	txn := cs.db.BlobTxn(true) // true for read-write transaction
	defer txn.Rollback() // Ensure transaction is rolled back if not committed

	// Set the serialized state using the blob transaction from the main transaction.
	err = txn.Blob().Set(cs.cursorKey, serializedState)
	if err != nil {
		return fmt.Errorf("failed to set cursor state in BadgerDB transaction: %w", err)
	}

	// Commit the transaction.
	if err := txn.Commit(); err != nil {
		return fmt.Errorf("failed to commit cursor state transaction: %w", err)
	}

	return nil
}

// GetCursor retrieves the last stored slot number and block hash from BadgerDB using a database transaction.
// It deserializes the stored data into a CursorState struct.
// If no cursor state is found, it returns a zero-valued CursorState and nil error.
func (cs *CursorStore) GetCursor() (CursorState, error) {
	var state CursorState

	// Use a read-only blob-only transaction from the main database.
	txn := cs.db.BlobTxn(false) // false for read-only transaction
	defer txn.Discard() // Discard the read-only transaction

	item, err := txn.Blob().Get(cs.cursorKey)
	if err != nil {
		// If the key is not found, it means no cursor state has been saved yet.
		if err == badgerdb.ErrKeyNotFound {
			return CursorState{}, nil // Return zero-valued state and nil error
		}
		return CursorState{}, fmt.Errorf("failed to get cursor state from BadgerDB transaction: %w", err)
	}

	// Retrieve the value and deserialize it.
	err = item.Value(func(val []byte) error {
		// TODO: Ensure the deserialization method matches the serialization method used in UpdateCursor.
		return json.Unmarshal(val, &state)
	})
	if err != nil {
		return CursorState{}, fmt.Errorf("failed to deserialize cursor state: %w", err)
	}

	return state, nil // Return the retrieved state
}
