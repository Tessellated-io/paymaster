package main

import (
	"encoding/base64"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/tessellated-io/mail-in-rebates/paymaster/codec"
	"github.com/tessellated-io/mail-in-rebates/paymaster/crypto"
	"github.com/tessellated-io/mail-in-rebates/paymaster/skip"
	"github.com/tessellated-io/pickaxe/chains"
	pacrypto "github.com/tessellated-io/pickaxe/crypto"
	"github.com/tessellated-io/pickaxe/tx"
)

func main() {
	cdc := codec.GetCodec()

	// Get data from Skip API
	offlineRegistry := chains.NewOfflineChainRegistry()
	skipClient := skip.NewSkipClient(offlineRegistry, cdc)

	pubKeyBytes, _ := base64.StdEncoding.DecodeString("AtUYI/oVCa6LV2k6uXPF1u483KrICYZlOQhfer+weTu1")
	pubKey := crypto.Secp256k1PublicKeyFromBytes(pubKeyBytes)

	sourceChainID := "cosmoshub-4"
	senderAdddress := "cosmos15qth07rmamcue638q4fvzfrg9ra6eykngvzzd2"
	ibcXferMessage, err := skipClient.GetMessages(
		senderAdddress,
		pubKey,
		"axelar1jw7a28g98q3e7ul9f78cuzxnaw67dax8znaz9s",
		"100000",
		"uatom",
		sourceChainID,
		"uaxl",
		"axelar-dojo-1",
	)
	if err != nil {
		panic(err)
	}

	// Get a signer
	mnemonic := "TODO"
	keyPair := pacrypto.NewCosmosKeyPairFromMnemonic(mnemonic)

	grpcRes, err := tx.SendMessages(
		[]sdk.Msg{ibcXferMessage},
		"cosmos",
		sourceChainID,
		cdc,
		keyPair,
		txtypes.BroadcastMode_BROADCAST_MODE_SYNC,
		"cosmos-validator.tessageo.net:9090",
		1.1,
		0.0, // TODO
		"TODO",
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Codespace: %s \n", grpcRes.TxResponse.Codespace)
	fmt.Printf("Info: %s\n", grpcRes.TxResponse.Info)
	fmt.Printf("Info: %s\n", grpcRes.TxResponse.TxHash)
	fmt.Printf("Code: %d\n", grpcRes.TxResponse.Code) // Should be `0` if the tx is successful
	fmt.Printf("Logs: %s\n", grpcRes.TxResponse.Logs)

	if grpcRes.TxResponse.Code == 0 {

		fmt.Println("🎉 🎉 🎉")
		fmt.Printf("🎉 Transaction Sent\n")
		fmt.Println("🎉 🎉 🎉")
		return
	}
	fmt.Printf("Failed to send for height %d. Error: %s\n", 123, grpcRes.TxResponse.RawLog)
}
