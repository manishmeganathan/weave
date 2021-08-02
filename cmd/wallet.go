package cmd

import (
	"fmt"

	"github.com/manishmeganathan/weave/utils"
	"github.com/manishmeganathan/weave/wallet"
	"github.com/spf13/cobra"
)

// walletCmd represents the 'wallet' command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Interact with the Weave wallet",
	Long:  `Interact with the Weave wallet`,
	// Run: func(cmd *cobra.Command, args []string) {},
}

// wallet_getCmd represents the 'wallet get' command
var wallet_getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current wallet address",
	Long:  `Get the current wallet address`,
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.ReadConfigFile()
		fmt.Println(config.JBOK.Default)
	},
}

// wallet_setCmd represents the 'wallet set' command
var wallet_setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the current wallet address",
	Long:  `Set the current wallet address`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if args has elements
		if len(args) == 0 {
			fmt.Println("[error] address not provided.")
			return
		}

		// Retrieve the address from the first argument
		address := args[0]
		// Create a new JBOK object
		jbok := wallet.NewJBOK()
		// Check if the address is registered with the JBOK
		if !jbok.CheckWallet(address) {
			fmt.Println("[error] provided address does not exist in the JBOK.")
			return
		}

		// Read the configuration file into an object
		config := utils.ReadConfigFile()
		// Set the wallet address to the config
		config.JBOK.Default = address
		// Write the configuration to a file
		config.WriteConfigFile()
		// Print the confirmation
		fmt.Println("[success] weave wallet address set to:", address)
	},
}

// wallet_listCmd represents the 'wallet list' command
var wallet_listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all wallet addresses in the JBOK",
	Long:  `List all wallet addresses in the JBOK`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new JBOK object
		jbok := wallet.NewJBOK()
		addrs := jbok.GetAddresses()

		fmt.Println("--JBOK-Wallet-Addresses--")
		for index, addr := range addrs {
			fmt.Printf("%d] %s\n", index, addr)
		}

	},
}

func init() {
	// Add wallet command to root
	rootCmd.AddCommand(walletCmd)
	// Add get command to wallet
	walletCmd.AddCommand(wallet_getCmd)
	// Add set command to wallet
	walletCmd.AddCommand(wallet_setCmd)
	// Add list command to wallet
	walletCmd.AddCommand(wallet_listCmd)
}
