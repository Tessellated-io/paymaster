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
	ChainID string
	Grpc    string

	// What address do I care about
	Address   string
	Denom     string
	MinAmount json.Number

	// Payout configurations
	topUpAmount      json.Number
	rateLimitSeconds int
}
