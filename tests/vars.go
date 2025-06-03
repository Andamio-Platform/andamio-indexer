package tests

import "github.com/Andamio-Platform/andamio-indexer/config"

const (
	TEST_ADDRESS              string = "addr_test1qpnwndva8xzjqm5djd8fcc23e0a24pnvahpkgptutjvs0xhjjpekkfaj2f2myflapky4ahcgwpqkcmjzta549k6ajmfq7hfan9"
	EXPECTED_TX_HASH          string = "cc2b9a785f8df26e55ca643488e4531f1e74855b0e775b7c9a2f3e6ac2a14d5c"
	EXPECTED_OUTPUT_TX_HASH   string = "cc2b9a785f8df26e55ca643488e4531f1e74855b0e775b7c9a2f3e6ac2a14d5c"
	EXPECTED_OUTPUT_INDEX     uint32 = 6
	TEST_ASSET_FINGERPRINT    string = "asset1j37m0xchc63d8ut3c53mw5ymqy908jamkghs88"
	TEST_ASSET_POLICY_ID      string = "c76c35088ac826c8a0e6947c8ff78d8d4495789bc729419b3a334305"
	TEST_ASSET_TOKEN_NAME     string = "222MIxAxIM"
	TEST_ASSET_TOKEN_NAME_HEX string = "3232324d49784178494d"
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
