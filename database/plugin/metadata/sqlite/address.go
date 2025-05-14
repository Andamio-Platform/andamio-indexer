package sqlite

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"gorm.io/gorm"
)

// AddAddress adds a new address to the database
func (d *MetadataStoreSqlite) AddAddress(address string, txn *gorm.DB) error {
	newAddress := models.Address{
		Address: address,
	}
	result := txn.Create(&newAddress)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAddress returns the address from eventCtx.TransactionHash the database
func (d *MetadataStoreSqlite) GetAddress(address string, txn *gorm.DB) (string, error) {
	var addr models.Address
	result := txn.Where("address = ?", address).First(&addr)
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
	var addresses []models.Address
	result := txn.Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}

	var addrList []string
	for _, addr := range addresses {
		addrList = append(addrList, addr.Address)
	}

	return addrList, nil
}
