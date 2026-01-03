package config

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Validate validates the configuration.
func Validate(config *Config) error {
	if config.Version == "" {
		return errors.New("version field is required")
	}

	if len(config.Profiles) == 0 {
		return errors.New("at least one profile must be defined")
	}

	// routesが定義されていない場合
	if len(config.Routes) == 0 {
		return errors.New("routes must have at least one route")
	}

	// ルート名のバリデーション
	for routeName, route := range config.Routes {
		// ルート名が空
		if routeName == "" {
			return errors.New("route name cannot be empty")
		}

		// ルート名が有効なファイル名か
		if !isValidFileName(routeName) {
			return fmt.Errorf("route name '%s' contains invalid characters. Use only alphanumeric, hyphens, and underscores", routeName)
		}

		// ルートが空
		if len(route) == 0 {
			return fmt.Errorf("route '%s' must have at least one step", routeName)
		}

		// プロファイル参照チェック
		for i, step := range route {
			if _, ok := config.Profiles[step.Profile]; !ok {
				return fmt.Errorf("route '%s': profile '%s' not found (step %d)", routeName, step.Profile, i)
			}
		}

		// 2段目以降のpassword_promptチェック
		for i, step := range route {
			if i == 0 {
				continue // 1段目はconnectコマンドを使用するためpassword_prompt不要
			}
			profile := config.Profiles[step.Profile]
			if profile.Auth.Type == "password" && profile.Auth.PasswordPrompt == "" {
				return fmt.Errorf("route '%s': profile '%s': password_prompt is required for password auth in route step %d", routeName, step.Profile, i+1)
			}
		}
	}

	// プロファイル設定チェック
	for name, profile := range config.Profiles {
		// prompt_marker必須チェック
		if profile.PromptMarker == "" {
			return fmt.Errorf("profile '%s': prompt_marker is required", name)
		}

		// 認証設定チェック
		if err := validateAuth(profile.Auth); err != nil {
			return fmt.Errorf("invalid auth in profile '%s': %w", name, err)
		}

		// keyfile認証でpassword_promptが設定されている場合はエラー
		if profile.Auth.Type == "keyfile" && profile.Auth.PasswordPrompt != "" {
			return fmt.Errorf("profile '%s': password_prompt should not be set for keyfile auth", name)
		}

		// password_promptにシングルクォートが含まれる場合はエラー（TTLインジェクション対策）
		if profile.Auth.PasswordPrompt != "" && strings.Contains(profile.Auth.PasswordPrompt, "'") {
			return fmt.Errorf("profile '%s': password_prompt cannot contain single quotes", name)
		}
	}

	// 注: options.auto_disconnectには明示的なバリデーションは不要です。
	// YAMLパーサー（gopkg.in/yaml.v3）が自動的にboolean型を検証し、
	// 不正な値（文字列、数値など）はこの地点に到達する前にパースエラーになります。

	return nil
}

func validateAuth(auth *Auth) error {
	if auth == nil {
		return errors.New("auth is required")
	}

	switch auth.Type {
	case "password":
		// Value, Env, Prompt のいずれか1つが必要
		if auth.Value == "" && auth.Env == "" && !auth.Prompt {
			return errors.New("password auth requires 'value', 'env', or 'prompt'")
		}
	case "keyfile":
		if auth.Path == "" {
			return errors.New("keyfile auth requires 'path'")
		}
	default:
		return fmt.Errorf("invalid auth type: %s (must be 'password' or 'keyfile')", auth.Type)
	}

	return nil
}

// isValidFileName はファイル名として有効な文字列かチェック
func isValidFileName(name string) bool {
	// 英数字、ハイフン、アンダースコアのみ許可
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched
}
