package config

import "encoding/json"

// UserConfig is the top level config provided by a user of paymaaster
type UserConfig struct {
	// Account that will pay us
	Mnemonic string `yaml:"mnemonic"`

	// How often the check is run
	RunIntervalSeconds int `yaml:"runIntervalSeconds"`

	// Accounts to manage
	Accounts []UserAccountConfig `yaml:"userAccountConfig"`

	// TODO: some sort of alerting?
}

// UserAccountConfig is configuration for an account to monitor.
type UserAccountConfig struct {
	// Where to find data
	ChainID string `yaml:"chainID"`
	Grpc    string `yaml:"grpc"`

	// What address do I care about
	Address   string      `yaml:"address"`
	Denom     string      `yaml:"denom"`
	MinAmount json.Number `yaml:"minAmount"`

	// Payout configurations
	TopUpAmount      json.Number `yaml:"topUpAmount"`
	RateLimitSeconds int         `yaml:"rateLimitSeconds"`
}
