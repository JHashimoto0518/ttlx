# 技術仕様書

## 概要

本ドキュメントでは、ttlx（Tera Term Language eXtended）の技術仕様を定義します。実装言語、開発ツール、技術的制約、パフォーマンス要件を明確にし、実装の技術的基盤を提供します。

---

## テクノロジースタック

### 実装言語

**選定言語: Go (Golang)**

#### 選定理由

| 要件 | Go | Python | Rust |
|------|:--:|:------:|:----:|
| シングルバイナリ配布 | ✅ 容易 | ❌ ランタイム必要 | ✅ 容易 |
| クロスコンパイル | ✅ 標準機能 | ❌ 複雑 | ✅ 可能 |
| 開発速度 | ✅ 高速 | ✅ 非常に高速 | ❌ 学習コスト高 |
| 実行速度 | ✅ 高速 | ❌ 遅い | ✅ 非常に高速 |
| エコシステム | ✅ 豊富 | ✅ 非常に豊富 | ⚠️ 成長中 |
| 保守性 | ✅ シンプル | ✅ 高い | ⚠️ 複雑 |

**総合評価**:
- **Go**: クロスプラットフォーム配布が容易で、ユーザーがランタイムをインストールする必要がない
- シングルバイナリで配布できるため、インストールが簡単（ダウンロード → 実行）
- 標準ライブラリが充実しており、外部依存を最小化できる
- 開発速度と実行速度のバランスが良い
- YAMLパーサー、CLIフレームワークのエコシステムが成熟

#### バージョン
- **Go 1.21以降** を使用
- モジュール管理に Go Modules を使用

---

### 主要ライブラリ

#### YAML パーサー
- **ライブラリ**: `gopkg.in/yaml.v3`
- **理由**: Go の標準的な YAML ライブラリ、パフォーマンスが良い、メンテナンスが活発

#### CLI フレームワーク
- **ライブラリ**: `github.com/spf13/cobra`
- **理由**: Go で最も人気のある CLI フレームワーク、サブコマンド対応、ヘルプ生成が容易
- **補助**: `github.com/spf13/viper`（設定ファイル管理、環境変数読み込み）

#### テストフレームワーク
- **標準**: Go の標準 `testing` パッケージ
- **補助**: `github.com/stretchr/testify`（アサーション、モック）

#### 差分計算（Phase 2）
- **ライブラリ**: `github.com/sergi/go-diff`
- **理由**: unified diff 形式の生成が容易

#### カラー出力
- **ライブラリ**: `github.com/fatih/color`
- **理由**: クロスプラットフォーム対応、シンプルなAPI

---

### 開発ツール

#### バージョン管理
- **Git**: ソースコード管理
- **GitHub**: リポジトリホスティング、Issue管理、CI/CD

#### コード品質管理
- **golangci-lint**: 包括的なリンター（複数のリンターを統合）
  - `go vet`: 標準の静的解析
  - `staticcheck`: 高度な静的解析
  - `errcheck`: エラーハンドリングチェック
  - `gofmt`: コードフォーマット
  - `goimports`: import文の整理

#### ビルド・パッケージング
- **Go標準ツール**: `go build`, `go test`, `go mod`
- **GoReleaser**: マルチプラットフォームビルド、リリース自動化
  - Windows (amd64, arm64)
  - macOS (amd64, arm64)
  - Linux (amd64, arm64)

#### CI/CD
- **GitHub Actions**: 自動テスト、ビルド、リリース
- **ワークフロー**:
  - プルリクエスト時: リント、テスト実行
  - main ブランチマージ時: テスト、ビルド
  - タグプッシュ時: リリースビルド、GitHub Releases への公開

---

## プロジェクト構造

```
ttlx/
├── cmd/                    # CLIエントリーポイント
│   └── ttlx/
│       └── main.go
├── internal/               # 内部パッケージ（非公開）
│   ├── config/             # YAML設定の読み込み・解析
│   │   ├── loader.go
│   │   ├── model.go
│   │   └── validator.go
│   ├── generator/          # TTL生成
│   │   ├── generator.go
│   │   └── template.go
│   ├── differ/             # 差分計算（Phase 2）
│   │   └── differ.go
│   └── cli/                # CLIコマンド実装
│       ├── build.go
│       ├── validate.go
│       └── diff.go
├── pkg/                    # 公開パッケージ（将来のライブラリ化）
│   └── ttlx/
│       └── api.go
├── test/                   # テストデータ
│   ├── fixtures/
│   │   ├── valid.yml
│   │   ├── invalid.yml
│   │   └── expected.ttl
│   └── integration/
├── docs/                   # ドキュメント
├── .github/
│   └── workflows/          # GitHub Actions
│       ├── test.yml
│       └── release.yml
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── .goreleaser.yml
```

---

## 開発手法

### コーディング規約

#### 命名規則
- **パッケージ名**: 小文字、単数形（`config`, `generator`）
- **ファイル名**: スネークケース（`yaml_loader.go`）
- **型名**: パスカルケース（`ConfigModel`, `TTLGenerator`）
- **関数名**: パスカルケース（公開）、キャメルケース（非公開）
- **変数名**: キャメルケース
- **定数名**: パスカルケースまたは大文字スネークケース

#### コメント
- 公開関数・型には必ずドキュメントコメントを記述
- コメントは英語または日本語（統一する）
- 複雑なロジックには説明コメントを追加

#### エラーハンドリング
- エラーは明示的に処理する（`errcheck` で検証）
- エラーメッセージは具体的で、ユーザーが理解できる内容にする
- エラーのラップには `fmt.Errorf` と `%w` を使用

```go
if err != nil {
    return fmt.Errorf("failed to load YAML file %s: %w", path, err)
}
```

---

### テスト戦略

#### ユニットテスト
- **カバレッジ目標**: 80%以上
- **テストファイル**: `*_test.go`
- **テーブル駆動テスト**: 複数のケースを効率的にテスト

```go
func TestValidator_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {"valid config", validConfig, false},
        {"missing profiles", missingProfilesConfig, true},
        // ...
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

#### 統合テスト
- `test/integration/` に配置
- 実際のYAMLファイルからTTL生成をテスト
- 生成されたTTLが期待通りの内容か検証

#### E2Eテスト（Phase 2以降）
- CLIコマンドの実行をテスト
- 標準出力・エラー出力の検証
- 終了コードの検証

---

### CI/CDパイプライン

#### プルリクエスト時
```yaml
# .github/workflows/test.yml
name: Test

on: [pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Lint
        run: golangci-lint run
      - name: Test
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Coverage
        run: go tool cover -func=coverage.out
```

#### リリース時
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 技術的制約と要件

### プラットフォーム要件

#### サポート対象
- **Windows**: 10以降（amd64, arm64）
- **macOS**: 11 (Big Sur) 以降（amd64, arm64）
- **Linux**: 主要ディストリビューション（amd64, arm64）
  - Ubuntu 20.04以降
  - Debian 11以降
  - CentOS/RHEL 8以降

#### Tera Term互換性
- **対象バージョン**: Tera Term 4.106以降
- 生成されるTTLスクリプトは、上記バージョンで動作することを保証
- Tera Term 5.x系もサポート

---

### セキュリティ要件

#### 機密情報の扱い
1. **パスワード平文記述の警告**
   - YAML設定でパスワードを直接記述した場合、警告を表示
   - 推奨: 環境変数または実行時入力

2. **環境変数からの読み込み**
   - 環境変数名を指定することで、パスワードを外部化

3. **ファイルパーミッション警告**（Phase 2）
   - 生成されたTTLにパスワードが含まれる場合、ファイルパーミッションを確認
   - 推奨: `chmod 600`

4. **ログファイルの扱い**
   - ログファイルに機密情報を出力しない
   - デバッグモード時も、パスワードはマスキング

#### 依存関係の脆弱性管理
- **Dependabot**: 依存ライブラリの脆弱性を自動検出
- **定期的な更新**: 依存ライブラリを最新に保つ

---

### パフォーマンス要件

#### TTL生成速度
- **目標**: 100KB以下のYAMLファイルを1秒以内に処理
- **最大**: 1MB以下のYAMLファイルを5秒以内に処理

#### メモリ使用量
- **目標**: 通常の使用で50MB以下
- **最大**: 100MBを超えない

#### バイナリサイズ
- **目標**: 10MB以下（圧縮前）
- Go のシングルバイナリは比較的大きいが、許容範囲

#### 起動時間
- **目標**: コマンド実行から結果表示まで100ms以内（小さいファイル）

---

## 配布・インストール

### 配布形態

#### GitHub Releases
- **形式**: バイナリ（各プラットフォーム用）
- **命名**: `ttlx-<version>-<os>-<arch>.tar.gz`
  - 例: `ttlx-1.0.0-linux-amd64.tar.gz`
- **チェックサム**: SHA256ハッシュを提供

#### Homebrew（macOS/Linux）（Phase 2以降）
```bash
brew install ttlx
```

#### Scoop（Windows）（Phase 2以降）
```powershell
scoop install ttlx
```

#### Go Install（開発者向け）
```bash
go install github.com/your-org/ttlx/cmd/ttlx@latest
```

---

### インストール手順

#### 手動インストール（推奨）
1. GitHub Releases から対応するバイナリをダウンロード
2. 解凍: `tar -xzf ttlx-1.0.0-linux-amd64.tar.gz`
3. パスに配置: `sudo mv ttlx /usr/local/bin/`
4. 実行権限付与: `chmod +x /usr/local/bin/ttlx`
5. 確認: `ttlx --version`

#### Windowsの場合
1. ZIPファイルをダウンロード・解凍
2. `ttlx.exe` を適切なフォルダに配置
3. 環境変数 `PATH` にフォルダを追加
4. 確認: `ttlx --version`

---

## 拡張性と将来の展望

### ライブラリ化（Phase 3以降）
- `pkg/ttlx` パッケージを公開
- 他のGoプログラムから ttlx の機能を利用可能に

```go
import "github.com/your-org/ttlx/pkg/ttlx"

config, err := ttlx.LoadConfig("config.yml")
ttl, err := ttlx.Generate(config)
```

### プラグインシステム（Phase 3以降）
- カスタムジェネレーターの追加
- カスタムバリデーターの追加
- Go Plugin または WebAssembly ベース

### Language Server Protocol (LSP)（Phase 3以降）
- VSCode、Vim、Emacsなどのエディタでの自動補完
- リアルタイムバリデーション
- ホバー時のドキュメント表示

---

## 技術的負債の管理

### 定期的なリファクタリング
- コードレビュー時に技術的負債を識別
- 四半期ごとにリファクタリングタスクを計画

### 依存関係の更新
- 月次で依存ライブラリを確認
- セキュリティパッチは即座に適用

### パフォーマンス最適化
- ベンチマークテストの実施
- プロファイリングによるボトルネック特定

---

## 制約事項

### 技術的制約
1. **TTL生成のみ**: ttlx は TTL を生成するだけで、SSH接続機能は実装しない
2. **Tera Term依存**: 生成されたTTLは、Tera Termでの実行を前提とする
3. **オフライン動作**: インターネット接続不要で動作する

### 設計上の制約
1. **YAML形式**: 設定ファイルはYAML形式のみサポート（JSON、TOML等は非対応）
2. **シンプルさ優先**: 複雑な機能よりも、シンプルで使いやすいことを優先

---

## まとめ

本技術仕様書では、ttlx の実装言語として **Go** を選定し、主要ライブラリ、開発ツール、テスト戦略、CI/CDパイプラインを定義しました。クロスプラットフォーム対応、シンプルなインストール、高速な実行を実現するための技術的基盤が整いました。

次のステップとして、リポジトリ構造定義書（`repository-structure.md`）で具体的なファイル配置とディレクトリ構成を定義します。
