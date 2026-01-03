# ユビキタス言語定義（用語集）

## 概要

本ドキュメントでは、ttlx プロジェクトで使用する用語を定義します。すべての開発者、ドキュメント、コード、コミュニケーションにおいて、これらの用語を統一して使用します。

---

## ドメイン用語

### Tera Term 関連

#### Tera Term
**日本語**: テラターム

**定義**: Windows用のオープンソース端末エミュレータ。SSH、Telnet、シリアル接続をサポート。

**関連用語**: Tera Term Macro, TTL

---

#### TTL (Tera Term Language)
**日本語**: ティーティーエル

**定義**: Tera Term のマクロ言語。拡張子は `.ttl`。接続、コマンド実行、ファイル転送などを自動化できる。

**例**:
```ttl
connect 'example.com:22 /ssh'
sendln 'ls -la'
```

**関連用語**: Tera Term Macro

---

#### Tera Term Macro
**日本語**: テラタームマクロ

**定義**: TTLで記述されたスクリプト。Tera Termで実行することで、一連の操作を自動化する。

**同義語**: TTLスクリプト、マクロ

---

### SSH 関連

#### SSH (Secure Shell)
**日本語**: エスエスエイチ

**定義**: 暗号化された通信でリモートサーバーに接続するプロトコル。

**関連用語**: SSH接続、公開鍵認証、パスワード認証

---

#### 踏み台サーバー（Bastion Host）
**日本語**: ふみだいサーバー

**英語**: Bastion Host, Jump Server

**定義**: セキュリティのために、目的のサーバーへのアクセスを中継するサーバー。外部からは踏み台サーバーにのみ接続可能で、踏み台サーバー経由で内部サーバーにアクセスする。

**コード上の表記**: `bastion`, `jump`

**例**:
```
外部 → 踏み台サーバー → 目的のサーバー
```

---

#### 多段SSH（Multi-hop SSH）
**日本語**: ただんエスエスエイチ

**英語**: Multi-hop SSH, SSH Chaining

**定義**: 複数のサーバーを経由してSSH接続すること。踏み台サーバーを介して目的のサーバーに接続する場合など。

**同義語**: SSH多段接続、SSHチェイニング

**例**:
```
ローカル → 踏み台1 → 踏み台2 → 目的サーバー
```

---

#### 公開鍵認証（Public Key Authentication）
**日本語**: こうかいかぎにんしょう

**英語**: Public Key Authentication, Key-based Authentication

**定義**: SSH接続時に、パスワードの代わりに公開鍵/秘密鍵ペアを使用する認証方式。

**コード上の表記**: `keyfile`, `key-based`

**関連ファイル**: `~/.ssh/id_rsa`, `~/.ssh/id_ed25519`

---

#### パスワード認証（Password Authentication）
**日本語**: パスワードにんしょう

**英語**: Password Authentication

**定義**: SSH接続時に、ユーザー名とパスワードを使用する認証方式。

**コード上の表記**: `password`

---

### ttlx プロジェクト固有用語

#### ttlx
**読み方**: ティーティーエルエックス

**定義**: Tera Term Language eXtended の略。YAML設定ファイルからTTLスクリプトを生成するCLIツール。

**フルネーム**: Tera Term Language eXtended

---

#### プロファイル（Profile）
**日本語**: プロファイル

**英語**: Profile

**定義**: SSH接続先の情報（ホスト名、ユーザー名、認証方式など）をまとめた設定。YAML設定ファイルの `profiles` セクションで定義される。

**コード上の表記**: `Profile` (構造体)

**例**:
```yaml
profiles:
  bastion:
    host: bastion.example.com
    user: user1
    prompt_marker: "$ "
    auth:
      type: password
```

---

#### プロンプトマーカー（Prompt Marker）
**日本語**: プロンプトマーカー

**英語**: Prompt Marker

**定義**: プロンプトを識別するための文字列。TTL の `wait` コマンドで使用され、コマンド実行後にプロンプトが表示されたことを検出するために使用される。

**コード上の表記**: `PromptMarker` (フィールド), `prompt_marker` (YAML)

**例**:
```yaml
prompt_marker: "$ "   # 一般ユーザーのプロンプト
prompt_marker: "# "   # rootのプロンプト
prompt_marker: "mysql> "  # MySQLのプロンプト
```

---

#### ルート（Route）
**日本語**: ルート、経路

**英語**: Route

**定義**: 接続する順序を定義したもの。プロファイルのリストで、多段SSH接続の経路を表現する。

**コード上の表記**: `Route` (構造体), `RouteStep`

**例**:
```yaml
route:
  - profile: bastion
  - profile: target
```

---

#### ルートステップ（Route Step）
**日本語**: ルートステップ

**英語**: Route Step

**定義**: ルートの各段階。1つのプロファイルと、その接続後に実行するコマンドで構成される。

**コード上の表記**: `RouteStep` (構造体)

**例**:
```yaml
- profile: bastion
  commands:
    - su - root
```

---

#### グローバルプロファイル（Global Profile）
**日本語**: グローバルプロファイル

**英語**: Global Profile

**定義**: 複数の設定ファイルで共有されるプロファイル定義。通常 `~/.ttlx/profiles.yml` に保存される。

**Phase**: Phase 2

**関連用語**: ローカルプロファイル、プロファイルのインポート

---

#### メタデータファイル（Metadata File）
**日本語**: メタデータファイル

**英語**: Metadata File

**定義**: TTL生成時の設定スナップショット。差分表示機能で、前回の生成内容と比較するために使用される。JSON形式で `.ttlx/*.meta.json` に保存される。

**Phase**: Phase 2

**関連用語**: 差分表示

---

## 技術用語

### YAML
**読み方**: ヤムル

**正式名称**: YAML Ain't Markup Language

**定義**: 人間が読みやすいデータシリアライゼーション形式。ttlxでは設定ファイルのフォーマットとして使用。

**拡張子**: `.yml`, `.yaml`

---

### CLI (Command Line Interface)
**読み方**: シーエルアイ

**日本語**: コマンドラインインターフェース

**定義**: テキストベースでコマンドを実行するインターフェース。ttlxはCLIツール。

**対義語**: GUI (Graphical User Interface)

---

### バリデーション（Validation）
**日本語**: バリデーション、検証

**英語**: Validation

**定義**: データの妥当性を検証すること。ttlxでは、YAML設定ファイルの構文や内容が正しいかチェックする。

**コマンド**: `ttlx validate`

**コード上の表記**: `Validator`, `Validate()`

---

### ジェネレーター（Generator）
**日本語**: ジェネレーター、生成器

**英語**: Generator

**定義**: TTLスクリプトを生成するコンポーネント。YAML設定を読み込み、TTLコードを出力する。

**コード上の表記**: `Generator`, `Generate()`

---

### 差分（Diff）
**日本語**: さぶん

**英語**: Difference, Diff

**定義**: 2つのファイルやデータの違い。ttlxでは、設定変更前後のTTLスクリプトの差分を表示する。

**コマンド**: `ttlx diff`

**コード上の表記**: `Differ`, `Diff()`

---

### トランスパイル（Transpile）
**日本語**: トランスパイル

**英語**: Transpile

**定義**: ある言語のソースコードを、別の言語のソースコードに変換すること。ttlxでは、YAMLをTTLに変換する。

**関連用語**: コンパイル、変換

---

## コマンド・機能用語

### build
**コマンド**: `ttlx build <config.yml>`

**定義**: YAML設定ファイルからTTLスクリプトを生成するコマンド。

**例**:
```bash
ttlx build production-db.yml -o output.ttl
```

---

### validate
**コマンド**: `ttlx validate <config.yml>`

**定義**: YAML設定ファイルのバリデーションを実行するコマンド。

**例**:
```bash
ttlx validate production-db.yml
```

---

### diff
**コマンド**: `ttlx diff <config.yml>`

**Phase**: Phase 2

**定義**: 設定変更前後のTTL差分を表示するコマンド。

**例**:
```bash
ttlx diff production-db.yml
```

---

### init
**コマンド**: `ttlx init`

**Phase**: Phase 2

**定義**: 対話的に設定ファイルを作成するコマンド。

**例**:
```bash
ttlx init
? 踏み台サーバーのホスト名: bastion.example.com
```

---

## データ構造・コード用語

### Config
**型**: 構造体

**定義**: YAML設定ファイル全体を表すデータモデル。

**主要フィールド**:
- `Version`: バージョン番号
- `Profiles`: プロファイルのマップ
- `Route`: 接続ルート
- `Options`: オプション設定

**パッケージ**: `internal/config`

---

### Profile
**型**: 構造体

**定義**: SSH接続先の情報を表すデータモデル。

**主要フィールド**:
- `Host`: ホスト名
- `Port`: ポート番号
- `User`: ユーザー名
- `PromptMarker`: プロンプト識別文字列
- `Auth`: 認証設定

**パッケージ**: `internal/config`

---

### RouteStep
**型**: 構造体

**定義**: ルートの各段階を表すデータモデル。

**主要フィールド**:
- `Profile`: プロファイル名
- `Commands`: 実行するコマンドのリスト

**パッケージ**: `internal/config`

---

### Auth
**型**: 構造体

**定義**: 認証設定を表すデータモデル。

**主要フィールド**:
- `Type`: 認証タイプ（`password` | `keyfile`）
- `Value`: パスワード（直接記述）
- `Env`: 環境変数名
- `Prompt`: 実行時入力フラグ
- `Path`: 秘密鍵ファイルパス

**パッケージ**: `internal/config`

---

## 英語・日本語対応表

| 英語 | 日本語 | コード表記 | 備考 |
|------|--------|-----------|------|
| Profile | プロファイル | `Profile` | 接続設定 |
| Prompt Marker | プロンプトマーカー | `PromptMarker` | プロンプト識別文字列 |
| Route | ルート、経路 | `Route` | 接続順序 |
| Route Step | ルートステップ | `RouteStep` | ルートの各段階 |
| Bastion Host | 踏み台サーバー | `bastion` | 中継サーバー |
| Multi-hop SSH | 多段SSH | - | 複数サーバー経由 |
| Password Authentication | パスワード認証 | `password` | - |
| Public Key Authentication | 公開鍵認証 | `keyfile` | - |
| Validation | バリデーション、検証 | `Validate()` | - |
| Generator | ジェネレーター、生成器 | `Generator` | - |
| Diff | 差分 | `Diff()` | - |
| Metadata | メタデータ | `Metadata` | Phase 2 |
| Global Profile | グローバルプロファイル | - | Phase 2 |
| Template | テンプレート | `Template` | Phase 3 |

---

## 略語

| 略語 | 正式名称 | 説明 |
|------|----------|------|
| TTL | Tera Term Language | Tera Termのマクロ言語 |
| SSH | Secure Shell | 暗号化通信プロトコル |
| YAML | YAML Ain't Markup Language | データ記述言語 |
| CLI | Command Line Interface | コマンドライン操作 |
| PR | Pull Request | GitHubのプルリクエスト |
| CI/CD | Continuous Integration / Continuous Deployment | 継続的インテグレーション/デプロイ |
| MVP | Minimum Viable Product | 最小実用製品 |

---

## フェーズ別用語

### Phase 1 (MVP)
- プロファイル
- ルート
- ルートステップ
- バリデーション
- ジェネレーター
- build コマンド
- validate コマンド

### Phase 2
- グローバルプロファイル
- メタデータファイル
- 差分（Diff）
- diff コマンド
- init コマンド
- 変数展開

### Phase 3（将来）
- テンプレート
- Language Server Protocol (LSP)
- プラグインシステム

---

## 用語の使い方

### ドキュメント
- **ユーザー向け**: 日本語用語を使用
  - 例: 「踏み台サーバー経由で接続します」
- **開発者向け**: 英語用語を優先、必要に応じて日本語を併記
  - 例: 「`Profile` 構造体（プロファイル）」

### コード
- **変数名・関数名**: 英語のみ
  - 例: `profile`, `routeStep`, `LoadConfig()`
- **コメント**: 英語または日本語（統一する）
  - 例: `// Load YAML configuration file` または `// YAML設定ファイルを読み込む`

### コミットメッセージ・PR
- **英語を推奨**、日本語も可
  - 例: `feat(cli): add diff command`
  - 例: `fix(generator): 特殊文字のエスケープ処理を修正`

---

## まとめ

本用語集では、ttlx プロジェクトで使用する用語を統一的に定義しました。すべてのコミュニケーション、ドキュメント、コードにおいて、これらの用語を一貫して使用することで、チーム内の認識の齟齬を防ぎます。

新しい概念や用語が登場した場合は、この用語集を更新し、チーム全体で共有してください。
