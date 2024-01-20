/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "workito",
	Short: "Workito is a small CLI tool to send commands to a b7s worker node",
	Long: `Workito is a small CLI tool to send commands to a b7s worker node.
	
At the moment it supports instructing the worker node to install or to execute 
a function. It can be used to issue commands to the worker node without having 
a head node running alongside it, allowing simpler setups and easier testing.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	PeerAddress string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&PeerAddress, "address", "a", "", "multiaddress of the worker node")
}
