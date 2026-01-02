// Package cli provides the command-line interface for ttlx.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JHashimoto0518/ttlx/internal/config"
	"github.com/JHashimoto0518/ttlx/internal/generator"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <config.yml>",
	Short: "Generate TTL script from YAML configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := args[0]
		outputPath, err := cmd.Flags().GetString("output")
		if err != nil {
			return fmt.Errorf("failed to get output flag: %w", err)
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return fmt.Errorf("failed to get dry-run flag: %w", err)
		}

		// 1. 設定読み込み
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// 2. バリデーション
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		// 3. TTL生成
		ttl, err := generator.Generate(cfg, filepath.Base(configPath))
		if err != nil {
			return fmt.Errorf("failed to generate TTL: %w", err)
		}

		// 4. 出力
		if dryRun {
			fmt.Println(ttl)
			return nil
		}

		if outputPath == "" {
			outputPath = strings.TrimSuffix(configPath, filepath.Ext(configPath)) + ".ttl"
		}

		if err := os.WriteFile(outputPath, []byte(ttl), 0o644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		fmt.Printf("Generated: %s\n", outputPath)
		return nil
	},
}

func init() {
	buildCmd.Flags().StringP("output", "o", "", "Output file path")
	buildCmd.Flags().Bool("dry-run", false, "Print to stdout instead of file")
}
