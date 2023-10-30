/*
Copyright © 2023 Tessellated <tessellated.io>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start paymaster",
	Long:  `Starts paymaster with the given configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		panic("TODO: Implement me")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
