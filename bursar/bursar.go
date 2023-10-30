package bursar

import (
	"fmt"
	"time"

	"github.com/tessellated-io/paymaster/crypto"
	"github.com/tessellated-io/paymaster/skip"
	"github.com/tessellated-io/pickaxe/arrays"
	"github.com/tessellated-io/pickaxe/chains"
	"github.com/tessellated-io/pickaxe/coding"
	pacrypto "github.com/tessellated-io/pickaxe/crypto"
	"github.com/tessellated-io/pickaxe/tx"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// Time between disbursement
// TODO: Remove
const disbursementCooldown = 1 * time.Hour

// Dispenses funds
type Bursar struct {
	cdc               *codec.ProtoCodec
	offlineRegistry   *chains.OfflineChainRegistry
	skipClient        *skip.SkipClient
	bytesSigner       *pacrypto.KeyPair
	knownPublicKeys   []cryptotypes.PubKey
	lastDisbursements map[string]time.Time

	sourceChainID            string
	senderAddress            string
	sendDenom                string
	sendAmount               string
	sourceChainAddressPrefix string
	sourceChainGrpcUri       string
}

func NewBursar(
	cdc *codec.ProtoCodec,
	offlineRegistry *chains.OfflineChainRegistry,
	skipClient *skip.SkipClient,
	bytesSigner *pacrypto.KeyPair,
) *Bursar {
	knownPublicKeysHex := []string{"02d51823fa1509ae8b57693ab973c5d6ee3cdcaac809866539085f7abfb0793bb5"}
	knownPublicKeys := arrays.Map(arrays.Map(knownPublicKeysHex, coding.UnsafeHexToBytes), crypto.Secp256k1PublicKeyFromBytes)

	// TODO: stop hardcoding these
	sourceChainID := "cosmoshub-4"
	sourceChainAddressPrefix := "cosmos"
	senderAddress := bytesSigner.GetAddress(sourceChainAddressPrefix)
	fmt.Printf("Bursar online, disbursing from %s\n", senderAddress)
	sendDenom := "uatom"
	sendAmount := "500000"
	sourceChainGrpcUri := "cosmos-validator.tessageo.net:9090"

	return &Bursar{
		cdc:               cdc,
		offlineRegistry:   offlineRegistry,
		skipClient:        skipClient,
		bytesSigner:       bytesSigner,
		knownPublicKeys:   knownPublicKeys,
		lastDisbursements: make(map[string]time.Time),

		senderAddress:            senderAddress,
		sourceChainAddressPrefix: sourceChainAddressPrefix,
		sourceChainID:            sourceChainID,
		sendDenom:                sendDenom,
		sendAmount:               sendAmount,
		sourceChainGrpcUri:       sourceChainGrpcUri,
	}
}

// ==============
// PUBLIC API

func (b *Bursar) SendFunds(targetAddress, prefix string) (string, error) {
	// Auth checks
	// 1. Is known address
	// if !b.isKnownAddress(targetAddress, prefix) {
	// 	return "", fmt.Errorf("unauthorized address %s", targetAddress)
	// }

	// 2. Only give out funds once per hour
	if b.isInCooldown(targetAddress) {
		return "", fmt.Errorf("address %s has received funds too recently. Please wait a moment before trying again", targetAddress)
	}

	// 3. Create an API message
	destChainID := b.offlineRegistry.AccountPrefixToData[prefix].ChainID
	destChainDenom := b.offlineRegistry.AccountPrefixToData[prefix].NativeToken
	publicKey := b.bytesSigner.GetPublicKey()
	ibcXferMessage, err := b.skipClient.GetMessages(
		b.senderAddress,
		publicKey,
		targetAddress,
		b.sendAmount,
		b.sendDenom,
		b.sourceChainID,
		destChainDenom,
		destChainID,
	)
	if err != nil {
		return "", err
	}

	// 3. Send funds
	grpcRes, err := tx.SendMessages(
		[]sdk.Msg{ibcXferMessage},
		b.sourceChainAddressPrefix,
		b.sourceChainID,
		b.cdc,
		b.bytesSigner,
		txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
		b.sourceChainGrpcUri,
		1.1,
		0.0, // TODO
		"TODO",
	)
	if err != nil {
		return "", err
	}

	fmt.Printf(
		"sent %s%s to %s with code %d (%s) in %s\n",
		b.sendAmount,
		b.sendDenom,
		targetAddress,
		grpcRes.TxResponse.Code,
		grpcRes.TxResponse.RawLog,
		grpcRes.TxResponse.TxHash,
	)

	// 4. If successful, notate time to activate cooldown
	if grpcRes.TxResponse.Code == 0 {
		b.lastDisbursements[targetAddress] = time.Now()
		return grpcRes.TxResponse.TxHash, nil
	}
	return "", fmt.Errorf("failed to send for height %d. Error: %s", 123, grpcRes.TxResponse.RawLog)
}

// ==============
// HELPERS

// Check that we haven't sent funds recently.
func (b *Bursar) isInCooldown(address string) bool {
	lastDisbursement, hasPrevious := b.lastDisbursements[address]
	if !hasPrevious {
		// have not sent money yet
		return false
	}

	now := time.Now()
	nextDisburment := lastDisbursement.Add(disbursementCooldown)
	return now.Before(nextDisburment)
}

// TODO
// // Usage: isKnownAddress("osmo15qth07rmamcue638q4fvzfrg9ra6eyknqh3jmc", "osmo")
// func (b *Bursar) isKnownAddress(addressToCheck, prefix string) bool {
// 	for _, knownPublicKey := range b.knownPublicKeys {
// 		computedAddress, err := crypto.PubKeyToAddress(knownPublicKey, prefix)
// 		if err != nil {
// 			fmt.Printf(
// 				"Warning: Could not formulate an address for public key and prefix, skipping. Error: %s (pub_key_hex=%s, prefix=%s)",
// 				err,
// 				hex.EncodeToString(knownPublicKey.Bytes()),
// 				prefix,
// 			)
// 			continue
// 		}

// 		if strings.EqualFold(computedAddress, addressToCheck) {
// 			return true
// 		}
// 	}
// 	return false
// }
