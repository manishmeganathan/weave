package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weave",
	Short: "Weave CLI",
	Long: `
CLI tool for the Weave Network. Can be used to interact 
with the network, perform transactions and manage wallets.`,

	// Run: func(cmd *cobra.Command, args []string) {},
}

// CLI Entypoint.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Seed the random library
	rand.Seed(time.Now().UTC().UnixNano())
	// Setup logger
	intializelogger(4)
}

// A function that initializes the log level, formatter and output
func intializelogger(loglevel logrus.Level) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the specified level or above.
	logrus.SetLevel(loglevel)
}
