package cmd

import (
	"fmt"

	"github.com/manishmeganathan/animus/blockchain"
	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send tokens",
	Long: `Sends tokens from one address to another

Command requires the 'sender', 'reciever' and 'amount' flags`,

	Run: func(cmd *cobra.Command, args []string) {
		sender, _ := cmd.Flags().GetString("sender")
		receiver, _ := cmd.Flags().GetString("receiver")
		amount, _ := cmd.Flags().GetInt("amount")

		chain, err := blockchain.AnimateBlockChain()
		if err != nil {
			fmt.Println("Animus Blockchain does not exist! Use 'animus chain create' to create one.")
			chain.Database.Close()
			return
		}

		defer chain.Database.Close()

		tx := blockchain.NewTransaction(sender, receiver, amount, chain)
		chain.AddBlock([]*blockchain.Transaction{tx})
		fmt.Println("Success!")
	},
}

func init() {
	// Create the 'send' command
	rootCmd.AddCommand(sendCmd)

	// Add the 'sender' flag to the 'send' command and mark it as required
	sendCmd.Flags().StringP("sender", "s", "", "address of sender")
	sendCmd.MarkFlagRequired("sender")

	// Add the 'reciever' flag to the 'send' command and mark it as required
	sendCmd.Flags().StringP("reciever", "r", "", "address of reciever")
	sendCmd.MarkFlagRequired("reciever")

	// Add the 'amount' flag to the 'send' command and mark it as required
	sendCmd.Flags().IntP("amount", "a", 0, "amount to send")
	sendCmd.MarkFlagRequired("amount")
}