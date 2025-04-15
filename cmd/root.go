package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "netbird-mgmt",
	Short:        "",
	Long:         "",
	Version:      "",
	SilenceUsage: true,
}

// Execute is the entry point for all commands.
func Execute() error {
	return rootCmd.Execute()
}
