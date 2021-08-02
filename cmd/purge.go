package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/manishmeganathan/weave/utils"
	"github.com/manishmeganathan/weave/wallet"
	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge Weave resources",
	Long:  `Purge Weave resources such as the JBOK data, database files or the config file.`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// purge_jbokCmd represents the 'purge jbok' command
var purge_jbokCmd = &cobra.Command{
	Use:   "jbok",
	Short: "Purge the Weave JBOK file",
	Long:  `Purge the Weave JBOK file`,
	Run: func(cmd *cobra.Command, args []string) {
		// Purge JBOK data
		wallet.PurgeJBOK()
	},
}

// purge_walletCmd represents the 'purge wallet' command
var purge_walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Purge a Weave wallet",
	Long:  `Purge a Weave wallet`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if args has elements
		if len(args) == 0 {
			fmt.Println("[error] wallet address not provided.")
			return
		}

		// Get the wallet address to purge
		walletres := args[0]
		// Create a new JBOK object
		jbok := wallet.NewJBOK()
		// Remove the wallet from the JBOK
		jbok.RemoveWallet(walletres)
	},
}

// purge_configCmd represents the 'purge config' command
var purge_configCmd = &cobra.Command{
	Use:   "config",
	Short: "Purge the Weave config file",
	Long:  `Purge the Weave config file`,
	Run: func(cmd *cobra.Command, args []string) {
		// Purge Config data
		utils.RemoveConfigFile()
	},
}

// purge_dbCmd represents the 'purge db' command
var purge_dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Purge the Weave database files",
	Long:  `Purge the Weave database files`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the path the db directory
		dbdir := filepath.Join(utils.ConfigDirectory(), "db")
		// Clear the db directory
		utils.ClearDirectory(dbdir)
	},
}

func init() {
	// Add purge command to root
	rootCmd.AddCommand(purgeCmd)
	// Add jbok command to purge
	purgeCmd.AddCommand(purge_jbokCmd)
	// Add wallet command to purge
	purgeCmd.AddCommand(purge_walletCmd)
	// Add config command to purge
	purgeCmd.AddCommand(purge_configCmd)
	// Add db command to purge
	purgeCmd.AddCommand(purge_dbCmd)
}
