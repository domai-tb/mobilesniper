package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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
	Use:   "services <network range or single IP>",
	Short: "Perform a port scan on a given network",
	Long:  `This command performs a port scan with service discovery on a given network range by nmap`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		_, err := utils.GetIPsInCIDR(cidrOrIP)
		if err != nil {
			panic(err)
		}

		bar, _ := NewProgressBar(1, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup
		defer wg.Wait()

		targetChan := make(chan models.Target)

		wg.Add(1)
		go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &wg, maxConcurrency)

		for target := range targetChan {
			bar.ChangeMax(bar.GetMax() + 1)
			log.Println(target)

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}

var nfsCmd = &cobra.Command{
	Use:   "nf <network range>",
	Short: "Probe network to identify specific network function",
	Long:  "This command enumerates 3GPP OpenAPI definition to identifies 5G network function within the given network range or IP address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		openapiPath, _ := cmd.Flags().GetString("openapi")

		files, _ := os.ReadDir(openapiPath)
		bar, _ := NewProgressBar(0, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup
		defer wg.Wait()

		targetChan := make(chan models.Target)
		nfrChan := make(chan models.NetworkFunctionResult)

		mayIpAndPort := strings.Split(cidrOrIP, ":")
		if len(mayIpAndPort) == 2 {
			port, _ := strconv.Atoi(mayIpAndPort[1])
			target := models.Target{
				IP:       mayIpAndPort[0],
				Port:     port,
				Protocol: "tcp",
			}

			go func() {
				targetChan <- target
				close(targetChan)
			}()

		} else {
			_, err := utils.GetIPsInCIDR(cidrOrIP)
			if err != nil {
				panic(err)
			}

			wg.Add(1)
			go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &wg, maxConcurrency)
		}

		semaphore := make(chan struct{}, maxConcurrency)

		for target := range targetChan {
			bar.ChangeMax(bar.GetMax() + len(files))

			if utils.IsHTTPorHTTPS(target.IP, target.Port) {

				wg.Add(1)
				semaphore <- struct{}{} // add to channel

				go func() {
					defer wg.Done()
					defer func() { <-semaphore }() // remove from channel

					log.Printf("Start identifing NFs on %s:%d", target.IP, target.Port)

					err := filepath.Walk(openapiPath, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}

						openapi, err := models.ValidateOpenAPIFile(path)
						if err != nil {
							log.Printf("Invalid OpenAPI file: %s - %v\n", path, err)
							return nil // Continue processing other files
						}

						enum.TestOfNetworkFunction(target.IP, target.Port, openapi, nfrChan)

						return nil
					})

					if err != nil {
						log.Panicf("error walking the path %q: %v", openapiPath, err)
					}

					log.Printf("Finish identifing NFs on %s:%d", target.IP, target.Port)

					go func() {
						wg.Wait()
						close(nfrChan)
					}()
				}()

				for nfr := range nfrChan {
					if nfr.Accuracy > 0 {
						log.Println(nfr.String())
					}

					time.Sleep(100 * time.Millisecond)
					bar.Add(1)
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
