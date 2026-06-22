package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// pingListCmd represents the list command
var pingListCmd = &cobra.Command{
	Use:   "list <node>",
	Args:  cobra.ExactArgs(1),
	Short: "List existing ping log files on a host",
	Long:  `No longer description available`,
	Run: func(cmd *cobra.Command, args []string) {
		handler.HandlePingList(args)
	},
}

func init() {
	PingCmd.AddCommand(pingListCmd)
}
