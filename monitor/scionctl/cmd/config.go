/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config <command>",
	Short: "Manage SCION node configuration",
	Long:  `Configure various aspects of SCION nodes ISD lists and AS lists`,
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
}
