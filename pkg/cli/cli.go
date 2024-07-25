package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var maxConcurrency int // maximum number of concurrent Go-routines

var rootCmd = &cobra.Command{
	Use:   "mobilesniper",
	Short: "A pentesting tool for 5G mobile networks.",
	Long: `MobileSniper is a CLI application for performing
		   various pentesting tasks specialicied on 5G mobile networks.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(enumCmd)

	rootCmd.PersistentFlags().IntVarP(
		&maxConcurrency, "max-goroutines", "c", 256, "Maximum number of concurrent Go-routines",
	)
}
