/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

// https://github.com/spf13/cobra-cli/blob/main/README.md

package cmd

import (
	"fmt"
	"os"

	"scionctl/internal/config"
	"scionctl/internal/pprinter"

	"github.com/spf13/cobra"
)

var (
	// Global manager variable for nodes
	CmdNodeManager *config.NodeManager
	cfgFile        string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "scionctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scionctl called")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/home/nodeconfig.yaml", "config file (default is /home/nodeconfig.yaml)")


	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {

	// If a config file is found, read it in.
	if cfgFile != "" {
		// Use config file from the flag.
		err := config.InitializeManager(cfgFile)

		if err != nil {
			fmt.Printf("Error fetching config: %s\n", cfgFile)
			pprinter.PrintError(err)
			os.Exit(1)
		}
	}
}
