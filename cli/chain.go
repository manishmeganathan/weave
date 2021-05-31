package cli

import (
	"fmt"
	"strconv"

	"github.com/manishmeganathan/animus/blockchain"
	"github.com/spf13/cobra"
)

// chainCmd represents the 'chain' command
var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Chain Toolset",
	Long:  `Chain Toolset - A set of tools to manipulate the chain of the Animus Blockchain`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please invoke a chain command. Use 'animus chain --help' for usage instructions")
	},
}

// chainShowCmd represents the 'chain show' command
var chainShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the chain of Blocks",
	Long: `Show the chain of Blocks on the Animus Blockchain.

Prints the Data, Signature and Validation of the 
every block on the Animus Blockchain`,

	Run: func(cmd *cobra.Command, args []string) {

		chain := blockchain.NewBlockChain()
		defer chain.Database.Close()

		iter := blockchain.NewIterator(chain)

		fmt.Print("\n---- ANIMUS BLOCKCHAIN ----\n\n")
		for {
			block := iter.Next()

			fmt.Printf("Previous Hash: %x\n", block.PrevHash)
			fmt.Printf("Data in Block: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Printf("PoW: %s\n\n", strconv.FormatBool(block.Validate()))

			if len(block.PrevHash) == 0 {
				break
			}
		}
		fmt.Print("---- END OF BLOCKCHAIN ----\n\n")
	},
}

func init() {
	rootCmd.AddCommand(chainCmd)
	chainCmd.AddCommand(chainShowCmd)
}
