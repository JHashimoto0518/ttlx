# ttlx - Tera Term Language eXtended

**Read this in other languages: 日本語 | [English](README.en.md)**

> [!WARNING]
> 本プロジェクトは開発中です。APIや機能は予告なく変更される可能性があります。
> 本番環境での使用は推奨しません。

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

ttlxは、YAML設定ファイルからTera Termマクロ（TTL）スクリプトを生成するツールです。
共通の接続設定やコマンドをプロファイル化することで、複数のスクリプト間での再利用を可能にし、保守性を向上させます。

## Tera Term Language機能の対応状況

| TTL機能カテゴリ | 対応状況 | 説明 |
|----------------|---------|------|
| **SSH接続** | ✅ 対応 | 多段SSH接続（踏み台サーバー経由） |
| **認証** | ✅ 対応 | パスワード認証（環境変数/実行時入力/直接指定）<br>公開鍵認証 |
| **コマンド実行** | ✅ 対応 | 接続後の任意コマンド実行 |
| **エラーハンドリング** | ✅ 対応 | タイムアウト処理、接続失敗時の処理 |
| **ファイル転送** | 🔄 未対応 | 将来対応予定 |
| **ダイアログ表示** | ⚠️ 部分対応 | パスワード入力、エラーメッセージのみ |
| **変数操作** | ⚠️ 部分対応 | 環境変数読み込みのみ |
| **ループ・分岐** | 🔄 未対応 | 将来対応予定 |

## 特徴

- 📝 **YAML設定**: シンプルで読みやすいYAML形式でSSHルートとコマンドを定義
- 🔐 **複数の認証方式**: パスワード認証と公開鍵認証をサポート
- 🔗 **多段SSH接続**: 踏み台サーバーやプロキシサーバー経由の接続を自動化
- ✅ **バリデーション**: わかりやすいエラーメッセージ付きの設定検証機能
- 🎯 **型安全**: Goの型システムを活用した堅牢なコード生成
- 🧪 **高品質**: テストカバレッジ97.8%

## インストール

### ソースからビルド

```bash
git clone https://github.com/JHashimoto0518/ttlx.git
cd ttlx
go build -o ttlx cmd/ttlx/main.go
```

### go installを使用

```bash
go install github.com/JHashimoto0518/ttlx/cmd/ttlx@latest
```

## クイックスタート

### 1. YAML設定ファイルを作成

```yaml
version: "1.0"

profiles:
  bastion:
    host: bastion.example.com
    user: user1
    auth:
      type: password
      prompt: true

  target:
    host: 10.0.0.50
    user: user2
    auth:
      type: password
      env: TARGET_PASSWORD

route:
  - profile: bastion
    commands:
      - echo "踏み台サーバーに接続しました"

  - profile: target
    commands:
      - ps aux
      - df -h
```

### 2. TTLスクリプトを生成

```bash
ttlx build config.yml
```

これにより、Tera Termで実行可能な `config.ttl` が生成されます。

### 3. 設定ファイルを検証

```bash
ttlx validate config.yml
```

## 設定方法

### プロファイル設定

各プロファイルはSSH接続先を定義します：

```yaml
profiles:
  server_name:
    host: ホスト名またはIP    # 必須
    port: 22                  # オプション、デフォルト: 22
    user: ユーザー名           # 必須
    auth:                     # 必須
      type: password|keyfile
      # ... 認証方式固有の設定
```

### 認証方式

#### パスワード認証

```yaml
auth:
  type: password
  # 以下のいずれかを選択:
  prompt: true              # 実行時にパスワードを入力
  env: ENV_VAR_NAME        # 環境変数から読み込み
  value: "password"        # パスワードを直接記述（非推奨）
```

#### 公開鍵認証

```yaml
auth:
  type: keyfile
  path: ~/.ssh/id_rsa      # 秘密鍵ファイルのパス
```

### ルート設定

SSH接続の順序を定義します：

```yaml
route:
  - profile: bastion       # 1段目
    commands:              # 実行するコマンド（オプション）
      - su - root
      - cd /var/log

  - profile: target        # 2段目
    commands:
      - ps aux
```

### グローバルオプション

```yaml
options:
  timeout: 30              # 接続タイムアウト（秒）、デフォルト: 30
  retry: 3                 # リトライ回数（未実装）
  log: true                # ログ有効化（未実装）
  log_file: /tmp/ttlx.log  # ログファイルパス（未実装）
```

## CLIコマンド

### build

YAML設定からTTLスクリプトを生成：

```bash
ttlx build <config.yml> [フラグ]

フラグ:
  -o, --output string   出力ファイルパス（デフォルト: <config>.ttl）
      --dry-run         ファイルではなく標準出力に出力
```

### validate

YAML設定を検証：

```bash
ttlx validate <config.yml>
```

### version

バージョン情報を表示：

```bash
ttlx version
```

## サンプル

より多くのサンプルは [test/fixtures/valid](test/fixtures/valid) ディレクトリを参照してください：

- [simple.yml](test/fixtures/valid/simple.yml) - 基本的な2段SSH接続
- [full.yml](test/fixtures/valid/full.yml) - 全機能を使用した設定例

## 開発

### 必要要件

- Go 1.21以降
- golangci-lint（リント用）

### ビルド

```bash
go build -o ttlx cmd/ttlx/main.go
```

### テスト

```bash
# 全テスト実行
go test ./...

# カバレッジ付きテスト実行
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# リント実行
golangci-lint run
```

### プロジェクト構造

```
ttlx/
├── cmd/ttlx/          # CLIエントリーポイント
├── internal/
│   ├── cli/           # CLIコマンド
│   ├── config/        # 設定処理
│   └── generator/     # TTL生成
├── test/
│   ├── fixtures/      # テストデータ
│   └── integration/   # 統合テスト
└── docs/              # ドキュメント
```

## コントリビューション

コントリビューションを歓迎します！詳細は [CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 謝辞

- [Tera Term Project](https://teratermproject.github.io/) - このツールがスクリプトを生成する対象のターミナルエミュレーター
- [Cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [go-yaml](https://github.com/go-yaml/yaml) - YAMLパーサー

## サポート

- 🐛 [バグを報告](https://github.com/JHashimoto0518/ttlx/issues)
- 💡 [機能をリクエスト](https://github.com/JHashimoto0518/ttlx/issues)
- 📖 [ドキュメント](docs/)
