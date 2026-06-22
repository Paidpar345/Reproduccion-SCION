/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// pingStatusCmd represents the list command
var pingStatusCmd = &cobra.Command{
	Use:   "status <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Status of ping process on node",
	Long:  `No longer description available`,
	Run: func(cmd *cobra.Command, args []string) {
		handler.HandlePingStatus(args)
	},
}

func init() {
	PingCmd.AddCommand(pingStatusCmd)
}
