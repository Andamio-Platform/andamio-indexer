package database

// GetAllAddresses returns all addresses from the database
func (d *Database) GetAllAddresses() ([]string, error) {
	txn := d.Transaction(false)
	defer txn.Discard()

	addrList, err := d.metadata.GetAllAddresses(txn.Metadata())
	if err != nil {
		return nil, err
	}

	return addrList, nil
}
