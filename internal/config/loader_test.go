package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Success(t *testing.T) {
	t.Run("loads simple configuration", func(t *testing.T) {
		cfg, err := LoadConfig("../../test/fixtures/valid/simple.yml")
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.NotEmpty(t, cfg.Version)
		assert.NotEmpty(t, cfg.Profiles)
		assert.NotEmpty(t, cfg.Routes)

		// デフォルト値が設定されていることを確認
		assert.NotNil(t, cfg.Options)
		assert.Equal(t, 30, cfg.Options.Timeout)

		// プロファイルのデフォルトポートが設定されていることを確認
		for _, profile := range cfg.Profiles {
			assert.NotZero(t, profile.Port)
		}
	})

	t.Run("loads full configuration", func(t *testing.T) {
		cfg, err := LoadConfig("../../test/fixtures/valid/full.yml")
		require.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.NotEmpty(t, cfg.Version)
		assert.NotEmpty(t, cfg.Profiles)
		assert.NotEmpty(t, cfg.Routes)

		// カスタム値が保持されていることを確認
		assert.NotNil(t, cfg.Options)
		assert.Equal(t, 60, cfg.Options.Timeout)
		assert.Equal(t, 3, cfg.Options.Retry)
		assert.True(t, cfg.Options.Log)

		// プロファイルのポートが設定されていることを確認
		for _, profile := range cfg.Profiles {
			assert.NotZero(t, profile.Port)
		}
	})
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestLoadConfig_SyntaxError(t *testing.T) {
	_, err := LoadConfig("../../test/fixtures/invalid/syntax-error.yml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML")
}

func TestLoadConfig_PreservesCustomValues(t *testing.T) {
	cfg, err := LoadConfig("../../test/fixtures/valid/full.yml")
	require.NoError(t, err)

	// カスタムポートが保持されていることを確認
	assert.Equal(t, 22, cfg.Profiles["bastion"].Port)
	assert.Equal(t, 2222, cfg.Profiles["target"].Port)

	// カスタムタイムアウトが保持されていることを確認
	assert.Equal(t, 60, cfg.Options.Timeout)
}
