/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "Execute a function",
	Long:  `Execute sends an execution request to the worker node and waits for a result.`,
	Run:   runExecute,
}

type ExecuteFlags struct {
	FunctionID string
	Method     string
	Entry      string

	Parameters      []string
	EnvironmentVars []string

	Nodes       uint
	Stdin       string
	Permissions []string
}

var ExecuteCmdFlags ExecuteFlags

func init() {
	rootCmd.AddCommand(executeCmd)

	flags := executeCmd.Flags()

	flags.StringVarP(&ExecuteCmdFlags.FunctionID, "function-id", "f", "", "function to execute")
	flags.StringVarP(&ExecuteCmdFlags.Method, "method", "m", "", "method to use")
	flags.StringVarP(&ExecuteCmdFlags.Entry, "entry", "t", "", "entry point")

	flags.StringSliceVarP(&ExecuteCmdFlags.EnvironmentVars, "env-var", "e", nil, "environment variables to set, in name=val format")
	flags.StringSliceVarP(&ExecuteCmdFlags.Parameters, "parameter", "p", nil, "parameters to set")
	flags.StringSliceVarP(&ExecuteCmdFlags.Permissions, "permissions", "r", nil, "permissions to set")

	flags.UintVarP(&ExecuteCmdFlags.Nodes, "nodes", "n", 1, "how many nodes should execute this request")
	flags.StringVar(&ExecuteCmdFlags.Stdin, "stdin", "", "standard input for the execution")
}
