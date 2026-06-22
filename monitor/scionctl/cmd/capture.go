// Keeping this in one file for now

package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

var CaptureCmd = &cobra.Command{
	Use:   "capture <command>",
	Short: "Manage tcpdump capture on a node",
	Long:  "No long description available atm",
}

var captureStartCmd = &cobra.Command{
	Use:   "start <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Starts a tcpdump capture on the on the specified node (eth0)", // Todo keep this updated
	Long:  "No long description available atm",

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleCaptureStart(args)
	},
}

var captureStopCmd = &cobra.Command{
	Use:   "stop <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Stops a running tcpdump capture on the specified node", // Todo keep this updated
	Long:  "No long description available atm",

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleCaptureStop(args)
	},
}

var captureStatusCmd = &cobra.Command{
	Use:   "start <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Starts a tcpdump capture on the default interface eth0", // Todo keep this updated
	Long:  "No long description available atm",

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleCaptureStatus(args)
	},
}

var captureListCmd = &cobra.Command{
	Use:   "list <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Returns available .pcap files on the node", // Todo keep this updated
	Long:  "No long description available atm",

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleCaptureList(args)
	},
}

func init() {
	RootCmd.AddCommand(CaptureCmd)
	CaptureCmd.AddCommand(captureStartCmd)
	CaptureCmd.AddCommand(captureStopCmd)
	CaptureCmd.AddCommand(captureStatusCmd)
	CaptureCmd.AddCommand(captureListCmd)
}
