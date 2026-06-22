/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

var (
	scionpingCount int
)

// PingStartCmd represents the start ping command
var scionPingStartCmd = &cobra.Command{
	Use:   "start <source> <target>",
	Args:  cobra.ExactArgs(2),
	Short: "Start a ping operation between two nodes",
	Long:  `Initiate a ping operation between two SCION nodes`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleScionPingStart(args, scionpingCount)
	},
}

func init() {
	ScionPingCmd.AddCommand(scionPingStartCmd)
	scionPingStartCmd.Flags().IntVarP(&scionpingCount, "count", "c", -1, "Number of scion ping packets to send")
}
