# 初回実装タスクリスト

## 概要

本ドキュメントでは、ttlx の初回実装（Phase 1 MVP）における具体的なタスクを列挙し、進捗を管理します。

---

## フェーズ1: 基盤構築

### プロジェクト初期化

- [ ] Go モジュール初期化
  - `go mod init github.com/your-org/ttlx`
  - `go.mod`, `go.sum` 作成

- [ ] 依存ライブラリ追加
  - `gopkg.in/yaml.v3` (YAML パーサー)
  - `github.com/spf13/cobra` (CLI フレームワーク)
  - `github.com/stretchr/testify` (テストライブラリ)

- [ ] ディレクトリ構造作成
  - `cmd/ttlx/`
  - `internal/config/`
  - `internal/generator/`
  - `internal/cli/`
  - `internal/errors/`
  - `test/fixtures/valid/`
  - `test/fixtures/invalid/`
  - `test/fixtures/expected/`
  - `test/integration/`

- [ ] `.gitignore` 作成
  - ビルド成果物、テスト成果物、IDE設定を除外

- [ ] `.golangci.yml` 作成
  - golangci-lint 設定

---

### データモデル定義

- [ ] `internal/config/model.go` 作成
  - `Config` 構造体定義
  - `Profile` 構造体定義
  - `Auth` 構造体定義
  - `RouteStep` 構造体定義
  - `Options` 構造体定義
  - `SetDefaults()` メソッド実装

- [ ] `internal/config/model_test.go` 作成
  - `SetDefaults` のテスト

---

### YAML読み込み

- [ ] `internal/config/loader.go` 作成
  - `LoadConfig(path string) (*Config, error)` 実装
  - ファイル読み込み
  - YAML パース
  - デフォルト値設定

- [ ] `internal/config/loader_test.go` 作成
  - 正常系: 有効なYAMLファイルの読み込み
  - 異常系: ファイル不存在
  - 異常系: YAML構文エラー

- [ ] テストフィクスチャ作成
  - `test/fixtures/valid/simple.yml` (最小構成)
  - `test/fixtures/valid/full.yml` (全機能使用)
  - `test/fixtures/invalid/syntax-error.yml`

---

## フェーズ2: バリデーション

### バリデーター実装

- [ ] `internal/config/validator.go` 作成
  - `Validate(config *Config) error` 実装
  - バージョンチェック
  - プロファイル存在チェック
  - ルート存在チェック
  - プロファイル参照チェック
  - 認証設定チェック (`validateAuth`)

- [ ] `internal/config/validator_test.go` 作成
  - 正常系: 有効な設定のバリデーション成功
  - 異常系: バージョン欠落
  - 異常系: プロファイル欠落
  - 異常系: ルート欠落
  - 異常系: プロファイル参照不正
  - 異常系: 認証設定不正 (各パターン)

- [ ] テストフィクスチャ追加
  - `test/fixtures/invalid/missing-version.yml`
  - `test/fixtures/invalid/missing-profiles.yml`
  - `test/fixtures/invalid/invalid-profile-ref.yml`
  - `test/fixtures/invalid/invalid-auth.yml`

---

## フェーズ3: TTL生成（パスワード認証）

### テンプレート定義

- [ ] `internal/generator/template.go` 作成
  - ヘッダーテンプレート定義
  - 変数定義テンプレート
  - 接続テンプレート（最初のステップ）
  - SSH コマンドテンプレート（2番目以降）
  - パスワード認証テンプレート
  - コマンド実行テンプレート
  - エラーハンドリングテンプレート
  - 終了処理テンプレート

---

### ジェネレーター実装

- [ ] `internal/generator/generator.go` 作成
  - `Generate(config *Config, sourceFile string) (string, error)` 実装
  - ヘッダー生成
  - 変数定義生成
  - 接続処理生成（パスワード認証のみ）
    - 最初のステップ: `connect` コマンド
    - 2番目のステップ: `sendln 'ssh ...'` コマンド
  - パスワード認証処理生成
    - 環境変数から読み込み
    - 実行時入力 (`passwordbox`)
  - コマンド実行処理生成
  - エラーハンドリングコード生成
  - 正常終了・異常終了処理生成

- [ ] `internal/generator/generator_test.go` 作成
  - 正常系: 最小構成のTTL生成
  - 正常系: フル機能のTTL生成
  - 各コンポーネントのテスト（ヘッダー、変数、接続、コマンド）

- [ ] 期待されるTTL作成
  - `test/fixtures/expected/simple.ttl`
  - `test/fixtures/expected/full.ttl`

---

## フェーズ4: 公開鍵認証対応

### 公開鍵認証のTTL生成

- [ ] `internal/generator/generator.go` 更新
  - 公開鍵認証処理生成
  - `connect` コマンドに `/keyfile=<path>` オプション追加

- [ ] テストフィクスチャ追加
  - `test/fixtures/valid/keyfile-auth.yml`
  - `test/fixtures/expected/keyfile-auth.ttl`

- [ ] `internal/generator/generator_test.go` 更新
  - 公開鍵認証のテスト追加

---

## フェーズ5: CLIコマンド

### CLI基盤構築

- [ ] `internal/cli/root.go` 作成
  - ルートコマンド定義
  - `Execute()` 関数実装

- [ ] `cmd/ttlx/main.go` 作成
  - `main()` 関数実装
  - `cli.Execute()` 呼び出し

---

### build コマンド

- [ ] `internal/cli/build.go` 作成
  - `buildCmd` 定義
  - コマンドライン引数解析
  - `LoadConfig` 呼び出し
  - `Validate` 呼び出し
  - `Generate` 呼び出し
  - ファイル出力
  - `--output` フラグ実装
  - `--dry-run` フラグ実装

- [ ] 動作確認
  - `go run cmd/ttlx/main.go build test/fixtures/valid/simple.yml`
  - 出力ファイルの確認

---

### validate コマンド

- [ ] `internal/cli/validate.go` 作成
  - `validateCmd` 定義
  - `LoadConfig` 呼び出し
  - `Validate` 呼び出し
  - 成功/失敗メッセージ表示

- [ ] 動作確認
  - `go run cmd/ttlx/main.go validate test/fixtures/valid/simple.yml`
  - `go run cmd/ttlx/main.go validate test/fixtures/invalid/missing-version.yml`

---

### version コマンド

- [ ] `internal/cli/version.go` 作成
  - `versionCmd` 定義
  - バージョン情報表示

- [ ] バージョン定数定義
  - `internal/cli/version.go` に `Version` 定数

- [ ] 動作確認
  - `go run cmd/ttlx/main.go version`

---

## フェーズ6: 統合・テスト

### 統合テスト

- [ ] `test/integration/build_test.go` 作成
  - `TestBuild_SimpleConfig` 実装
  - `TestBuild_FullConfig` 実装
  - `TestBuild_InvalidConfig` 実装

- [ ] 統合テスト実行
  - `go test ./test/integration/...`

---

### カバレッジ確認

- [ ] テストカバレッジ測定
  - `go test -coverprofile=coverage.out ./...`
  - `go tool cover -html=coverage.out`

- [ ] カバレッジ80%以上確認
  - 不足している部分のテスト追加

---

### リント・フォーマット

- [ ] golangci-lint 実行
  - `golangci-lint run`
  - エラー修正

- [ ] gofmt 実行
  - `gofmt -w .`

- [ ] goimports 実行
  - `goimports -w .`

---

### ビルド・動作確認

- [ ] ローカルビルド
  - `go build -o ttlx cmd/ttlx/main.go`

- [ ] `ttlx build` 動作確認
  - `./ttlx build test/fixtures/valid/simple.yml`
  - 生成されたTTLファイルの確認

- [ ] `ttlx validate` 動作確認
  - `./ttlx validate test/fixtures/valid/simple.yml`
  - `./ttlx validate test/fixtures/invalid/missing-version.yml`

---

### 手動テスト（Tera Term）

- [ ] Tera Term での実行テスト
  - 生成されたTTLをTera Termで実行
  - 2段SSH接続が成功することを確認
  - コマンドが実行されることを確認

- [ ] パスワード認証テスト
  - 環境変数からのパスワード読み込み確認
  - 実行時入力（`passwordbox`）の確認

- [ ] 公開鍵認証テスト
  - 秘密鍵ファイルを使った接続確認

---

## ドキュメント整備

### README.md 作成

- [ ] `README.md` 作成
  - プロジェクト概要
  - インストール方法
  - 使い方（基本的な例）
  - ライセンス情報

---

### CONTRIBUTING.md 作成

- [ ] `CONTRIBUTING.md` 作成
  - コントリビューション方法
  - 開発環境のセットアップ
  - テストの実行方法

---

### CHANGELOG.md 作成

- [ ] `CHANGELOG.md` 作成
  - v1.0.0 の変更内容

---

## CI/CD セットアップ

### GitHub Actions

- [ ] `.github/workflows/test.yml` 作成
  - PR時のテスト実行
  - golangci-lint 実行
  - カバレッジレポート

- [ ] `.github/workflows/release.yml` 作成（将来用）
  - タグプッシュ時のリリースビルド
  - GoReleaser 設定

---

## リリース準備

### バージョン 1.0.0 リリース

- [ ] すべてのテストが通過
  - `go test ./...`

- [ ] リントエラーなし
  - `golangci-lint run`

- [ ] カバレッジ80%以上

- [ ] ドキュメント完成
  - README.md
  - CONTRIBUTING.md
  - CHANGELOG.md

- [ ] Gitタグ作成
  - `git tag v1.0.0`
  - `git push origin v1.0.0`

- [ ] GitHub Release 作成
  - リリースノート記載
  - ビルド済みバイナリ添付（手動またはGoReleaser）

---

## 進捗サマリー

### 完了したフェーズ
- [x] フェーズ1: 基盤構築
- [x] フェーズ2: バリデーション
- [x] フェーズ3: TTL生成（パスワード認証）
- [x] フェーズ4: 公開鍵認証対応
- [x] フェーズ5: CLIコマンド
- [x] フェーズ6: 統合・テスト（手動テストを除く）
- [x] ドキュメント整備
- [x] CI/CD セットアップ
- [ ] リリース準備（タグ作成・GitHub Release作成）

---

## 備考

### タスク管理のルール

- タスクは上から順に実施する
- 各タスク完了時にチェックボックスをチェック
- ブロッカーが発生した場合は、本ファイルに記録
- 新しいタスクが発生した場合は、適切なセクションに追加

### 完了条件

各タスクの完了条件：
- コードが実装されている
- テストが実装され、成功している
- リントエラーがない
- 必要に応じてドキュメントが更新されている

---

## まとめ

本タスクリストでは、ttlx の初回実装を段階的に進めるための具体的なタスクを定義しました。各タスクを順次完了させることで、確実にPhase 1 MVPを完成させることができます。
