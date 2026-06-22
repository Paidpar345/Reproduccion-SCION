package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// configISDListCmd represents the isdlist subcommand
var configIsdListCmd = &cobra.Command{
	Use:   "isdlist <node> <ISD16> <ISD17> ... | delete",
	Args:  cobra.MinimumNArgs(2),
	Short: "Configure ISD blacklist for a SCION node",
	Long: `Configure the ISD blacklist for a SCION node's path policy.
	
Examples:
  # Blacklist specific ISDs (provide ISD numbers only)
  scionctl config isdlist scion11 16
  
  # Clear the blacklist
  scionctl config isdlist scion11 delete`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleConfigISDList(args)
	},
}

func init() {
	ConfigCmd.AddCommand(configIsdListCmd)
}