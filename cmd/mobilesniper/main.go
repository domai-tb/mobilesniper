package main

import (
	"log"
	"sync"

	"github.com/awareseven/mobilesniper/pkg/enum"
)

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	targetChan := make(chan enum.Target)

	wg.Add(1)
	go enum.DiscoverOpenPorts("10.13.37.0/24", targetChan, &wg)

	go func() {
		for target := range targetChan {
			log.Printf("Found open port: %s:%d\n", target.IP, target.Port)
			// A HTTP(S) port indicates a network function throught theire REST-API design.
			// TODO: Identify network function based on swagger documentation
			// log.Printf("%s:%d speaks HTTP(S): %s", target.IP, target.Port, enum.IsHTTPorHTTPS(target.IP, target.Port))
		}
	}()
}
