# 開発ガイドライン

## 概要

本ドキュメントでは、ttlx プロジェクトの開発における規約とベストプラクティスを定義します。すべての開発者がこのガイドラインに従うことで、コードの一貫性と保守性を確保します。

---

## コーディング規約

### 基本方針

1. **Effective Go に従う**
   - 公式ドキュメント: https://go.dev/doc/effective_go
   - Go の標準的な書き方を優先

2. **シンプルさを重視**
   - 過度な抽象化を避ける
   - 明確で読みやすいコードを書く

3. **ツールによる自動整形**
   - `gofmt` でフォーマット
   - `goimports` で import 文を整理

---

### フォーマット

#### インデント
- **タブ文字**を使用（Go の標準）
- `gofmt` が自動的に整形

#### 行の長さ
- **推奨**: 100文字以内
- **最大**: 120文字
- 長い行は適切に改行

#### 空行
- 関数間に1行の空行
- ロジックの区切りに空行を挿入

```go
// Good
func Foo() {
    // ...
}

func Bar() {
    // ...
}

// Bad
func Foo() {
    // ...
}
func Bar() {
    // ...
}
```

---

### 命名規則

#### パッケージ名
- **小文字のみ**、単数形
- 短く、簡潔に
- アンダースコア不可

```go
// Good
package config
package generator

// Bad
package configLoader
package ttl_generator
```

#### ファイル名
- **スネークケース**
- 小文字、アンダースコア区切り

```
// Good
yaml_loader.go
ssh_generator.go

// Bad
yamlLoader.go
SSHGenerator.go
```

#### 型名
- **パスカルケース**（大文字始まり）
- 公開型は大文字、非公開型は小文字始まり

```go
// Good
type Config struct { ... }       // 公開型
type Profile struct { ... }      // 公開型
type validator struct { ... }    // 非公開型

// Bad
type config struct { ... }       // 小文字始まりは非公開
type PROFILE struct { ... }      // 全て大文字は避ける
```

#### 関数・メソッド名
- **パスカルケース**（公開）、**キャメルケース**（非公開）
- 動詞で始める

```go
// Good
func LoadConfig(path string) (*Config, error) { ... }  // 公開
func validateProfile(p *Profile) error { ... }         // 非公開

// Bad
func load_config(path string) { ... }                  // スネークケース不可
func ProfileValidate(p *Profile) { ... }               // 名詞始まり
```

#### 変数名
- **キャメルケース**
- 短いスコープでは短い名前（`i`, `err`）
- 長いスコープでは説明的な名前

```go
// Good
var configPath string
var err error
for i := 0; i < 10; i++ { ... }

// Bad
var config_path string
var e error
```

#### 定数名
- **パスカルケース** または **大文字スネークケース**

```go
// Good
const DefaultTimeout = 30
const MAX_RETRY_COUNT = 5

// Bad
const default_timeout = 30
```

#### インターフェース名
- 単一メソッドの場合、メソッド名 + `er`

```go
// Good
type Loader interface {
    Load(path string) error
}

type Generator interface {
    Generate(config *Config) (string, error)
}

// Bad
type ILoader interface { ... }
type LoaderInterface interface { ... }
```

---

### コメント

#### パッケージコメント
- パッケージの最初のファイルに記述
- `// Package <name> ...` の形式

```go
// Package config provides YAML configuration loading and validation.
package config
```

#### 公開型・関数のコメント
- **必須**: すべての公開型・関数にコメント
- 型名・関数名で始める

```go
// Config represents the ttlx YAML configuration.
type Config struct {
    Version  string
    Profiles map[string]*Profile
}

// LoadConfig loads and parses a YAML configuration file.
func LoadConfig(path string) (*Config, error) {
    // ...
}
```

#### 内部コメント
- 複雑なロジックには説明コメント
- 明白なことはコメント不要

```go
// Good
// Merge global profiles into the local configuration.
// Local profiles take precedence over global ones.
func mergeProfiles(local, global map[string]*Profile) {
    // ...
}

// Bad
// Loop through profiles
for _, p := range profiles {
    // ...
}
```

---

### エラーハンドリング

#### エラーチェック
- **すべてのエラーを明示的に処理**
- `errcheck` リンターで検証

```go
// Good
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("failed to read file %s: %w", path, err)
}

// Bad
data, _ := os.ReadFile(path)  // エラーを無視
```

#### エラーのラップ
- `fmt.Errorf` と `%w` を使用
- コンテキストを追加

```go
// Good
if err := validate(config); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// Bad
if err := validate(config); err != nil {
    return err  // コンテキストなし
}
```

#### カスタムエラー
- 特定のエラー型は `internal/errors` に定義

```go
// internal/errors/errors.go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}
```

---

### 関数設計

#### 関数の長さ
- **推奨**: 50行以内
- **最大**: 100行
- 長すぎる関数は分割

#### 引数の数
- **推奨**: 3個以内
- **最大**: 5個
- 多すぎる場合は構造体を使用

```go
// Good
type GenerateOptions struct {
    Timeout int
    Retry   int
    Log     bool
}

func Generate(config *Config, opts GenerateOptions) (string, error) {
    // ...
}

// Bad
func Generate(config *Config, timeout int, retry int, log bool, logFile string) {
    // ...
}
```

#### 戻り値
- エラーは最後の戻り値

```go
// Good
func LoadConfig(path string) (*Config, error) { ... }

// Bad
func LoadConfig(path string) (error, *Config) { ... }
```

---

## テスト規約

### テストファイル

#### 配置
- 実装ファイルと同じディレクトリ
- `*_test.go` の形式

```
internal/config/
├── loader.go
└── loader_test.go
```

#### パッケージ名
- **ホワイトボックステスト**: 同じパッケージ名（`package config`）
- **ブラックボックステスト**: `_test` サフィックス（`package config_test`）

---

### テストの書き方

#### テーブル駆動テスト
- 複数のケースを効率的にテスト

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  validConfig(),
            wantErr: false,
        },
        {
            name:    "missing version",
            config:  configWithoutVersion(),
            wantErr: true,
        },
        {
            name:    "invalid profile",
            config:  configWithInvalidProfile(),
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### テスト名
- `Test<FunctionName>` の形式
- サブテストは説明的な名前

```go
func TestLoadConfig(t *testing.T) {
    t.Run("valid YAML file", func(t *testing.T) { ... })
    t.Run("file not found", func(t *testing.T) { ... })
    t.Run("invalid YAML syntax", func(t *testing.T) { ... })
}
```

#### アサーション
- `testify/assert` または `testify/require` を使用

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFoo(t *testing.T) {
    result, err := Foo()
    require.NoError(t, err)           // エラー時にテスト中断
    assert.Equal(t, expected, result) // 比較
}
```

---

### カバレッジ

#### 目標
- **全体**: 80%以上
- **コアロジック**: 90%以上（`config`, `generator`）

#### カバレッジ確認
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

### モックとスタブ

#### ファイルI/O
- テスト用のフィクスチャを使用（`test/fixtures/`）

```go
func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("../../test/fixtures/valid/simple.yml")
    require.NoError(t, err)
    assert.Equal(t, "1.0", config.Version)
}
```

#### 外部依存
- インターフェースでモック化

```go
// 本番コード
type FileReader interface {
    ReadFile(path string) ([]byte, error)
}

// テストコード
type mockFileReader struct {
    data []byte
    err  error
}

func (m *mockFileReader) ReadFile(path string) ([]byte, error) {
    return m.data, m.err
}
```

---

## Git規約

### ブランチ戦略

#### ブランチ名
- **main**: メインブランチ（常にリリース可能）
- **feature/**: 機能開発（`feature/add-diff-command`）
- **bugfix/**: バグ修正（`bugfix/fix-validation`）
- **hotfix/**: 緊急修正（`hotfix/security-patch`）

#### ブランチ命名規則
```
<type>/<short-description>

例:
feature/add-init-command
bugfix/fix-yaml-parsing
hotfix/critical-security-fix
```

---

### コミットメッセージ

#### フォーマット
```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Type（必須）
- **feat**: 新機能
- **fix**: バグ修正
- **docs**: ドキュメント変更
- **style**: コードフォーマット（機能変更なし）
- **refactor**: リファクタリング
- **test**: テスト追加・修正
- **chore**: ビルド、CI/CD、依存関係更新

#### Scope（任意）
- 変更の範囲（`config`, `generator`, `cli`）

#### Subject（必須）
- 変更の簡潔な説明
- 命令形（"add" not "added"）
- 小文字始まり
- 末尾にピリオド不要

#### 例
```
feat(cli): add diff command

Implement the `ttlx diff` command to show differences
between current config and previously generated TTL.

Closes #42
```

```
fix(generator): escape special characters in passwords

Special characters in passwords were not properly escaped,
causing TTL generation to fail.

Fixes #58
```

```
docs: update README with installation instructions
```

---

### プルリクエスト規約

#### PRタイトル
- コミットメッセージと同じフォーマット

```
feat(cli): add diff command
fix(generator): handle special characters in passwords
```

#### PR説明
- **変更内容**: 何を変更したか
- **動機**: なぜ変更したか
- **影響範囲**: どこに影響するか
- **テスト**: どのようにテストしたか

#### テンプレート
```markdown
## 変更内容
- `ttlx diff` コマンドを追加
- メタデータファイルの保存・読み込み機能を実装

## 動機
設定ファイル変更前後のTTL差分を確認できるようにするため（#42）

## 影響範囲
- `internal/differ` パッケージを新規追加
- `internal/cli/diff.go` を追加
- `internal/generator` に `Metadata` 保存機能を追加

## テスト
- `internal/differ` のユニットテストを追加（カバレッジ 85%）
- `test/fixtures/` にテストデータを追加
- 手動テスト: `ttlx diff` コマンドの動作確認

## チェックリスト
- [ ] テストが通る（`go test ./...`）
- [ ] リントエラーがない（`golangci-lint run`）
- [ ] ドキュメント更新済み
- [ ] CHANGELOG.md に記載
```

---

### コードレビュー

#### レビュアーの責務
- コードの正確性を確認
- テストの十分性を確認
- 可読性・保守性を評価
- 設計の妥当性を検討

#### レビューコメント
- 建設的なフィードバック
- 具体的な改善案を提示

```markdown
# Good
`config.Validate()` で `profiles` が空の場合もエラーにすべきでは？
現状だと空のプロファイルでもパスしてしまいます。

# Bad
これはダメです。
```

---

## バージョニングとリリースポリシー

### バージョン番号の方針

本プロジェクトは **セマンティックバージョニング（Semantic Versioning）** に準拠します。

```
<major>.<minor>.<patch>[-<pre-release>]

例:
1.0.0       # 安定版
0.1.0-beta  # ベータ版
0.2.0-rc.1  # リリース候補
```

#### バージョンの意味

- **Major（メジャー）**: 後方互換性のない変更
- **Minor（マイナー）**: 後方互換性のある機能追加
- **Patch（パッチ）**: 後方互換性のあるバグ修正
- **Pre-release（プレリリース）**: 開発版・ベータ版・RC版など

---

### リリース前の後方互換性ポリシー

#### ベータ版（0.x.x-beta）における後方互換性

**重要**: ベータ版では後方互換性を担保する必要はありません。

**理由**:
- ベータ版は正式リリース前の開発版であり、仕様が流動的
- 機能の追加・変更・削除が頻繁に発生する
- ユーザーはベータ版使用時に破壊的変更があることを理解している

**適用例**:
```yaml
# v0.1.0-beta で route → routes に変更
# 後方互換性なし（旧仕様は削除）

# ❌ 非推奨（ベータ版では不要）
route:           # deprecated, use routes instead
  - profile: bastion

# ✅ 推奨（新仕様のみサポート）
routes:
  main:
    - profile: bastion
```

**開発時の指針**:
- ベータ版では設計のシンプルさを優先
- 後方互換性のためのレガシーコードを残さない
- CHANGELOG に破壊的変更（BREAKING CHANGE）を明記
- 移行方法を README に記載

---

### 安定版（1.0.0以降）における後方互換性

**重要**: 1.0.0以降の安定版では、後方互換性を厳密に管理します。

#### Major バージョンアップ時（破壊的変更）
- 後方互換性のない変更が許可される
- CHANGELOG に詳細な移行ガイドを記載
- 可能な限り非推奨警告を先に導入（1つ前のバージョンで警告）

#### Minor バージョンアップ時（機能追加）
- 後方互換性を維持
- 既存機能の動作を変更しない
- 新機能の追加のみ

#### Patch バージョンアップ時（バグ修正）
- 後方互換性を維持
- バグ修正のみ（新機能追加なし）

---

### バージョン管理の実践

#### バージョン番号の更新箇所

バージョン番号は以下のファイルで管理されます:

1. **`internal/cli/version.go`** - CLI バージョン表示用
   ```go
   const version = "0.1.0-beta"
   ```

2. **`internal/generator/generator.go`** - TTL ヘッダー生成用
   ```go
   const version = "0.1.0-beta"
   ```

3. **`CHANGELOG.md`** - 変更履歴

#### バージョンアップ手順

1. **バージョン番号の決定**
   - 変更内容に応じて適切なバージョンを決定
   - ベータ版: `0.x.x-beta`
   - 安定版: `x.y.z`

2. **コード内のバージョン更新**
   ```bash
   # 2箇所のバージョン定数を更新
   vim internal/cli/version.go
   vim internal/generator/generator.go
   ```

3. **CHANGELOG.md の更新**
   - `[Unreleased]` セクションを新バージョンに変更
   - 変更内容を `Added`, `Changed`, `Fixed` などに分類
   - 破壊的変更は `**BREAKING**:` プレフィックスで明示

4. **Git タグの作成**
   ```bash
   git tag -a v0.1.0-beta -m "Release version 0.1.0-beta"
   git push origin v0.1.0-beta
   ```

---

### リリースタイプ

#### 開発版（Development）
- バージョン: `0.0.x`
- 対象: 開発者のみ
- 後方互換性: 不要

#### ベータ版（Beta）
- バージョン: `0.x.x-beta`
- 対象: 早期アダプター、テスター
- 後方互換性: **不要**
- CHANGELOG に全ての破壊的変更を明記

#### リリース候補（Release Candidate）
- バージョン: `x.y.z-rc.N`
- 対象: 本番環境での最終検証
- 後方互換性: 必須（安定版と同等）

#### 安定版（Stable）
- バージョン: `1.0.0` 以降
- 対象: 全ユーザー
- 後方互換性: 厳密に管理

---

### CHANGELOG.md の記載ルール

#### フォーマット

```markdown
## [Unreleased]

### Note
Version 1.0.0 will be the first stable release. Currently in beta.

## [0.2.0-beta] - 2026-01-15

### Added
- 新機能の説明

### Changed
- **BREAKING**: 破壊的変更の説明
- その他の変更

### Fixed
- バグ修正の説明

### Deprecated
- 非推奨機能（将来削除予定）

### Removed
- 削除された機能

### Security
- セキュリティ修正
```

#### 破壊的変更の記載方法

```markdown
### Changed
- **BREAKING**: `route` field is no longer supported. Use `routes` instead.
  - Migration: Rename `route:` to `routes:` and add route names as keys.
  - Example:
    ```yaml
    # Before (v0.0.x)
    route:
      - profile: bastion

    # After (v0.1.0-beta)
    routes:
      main:
        - profile: bastion
    ```
```

---

## スタイリング規約

### gofmt に準拠
- すべてのコードは `gofmt` でフォーマット
- CI で自動チェック

### import 文の整理
- `goimports` で自動整理
- 標準ライブラリ → 外部ライブラリ → 内部パッケージ

```go
import (
    // 標準ライブラリ
    "fmt"
    "os"

    // 外部ライブラリ
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    // 内部パッケージ
    "github.com/your-org/ttlx/internal/config"
)
```

---

## CI/CDチェックリスト

### プルリクエスト時
- [ ] `go test ./...` が成功
- [ ] `golangci-lint run` がエラーなし
- [ ] カバレッジ 80% 以上
- [ ] コードレビュー承認

### リリース時
- [ ] すべてのテストが成功
- [ ] ドキュメント更新
- [ ] CHANGELOG.md 更新
- [ ] バージョンタグ作成（`v1.0.0`）

---

## ベストプラクティス

### DRY（Don't Repeat Yourself）
- 重複コードを避ける
- 共通処理は関数化

### YAGNI（You Aren't Gonna Need It）
- 今必要な機能のみを実装
- 将来の拡張を過度に考慮しない

### エラーメッセージ
- ユーザーフレンドリー
- 具体的で、解決方法を示唆

```go
// Good
return fmt.Errorf("profile 'bastion' not found in profiles section (line 15)")

// Bad
return fmt.Errorf("profile not found")
```

### パフォーマンス
- 過度な最適化を避ける
- まず正しく動くコードを書く
- ボトルネックが判明してから最適化

---

## まとめ

本開発ガイドラインでは、ttlx プロジェクトのコーディング規約、命名規則、テスト規約、Git規約を定義しました。すべての開発者がこのガイドラインに従うことで、一貫性のある高品質なコードを維持できます。

次のステップとして、ユビキタス言語定義（`glossary.md`）でプロジェクト内で使用する用語を統一します。
