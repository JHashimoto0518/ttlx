# 初回実装の設計書

## 概要

本ドキュメントでは、ttlx の初回実装（Phase 1 MVP）における具体的な設計を定義します。パッケージ構成、データモデル、コンポーネント設計、処理フローを明確にし、実装の指針とします。

---

## パッケージ構成

### 実装するパッケージ（Phase 1）

```
ttlx/
├── cmd/
│   └── ttlx/
│       └── main.go                # エントリーポイント
│
├── internal/
│   ├── config/                    # 設定ファイル処理
│   │   ├── model.go              # データモデル定義
│   │   ├── loader.go             # YAML読み込み
│   │   ├── validator.go          # バリデーション
│   │   ├── model_test.go
│   │   ├── loader_test.go
│   │   └── validator_test.go
│   │
│   ├── generator/                 # TTL生成
│   │   ├── generator.go          # メイン生成ロジック
│   │   ├── template.go           # テンプレート管理
│   │   ├── generator_test.go
│   │   └── template_test.go
│   │
│   ├── cli/                       # CLIコマンド
│   │   ├── root.go               # ルートコマンド
│   │   ├── build.go              # buildコマンド
│   │   ├── validate.go           # validateコマンド
│   │   ├── version.go            # versionコマンド
│   │   └── cli_test.go
│   │
│   └── errors/                    # エラー定義
│       ├── errors.go             # カスタムエラー型
│       └── messages.go           # エラーメッセージ
│
└── test/
    ├── fixtures/
    │   ├── valid/
    │   │   ├── simple.yml        # 最小構成
    │   │   └── full.yml          # 全機能使用
    │   ├── invalid/
    │   │   ├── missing-version.yml
    │   │   ├── invalid-profile.yml
    │   │   └── syntax-error.yml
    │   └── expected/
    │       ├── simple.ttl
    │       └── full.ttl
    └── integration/
        └── build_test.go
```

---

## データモデル設計

### Go 構造体定義

#### `internal/config/model.go`

```go
package config

// Config represents the entire YAML configuration.
type Config struct {
    Version  string              `yaml:"version"`
    Profiles map[string]*Profile `yaml:"profiles"`
    Route    []*RouteStep        `yaml:"route"`
    Options  *Options            `yaml:"options,omitempty"`
}

// Profile represents an SSH connection profile.
type Profile struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port,omitempty"` // デフォルト: 22
    User string `yaml:"user"`
    Auth *Auth  `yaml:"auth"`
}

// Auth represents authentication settings.
type Auth struct {
    Type   string `yaml:"type"` // "password" | "keyfile"
    Value  string `yaml:"value,omitempty"`  // パスワード直接記述
    Env    string `yaml:"env,omitempty"`    // 環境変数名
    Prompt bool   `yaml:"prompt,omitempty"` // 実行時入力
    Path   string `yaml:"path,omitempty"`   // 秘密鍵ファイルパス
}

// RouteStep represents a step in the connection route.
type RouteStep struct {
    Profile  string   `yaml:"profile"`
    Commands []string `yaml:"commands,omitempty"`
}

// Options represents global options.
type Options struct {
    Timeout int    `yaml:"timeout,omitempty"`
    Retry   int    `yaml:"retry,omitempty"`
    Log     bool   `yaml:"log,omitempty"`
    LogFile string `yaml:"log_file,omitempty"`
}
```

#### デフォルト値の設定

```go
// SetDefaults sets default values for the config.
func (c *Config) SetDefaults() {
    for _, profile := range c.Profiles {
        if profile.Port == 0 {
            profile.Port = 22
        }
    }

    if c.Options == nil {
        c.Options = &Options{}
    }
    if c.Options.Timeout == 0 {
        c.Options.Timeout = 30
    }
}
```

---

## コンポーネント設計

### 1. Config Loader (`internal/config/loader.go`)

**責務**: YAML設定ファイルの読み込みと解析

**公開関数**:
```go
// LoadConfig loads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error)
```

**実装ステップ**:
1. ファイル読み込み（`os.ReadFile`）
2. YAMLパース（`yaml.Unmarshal`）
3. デフォルト値の設定
4. 基本的なバリデーション呼び出し

**エラーハンドリング**:
- ファイル不存在: `os.IsNotExist` でチェック、分かりやすいエラーメッセージ
- YAML構文エラー: `yaml` パッケージのエラーをラップ、行番号を含める

---

### 2. Validator (`internal/config/validator.go`)

**責務**: 設定の妥当性検証

**公開関数**:
```go
// Validate validates the configuration.
func Validate(config *Config) error
```

**バリデーションルール**:
1. `version` フィールドが存在するか
2. `profiles` に少なくとも1つのプロファイルがあるか
3. `route` に少なくとも1つのステップがあるか
4. 各ルートステップの `profile` が存在するか
5. 認証設定が妥当か（`validateAuth`）
6. ホスト名が妥当か（基本的なチェック）

**実装例**:
```go
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

    // 認証設定チェック
    for name, profile := range config.Profiles {
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
        // 平文パスワード警告（実装時に追加）
    case "keyfile":
        if auth.Path == "" {
            return errors.New("keyfile auth requires 'path'")
        }
    default:
        return fmt.Errorf("invalid auth type: %s (must be 'password' or 'keyfile')", auth.Type)
    }

    return nil
}
```

---

### 3. TTL Generator (`internal/generator/generator.go`)

**責務**: Config から TTL スクリプトを生成

**公開関数**:
```go
// Generate generates a TTL script from the configuration.
func Generate(config *config.Config, sourceFile string) (string, error)
```

**実装ステップ**:
1. ヘッダーコメント生成
2. 変数定義生成
3. 各ルートステップの処理生成
   - 接続処理
   - コマンド実行処理
4. エラーハンドリングコード生成
5. 正常終了・異常終了処理生成

**TTL構造**:
```ttl
; ========================================
; Generated by ttlx v1.0.0
; Source: config.yml
; Generated at: 2026-01-01 10:00:00
; ========================================

; === 変数定義 ===
timeout = 30

; === ステップ1: bastion ===
:CONNECT_BASTION
connect 'bastion.example.com:22 /ssh /auth=password /user=user1'
if result <> 0 then
    goto :ERROR_CONNECT_BASTION
endif
wait '$' timeout
if result = 0 then
    goto :TIMEOUT_BASTION
endif

; パスワード認証の場合
passwordbox 'Enter password for bastion:' 'password'
sendln password

; === ステップ1: コマンド実行 ===
sendln 'su - root'
wait '#' timeout

; === ステップ2: target ===
sendln 'ssh user2@10.0.0.50'
wait 'password:' timeout
if result = 0 then
    goto :TIMEOUT_TARGET
endif
; ...

:SUCCESS
closett
end

:ERROR_CONNECT_BASTION
messagebox 'Failed to connect to bastion' 'Error'
goto :CLEANUP

:TIMEOUT_BASTION
messagebox 'Connection timeout: bastion' 'Error'
goto :CLEANUP

:CLEANUP
closett
end
```

---

### 4. Template Manager (`internal/generator/template.go`)

**責務**: TTLコード片のテンプレート管理

**テンプレート例**:
```go
const (
    // ヘッダーテンプレート
    headerTemplate = `; ========================================
; Generated by ttlx %s
; Source: %s
; Generated at: %s
; ========================================
`

    // 接続テンプレート（最初のステップ）
    connectTemplate = `; === ステップ%d: %s ===
:CONNECT_%s
connect '%s:%d /ssh /auth=%s /user=%s'
if result <> 0 then
    goto :ERROR_CONNECT_%s
endif
wait '$' timeout
if result = 0 then
    goto :TIMEOUT_%s
endif
`

    // SSH コマンドテンプレート（2番目以降のステップ）
    sshTemplate = `; === ステップ%d: %s ===
sendln 'ssh %s@%s'
wait 'password:' timeout
if result = 0 then
    goto :TIMEOUT_%s
endif
`
)
```

---

### 5. CLI Commands (`internal/cli/`)

#### `root.go`
```go
package cli

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "ttlx",
    Short: "Tera Term Language eXtended - YAML to TTL generator",
    Long:  `ttlx generates Tera Term macro (TTL) scripts from YAML configuration files.`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(buildCmd)
    rootCmd.AddCommand(validateCmd)
    rootCmd.AddCommand(versionCmd)
}
```

#### `build.go`
```go
var buildCmd = &cobra.Command{
    Use:   "build <config.yml>",
    Short: "Generate TTL script from YAML configuration",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        configPath := args[0]
        outputPath, _ := cmd.Flags().GetString("output")
        dryRun, _ := cmd.Flags().GetBool("dry-run")

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
        ttl, err := generator.Generate(cfg, configPath)
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

        if err := os.WriteFile(outputPath, []byte(ttl), 0644); err != nil {
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
```

#### `validate.go`
```go
var validateCmd = &cobra.Command{
    Use:   "validate <config.yml>",
    Short: "Validate YAML configuration",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
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
```

---

## 処理フロー詳細

### `ttlx build` の実装フロー

```
1. コマンドライン引数解析
   ↓
2. YAML設定ファイル読み込み (config.LoadConfig)
   ↓
3. デフォルト値設定 (Config.SetDefaults)
   ↓
4. バリデーション (config.Validate)
   ↓
5. TTL生成 (generator.Generate)
   ├─ ヘッダー生成
   ├─ 変数定義生成
   ├─ 各ルートステップ処理
   │  ├─ 接続処理生成
   │  └─ コマンド実行生成
   └─ エラーハンドリング生成
   ↓
6. 出力 (ファイル or 標準出力)
   ↓
7. 成功メッセージ表示
```

---

## エラーハンドリング戦略

### エラーの種類と対応

1. **ファイルI/Oエラー**
   - ファイル不存在: `ファイルが見つかりません: <path>`
   - 読み込みエラー: `ファイル読み込みエラー: <error>`
   - 書き込みエラー: `ファイル書き込みエラー: <error>`

2. **YAML構文エラー**
   - パースエラー: `YAML構文エラー (line X): <error>`

3. **バリデーションエラー**
   - 必須項目欠落: `<field> is required`
   - プロファイル不存在: `profile '<name>' not found in route step X`
   - 認証設定不正: `invalid auth in profile '<name>': <reason>`

4. **生成エラー**
   - テンプレートエラー: `TTL generation failed: <error>`

### エラーメッセージのフォーマット

```
Error: <main error message>

<detailed explanation>

Suggestion: <how to fix>
```

**例**:
```
Error: profile 'target' not found in route step 2

The profile 'target' is referenced in the route but not defined in the profiles section.

Suggestion: Add a profile named 'target' to the profiles section.
```

---

## テスト戦略

### ユニットテスト

#### `internal/config`
- `LoadConfig`: 正常系、YAML構文エラー、ファイル不存在
- `Validate`: 各バリデーションルール

#### `internal/generator`
- `Generate`: 最小構成、フル機能、エラーケース
- テンプレート生成: 各コンポーネント

#### `internal/cli`
- コマンド実行: 各コマンド、フラグの動作

### 統合テスト

#### `test/integration/build_test.go`
```go
func TestBuild_SimpleConfig(t *testing.T) {
    // 1. YAML読み込み
    cfg, err := config.LoadConfig("../fixtures/valid/simple.yml")
    require.NoError(t, err)

    // 2. バリデーション
    err = config.Validate(cfg)
    require.NoError(t, err)

    // 3. TTL生成
    ttl, err := generator.Generate(cfg, "simple.yml")
    require.NoError(t, err)

    // 4. 期待されるTTLと比較
    expected, err := os.ReadFile("../fixtures/expected/simple.ttl")
    require.NoError(t, err)

    assert.Contains(t, ttl, "connect 'bastion.example.com:22")
    assert.Contains(t, ttl, "ssh user@10.0.0.50")
}
```

### 手動テスト

1. 生成されたTTLをTera Termで実行
2. テスト環境での実際の接続確認

---

## 実装の優先順位

### フェーズ1: 基盤構築
1. プロジェクト初期化（`go mod init`）
2. ディレクトリ構造作成
3. データモデル定義（`internal/config/model.go`）
4. YAML読み込み（`internal/config/loader.go`）
5. 基本的なテスト

### フェーズ2: バリデーション
1. バリデーター実装（`internal/config/validator.go`）
2. テスト追加

### フェーズ3: TTL生成
1. テンプレート定義（`internal/generator/template.go`）
2. ジェネレーター実装（`internal/generator/generator.go`）
3. 接続処理生成（パスワード認証のみ）
4. コマンド実行生成
5. テスト

### フェーズ4: 公開鍵認証対応
1. 公開鍵認証のTTL生成
2. テスト

### フェーズ5: CLIコマンド
1. Cobra セットアップ
2. `build` コマンド実装
3. `validate` コマンド実装
4. `version` コマンド実装

### フェーズ6: 統合・テスト
1. 統合テスト実装
2. 手動テスト
3. バグ修正

---

## まとめ

本設計書では、ttlx の初回実装における具体的な設計を定義しました。Go の構造体定義、各コンポーネントの実装方針、処理フロー、エラーハンドリング戦略を明確にすることで、実装をスムーズに進めることができます。

次のステップとして、タスクリスト（`tasklist.md`）で具体的な実装タスクを列挙し、進捗管理を行います。
