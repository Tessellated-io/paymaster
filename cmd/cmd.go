package main

import (
	"github.com/spf13/cobra"
	"github.com/tessellated-io/mail-in-rebates/paymaster/bursar"
	"github.com/tessellated-io/mail-in-rebates/paymaster/codec"
	"github.com/tessellated-io/mail-in-rebates/paymaster/server"
	"github.com/tessellated-io/mail-in-rebates/paymaster/skip"
	"github.com/tessellated-io/mail-in-rebates/paymaster/tracker"
	"github.com/tessellated-io/pickaxe/chains"
	pacrypto "github.com/tessellated-io/pickaxe/crypto"
)

// TODO: Separate these out into different commands
// TODO: Enable persistence on state.
var rootCmd = &cobra.Command{
	Use:   "paymaster",
	Short: "Paymaster helps distribute payments to crypto wallets",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		mnemonic := "TODO"
		keyPair := pacrypto.NewCosmosKeyPairFromMnemonic(mnemonic)

		cdc := codec.GetCodec()

		offlineRegistry := chains.NewOfflineChainRegistry()
		skipClient := skip.NewSkipClient(offlineRegistry, cdc)

		addressTracker, err := tracker.NewAddressTracker("/home/ubuntu/paymaster.csv")
		if err != nil {
			panic(err)
		}
		addressTracker.AddAddress("Test test")

		bursar := bursar.NewBursar(
			cdc,
			offlineRegistry,
			skipClient,
			keyPair,
		)
		server.StartPaymasterServer(bursar, addressTracker, port)
	},
}

func init() {
	// Define the 'port' flag and set it as optional
	rootCmd.Flags().Int("port", 8080, "Port number (optional)")
}

func main() {
	rootCmd.Execute()
}
