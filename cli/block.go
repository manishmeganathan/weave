package cli

import (
	"fmt"
	"math/rand"

	"github.com/manishmeganathan/animus/blockchain"
	"github.com/spf13/cobra"
)

// blockCmd represents the 'block' command
var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "Block Toolset",
	Long:  `Block Toolset - A set of tools to manipulate blocks on the Animus Blockchain`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please invoke a block command. Use 'animus block --help' for usage instructions")
	},
}

// blockAddCmd represents the 'block add' command
var blockAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Block",
	Long: `Add a Block to the Animus Blockchain

The string provided in the data flag of the command is used as the block data
to create and sign a block on the Animus Blockchain. If the data flag is not set, 
a randomly generated 20 character alpha-numeric string is used instead.`,

	Run: func(cmd *cobra.Command, args []string) {
		data, _ := cmd.Flags().GetString("data")

		if data == "" {
			data = generaterandomtext()
		}

		chain := blockchain.NewBlockChain()
		defer chain.Database.Close()

		chain.AddBlock(data)
		fmt.Println("Block Added!")
	},
}

func init() {
	// Create the 'block' command
	rootCmd.AddCommand(blockCmd)
	// Create the 'add' subcommand to 'block'
	blockCmd.AddCommand(blockAddCmd)
	// Add the 'data' flag to the 'block add' command
	blockAddCmd.Flags().StringP("data", "d", "", "data of the block")
}

// A function to generate a random 20 character case-sensitive alpha-numeric string.
func generaterandomtext() string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, 20)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}
