package config

import (
	"context"
	"time"

	chainregistry "github.com/tessellated-io/pickaxe/cosmos/chain-registry"
	"github.com/tessellated-io/pickaxe/log"
	"github.com/tessellated-io/pickaxe/util"
	r "github.com/tessellated-io/router/router"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GlobalConfig struct {
	Mnemonic    string
	RunInterval time.Duration
}

type AccountConfig struct {
	Address string
	MinCoin sdk.Coin

	TopUpAmount sdk.Coin

	RateLimit time.Duration
}

func GetConfig(ctx context.Context, filename string, log *log.Logger) (*GlobalConfig, []*AccountConfig, r.Router, error) {
	// Get data from the file
	fileConfig, err := parseConfig(filename)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a global config
	globalConfig := &GlobalConfig{
		Mnemonic:    fileConfig.Mnemonic,
		RunInterval: time.Duration(fileConfig.RunIntervalSeconds) * time.Second,
	}

	// Create and route account configs for each account
	router, err := r.NewRouter(nil)
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
		minCoin := sdk.NewCoin(fileAccountConfig.Denom, sdk.NewIntFromBigInt(minAmount))

		topUpAmount, err := util.NumberToBigInt(fileAccountConfig.TopUpAmount)
		if err != nil {
			return nil, nil, nil, err
		}
		topUpCoin := sdk.NewCoin(fileAccountConfig.Denom, sdk.NewIntFromBigInt(topUpAmount))

		rateLimit := time.Duration(fileAccountConfig.RateLimitSeconds) * time.Second

		accountConfig := &AccountConfig{
			Address:     fileAccountConfig.Address,
			MinCoin:     minCoin,
			TopUpAmount: topUpCoin,
			RateLimit:   rateLimit,
		}
		accountConfigs = append(accountConfigs, accountConfig)

		// TODO: hit cache for a chain client
		// TODO: hit cache for a chain config?
		registryClient := chainregistry.NewRegistryClient(log)
		chainName, err := registryClient.ChainNameForChainID(ctx, fileAccountConfig.ChainID)
		if err != nil {
			return nil, nil, nil, err
		}

		registryChainInfo, err := registryClient.GetChainInfo(ctx, chainName)
		if err != nil {
			return nil, nil, nil, err
		}

		// Create and add chain to router
		chain, err := r.NewChain(
			fileAccountConfig.ChainID,
			fileAccountConfig.ChainID, // TODO: use a name here when we have chain registry
			registryChainInfo.Bech32Prefix,
			&fileAccountConfig.Grpc,
		)
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
