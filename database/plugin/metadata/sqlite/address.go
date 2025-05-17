package sqlite

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// AddAddress adds a new address to the database
func (d *MetadataStoreSqlite) AddAddress(txn *gorm.DB, address string) error {
	db := txn
	if db == nil {
		db = d.db
	}
	newAddress := models.Address{
		Address: address,
	}
	result := db.Create(&newAddress)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAddress returns the address from eventCtx.TransactionHash the database
func (d *MetadataStoreSqlite) GetAddress(txn *gorm.DB, address string) (string, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var addr models.Address
	result := db.Where("address = ?", address).First(&addr)
	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}
	return addr.Address, nil
}

// GetAllAddresses returns all addresses from the database
func (d *MetadataStoreSqlite) GetAllAddresses(txn *gorm.DB) ([]string, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var addresses []models.Address
	result := db.Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}

	var addrList []string
	for _, addr := range addresses {
		addrList = append(addrList, addr.Address)
	}

	return addrList, nil
}

// GetTxInputsByAddress retrieves transaction inputs for a given address with pagination support.
func (d *MetadataStoreSqlite) GetTxInputsByAddress(txn *gorm.DB, address string, limit, offset int) ([]models.TransactionInput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var inputs []models.TransactionInput
	query := db.Where("address = ?", []byte(address))

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Find(&inputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return inputs, nil
}

// GetTxOutputsByAddress retrieves transaction outputs for a given address with pagination support.
func (d *MetadataStoreSqlite) GetTxOutputsByAddress(txn *gorm.DB, address string, limit, offset int) ([]models.TransactionOutput, error) {
	db := txn
	if db == nil {
		db = d.db
	}
	var outputs []models.TransactionOutput
	query := db.Where("address = ?", []byte(address))

	if limit > 0 || offset >= 0 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Find(&outputs)
	if result.Error != nil {
		return nil, result.Error
	}
	return outputs, nil
}

// RemoveAddress removes an address from the database
func (d *MetadataStoreSqlite) RemoveAddress(txn *gorm.DB, address string) error {
	db := txn
	if db == nil {
		db = d.db
	}
	result := db.Where("address = ?", address).Delete(&models.Address{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
