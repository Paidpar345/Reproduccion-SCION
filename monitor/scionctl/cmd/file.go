/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

var (
	src string
)

// startping represents the ping command
func FileCmd(parentType string) *cobra.Command {
	return &cobra.Command{
		Use:   "file <node> <name>",
		Args:  cobra.ExactArgs(2),
		Short: "Fetch a file from a node",
		Long:  `Returns a file in stdout that can be piped to a new file, e.g. <command> > newfile.pcap`,

		Run: func(cmd *cobra.Command, args []string) {
			src = parentType
			handler.HandleGetFile(args, src)
		},
	}

}

func init() {
	ScionPingCmd.AddCommand(FileCmd("scionping"))
	PingCmd.AddCommand(FileCmd("ping"))
	CaptureCmd.AddCommand(FileCmd("capture"))
}
