package config

import (
	"context"
	"log"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tessellated-io/pickaxe/util"
	"github.com/tessellated-io/router/router"
	r "github.com/tessellated-io/router/router"
)

type GlobalConfig struct {
	Mnemonic           string
	RunIntervalSeconds time.Duration
}

type AccountConfig struct {
	Address string
	minCoin sdk.Coin

	topUpAmount sdk.Coin

	rateLimiteSeconds time.Duration
}

func GetConfig(ctx context.Context, filename string, log *log.Logger) (*GlobalConfig, []*AccountConfig, router.Router, error) {
	// Get data from the file
	fileConfig, err := parseConfig(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a global config
	globalConfig := &GlobalConfig{
		Mnemonic:           fileConfig.Mnemonic,
		RunIntervalSeconds: time.Duration(fileConfig.RunIntervalSeconds) * time.Second,
	}

	// Create and route account configs for each account
	router, err := router.NewRouter(nil)
	if err != nil {
		return nil, nil, nil, err
	}
	accountConfigs := []*AccountConfig{}
	for _, fileAccountConfig := range fileConfig.Accounts {
		// Create and add account config
		minAmount, err := util.NumberToBigInt(fileAccountConfig.MinAmount)
		if err != nil {
			return nil, nil, nil, err
		}
		minCoin := sdk.NewCoin(fileAccountConfig.Denom, minAmount)

		accountConfig := &AccountConfig{
			Address: fileAccountConfig.Address,
			minCoin: minCoin,
		}
		accountConfigs = append(accountConfigs, accountConfig)

		// Create and add chain to router
		chain, err := r.NewChain(fileAccountConfig.ChainID, fileAccountConfig.ChainID, &fileAccountConfig.Grpc)
		if err != nil {
			return nil, nil, nil, err
		}

		err = router.AddChain(chain)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return globalConfig, accountConfigs, router, nil
}
