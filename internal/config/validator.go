package config

import (
	"errors"
	"fmt"
)

// Validate validates the configuration.
func Validate(config *Config) error {
	if config.Version == "" {
		return errors.New("version field is required")
	}

	if len(config.Profiles) == 0 {
		return errors.New("at least one profile must be defined")
	}

	if len(config.Route) == 0 {
		return errors.New("route must have at least one step")
	}

	// プロファイル参照チェック
	for i, step := range config.Route {
		if _, ok := config.Profiles[step.Profile]; !ok {
			return fmt.Errorf("profile '%s' not found (route step %d)", step.Profile, i)
		}
	}

	// 2段目以降のルートステップでパスワード認証の場合、password_prompt必須
	for i, step := range config.Route {
		if i == 0 {
			continue // 1段目はconnectコマンドを使用するためpassword_prompt不要
		}
		profile := config.Profiles[step.Profile]
		if profile.Auth.Type == "password" && profile.Auth.PasswordPrompt == "" {
			return fmt.Errorf("profile '%s': password_prompt is required for password auth in route step %d", step.Profile, i+1)
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
	}

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
