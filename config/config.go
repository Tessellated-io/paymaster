package config

import (
	"context"
	"log"
	"time"

	"github.com/tessellated-io/pickaxe/util"
	r "github.com/tessellated-io/router/router"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GlobalConfig struct {
	Mnemonic    string
	RunInterval time.Duration
}

type AccountConfig struct {
	ChainID string
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
