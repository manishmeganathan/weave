package cmd

import (
	"math/rand"
	"time"

	"github.com/manishmeganathan/weave/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weave",
	Short: "Weave CLI",
	Long: `
CLI tool for the Weave Network. Can be used to interact 
with the network, perform transactions and manage wallets.`,

	// Run: func(cmd *cobra.Command, args []string) {},
}

// CLI Entypoint.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Seed the random library
	rand.Seed(time.Now().UTC().UnixNano())
	// Setup logger
	utils.LogInitialize(4)
}
