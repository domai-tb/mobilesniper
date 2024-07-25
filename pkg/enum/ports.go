package enum

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/tomsteele/go-nmap"

	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

func DiscoverOpenPorts(targetNet string, targetChan chan<- models.Target, wg *sync.WaitGroup, maxConcurrency int,
) (*[]models.Target, error) {
	retVal := make([]models.Target, 0)
	defer wg.Done()

	ips, err := utils.GetIPsInCIDR(targetNet)
	if err != nil {
		return nil, fmt.Errorf("%s is not a network in valid CIDR-notation", targetNet)
	}

	semaphore := make(chan struct{}, maxConcurrency)

	for _, ip := range ips {
		wg.Add(1)
		semaphore <- struct{}{} // add to channel

		go func(ip string) {
			defer wg.Done()
			defer func() { <-semaphore }() // remove from channel

			// log.Printf("Start scanning %s", ip)
			cmd := exec.Command(
				"nmap", "--scan-delay", "0", "--max-scan-delay", "20ms",
				"-T5", "-Pn", "-oX", "-", "-p-", "-sV", ip,
			)
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			if err != nil {
				log.Printf("error running nmap for %s: %v", ip, err)
				return
			}
			log.Printf("Finish scanning %s", ip)

			var nmapRun nmap.NmapRun
			err = xml.Unmarshal(out.Bytes(), &nmapRun)
			if err != nil {
				log.Printf("error parsing nmap XML output for %s: %v", ip, err)
				return
			}

			for _, host := range nmapRun.Hosts {
				if len(host.Ports) > 0 {
					for _, port := range host.Ports {
						if port.State.State == "open" {
							retVal = append(retVal, models.Target{IP: ip, Port: int(port.PortId)})
							targetChan <- models.Target{
								IP: ip, Port: int(port.PortId), Protocol: port.Protocol,
								Service: models.Service{
									Name: port.Service.Name, Version: port.Service.Version,
									Product: port.Service.Product,
								}}
						}
					}
				}
			}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(targetChan)
	}()

	return &retVal, err
}
