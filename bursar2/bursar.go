package bursar2

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tessellated-io/mail-in-rebates/paymaster/config"
	"github.com/tessellated-io/pickaxe/cosmos/rpc"
	r "github.com/tessellated-io/router/router"
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
}

// Ensure bursar is a Bursar
var _ Bursar = (*bursar)(nil)

func NewBursar(accounts []*config.AccountConfig) (Bursar, error) {
	return &bursar{
		accounts:         accounts,
		rateLimitTracker: make(map[string]time.Time),
	}, nil
}

// Bursar Interface

// TODO: Add logging
func (b *bursar) PollForTopUps(ctx context.Context) {
	for _, account := range b.accounts {
		// TODO: Incorporate router

		// Create a gRPC client
		// TODO: We can probably create these up front to avoid connection overhead
		endpoint, err := b.router.GetGrpcEndpoint()
		if err != nil {
			return nil
		}
		rpcClient := rpc.NewGRpcClient())

		// TODO: retries
		balance, err := rpcClient.GetBalance(ctx, account.Address, account.MinCoin.Denom)
		if err != nil {
			return err
		}

		// Check balance
		if balance.LessThank(account.MinCoin.Amount) {
			txHash, err := b.SendFunds(account.TopUpAmount, account.Address, account.RateLimit)
			if err != nil {
				return err
			}
		}
	}
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
