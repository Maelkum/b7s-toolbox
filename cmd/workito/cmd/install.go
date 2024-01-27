/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a function",
	Long:  `Instruct the worker node to install the function(s) based on the specified CIDs.`,
	Example: `  workito install --address <node-multiaddress> <CID>
  workito install --address <node-multiaddress> <CID1> <CID2> <CID2>`,
	Run: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
}
