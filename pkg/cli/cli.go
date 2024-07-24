package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mobilesniper",
	Short: "A pentesting tool for 5G mobile networks.",
	Long:  `MobileSniper is a CLI application for performing various pentesting tasks specialicied on 5G mobile networks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(enumCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// You can implement configuration loading here if needed
}
