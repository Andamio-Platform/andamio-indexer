package indexer

import (
	"github.com/Andamio-Platform/andamio-indexer/indexer/filters"
)

func InitAdder() error {
	err := filters.FilterByAndamioAddresses()
	if err != nil {
		return err
	}

	err = filters.FilterByAndamioPolicyIDs()
	if err != nil {
		return err
	}

	return nil
}
