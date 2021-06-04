package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/manishmeganathan/animus/blockchain"
	"github.com/manishmeganathan/animus/wallet"
	"github.com/spf13/cobra"
)

// chainCmd represents the 'chain' command
var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Create/View Animus Blockchains",
	Long:  `Create/View Animus Blockchains`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please invoke a chain command. Use 'animus chain --help' for usage instructions")
	},
}

// chainShowCmd represents the 'chain show' command
var chainShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the chain of Blocks",
	Long: `Show the chain of Blocks on the Animus Blockchain.

Prints the Hash of the previous block, current block and validation of
each block on the Animus Blockchain`,

	Run: func(cmd *cobra.Command, args []string) {

		chain, err := blockchain.AnimateBlockChain()
		if err != nil {
			fmt.Println("Animus Blockchain does not exist! Use 'animus chain create' to create one.")
			chain.Database.Close()
			return
		}

		defer chain.Database.Close()
		iter := blockchain.NewIterator(chain)

		fmt.Print("\n---- ANIMUS BLOCKCHAIN ----\n\n")
		for {
			block := iter.Next()

			fmt.Printf("Previous Hash: %x\n", block.PrevHash)
			fmt.Printf("Block Hash: %x\n", block.Hash)
			fmt.Printf("Signed: %s\n\n", strconv.FormatBool(block.Validate()))

			for _, txn := range block.Transactions {
				fmt.Println(txn.String())
			}

			fmt.Println()

			if len(block.PrevHash) == 0 {
				break
			}
		}
		fmt.Print("---- END OF BLOCKCHAIN ----\n\n")
	},
}

// chainCreateCmd represents the 'chain create' command
var chainCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Blockchain",
	Long:  `Create a new Blockchain with a coinbase transaction for the specified miner`,

	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")

		if !wallet.ValidateWalletAddress(address) {
			log.Panic("Invalid Address!")
		}

		chain, err := blockchain.SeedBlockChain(address)
		if err != nil {
			fmt.Println("Animus Blockchain already exists. Use 'animus chain show' to view it.")
			chain.Database.Close()
			return
		}

		defer chain.Database.Close()

		fmt.Println("Animus Blockchain Created!")
	},
}

func init() {
	// Create the 'chain' command
	rootCmd.AddCommand(chainCmd)

	// Create the 'show' subcommand to 'chain'
	chainCmd.AddCommand(chainShowCmd)

	// Create the 'create' subcommand to 'chain'
	chainCmd.AddCommand(chainCreateCmd)
	// Add the 'address' flag to the 'chain create' command and mark it as required
	chainCreateCmd.Flags().StringP("address", "a", "", "address of the block miner")
	chainCreateCmd.MarkFlagRequired("address")
}
