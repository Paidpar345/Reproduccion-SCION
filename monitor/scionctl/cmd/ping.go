/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// startping represents the ping command
var PingCmd = &cobra.Command{
	Use:   "ping <command>",
	Short: "Manage ping operations between nodes",
	Long:  `Start, stop, and monitor ping operations between SCION nodes`,

	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func init() {
	RootCmd.AddCommand(PingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
