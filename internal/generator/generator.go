// Package generator generates TTL scripts from configuration data.
package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/JHashimoto0518/ttlx/internal/config"
)

const version = "1.0.0"

// Generate generates a TTL script from the configuration.
func Generate(cfg *config.Config, sourceFile string) (string, error) {
	var sb strings.Builder

	// ヘッダー生成
	sb.WriteString(generateHeader(sourceFile))

	// 変数定義生成
	sb.WriteString(generateVariables(cfg))

	// ルートステップごとの処理生成
	errorLabels := make([]string, 0)
	for i, step := range cfg.Route {
		profile := cfg.Profiles[step.Profile]
		upperProfileName := strings.ToUpper(step.Profile)

		if i == 0 {
			// 最初のステップ: connect コマンド
			sb.WriteString(generateConnect(i+1, step.Profile, upperProfileName, profile))

			// パスワード認証処理（connect コマンドに含まれていない場合のみ）
			if profile.Auth.Type == "password" && profile.Auth.Value == "" {
				sb.WriteString(generatePasswordAuth(step.Profile, profile.Auth))
			}
		} else {
			// 2番目以降のステップ: ssh コマンド
			sb.WriteString(generateSSH(i+1, step.Profile, upperProfileName, profile))

			// パスワード認証処理
			if profile.Auth.Type == "password" {
				sb.WriteString(generatePasswordAuth(step.Profile, profile.Auth))
			}
		}

		// コマンド実行
		if len(step.Commands) > 0 {
			sb.WriteString(generateCommands(step.Commands, profile.PromptMarker, upperProfileName))
		}

		// エラーラベルを記録
		errorLabels = append(errorLabels, upperProfileName)
	}

	// 成功終了
	sb.WriteString(successTemplate)

	// エラーハンドリング生成
	sb.WriteString(generateErrorHandling(errorLabels, cfg.Route))

	return sb.String(), nil
}

func generateHeader(sourceFile string) string {
	now := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf(headerTemplate, version, sourceFile, now)
}

func generateVariables(cfg *config.Config) string {
	timeout := 30
	if cfg.Options != nil && cfg.Options.Timeout > 0 {
		timeout = cfg.Options.Timeout
	}
	return fmt.Sprintf(variablesTemplate, timeout)
}

func generateConnect(stepNum int, profileName, upperProfileName string, profile *config.Profile) string {
	authType := profile.Auth.Type
	keyfileOption := ""
	passwordOption := ""

	if authType == "keyfile" {
		keyfileOption = fmt.Sprintf(" /keyfile=%s", profile.Auth.Path)
	} else if authType == "password" && profile.Auth.Value != "" {
		// パスワードが直接指定されている場合は connect コマンドに含める
		passwordOption = fmt.Sprintf(" /passwd=%s", profile.Auth.Value)
	}

	return fmt.Sprintf(
		connectTemplate,
		stepNum,
		profileName,
		upperProfileName,
		profile.Host,
		profile.Port,
		authType,
		profile.User,
		keyfileOption,
		passwordOption,
		upperProfileName,
		profile.PromptMarker,
		upperProfileName,
	)
}

func generateSSH(stepNum int, profileName, upperProfileName string, profile *config.Profile) string {
	return fmt.Sprintf(
		sshTemplate,
		stepNum,
		profileName,
		profile.User,
		profile.Host,
		profile.Port,
		upperProfileName,
	)
}

func generatePasswordAuth(profileName string, auth *config.Auth) string {
	if auth.Env != "" {
		return fmt.Sprintf(passwordEnvTemplate, auth.Env)
	}
	if auth.Prompt {
		return fmt.Sprintf(passwordPromptTemplate, profileName)
	}
	if auth.Value != "" {
		return fmt.Sprintf(passwordValueTemplate, auth.Value)
	}
	return ""
}

func generateCommands(commands []string, prompt, upperProfileName string) string {
	var sb strings.Builder
	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf(commandTemplate, cmd, cmd, prompt, upperProfileName))
	}
	return sb.String()
}

func generateErrorHandling(errorLabels []string, route []*config.RouteStep) string {
	var sb strings.Builder

	for i, label := range errorLabels {
		profileName := route[i].Profile

		// 接続エラー（最初のステップのみ）
		if i == 0 {
			sb.WriteString(fmt.Sprintf(errorConnectTemplate, label, profileName))
		}

		// タイムアウトエラー
		sb.WriteString(fmt.Sprintf(errorTimeoutTemplate, label, profileName))
	}

	// クリーンアップ
	sb.WriteString(cleanupTemplate)

	return sb.String()
}
