package enum

import (
	"encoding/xml"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/awareseven/mobilesniper/pkg/models"
	soap "github.com/awareseven/mobilesniper/pkg/models/soap"
	"github.com/awareseven/mobilesniper/pkg/utils"
)

var sdcCmd = &cobra.Command{
	Use:   "sdc <network interface>",
	Short: "Enumerate SDC devices.",
	Long:  `This command performs a SDC discovery process via UDP multicast to discover medical devices.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ifi := args[0]
		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover SDC Devices: %s", ifi))
		defer bar.Finish()

		var wg sync.WaitGroup
		sdcChan := make(chan soap.ProbeMatch)

		wg.Add(1)
		go enum.DiscoverSDCDevices(ifi, sdcChan, &wg, core.Verbose)

		go func() {
			wg.Wait()
			close(sdcChan) // Ensure channel closure after all operations are complete
		}()

		for probeMatch := range sdcChan {
			bar.ChangeMax(bar.GetMax() + 1)

			// call HTTP server with a "Get" message that should return device information
			soapGetMsg := soap.NewSOAPMessage(
				"http://schemas.xmlsoap.org/ws/2004/09/transfer/Get",
				probeMatch.GetXAddrs(),
				soap.NewGetSOAPBody(),
			)

			// Marshal the Probe message to XML
			xmlBytes, err := soapGetMsg.XMLMarshal()
			if err != nil {
				log.Println(err)
				continue // continue with next received probe match
			}

			// actually perform the HTTP request with given porbe match address
			_, bodyBytes, err := utils.DoSdcXClientSOAPPost(probeMatch.GetXAddrs(), xmlBytes, core.Verbose)
			if err != nil {
				log.Println(err) // on mTLS enforced devices, this will return a mTLS mismatch error
				continue         // continue with next received probe match
			}

			// try to parse the returned SOAP message
			getResp := &soap.GetResponse{}
			err = xml.Unmarshal(bodyBytes, getResp)
			if err != nil {
				log.Println(err)
				continue // continue with next received probe match
			}

			// TODO: SDC communication to perform (more) malicious actions
			// -> successfully connection indicates that the device doesn't enfore mTLS
			// -> without mTLS there is no other authentication method (as soon as we know)
			device := models.CreateSDCDevicebyGetResponse(*getResp)
			log.Printf("Endpoint states it is %s", device.String())

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}
