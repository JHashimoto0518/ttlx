// Package config handles YAML configuration file loading, validation, and data models.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
	// ファイル読み込み
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// YAMLパース
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// デフォルト値設定
	cfg.SetDefaults()

	return &cfg, nil
}
