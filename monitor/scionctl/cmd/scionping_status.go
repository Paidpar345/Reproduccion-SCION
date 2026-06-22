/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// scionPingStatusCmd represents the list command
var scionPingStatusCmd = &cobra.Command{
	Use:   "status <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Status of ping process on node",
	Long:  `No longer description available`,
	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleScionPingStatus(args)
	},
}

func init() {
	ScionPingCmd.AddCommand(scionPingStatusCmd)
}
