# 設計書: 複数ルート対応

## 概要

`route`（単数形）から`routes`（複数形）への移行により、1つのYAMLファイルで複数のルートを定義できるようにします。各ルートに名前を付け、個別のTTLファイルを生成します。

## アーキテクチャ

### データフロー

```
YAML (routes定義)
    ↓
Parser (gopkg.in/yaml.v3)
    ↓
Config構造体 (Routes: map[string][]RouteStep)
    ↓
Validator (ルート名、各ルート内容の検証)
    ↓
Generator (ルートごとにTTL生成)
    ↓
map[string]string (ルート名 → TTL内容)
    ↓
CLI (ファイル書き込み)
    ↓
複数の.ttlファイル
```

## データ構造の変更

### 1. Config構造体 (`internal/config/model.go`)

#### 変更前

```go
type Config struct {
    Version  string               `yaml:"version"`
    Profiles map[string]*Profile  `yaml:"profiles"`
    Route    []RouteStep          `yaml:"route"`  // 単一ルート
    Options  *Options             `yaml:"options,omitempty"`
}

type RouteStep struct {
    Profile  string   `yaml:"profile"`
    Commands []string `yaml:"commands,omitempty"`
}
```

#### 変更後

```go
type Config struct {
    Version  string                       `yaml:"version"`
    Profiles map[string]*Profile          `yaml:"profiles"`
    Routes   map[string][]RouteStep       `yaml:"routes"`  // 複数ルート（マップ）
    Options  *Options                     `yaml:"options,omitempty"`
}

type RouteStep struct {
    Profile  string   `yaml:"profile"`
    Commands []string `yaml:"commands,omitempty"`
}
```

**設計上のポイント:**
- `Routes`をマップに変更（キー: ルート名、値: RouteStepのスライス）

### 2. Generator の返り値 (`internal/generator/generator.go`)

#### 変更前

```go
func Generate(cfg *config.Config) (string, error)
```

#### 変更後

```go
// GenerateAll は全ルートのTTLを生成します
func GenerateAll(cfg *config.Config) (map[string]string, error)

// generateRoute は単一ルートのTTLを生成します（内部関数）
func generateRoute(cfg *config.Config, routeName string, route []config.RouteStep) (string, error)
```

**設計上のポイント:**
- `GenerateAll`が全ルートをループしてTTL生成
- `generateRoute`は既存の`Generate`関数のロジックを再利用（リファクタリング）
- 返り値は `map[string]string`（ルート名 → TTL内容）

## コンポーネントの変更

### 1. internal/config/model.go

**変更内容:**
- `Route []RouteStep` → `Routes map[string][]RouteStep`
- `DeprecatedRoute []RouteStep` を追加（旧仕様検出用）

**影響範囲:**
- `SetDefaults()`: Routesの各ルートに対してデフォルト設定を適用
- YAMLパース処理（自動）

### 2. internal/config/validator.go

**変更内容:**
- `route`の検証 → `routes`の検証に変更
- 新しいバリデーションロジックを追加

**追加バリデーション:**

```go
func Validate(config *Config) error {
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

        // 2段目以降のpassword_promptチェック（既存ロジック）
        for i, step := range route {
            if i == 0 {
                continue
            }
            profile := config.Profiles[step.Profile]
            if profile.Auth.Type == "password" && profile.Auth.PasswordPrompt == "" {
                return fmt.Errorf("route '%s': profile '%s': password_prompt is required for password auth in route step %d", routeName, step.Profile, i+1)
            }
        }
    }

    // 既存のプロファイル設定チェック（変更なし）
    // ...

    return nil
}

// isValidFileName はファイル名として有効な文字列かチェック
func isValidFileName(name string) bool {
    // 英数字、ハイフン、アンダースコアのみ許可
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
    return matched
}
```

### 3. internal/generator/generator.go

**変更内容:**
- `Generate()` を `GenerateAll()` にリファクタリング
- 単一ルート生成ロジックを `generateRoute()` に分離
- 複数ルートのループ処理を追加

**実装例:**

```go
// GenerateAll は全ルートのTTLを生成します
func GenerateAll(cfg *config.Config) (map[string]string, error) {
    result := make(map[string]string)

    for routeName, route := range cfg.Routes {
        ttl, err := generateRoute(cfg, routeName, route)
        if err != nil {
            return nil, fmt.Errorf("failed to generate TTL for route '%s': %w", routeName, err)
        }
        result[routeName] = ttl
    }

    return result, nil
}

// generateRoute は単一ルートのTTLを生成します
func generateRoute(cfg *config.Config, routeName string, route []config.RouteStep) (string, error) {
    var sb strings.Builder

    // ヘッダー（ルート名を含める）
    sb.WriteString(generateHeader(cfg, routeName))

    // 変数定義
    sb.WriteString(generateVariables(cfg.Options))

    // 各ステップの生成（既存ロジックを使用、cfg.Route → route引数に変更）
    for i, step := range route {
        profile := cfg.Profiles[step.Profile]

        if i == 0 {
            // 1段目: connect コマンド
            sb.WriteString(generateFirstStep(step.Profile, profile, step.Commands))
        } else {
            // 2段目以降: SSH コマンド
            sb.WriteString(generateSubsequentStep(i, step.Profile, profile, step.Commands))
        }
    }

    // 成功処理
    autoDisconnect := false
    if cfg.Options != nil && cfg.Options.AutoDisconnect != nil {
        autoDisconnect = *cfg.Options.AutoDisconnect
    }

    if autoDisconnect {
        sb.WriteString(generateAutoDisconnect(len(route)))
    } else {
        sb.WriteString(successKeepAliveTemplate)
    }

    // エラーハンドリング（既存ロジックを使用、cfg.Route → route引数に変更）
    sb.WriteString(generateErrorHandling(route, cfg.Profiles))

    return sb.String(), nil
}

// generateHeader はTTLヘッダーを生成（ルート名を追加）
func generateHeader(cfg *config.Config, routeName string) string {
    return fmt.Sprintf(`; ========================================
; Generated by ttlx %s
; Route: %s
; Generated at: %s
; ========================================

`, version, routeName, time.Now().Format("2006-01-02 15:04:05"))
}
```

**既存関数の変更:**
- `generateFirstStep()`, `generateSubsequentStep()`, `generateErrorHandling()` は引数を調整
- これらの関数は `cfg.Route` に依存しないように変更

### 4. internal/cli/build.go

**変更内容:**
- `generator.Generate()` → `generator.GenerateAll()` に変更
- 単一ファイル出力 → 複数ファイル出力
- `-o`オプションの扱いを変更（ファイル名 → ディレクトリ名）

**実装例:**

```go
func runBuild(cmd *cobra.Command, args []string) error {
    configPath := args[0]

    // YAML読み込み
    cfg, err := config.Load(configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    // バリデーション
    if err := config.Validate(cfg); err != nil {
        return fmt.Errorf("validation error: %w", err)
    }

    // TTL生成
    ttls, err := generator.GenerateAll(cfg)
    if err != nil {
        return fmt.Errorf("failed to generate TTL: %w", err)
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

    // dry-runモード
    if dryRun {
        for routeName, ttl := range ttls {
            fmt.Printf("=== %s.ttl ===\n", routeName)
            fmt.Println(ttl)
            fmt.Println()
        }
        return nil
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
}
```

**変更ポイント:**
- `-o`オプションはディレクトリパスとして扱う（従来はファイルパス）
- 複数ファイル生成時のメッセージを分かりやすく
- ファイル一覧をソートして表示

### 5. internal/cli/validate.go

**変更内容:**
- バリデーションロジック自体は`validator.go`に集約されているため、大きな変更なし
- エラーメッセージの表示のみ

## テストの変更

### 1. internal/config/model_test.go

**追加テストケース:**
- `TestConfig_Parse_Routes`: 複数ルートのパーステスト

### 2. internal/config/validator_test.go

**追加テストケース:**
- `TestValidate_Routes_Empty`: routesが空の場合
- `TestValidate_Routes_InvalidRouteName`: 無効なルート名
- `TestValidate_Routes_EmptyRoute`: ルートが空の場合
- `TestValidate_Routes_ProfileNotFound`: プロファイル参照エラー（ルートごと）

### 3. internal/generator/generator_test.go

**変更内容:**
- `TestGenerate` → `TestGenerateAll`
- 複数ルート生成のテストケース追加
- 各ルートのTTL内容検証

### 4. test/fixtures/valid/

**既存フィクスチャの移行:**

すべての `.yml` ファイルで `route:` を `routes:` に変更し、ルート名を追加

**変更例:**

```yaml
# 変更前 (simple.yml)
route:
  - profile: bastion
  - profile: target

# 変更後 (simple.yml)
routes:
  simple-connection:
    - profile: bastion
    - profile: target
```

**新規フィクスチャ:**
- `multiple-routes.yml`: 複数ルート定義の例

### 5. test/fixtures/invalid/

**追加フィクスチャ:**
- `empty-routes.yml`: routesが空（エラー確認用）
- `invalid-route-name.yml`: 無効なルート名（エラー確認用）

## 影響範囲分析

### 変更が必要なファイル

| ファイル | 変更内容 | 影響度 |
|---------|---------|--------|
| `internal/config/model.go` | Route → Routes | 高 |
| `internal/config/validator.go` | バリデーションロジック全面改修 | 高 |
| `internal/generator/generator.go` | Generate → GenerateAll, ループ処理 | 高 |
| `internal/cli/build.go` | 複数ファイル出力処理 | 中 |
| `internal/cli/validate.go` | エラーメッセージ表示のみ | 低 |
| `internal/config/model_test.go` | テストケース追加 | 中 |
| `internal/config/validator_test.go` | テストケース追加 | 中 |
| `internal/generator/generator_test.go` | テストケース修正・追加 | 高 |
| `test/fixtures/valid/*.yml` | route → routes 移行 | 高 |
| `test/fixtures/invalid/*.yml` | 新規エラーケース追加 | 中 |
| `docs/functional-design.md` | YAML仕様更新 | 中 |
| `README.md` / `README.en.md` | 使用例更新 | 中 |

### 破壊的変更のポイント

1. **既存のYAMLファイルはすべて動作しなくなる**
   - `route:` → `routes:` に変更必須
   - 旧仕様使用時は分かりやすいエラーメッセージを表示

2. **`-o`オプションの挙動変更**
   - 従来: 出力ファイル名
   - 新仕様: 出力ディレクトリ名

3. **API変更（外部に公開していないため影響なし）**
   - `generator.Generate()` → `generator.GenerateAll()`

## フィクスチャの移行手順

### test/fixtures/valid/ の移行

すべてのYAMLファイルで `route:` を `routes:` に変換し、適切なルート名を付けます。

**変換例:**

```yaml
# 変更前 (simple.yml)
route:
  - profile: bastion
  - profile: target

# 変更後 (simple.yml)
routes:
  simple-connection:
    - profile: bastion
    - profile: target
```

**ルート名の命名規則:**
- ファイル名をベースにする（例: `simple.yml` → `simple-connection`）
- 複数ルートを定義する場合は意味のある名前を付ける

## リスク分析

### 中リスク

- **テストカバレッジの低下**
  - 軽減策: 既存テストを移行し、新規テストを追加

- **バグの混入**
  - 軽減策: 段階的な実装とテスト実行

### 低リスク

- **パフォーマンスの劣化**
  - 影響: 複数ルート生成でもループ処理のみ、O(n)で問題なし

## 実装順序

1. **Config構造体の変更** (`internal/config/model.go`)
2. **Validatorの変更** (`internal/config/validator.go`)
3. **Generatorの変更** (`internal/generator/generator.go`)
4. **CLIの変更** (`internal/cli/build.go`)
5. **テストの追加・修正**
6. **フィクスチャの移行**
7. **ドキュメント更新**
8. **統合テスト**

## 未解決の問題

- なし（現時点で明確な実装方針が定まっている）
