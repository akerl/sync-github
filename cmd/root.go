package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "syncgithub",
	Short:         "Tool for managing GitHub clones",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute function is the entrypoint for the CLI
func Execute() error {
	return rootCmd.Execute()
}
