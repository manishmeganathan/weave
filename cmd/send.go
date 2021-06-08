package cmd

import (
	"fmt"
	"log"

	"github.com/manishmeganathan/blockweave/blockchain"
	"github.com/manishmeganathan/blockweave/wallet"
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

		if !wallet.ValidateWalletAddress(sender) {
			log.Panic("Invalid Sender Address!")
		}

		if !wallet.ValidateWalletAddress(receiver) {
			log.Panic("Invalid Receiver Address!")
		}

		chain, err := blockchain.AnimateBlockChain()
		if err != nil {
			fmt.Println("Animus Blockchain does not exist! Use 'animus chain create' to create one.")
			chain.Database.Close()
			return
		}

		defer chain.Database.Close()

		tx := blockchain.NewTransaction(sender, receiver, amount, chain)
		coinbase := blockchain.NewCoinbaseTransaction(sender)

		block := chain.AddBlock([]*blockchain.Transaction{coinbase, tx})
		chain.UpdateUTXOS(block)

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
	sendCmd.Flags().StringP("receiver", "r", "", "address of reciever")
	sendCmd.MarkFlagRequired("receiver")

	// Add the 'amount' flag to the 'send' command and mark it as required
	sendCmd.Flags().IntP("amount", "a", 0, "amount to send")
	sendCmd.MarkFlagRequired("amount")
}
