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

package sqlite

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Register plugin
func init() {
	plugin.Register(
		plugin.PluginEntry{
			Type: plugin.PluginTypeMetadata,
			Name: "sqlite",
		},
	)
}

// MetadataStoreSqlite stores all data in sqlite. Data may not be persisted
type MetadataStoreSqlite struct {
	dataDir     string
	db          *gorm.DB
	logger      *slog.Logger
	timerVacuum *time.Timer
}

// New creates a new database
func New(
	dataDir string,
	logger *slog.Logger,
) (*MetadataStoreSqlite, error) {
	var metadataDb *gorm.DB
	var err error
	if dataDir == "" {
		// No dataDir, use in-memory config
		metadataDb, err = gorm.Open(
			sqlite.Open("file::memory:?cache=shared"),
			&gorm.Config{
				Logger:                 gormlogger.Discard,
				SkipDefaultTransaction: true,
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		// Make sure that we can read data dir, and create if it doesn't exist
		if _, err := os.Stat(dataDir); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("failed to read data dir: %w", err)
			}
			// Create data directory
			if err := os.MkdirAll(dataDir, fs.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create data dir: %w", err)
			}
		}
		// Open sqlite DB
		metadataDbPath := filepath.Join(
			dataDir,
			"metadata.sqlite",
		)
		// WAL journal mode, disable sync on write, increase cache size to 50MB (from 2MB)
		metadataConnOpts := "_pragma=journal_mode(WAL)&_pragma=sync(OFF)&_pragma=cache_size(-50000)"
		metadataDb, err = gorm.Open(
			sqlite.Open(
				fmt.Sprintf("file:%s?%s", metadataDbPath, metadataConnOpts),
			),
			&gorm.Config{
				Logger:                 gormlogger.Discard,
				SkipDefaultTransaction: true,
			},
		)
		if err != nil {
			return nil, err
		}
	}
	db := &MetadataStoreSqlite{
		db:      metadataDb,
		dataDir: dataDir,
		logger:  logger,
	}
	if err := db.init(); err != nil {
		// MetadataStoreSqlite is available for recovery, so return it with error
		return db, err
	}
	// Create table schemas
	db.logger.Debug(fmt.Sprintf("creating table: %#v", &CommitTimestamp{}))
	if err := db.db.AutoMigrate(&CommitTimestamp{}); err != nil {
		return db, err
	}
	for _, model := range models.MigrateModels {
		db.logger.Debug(fmt.Sprintf("creating table: %#v", model))
		if err := db.db.AutoMigrate(model); err != nil {
			return db, err
		}
	}
	return db, nil
}

func (d *MetadataStoreSqlite) init() error {
	if d.logger == nil {
		// Create logger to throw away logs
		// We do this so we don't have to add guards around every log operation
		d.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	// Configure tracing for GORM
	if err := d.db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return err
	}
	// Schedule daily database vacuum to free unused space
	d.scheduleDailyVacuum()
	return nil
}

func (d *MetadataStoreSqlite) runVacuum() error {
	if d.dataDir == "" {
		return nil
	}
	if result := d.DB().Raw("VACUUM"); result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *MetadataStoreSqlite) scheduleDailyVacuum() {
	if d.timerVacuum != nil {
		d.timerVacuum.Stop()
	}
	daily := time.Duration(24) * time.Hour
	f := func() {
		d.logger.Debug(
			"running vacuum on sqlite metadata database",
		)
		// schedule next run
		defer d.scheduleDailyVacuum()
		if err := d.runVacuum(); err != nil {
			d.logger.Error(
				"failed to free unused space in metadata store",
				"component", "database",
			)
		}
	}
	d.timerVacuum = time.AfterFunc(daily, f)
}

// AutoMigrate wraps the gorm AutoMigrate
func (d *MetadataStoreSqlite) AutoMigrate(txn *gorm.DB, dst ...interface{}) error {
	db := txn
	if db == nil {
		db = d.DB()
	}
	return db.AutoMigrate(dst...)
}

// Close gets the database handle from our MetadataStore and closes it
func (d *MetadataStoreSqlite) Close() error {
	// get DB handle from gorm.DB
	db, err := d.DB().DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// Create creates a record
func (d *MetadataStoreSqlite) Create(txn *gorm.DB, value interface{}) *gorm.DB {
	db := txn
	if db == nil {
		db = d.DB()
	}
	return db.Create(value)
}

// DB returns the database handle
func (d *MetadataStoreSqlite) DB() *gorm.DB {
	return d.db
}

// First returns the first DB entry
func (d *MetadataStoreSqlite) First(txn *gorm.DB, args interface{}) *gorm.DB {
	db := txn
	if db == nil {
		db = d.DB()
	}
	return db.First(args)
}

// Order orders a DB query
func (d *MetadataStoreSqlite) Order(txn *gorm.DB, args interface{}) *gorm.DB {
	db := txn
	if db == nil {
		db = d.DB()
	}
	return db.Order(args)
}

// Transaction creates a gorm transaction
func (d *MetadataStoreSqlite) Transaction() *gorm.DB {
	return d.DB().Begin()
}

// Where constrains a DB query
func (d *MetadataStoreSqlite) Where(
	txn *gorm.DB,
	query interface{},
	args ...interface{},
) *gorm.DB {
	db := txn
	if db == nil {
		db = d.DB()
	}
	return db.Where(query, args...)
}

// GetUniqueAddressesCount retrieves the total count of unique addresses from all transaction inputs and outputs, excluding a list of addresses.
func (d *MetadataStoreSqlite) GetUniqueAddressesCount(
	txn *gorm.DB,
	excludedAddresses []string,
) (int64, error) {
	db := txn
	if db == nil {
		db = d.DB()
	}

	uniqueAddresses := make(map[string]struct{})

	// Get addresses from TransactionInputs
	var inputAddresses []string
	if err := db.Model(&models.TransactionInput{}).Distinct("address").Pluck("address", &inputAddresses).Error; err != nil {
		return 0, err
	}
	for _, addr := range inputAddresses {
		uniqueAddresses[addr] = struct{}{}
	}

	// Get addresses from TransactionOutputs
	var outputAddresses []string
	if err := db.Model(&models.TransactionOutput{}).Distinct("address").Pluck("address", &outputAddresses).Error; err != nil {
		return 0, err
	}
	for _, addr := range outputAddresses {
		uniqueAddresses[addr] = struct{}{}
	}

	// Apply exclusion filter
	filteredAddressesCount := 0
	for addr := range uniqueAddresses {
		isExcluded := false
		for _, excludedAddr := range excludedAddresses {
			if addr == excludedAddr {
				isExcluded = true
				break
			}
		}
		if !isExcluded {
			filteredAddressesCount++
		}
	}

	return int64(filteredAddressesCount), nil
}

// GetTotalTransactionFees retrieves the total sum of all transaction fees.
func (d *MetadataStoreSqlite) GetTotalTransactionFees(txn *gorm.DB) (uint64, error) {
	db := txn
	if db == nil {
		db = d.DB()
	}

	var totalFees uint64
	result := db.Model(&models.Transaction{}).Select("SUM(fee)").Row().Scan(&totalFees)
	if result != nil {
		return 0, result
	}

	return totalFees, nil
}
