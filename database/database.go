// Copyright 2025 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"errors"
	"log/slog"
	"os"
	"sync"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/blob"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata"
)

// Database represents our data storage services
type Database struct {
	logger   *slog.Logger
	blob     blob.BlobStore
	metadata metadata.MetadataStore
	dataDir  string
}

var (
	globalDB *Database
	globalMu sync.RWMutex
)

// Blob returns the underling blob store instance
func (d *Database) Blob() blob.BlobStore {
	return d.blob
}

// DataDir returns the path to the data directory used for storage
func (d *Database) DataDir() string {
	return d.dataDir
}

// Logger returns the logger instance
func (d *Database) Logger() *slog.Logger {
	if d.logger != nil {
		return d.logger
	}
	return slog.Default()
}

// Metadata returns the underlying metadata store instance
func (d *Database) Metadata() metadata.MetadataStore {
	return d.metadata
}

// Transaction starts a new database transaction and returns a handle to it
func (d *Database) Transaction(readWrite bool) *Txn {
	slog.Debug("Starting new database transaction.", "readWrite", readWrite)
	return NewTxn(d, readWrite)
}

// BlobTxn starts a new blob-only database transaction and returns a handle to it
func (d *Database) BlobTxn(readWrite bool) *Txn {
	slog.Debug("Starting new blob-only transaction.", "readWrite", readWrite)
	return NewBlobOnlyTxn(d, readWrite)
}

// MetadataTxn starts a new metadata-only database transaction and returns a handle to it
func (d *Database) MetadataTxn(readWrite bool) *Txn {
	slog.Debug("Starting new metadata-only transaction.", "readWrite", readWrite)
	return NewMetadataOnlyTxn(d, readWrite)
}

// Close cleans up the database connections
func (d *Database) Close() error {
	slog.Info("Closing database connections...")
	var err error
	// Close metadata
	slog.Info("Closing metadata store...")
	metadataErr := d.Metadata().Close()
	err = errors.Join(err, metadataErr)
	if metadataErr == nil {
		slog.Info("Metadata store closed successfully.")
	} else {
		slog.Error("Failed to close metadata store.", "error", metadataErr)
	}
	// Close blob
	slog.Info("Closing blob store...")
	blobErr := d.Blob().Close()
	err = errors.Join(err, blobErr)
	if blobErr == nil {
		slog.Info("Blob store closed successfully.")
	} else {
		slog.Error("Failed to close blob store.", "error", blobErr)
	}
	slog.Info("Database connections closed.")
	return err
}

func (d *Database) init() error {
	slog.Info("Initializing database...")
	if d.logger == nil {
		// Create logger to throw away logs
		// We do this so we don't have to add guards around every log operation
		d.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	// Check commit timestamp
	slog.Info("Checking commit timestamp...")
	if err := d.checkCommitTimestamp(); err != nil {
		slog.Error("Commit timestamp check failed.", "error", err)
		return err
	}
	slog.Info("Commit timestamp check successful.")
	slog.Info("Database initialized.")
	return nil
}

// New creates a new database instance with optional persistence using the provided data directory
func New(
	logger *slog.Logger,
	dataDir string,
) (*Database, error) {
	slog.Info("Creating new database instance...", "dataDir", dataDir)
	slog.Info("Initializing metadata store (sqlite)...")
	metadataDb, err := metadata.New("sqlite", dataDir, logger)
	if err != nil {
		slog.Error("Failed to initialize metadata store.", "error", err)
		return nil, err
	}
	slog.Info("Metadata store initialized successfully.")

	slog.Info("Initializing blob store (badger)...")
	blobDb, err := blob.New("badger", dataDir, logger)
	if err != nil {
		slog.Error("Failed to initialize blob store.", "error", err)
		return nil, err
	}
	slog.Info("Blob store initialized successfully.")

	db := &Database{
		logger:   logger,
		blob:     blobDb,
		metadata: metadataDb,
		dataDir:  dataDir,
	}
	if err := db.init(); err != nil {
		// Database is available for recovery, so return it with error
		slog.Error("Database initialization failed.", "error", err)
		return db, err
	}
	slog.Info("New database instance created successfully.")
	return db, nil
}

// SetGlobalDB sets the global database instance
func SetGlobalDB(db *Database) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalDB = db
}

// GetUniqueAddressesCount retrieves the total count of unique addresses from all transactions, excluding a list of addresses.
func (d *Database) GetUniqueAddressesCount(excludedAddresses []string) (int64, error) {
	txn := d.MetadataTxn(false)
	defer txn.Rollback()
	return d.Metadata().GetUniqueAddressesCount(txn.Metadata(), excludedAddresses)
}

// GetTotalTransactionFees retrieves the total sum of all transaction fees.
func (d *Database) GetTotalTransactionFees() (uint64, error) {
	txn := d.MetadataTxn(false)
	defer txn.Rollback()
	return d.Metadata().GetTotalTransactionFees(txn.Metadata())
}


// GetGlobalDB returns the global database instance
func GetGlobalDB() *Database {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDB
}
