package bursar2

import (
	"context"
	"fmt"
	"time"

	"github.com/tessellated-io/mail-in-rebates/paymaster/config"
	"github.com/tessellated-io/pickaxe/cosmos/rpc"
	"github.com/tessellated-io/pickaxe/log"
	r "github.com/tessellated-io/router/router"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Lock on addresses to prevent multi thread attacks
// TODO: package name

// Bursar implements information for polling account balances and disbursing funds.
type Bursar interface{}

type bursar struct {
	// Configs for accounts to check.
	accounts []*config.AccountConfig
	router   r.Router

	// Map of address to last successful send time.
	// TODO: Move this to persistence
	rateLimitTracker map[string]time.Time

	cdc *codec.ProtoCodec
	log log.Logger
}

// Ensure bursar is a Bursar
var _ Bursar = (*bursar)(nil)

func NewBursar(accounts []*config.AccountConfig, router r.Router, cdc *codec.ProtoCodec, log log.Logger) (Bursar, error) {
	return &bursar{
		accounts:         accounts,
		rateLimitTracker: make(map[string]time.Time),
		router:           router,
		cdc:              cdc,
		log:              log,
	}, nil
}

// Bursar Interface

// TODO: Add logging
func (b *bursar) PollForTopUps(ctx context.Context) error {
	for _, account := range b.accounts {
		// TODO: Incorporate router

		// Create a gRPC client
		// TODO: We can probably create these up front to avoid connection overhead
		grpcEndpoint, err := b.router.GetGrpcEndpoint(account.ChainID)
		if err != nil {
			return nil
		}
		rpcClient, err := rpc.NewGRpcClient(grpcEndpoint, b.cdc, &b.log)
		if err != nil {
			return nil
		}

		// TODO: retries
		balance, err := rpcClient.GetBalance(ctx, account.Address, account.MinCoin.Denom)
		if err != nil {
			return err
		}

		// Check balance
		if balance.IsLT(account.MinCoin) {
			txHash, err := b.SendFunds(account.TopUpAmount, account.Address, account.RateLimit)
			if err != nil {
				return err
			}
			b.log.Info().Str("tx_hash", txHash).Str("target_address", account.Address).Msg(fmt.Sprintf("sent %s%s", account.TopUpAmount.Amount, account.TopUpAmount.Denom))
		}
		// TODO: Log?
	}
	return nil
}

// TODO: Public?
// SendFunds sends the given funds to the given address while respecting the rate limiting timer, returning a transaction hash or an error.
func (b *bursar) SendFunds(amount sdk.Coin, address string, rateLimit time.Duration) (string, error) {
	// Ensure we are not rate limited
	err := b.isRateLimited(address, rateLimit)
	if err != nil {
		return "", err
	}

	// Send funds
	txHash, err := b.sendFunds(amount, address)
	if err != nil {
		return "", err
	}

	// Update the last send time to enforce rate limit
	b.recordFundsSent(address)

	return txHash, nil
}

// Private Helpers

// Sends a transaction and returns a transaction hash and error on inclusion
func (b *bursar) sendFunds(amount sdk.Coin, address string) (string, error) {
	// TODO: retries

	// TODO
	return "", fmt.Errorf("TODO")
}

// Returns an error if the given address is rate limited.
func (b *bursar) isRateLimited(address string, rateLimit time.Duration) error {
	// If nothing is being tracked, there has never been a send and we are not rate limited
	lastDisbursement, isSet := b.rateLimitTracker[address]
	if !isSet {
		return nil
	}

	// Otherwise, check if now is past the cooldown period for that
	now := time.Now()
	nextEligibleTime := lastDisbursement.Add(rateLimit)

	if now.Before(nextEligibleTime) {
		return ErrTooSoonForTopUp
	}
	return nil
}

func (b *bursar) recordFundsSent(address string) {
	// TODO: possibly record a tx identifier here
	b.rateLimitTracker[address] = time.Now()
}
