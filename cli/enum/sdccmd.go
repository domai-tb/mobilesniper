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

		var wg sync.WaitGroup
		sdcChan := make(chan soap.ProbeMatch)

		wg.Add(1)
		go enum.DiscoverSDCDevices(ifi, sdcChan, &wg, core.Verbose)

		go func() {
			wg.Wait()
			close(sdcChan) // Ensure channel closure after all operations are complete
		}()

		for probeMatch := range sdcChan {

			soapGetMsg := soap.NewSOAPMessage(
				"http://schemas.xmlsoap.org/ws/2004/09/transfer/Get",
				probeMatch.GetXAddrs(),
				soap.NewGetSOAPBody(),
			)
			// Marshal the Probe message to XML
			xmlBytes, err := soapGetMsg.XMLMarshal()
			if err != nil {
				log.Fatalln(err)
			}

			_, bodyBytes, err := utils.DoSdcXClientSOAPPost(probeMatch.GetXAddrs(), xmlBytes, core.Verbose)
			if err != nil {
				log.Fatalln(err)
			}

			getResp := &soap.GetResponse{}
			err = xml.Unmarshal(bodyBytes, getResp)
			if err != nil {
				log.Fatalln(err)
			}

			device := models.CreateSDCDevicebyGetResponse(*getResp)
			log.Printf("Endpoint states it is %s", device.String())

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}

		bar.Finish() // Ensure progress bar finishes after all operations
	},
}
