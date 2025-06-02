package tests

import "github.com/Andamio-Platform/andamio-indexer/config"

const (
	TEST_ADDRESS            string = "addr_test1qpnwndva8xzjqm5djd8fcc23e0a24pnvahpkgptutjvs0xhjjpekkfaj2f2myflapky4ahcgwpqkcmjzta549k6ajmfq7hfan9"
	EXPECTED_TX_HASH        string = "4df3ebc0592b39124c5cc3a1cf680a5d7ac393531dd308e34ee499fbad7257e7"
	EXPECTED_INDEX          uint32 = 1
	EXPECTED_OUTPUT_TX_HASH string = "8a3a9c393bec05d40b73ed459a10a5c9c7a11f197c88d1aaca48080a2e48e7c5"
	EXPECTED_OUTPUT_INDEX   uint32 = 0
)

var (
	API_BASE_URL string
)

func init() {
	err := config.Load("../../config/config.json")
	if err != nil {
		panic(err)
	}
	API_BASE_URL = config.GetGlobalConfig().Indexer.APIBaseURL
}
