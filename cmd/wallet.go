package cmd

import (
	"fmt"

	"github.com/manishmeganathan/animus/wallet"
	"github.com/spf13/cobra"
)

// walletCmd represents the 'wallet' command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Create/List Animus Wallets",
	Long:  `Create/List Animus Wallets`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallet called")
	},
}

// walletCreateCmd represents the 'wallet create' command
var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Animus wallet",
	Long:  `Create a new Animus wallet`,

	Run: func(cmd *cobra.Command, args []string) {

		walletstore := wallet.NewWalletStore()
		address := walletstore.AddWallet()
		walletstore.Save()

		fmt.Println("Wallet Created!")
		fmt.Printf("Address of new wallet: %s\n", address)
	},
}

// walletShowCmd represents the 'wallet show' command
var walletShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all Animus wallets",
	Long:  `Show all Animus wallets`,

	Run: func(cmd *cobra.Command, args []string) {

		walletstore := wallet.NewWalletStore()
		addresses := walletstore.GetAddresses()

		for index, address := range addresses {
			fmt.Printf("%v] %s\n", index+1, address)
		}
	},
}

func init() {
	// Create the 'wallet' command
	rootCmd.AddCommand(walletCmd)

	// Create the 'create' subcommand to 'wallet'
	walletCmd.AddCommand(walletCreateCmd)

	// Create the 'show' subcommand to 'wallet'
	walletCmd.AddCommand(walletShowCmd)
}
