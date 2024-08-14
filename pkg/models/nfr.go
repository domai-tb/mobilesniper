package models

import "fmt"

type NetworkFunctionResult struct {
	Target          Target
	NetworkFunction string
	Accuracy        float64
}

func (nfr *NetworkFunctionResult) String() string {
	return fmt.Sprintf("- %v - %s (%f%%)", nfr.Target, nfr.NetworkFunction, nfr.Accuracy)
}
