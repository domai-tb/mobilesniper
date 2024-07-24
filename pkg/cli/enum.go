package cli

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/spf13/cobra"
)

var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Enumeration commands",
	Long:  `This command group contains commands related to enumeration.`,
}

var servicesCmd = &cobra.Command{
	Use:   "services <network range>",
	Short: "Perform a port scan on a given network",
	Long:  `This command performs a port scan with service discovery on a given network range by nmap`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cidr := args[0]
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("'%s' is not a valid CIDR network range", cidr)
		}
		fmt.Printf("Scanning network: %s\n", cidr)
		var wg sync.WaitGroup
		defer wg.Wait()

		targetChan := make(chan enum.Target)

		wg.Add(1)
		go enum.DiscoverOpenPorts(cidr, targetChan, &wg)

		go func() {
			for target := range targetChan {
				log.Println(target)
				// A HTTP(S) port indicates a network function throught theire REST-API design.
				// TODO: Identify network function based on swagger documentation
				// log.Printf("%s:%d speaks HTTP(S): %s", target.IP, target.Port, enum.IsHTTPorHTTPS(target.IP, target.Port))
			}
		}()
	},
}

func init() {
	enumCmd.AddCommand(servicesCmd)
}
