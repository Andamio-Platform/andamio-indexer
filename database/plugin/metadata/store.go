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

package metadata

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

type MetadataStore interface {
	// Database
	Close() error
	DB() *gorm.DB
	GetCommitTimestamp(txn *gorm.DB) (int64, error)
	SetCommitTimestamp(*gorm.DB, int64) error
	Transaction() *gorm.DB
	AutoMigrate(txn *gorm.DB, dst ...interface{}) error
	Create(txn *gorm.DB, value interface{}) *gorm.DB
	First(txn *gorm.DB, args interface{}) *gorm.DB
	Order(txn *gorm.DB, args interface{}) *gorm.DB
	Where(txn *gorm.DB, query interface{}, args ...interface{}) *gorm.DB

	// Address
	AddAddress(txn *gorm.DB, address string) error
	GetAddress(txn *gorm.DB, address string) (string, error)
	GetAllAddresses(txn *gorm.DB) ([]string, error)
	RemoveAddress(txn *gorm.DB, address string) error

	// Transaction
	SetTx(txn *gorm.DB, tx *models.Transaction) error
	GetTxByTxHash(txn *gorm.DB, txHash []byte) (*models.Transaction, error)
	GetTxsByBlockNumber(txn *gorm.DB, blockNumber uint64, limit, offset int) ([]models.Transaction, error)
	GetTxsBySlotRange(txn *gorm.DB, startSlot, endSlot uint64, limit, offset int) ([]models.Transaction, error)
	GetTxsByInputAddress(txn *gorm.DB, address string, limit, offset int) ([]models.TransactionInput, error)
	GetTxsByOutputAddress(txn *gorm.DB, address string, limit, offset int) ([]models.Transaction, error)
	GetTxsByAnyAddress(txn *gorm.DB, address string, limit, offset int) ([]models.Transaction, error)
	SetTxs(txn *gorm.DB, txs []*models.Transaction) error
	GetTxs(txn *gorm.DB, limit, offset int) ([]models.Transaction, error)
	CountTxs(txn *gorm.DB) (int64, error)
	DeleteTxByHash(txn *gorm.DB, txHash []byte) error
	DeleteTxsByBlockNumber(txn *gorm.DB, blockNumber uint64) error
	GetUniqueAddressesCount(txn *gorm.DB, excludedAddresses []string) (int64, error)
	GetTotalTransactionFees(txn *gorm.DB) (uint64, error)
	GetTxInputByUTxO(txn *gorm.DB, arg1 []byte, arg2 uint32) (*models.TransactionInput, error)
	GetTxOutputByUTxO(txn *gorm.DB, arg1 []byte, arg2 uint32) (*models.TransactionOutput, error)
	GetTxByID(txn *gorm.DB, arg1 uint) (*models.Transaction, error)
	GetTxInputByID(txn *gorm.DB, id uint) (*models.TransactionInput, error)
	GetTxOutputByID(txn *gorm.DB, id uint) (*models.TransactionOutput, error)

	// Added Getters and Setters for other types
	GetAssets(txn *gorm.DB, utxoID []byte, utxoIndex uint32) ([]models.Asset, error)
	SetAsset(txn *gorm.DB, asset *models.Asset) error
	CountUniqueAssets(txn *gorm.DB) (int64, error)

	GetDatum(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.Datum, error)
	GetDatumByHash(txn *gorm.DB, arg1 []byte) (*models.Datum, error)
	SetDatum(txn *gorm.DB, datum *models.Datum) error

	SetRedeemer(txn *gorm.DB, redeemer *models.Redeemer) error
	GetRedeemersByWitnessId(txn *gorm.DB, witnessID uint) ([]models.Redeemer, error)
	GetRedeemersByWitnessIdAndIndexAndTag(txn *gorm.DB, witnessID uint, index uint, tag []byte) (*models.Redeemer, error)
	GetRedeemersByWitnessIdAndTag(txn *gorm.DB, witnessID uint, tag []byte) ([]models.Redeemer, error)

	SetSimpleUTxO(txn *gorm.DB, utxo *models.SimpleUTxO) error
	GetSimpleUTxOByUTxO(txn *gorm.DB, utxoID []byte, utxoIndex uint32) (*models.SimpleUTxO, error)
	GetSimpleUTxOByID(txn *gorm.DB, arg1 []byte) ([]models.SimpleUTxO, error)
	GetSimpleUTxOsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.SimpleUTxO, error)
	GetSimpleUTxOByPrimaryKey(txn *gorm.DB, arg1 uint) (*models.SimpleUTxO, error)

	SetTransactionInput(txn *gorm.DB, input *models.TransactionInput) error
	GetTransactionInputsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.TransactionInput, error)
	GetTransactionInputsByAddress(txn *gorm.DB, address []byte) ([]models.TransactionInput, error)

	SetTransactionOutput(txn *gorm.DB, output *models.TransactionOutput) error
	GetTransactionOutputsByTransactionHash(txn *gorm.DB, transactionHash []byte) ([]models.TransactionOutput, error)
	GetTransactionOutputsByAddress(txn *gorm.DB, address []byte) ([]models.TransactionOutput, error)

	SetWitness(txn *gorm.DB, witness *models.Witness) error
	GetWitnessByTransactionHash(txn *gorm.DB, transactionHash []byte) (*models.Witness, error)
	GetWitnessByID(txn *gorm.DB, arg1 uint) (*models.Witness, error)

	// New functions for API endpoints
	GetTxInputsByAddress(txn *gorm.DB, address string, limit, offset int) ([]models.TransactionInput, error)
	GetTxOutputsByAddress(txn *gorm.DB, address string, limit, offset int) ([]models.TransactionOutput, error)
	GetTxsByPolicyId(txn *gorm.DB, policyId []byte, limit, offset int) ([]models.Transaction, error)
	GetTxsByTokenName(txn *gorm.DB, tokenName []byte, limit, offset int) ([]models.Transaction, error)
	GetTxsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.Transaction, error)
	GetTxsByPolicyIdAndTokenName(txn *gorm.DB, policyId []byte, tokenName []byte, limit, offset int) ([]models.Transaction, error)
	GetUTxOsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.SimpleUTxO, error)
	GetTransactionInputsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.TransactionInput, error)
	GetTransactionOutputsByAssetFingerprint(txn *gorm.DB, assetFingerprint []byte, limit, offset int) ([]models.TransactionOutput, error)
}

// For now, this always returns a sqlite plugin
func New(
	pluginName, dataDir string,
	logger *slog.Logger,
) (MetadataStore, error) {
	return sqlite.New(dataDir, logger)
}
