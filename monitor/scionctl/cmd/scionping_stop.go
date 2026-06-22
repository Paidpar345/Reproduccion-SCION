/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// ScionPingStopCmd represents the start ping command
var ScionPingStopCmd = &cobra.Command{
	Use:   "stop <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Stop a scion ping operation at a node",
	Long:  `Stop a scion ping operation at a SCION node`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleScionPingStop(args)
	},
}

func init() {
	ScionPingCmd.AddCommand(ScionPingStopCmd)
}
