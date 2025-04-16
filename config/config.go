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
	Database Database `json:"database"`
	Andamio  Andamio  `json:"andamio"`
}

type Network struct {
	Magic      uint32 `json:"magic"`
	SocketPath string `json:"socketPath"`
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
