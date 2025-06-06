package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	PayloadExpireAtMinute time.Duration = 15
)

var (
	GlobalConfig *Config
)

type Config struct {
	Network  Network  `json:"network"`
	Indexer  Indexer  `json:"indexer"`
	Database Database `json:"database"`
	Andamio  Andamio  `json:"andamio"`
}

type Network struct {
	Magic                      uint32 `json:"magic"`
	LocalCardanoNodeSocketPath string `json:"localCardanoNodesocketPath"`
	LocalCardanoNodeEndpoint   string `json:"localCardanoNodeEndpoint"`
	LocalKupoEndpoint          string `json:"localKupoEndpoint"`
	BlinklabKupoEndpoint       string `json:"blinklabKupoEndpoint"`
	CFCardanoNodeEndpoint      string `json:"CFCardanoNodeEndpoint"`
}

// "intercerptHash": "cd510710d2d680240540595aea3306750ad275e38ab4511eb10d2b5e02cc0186",
// "interceptSlot": 82402528,
// "intercerptHash": "6c7a5d8036284a4c1af79b62bd6702ac4bb23d33d2c09596272091688e986bcb",
// "interceptSlot": 88200000,

type Indexer struct {
	Host                  string `json:"host"`
	APIBaseURL            string `json:"APIBaseURL"`
	SwaggerURL            string `json:"swaggerURL"`
	IntercerptHash        string `json:"intercerptHash"`
	InterceptSlot         uint64 `json:"interceptSlot"`
	TrancactionCacheLimit int    `json:"trancactionCacheLimit"`
}

type Database struct {
	DatabaseDIR string `json:"databaseDir"`
}

type Andamio struct {
	GlobalAdmin           string                 `json:"globalAdmin"`
	GlobalStateRefMS      MintingContractConfig  `json:"globalStateRefMS"`
	GlobalStateS          SpendingContractConfig `json:"globalStateS"`
	GovernanceS           SpendingContractConfig `json:"governanceS"`
	IndexAdmin            string                 `json:"indexAdmin"`
	IndexMS               MintingContractConfig  `json:"indexMS"`
	IndexRefMS            MintingContractConfig  `json:"indexRefMS"`
	InstanceAdmin         string                 `json:"instanceAdmin"`
	InstanceMS            MintingContractConfig  `json:"instanceMS"`
	InstanceProvidedMS    MintingContractConfig  `json:"instanceProvidedMS"`
	InstanceProviderAdmin string                 `json:"instanceProviderAdmin"`
	ReferenceAddr         string                 `json:"referenceAddr"`
	StakingAdmin          string                 `json:"stakingAdmin"`
	StakingSH             string                 `json:"stakingSH"`
	V1GlobalStateObsTxRef string                 `json:"v1GlobalStateObsTxRef"`
}

type MintingContractConfig struct {
	MSCAddress  string `json:"mSCAddress"`
	MSCPolicyID string `json:"mSCPolicyID"`
	MSCTxRef    string `json:"mSCTxRef"`
}

type SpendingContractConfig struct {
	SCAddress string `json:"sCAddress"`
	SCTxRef   string `json:"sCTxRef"`
}

func Load(configFile string) error {

	if configFile != "" {
		buf, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %s", err)
		}
		err = json.Unmarshal(buf, &GlobalConfig)
		if err != nil {
			return fmt.Errorf("error parsing config file: %v", err)
		}
	}

	return nil
}

func GetGlobalConfig() *Config {
	return GlobalConfig
}

func (a *Andamio) GetAllAndamioPolicies() []string {
	var andamioPolicies []string
	andamioPolicies = append(andamioPolicies, a.GlobalStateRefMS.MSCPolicyID)
	andamioPolicies = append(andamioPolicies, a.IndexMS.MSCPolicyID)
	andamioPolicies = append(andamioPolicies, a.InstanceMS.MSCPolicyID) // every utxo with this token is trusted
	andamioPolicies = append(andamioPolicies, a.IndexRefMS.MSCPolicyID)
	andamioPolicies = append(andamioPolicies, a.InstanceProvidedMS.MSCPolicyID)
	return andamioPolicies
}

func (a *Andamio) GetAllAndamioAssetFingerprints() []string {
	var assetsFingersList []string
	assetsFingersList = append(assetsFingersList, a.GlobalAdmin)
	assetsFingersList = append(assetsFingersList, a.IndexAdmin)
	assetsFingersList = append(assetsFingersList, a.InstanceAdmin)
	assetsFingersList = append(assetsFingersList, a.InstanceProviderAdmin)
	assetsFingersList = append(assetsFingersList, a.StakingAdmin)

	return assetsFingersList
}

func (a *Andamio) GetAllAndamioAddresses() []string {
	var andamioAddr []string
	andamioAddr = append(andamioAddr, a.GlobalStateRefMS.MSCAddress)
	andamioAddr = append(andamioAddr, a.GlobalStateS.SCAddress)
	andamioAddr = append(andamioAddr, a.GovernanceS.SCAddress)
	andamioAddr = append(andamioAddr, a.IndexMS.MSCAddress)
	andamioAddr = append(andamioAddr, a.IndexRefMS.MSCAddress)
	andamioAddr = append(andamioAddr, a.InstanceMS.MSCAddress) // all instance vlidator at this address
	andamioAddr = append(andamioAddr, a.InstanceProvidedMS.MSCAddress)
	andamioAddr = append(andamioAddr, a.ReferenceAddr)

	return andamioAddr
}
