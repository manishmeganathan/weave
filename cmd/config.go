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
	Run:   func(cmd *cobra.Command, args []string) { print_config() },
}

// configShowCmd represents the 'config show' command
var config_showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the values from the configuration file.",
	Long:  `Display the values from the configuration file.`,
	Run:   func(cmd *cobra.Command, args []string) { print_config() },
}

// configResetCmd represents the 'config reset' command
var config_resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration file.",
	Long: `Reset the configuration file. If the file does not exist, it will be created.
Command expects a wallet address as an argument. This address must already be registered with the JBOK.`,
	Run: func(cmd *cobra.Command, args []string) { generate_config(args) },
}

// configGenerateCmd represents the 'config generate' command
var config_generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the configuration file.",
	Long: `Generate the configuration file. If the file already exists, it will be overwritten.
Command expects a wallet address as an argument. This address must already be registered with the JBOK.`,
	Run: func(cmd *cobra.Command, args []string) { generate_config(args) },
}

// A function to generate the configuration file
// Expects the wallet address as the first argument in a list of strings.
// The wallet address must be registered with the JBOK.
func generate_config(args []string) {
	if len(args) == 0 {
		fmt.Println("[error] default address not provided.")
		return
	}

	address := args[0]

	jbok := wallet.NewJBOK()
	if !jbok.CheckWallet(address) {
		fmt.Println("[error] provided address does not exist in the JBOK.")
		return
	}

	config := utils.GenerateConfigFile(false)
	config.JBOK.Default = address
	config.WriteConfigFile()
}

// A function that prints the configuration file values
func print_config() {
	config := utils.ReadConfigFile()
	config.PrintConfigFile()
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
