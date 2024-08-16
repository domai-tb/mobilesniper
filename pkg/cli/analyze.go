package cli

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyzing commands",
	Long:  `This command group contains commands related to offline analyzing.`,
}

var pcapCmd = &cobra.Command{
	Use:   "pcap",
	Short: "Analyze .pcap-file",
	Long:  "Analyze Wireshark traffic capturing (.pcap file) related to 5G communication.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		pcapFile := args[0]
		_, err := pcap.OpenOffline(pcapFile)

		if err != nil {
			panic(err)
		}

		bar, _ := NewProgressBar(1, fmt.Sprintf("Analyzing: %s", pcapFile))
		defer bar.Finish()

		log.Fatalln("PCAP Analyzing is currently unimplemented")
	},
}

func init() {
	analyzeCmd.Flags().String(
		"openapi", "assets/5GC-APIs", "Path to 3GPP OpenAPI definitions of 5G network functions",
	)

	analyzeCmd.AddCommand(pcapCmd)
}
