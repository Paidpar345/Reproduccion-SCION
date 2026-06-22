package cmd

import (
	"scionctl/cmd/handler"

	"github.com/spf13/cobra"
)

// configAsListCmd represents the aslist subcommand
var configAsListCmd = &cobra.Command{
	Use:   "aslist <node> <AS1> <AS2> ... | delete",
	Args:  cobra.MinimumNArgs(2),
	Short: "Configure AS blacklist for a SCION node",
	Long: `Configure the AS blacklist for a SCION node's path policy.
	
Examples:
  # Blacklist specific ASes (provide AS numbers only)
  scionctl config aslist scion11 12 13 15
  
  # Clear the blacklist
  scionctl config aslist scion11 delete`,

	Run: func(cmd *cobra.Command, args []string) {
		handler.HandleConfigASList(args)
	},
}

func init() {
	ConfigCmd.AddCommand(configAsListCmd)
}
