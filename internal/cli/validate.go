package cli

import (
	"fmt"

	"github.com/JHashimoto0518/ttlx/internal/config"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <config.yml>",
	Short: "Validate YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		configPath := args[0]

		// 1. 設定読み込み
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// 2. バリデーション
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Println("✓ Validation passed")
		return nil
	},
}
