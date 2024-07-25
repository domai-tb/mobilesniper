package models

import "fmt"

type NetworkFunctionResult struct {
	IP              string
	Port            int
	NetworkFunction string
	Accuracy        float64
}

func (nfr *NetworkFunctionResult) String() string {
	return fmt.Sprintf("%s:%d - %s (%f%%)", nfr.IP, nfr.Port, nfr.NetworkFunction, nfr.Accuracy)
}
