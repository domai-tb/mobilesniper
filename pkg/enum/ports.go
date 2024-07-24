package enum

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/tomsteele/go-nmap"
)

const (
	maxConcurrency = 256 // Maximum number of concurrent goroutines
)

func DiscoverOpenPorts(targetNet string, targetChan chan<- Target, wg *sync.WaitGroup) (targets *[]Target, err error) {
	retVal := make([]Target, 0)
	defer wg.Done()

	ips, err := getIPsInCIDR(targetNet)
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

			cmd := exec.Command("nmap", "-T5", "-Pn", "--scan-delay", "0", "--max-scan-delay", "20ms", "-oX", "-", "-p-", "-sV", ip)
			var out bytes.Buffer
			cmd.Stdout = &out

			// log.Printf("Start scanning %s", ip)
			err := cmd.Run()
			if err != nil {
				log.Printf("error running nmap for %s: %v", ip, err)
				return
			}
			// log.Printf("Finish scanning %s", ip)

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
							retVal = append(retVal, Target{Host: ip, IP: ip, Port: int(port.PortId)})
							targetChan <- Target{Host: ip, IP: ip, Port: int(port.PortId), Protocol: port.Protocol, Service: port.Service.Name}
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
