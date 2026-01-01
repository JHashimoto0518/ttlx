# リポジトリ構造定義書

## 概要

本ドキュメントでは、ttlx リポジトリのディレクトリ構成とファイル配置ルールを定義します。Go プロジェクトの標準的なレイアウトに従い、保守性と拡張性を確保します。

---

## ディレクトリ構成

```
ttlx/
├── cmd/                           # アプリケーションエントリーポイント
│   └── ttlx/
│       └── main.go               # メインエントリーポイント
│
├── internal/                      # 非公開パッケージ（プロジェクト固有）
│   ├── config/                    # 設定ファイルの読み込み・解析
│   │   ├── loader.go             # YAML読み込み、プロファイルマージ
│   │   ├── model.go              # データモデル定義
│   │   ├── validator.go          # バリデーションロジック
│   │   └── loader_test.go        # ユニットテスト
│   │
│   ├── generator/                 # TTL生成
│   │   ├── generator.go          # TTL生成メインロジック
│   │   ├── template.go           # テンプレート管理
│   │   ├── ssh.go                # SSH接続TTL生成
│   │   ├── command.go            # コマンド実行TTL生成
│   │   └── generator_test.go    # ユニットテスト
│   │
│   ├── differ/                    # 差分計算（Phase 2）
│   │   ├── differ.go             # 差分計算ロジック
│   │   ├── metadata.go           # メタデータ管理
│   │   └── differ_test.go        # ユニットテスト
│   │
│   ├── cli/                       # CLIコマンド実装
│   │   ├── root.go               # ルートコマンド
│   │   ├── build.go              # buildコマンド
│   │   ├── validate.go           # validateコマンド
│   │   ├── diff.go               # diffコマンド（Phase 2）
│   │   ├── init.go               # initコマンド（Phase 2）
│   │   └── version.go            # versionコマンド
│   │
│   └── errors/                    # エラー定義
│       ├── errors.go             # カスタムエラー型
│       └── messages.go           # エラーメッセージ定義
│
├── pkg/                           # 公開パッケージ（将来のライブラリ化）
│   └── ttlx/
│       ├── api.go                # 公開API
│       └── types.go              # 公開型定義
│
├── test/                          # テストデータ・統合テスト
│   ├── fixtures/                 # テストフィクスチャ
│   │   ├── valid/                # 正常系YAMLファイル
│   │   │   ├── simple.yml
│   │   │   ├── multi-hop.yml
│   │   │   ├── with-commands.yml
│   │   │   └── global-profiles.yml
│   │   ├── invalid/              # 異常系YAMLファイル
│   │   │   ├── missing-version.yml
│   │   │   ├── invalid-profile.yml
│   │   │   └── syntax-error.yml
│   │   └── expected/             # 期待されるTTL出力
│   │       ├── simple.ttl
│   │       ├── multi-hop.ttl
│   │       └── with-commands.ttl
│   │
│   └── integration/               # 統合テスト
│       └── build_test.go
│
├── docs/                          # ドキュメント
│   ├── product-requirements.md   # プロダクト要求定義書
│   ├── functional-design.md      # 機能設計書
│   ├── architecture.md           # 技術仕様書
│   ├── repository-structure.md   # リポジトリ構造定義書（本ドキュメント）
│   ├── development-guidelines.md # 開発ガイドライン
│   └── glossary.md               # ユビキタス言語定義
│
├── .github/                       # GitHub関連設定
│   ├── workflows/                # GitHub Actions
│   │   ├── test.yml              # PR時のテスト
│   │   ├── lint.yml              # リント
│   │   └── release.yml           # リリースビルド
│   ├── ISSUE_TEMPLATE/           # Issueテンプレート
│   │   ├── bug_report.md
│   │   └── feature_request.md
│   └── PULL_REQUEST_TEMPLATE.md  # PRテンプレート
│
├── scripts/                       # ビルド・開発スクリプト
│   ├── build.sh                  # ローカルビルドスクリプト
│   ├── test.sh                   # テスト実行スクリプト
│   └── install.sh                # インストールスクリプト
│
├── .gitignore                     # Git除外設定
├── .golangci.yml                  # golangci-lint設定
├── .goreleaser.yml                # GoReleaser設定
├── go.mod                         # Go モジュール定義
├── go.sum                         # 依存関係のチェックサム
├── LICENSE                        # ライセンス
├── README.md                      # プロジェクト概要
├── CONTRIBUTING.md                # コントリビューションガイド
└── CHANGELOG.md                   # 変更履歴

```

---

## ディレクトリの役割

### `/cmd`
**役割**: アプリケーションのエントリーポイント

- プロジェクトで実行可能なバイナリごとにサブディレクトリを作成
- `main.go` のみを配置し、ロジックは `internal/` に実装
- 複数のCLIツールを提供する場合は、`cmd/ttlx-server/` などを追加可能

**ファイル配置ルール**:
- `main.go`: `main()` 関数のみを実装
- ビジネスロジックは含めない（`internal/` を呼び出すのみ）

**例**:
```go
// cmd/ttlx/main.go
package main

import (
    "github.com/your-org/ttlx/internal/cli"
    "os"
)

func main() {
    if err := cli.Execute(); err != nil {
        os.Exit(1)
    }
}
```

---

### `/internal`
**役割**: プロジェクト固有の非公開パッケージ

- 他のプロジェクトからインポート不可（Go の仕様）
- アプリケーションのコアロジックを実装
- パッケージごとにサブディレクトリを作成

**サブディレクトリ**:

#### `/internal/config`
- YAML設定ファイルの読み込み・解析
- データモデル（構造体）の定義
- バリデーションロジック
- グローバルプロファイルのマージ

**主要ファイル**:
- `model.go`: `Config`, `Profile`, `RouteStep` などの構造体定義
- `loader.go`: YAML読み込み、パース
- `validator.go`: バリデーションロジック

#### `/internal/generator`
- TTL スクリプト生成ロジック
- テンプレート管理
- SSH接続、コマンド実行のTTL生成

**主要ファイル**:
- `generator.go`: メイン生成ロジック
- `template.go`: テンプレート文字列管理
- `ssh.go`: SSH接続TTL生成（多段対応）
- `command.go`: コマンド実行TTL生成

#### `/internal/differ`（Phase 2）
- 設定ファイル変更前後の差分計算
- メタデータファイルの管理

**主要ファイル**:
- `differ.go`: 差分計算ロジック
- `metadata.go`: メタデータの保存・読み込み

#### `/internal/cli`
- CLI コマンドの実装（cobra使用）
- コマンドライン引数の解析
- 各コマンドの実行ロジック

**主要ファイル**:
- `root.go`: ルートコマンド、共通フラグ
- `build.go`: `ttlx build` コマンド
- `validate.go`: `ttlx validate` コマンド
- `diff.go`: `ttlx diff` コマンド
- `init.go`: `ttlx init` コマンド（対話的セットアップ）
- `version.go`: `ttlx version` コマンド

#### `/internal/errors`
- カスタムエラー型の定義
- エラーメッセージの一元管理

**主要ファイル**:
- `errors.go`: エラー型定義
- `messages.go`: エラーメッセージ定数

---

### `/pkg`
**役割**: 外部からインポート可能な公開パッケージ

- 将来的にライブラリとして使用される可能性のあるコード
- Phase 1 では最小限、Phase 3 以降で拡充

**配置ルール**:
- 安定したAPIのみを公開
- 後方互換性を保つ必要がある

---

### `/test`
**役割**: テストデータと統合テスト

#### `/test/fixtures`
- テスト用のYAMLファイル、期待されるTTL出力
- `valid/`: 正常系のテストデータ
- `invalid/`: 異常系のテストデータ
- `expected/`: 期待される出力

**命名規則**:
- ファイル名は内容を表す（`simple.yml`, `multi-hop.yml`）
- 対応する期待出力は同名で `.ttl` 拡張子（`simple.ttl`）

#### `/test/integration`
- 統合テスト（E2E的なテスト）
- 実際のファイルからTTL生成をテスト

---

### `/docs`
**役割**: プロジェクトドキュメント

- 永続的なドキュメントを配置
- Markdown 形式
- 図表は Mermaid または画像（`docs/images/`）

**ファイル一覧**:
- `product-requirements.md`: プロダクト要求定義
- `functional-design.md`: 機能設計
- `architecture.md`: 技術仕様
- `repository-structure.md`: リポジトリ構造（本ドキュメント）
- `development-guidelines.md`: 開発ガイドライン
- `glossary.md`: 用語集

---

### `/.github`
**役割**: GitHub 固有の設定

#### `/workflows`
- GitHub Actions のワークフロー定義
- CI/CD パイプライン

**ファイル一覧**:
- `test.yml`: PR時の自動テスト
- `lint.yml`: golangci-lint 実行
- `release.yml`: タグプッシュ時のリリースビルド

#### `/ISSUE_TEMPLATE`
- Issue 作成時のテンプレート
- `bug_report.md`: バグ報告用
- `feature_request.md`: 機能要望用

---

### `/scripts`
**役割**: ビルド・開発補助スクリプト

- シェルスクリプト（`.sh`）
- クロスプラットフォーム対応が必要な場合は Makefile も検討

**ファイル例**:
- `build.sh`: ローカルでのビルド
- `test.sh`: テスト実行（カバレッジレポート付き）
- `install.sh`: ローカルインストール

---

## ファイル配置ルール

### Goファイルの命名規則

1. **パッケージ名 = ディレクトリ名**
   - `internal/config/` 配下のファイルは `package config`

2. **ファイル名はスネークケース**
   - `yaml_loader.go`, `ssh_generator.go`
   - 単語が1つの場合は小文字（`model.go`, `errors.go`）

3. **テストファイル**
   - `*_test.go` の形式
   - 例: `loader.go` → `loader_test.go`

4. **1ファイル1責務**
   - 1つのファイルには関連する機能のみを実装
   - ファイルが500行を超えたら分割を検討

---

### テストファイルの配置

1. **ユニットテスト**
   - 実装ファイルと同じディレクトリに配置
   - 例: `internal/config/loader.go` → `internal/config/loader_test.go`

2. **統合テスト**
   - `/test/integration/` に配置
   - ファイル名: `*_integration_test.go` または `*_test.go`

3. **テストデータ**
   - `/test/fixtures/` に配置
   - サブディレクトリで分類（`valid/`, `invalid/`, `expected/`）

---

### ドキュメントの配置

1. **永続的ドキュメント**
   - `/docs/` に配置
   - Markdown 形式（`.md`）

2. **作業単位のドキュメント**
   - `/.steering/[YYYYMMDD]-[開発タイトル]/` に配置
   - 例: `.steering/20250103-initial-implementation/`

3. **画像・図表**
   - `/docs/images/` に配置（必要な場合のみ）
   - PNG または SVG 形式

---

## 設定ファイルの管理

### `.gitignore`
無視するファイル・ディレクトリ：
```gitignore
# ビルド成果物
/bin/
/dist/
*.exe
*.dll
*.so
*.dylib

# テスト成果物
*.test
*.out
coverage.txt
coverage.html

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# 一時ファイル
*.tmp
*.log

# ttlx固有
.ttlx/
*.ttl
```

### `.golangci.yml`
golangci-lint の設定：
```yaml
linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode

linters-settings:
  errcheck:
    check-blank: true
  govet:
    check-shadowing: true

issues:
  exclude-use-default: false
```

### `.goreleaser.yml`
GoReleaser の設定（マルチプラットフォームビルド）：
```yaml
builds:
  - main: ./cmd/ttlx
    binary: ttlx
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: your-org
    name: ttlx
```

---

## バージョン管理

### セマンティックバージョニング
- **フォーマット**: `vMAJOR.MINOR.PATCH`
- **例**: `v1.0.0`, `v1.2.3`

### タグ命名規則
- Git タグ: `v1.0.0`
- リリースブランチ: `release/v1.0.x`

### ブランチ戦略
- **main**: 安定版、常にリリース可能
- **develop**: 開発ブランチ（Phase 1では不要、Phase 2以降で検討）
- **feature/**: 機能開発（`feature/add-diff-command`）
- **bugfix/**: バグ修正（`bugfix/fix-validation-error`）

---

## 依存関係の管理

### `go.mod`
- Go Modules を使用
- 依存ライブラリのバージョンを明示

### 依存ライブラリの追加
```bash
go get github.com/spf13/cobra@latest
go mod tidy
```

### 依存ライブラリの更新
```bash
go get -u ./...
go mod tidy
```

---

## まとめ

本リポジトリ構造定義書では、ttlx プロジェクトのディレクトリ構成、ファイル配置ルール、命名規則を定義しました。Go の標準的なレイアウトに従い、保守性と拡張性を確保しています。

次のステップとして、開発ガイドライン（`development-guidelines.md`）でコーディング規約、Git規約、テスト規約を定義します。
