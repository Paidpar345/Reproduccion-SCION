package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// configFileCmd represents the file subcommand
var configFileCmd = &cobra.Command{
	Use:   "file <node>",
	Args:  cobra.ExactArgs(1),
	Short: "Retrieve the path-policy configuration file from a SCION node",
	Long: `Retrieve and display the path-policy configuration file from a SCION node.
	
Examples:
  # Get the path-policy file from a node
  scionctl config file scion11`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleConfigFile(args)
	},
}

func init() {
	ConfigCmd.AddCommand(configFileCmd)
}
