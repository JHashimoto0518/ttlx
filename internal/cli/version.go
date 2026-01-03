package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0-beta"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("ttlx version %s\n", version)
	},
}
