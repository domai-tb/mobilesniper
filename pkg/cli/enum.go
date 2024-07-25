package cli

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
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
		log.Printf("Scanning network: %s\n", cidr)

		var wg sync.WaitGroup
		defer wg.Wait()

		targetChan := make(chan models.Target)

		wg.Add(1)
		go enum.DiscoverOpenPorts(cidr, targetChan, &wg, maxConcurrency)

		for target := range targetChan {
			log.Println(target)
		}
	},
}

var nfsCmd = &cobra.Command{
	Use:   "nf <network range>",
	Short: "Probe network to identify specific network function",
	Long: `This command enumerates 3GPP OpenAPI definition to identifies 5G network function within 
			the given network range.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cidr := args[0]
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Fatalf("'%s' is not a valid CIDR network range", cidr)
		}

		openapiPath, _ := cmd.Flags().GetString("netfunc-openapi-path")
		fmt.Printf("Scanning network: %s\n", cidr)

		var wg sync.WaitGroup
		defer wg.Wait()

		targetChan := make(chan models.Target)
		nfrChan := make(chan models.NetworkFunctionResult)

		wg.Add(1)
		go enum.DiscoverOpenPorts(cidr, targetChan, &wg, maxConcurrency)

		for target := range targetChan {
			if utils.IsHTTPorHTTPS(target.IP, target.Port) {
				go func() {
					log.Printf("Start identifing NFs on %s:%d", target.IP, target.Port)
					err := filepath.Walk(openapiPath, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						openapi, err := models.ValidateOpenAPIFile(path)
						if err != nil {
							fmt.Printf("Invalid OpenAPI file: %s - %v\n", path, err)
							return nil // Continue processing other files
						}

						enum.TestOfNetworkFunction(target.IP, target.Port, openapi, nfrChan)

						return nil
					})

					if err != nil {
						log.Panicf("error walking the path %q: %v", openapiPath, err)
					}

					log.Printf("Finish identifing NFs on %s:%d", target.IP, target.Port)
				}()

				for nfr := range nfrChan {
					log.Println(nfr)
				}
			}
		}
	},
}

func init() {
	nfsCmd.Flags().String("openapi", "assets/5GC-APIs", "Path to 3GPP OpenAPI definitions of 5G network functions")

	enumCmd.AddCommand(servicesCmd)
	enumCmd.AddCommand(nfsCmd)
}
