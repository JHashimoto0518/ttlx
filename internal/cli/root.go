package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ttlx",
	Short: "Tera Term Language eXtended - YAML to TTL generator",
	Long:  `ttlx generates Tera Term macro (TTL) scripts from YAML configuration files.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}
