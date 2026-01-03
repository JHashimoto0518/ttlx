// Package cli provides the command-line interface for ttlx.
package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
		ttls, err := generator.GenerateAll(cfg, filepath.Base(configPath))
		if err != nil {
			return fmt.Errorf("failed to generate TTL: %w", err)
		}

		// 4. 出力
		if dryRun {
			// dry-runモード: 各ルートのTTLを区切って表示
			routeNames := make([]string, 0, len(ttls))
			for routeName := range ttls {
				routeNames = append(routeNames, routeName)
			}
			sort.Strings(routeNames)

			for _, routeName := range routeNames {
				fmt.Printf("=== %s.ttl ===\n", routeName)
				fmt.Println(ttls[routeName])
				fmt.Println()
			}
			return nil
		}

		// 出力先ディレクトリの決定
		outputDir := "."
		if outputPath != "" {
			outputDir = outputPath
		}

		// ディレクトリが存在しない場合は作成
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// 各ルートのTTLファイルを書き込み
		var generatedFiles []string
		for routeName, ttl := range ttls {
			filename := filepath.Join(outputDir, routeName+".ttl")
			if err := os.WriteFile(filename, []byte(ttl), 0644); err != nil {
				return fmt.Errorf("failed to write TTL file '%s': %w", filename, err)
			}
			generatedFiles = append(generatedFiles, filename)
		}

		// 生成されたファイル一覧を表示
		sort.Strings(generatedFiles)
		fmt.Println("Generated TTL files:")
		for _, file := range generatedFiles {
			fmt.Printf("  - %s\n", file)
		}

		return nil
	},
}

func init() {
	buildCmd.Flags().StringP("output", "o", "", "Output directory path")
	buildCmd.Flags().Bool("dry-run", false, "Print to stdout instead of file")
}
