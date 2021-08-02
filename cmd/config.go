package cmd

import (
	"fmt"

	"github.com/manishmeganathan/weave/utils"
	"github.com/manishmeganathan/weave/wallet"
	"github.com/spf13/cobra"
)

// configCmd represents the 'config' command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View configuration information",
	Long:  `View configuration information`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read the configuration file into an object
		config := utils.ReadConfigFile()
		// Print the configuration file values
		config.PrintConfigFile()
	},
}

// config_showCmd represents the 'config show' command
var config_showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show values from the configuration file",
	Long: `Show values from the configuration file. 
Commands expects a value that represents the config type. 
Valid values are 'all', 'jbok', 'db', 'blocks' and 'state'.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Read the configuration file into an object
		config := utils.ReadConfigFile()

		// Check if args has elements
		if len(args) == 0 {
			fmt.Println("[error] no config value provided.")
			return
		}

		// Check the show value
		switch args[0] {
		case "all":
			// Print all the configuration file values
			config.PrintConfigFile()

		case "jbok":
			// Print the JBOK configuration file values
			fmt.Println()
			fmt.Println("----JBOK-Configuration----")
			fmt.Printf("JBOK File: %v\n", config.JBOK.File)
			fmt.Printf("JBOK Default: %v\n", config.JBOK.Default)
			fmt.Println()

		case "db":
			// Print the DB configuration file values
			fmt.Println()
			fmt.Println("----Database-Configuration----")
			fmt.Printf("DB Root: %v\n", config.DB.Root)
			fmt.Printf("DB State File: %v\n", config.DB.State.File)
			fmt.Printf("DB State Directory: %v\n", config.DB.State.Directory)
			fmt.Printf("DB Blocks File: %v\n", config.DB.Blocks.File)
			fmt.Printf("DB Blocks Directory: %v\n", config.DB.Blocks.Directory)
			fmt.Println()

		case "blocks":
			// Print the Database Blocks configuration file values
			fmt.Println()
			fmt.Println("----Database-Blocks-Configuration----")
			fmt.Printf("DB Blocks File: %v\n", config.DB.Blocks.File)
			fmt.Printf("DB Blocks Directory: %v\n", config.DB.Blocks.Directory)
			fmt.Println()

		case "state":
			// Print the Database State configuration file values
			fmt.Println()
			fmt.Println("----Database-State-Configuration----")
			fmt.Printf("DB State File: %v\n", config.DB.State.File)
			fmt.Printf("DB State Directory: %v\n", config.DB.State.Directory)
			fmt.Println()

		case "net":
			// Print the Network configuration file values
			fmt.Println()
			// fmt.Println("----Network-Configuration----")
			// fmt.Println()

		default:
			fmt.Println("[error] invalid config value provided.")
		}
	},
}

// config_resetCmd represents the 'config reset' command
var config_resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration file.",
	Long: `Reset the configuration file. If the file does not exist, it will be created.
Command expects a wallet address as an argument. This address must already be registered with the JBOK.`,
	Run: func(cmd *cobra.Command, args []string) { generate_config(args) },
}

// config_generateCmd represents the 'config generate' command
var config_generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the configuration file",
	Long: `Generate the configuration file. If the file already exists, it will be overwritten.
Command expects a wallet address as an argument. This address must already be registered with the JBOK.`,
	Run: func(cmd *cobra.Command, args []string) { generate_config(args) },
}

// A function to generate the configuration file
// Expects the wallet address as the first argument in a list of strings.
// The wallet address must be registered with the JBOK.
func generate_config(args []string) {
	// Check if args has elements
	if len(args) == 0 {
		fmt.Println("[error] default address not provided.")
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

	// Create a new default configuration
	// without writing the object to a file
	config := utils.GenerateConfigFile(false)
	// Set the wallet address to the config
	config.JBOK.Default = address
	// Write the configuration to a file
	config.WriteConfigFile()
}

func init() {
	// Add config command to root
	rootCmd.AddCommand(configCmd)
	// Add show command to config
	configCmd.AddCommand(config_showCmd)
	// Add reset command to config
	configCmd.AddCommand(config_resetCmd)
	// Add generate command to config
	configCmd.AddCommand(config_generateCmd)
}
