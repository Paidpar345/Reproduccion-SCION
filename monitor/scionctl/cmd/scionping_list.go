package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// scionPingListCmd represents the list command
var scionPingListCmd = &cobra.Command{
	Use:   "list <node>",
	Args:  cobra.ExactArgs(1),
	Short: "List existing ping log files on a host",
	Long:  `No longer description available`,
	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleScionPingList(args)
	},
}

func init() {
	ScionPingCmd.AddCommand(scionPingListCmd)
}
