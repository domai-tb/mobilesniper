package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var maxConcurrency int // maximum number of concurrent Go-routines

var rootCmd = &cobra.Command{
	Use:   "mobilesniper",
	Short: "A pentesting tool for 5G mobile networks.",
	Long:  "MobileSniper is a CLI application for performing various pentesting tasks specialicied on 5G mobile networks.",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewProgressBar(count int, desc string) (*progressbar.ProgressBar, *ProgressLogger) {
	bar := progressbar.NewOptions(count,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription(fmt.Sprintf("[red]%s[green]", desc)),
		progressbar.OptionSpinnerType(21),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "][green]",
		}),
	)

	// Create the custom logger and set it as the default logger
	progressLogger := NewProgressLogger(bar)
	log.SetOutput(progressLogger)

	// Periodically refresh the progress bar to update the elapsed time
	go func() {
		for {
			time.Sleep(25 * time.Millisecond)
			bar.RenderBlank() // Refresh the progress bar without changing progress
		}
	}()

	return bar, progressLogger
}

func init() {
	rootCmd.AddCommand(enumCmd)

	rootCmd.PersistentFlags().IntVarP(
		&maxConcurrency, "max-goroutines", "c", 256, "Maximum number of concurrent Go-routines",
	)
}
