package enum

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/enum"
)

var sdcCmd = &cobra.Command{
	Use:   "sdc <network interface>",
	Short: "Enumerate SDC devices.",
	Long:  `This command performs a SDC discovery process via UDP multicast to discover medical devices.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ifi := args[0]
		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover SDC Devices: %s", ifi))

		enum.DiscoverSDCDevices(ifi, core.Verbose)

		bar.Finish()
	},
}
