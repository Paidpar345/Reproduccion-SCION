/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// pingStopCmd represents the start ping command
var pingStopCmd = &cobra.Command{
	Use:   "stop <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Stop a ping operation at a node",
	Long:  `Stop a ping operation at a SCION node`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandlePingStop(args)
	},
}

func init() {
	PingCmd.AddCommand(pingStopCmd)
}
