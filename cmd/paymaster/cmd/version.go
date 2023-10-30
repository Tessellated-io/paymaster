/*
Copyright © 2023 Tessellated <tessellated.io>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Binary name
const (
	binaryName = "paymaster"
	binaryIcon = "💸"
)

// Version
var (
	PaymasterVersion string
	GoVersion        string
	GitRevision      string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version of paymaster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s  %s:\n", binaryIcon, binaryName)
		fmt.Printf("  - Version: %s\n", PaymasterVersion)
		fmt.Printf("  - Git Revision: %s\n", GitRevision)
		fmt.Printf("  - Go Version: %s\n", GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
