# 要件定義: 複数ルート対応

## 概要

現在のttlx仕様では、1つのYAMLファイルに1つの`route`しか定義できないため、同じプロファイル（例: bastion）を使用して異なるターゲットサーバーに接続するルートを作成する際に、プロファイルの再利用ができません。

この機能では、1つのYAMLファイルで複数のルートを定義できるようにし、プロファイルの再利用性を向上させます。

## ユーザーストーリー

### US-1: プロファイルを再利用した複数ルートの定義

**As a** ネットワーク管理者
**I want to** 同じ踏み台サーバー（bastion）を経由して複数のターゲットサーバーに接続するルートを1つのYAMLファイルで定義したい
**So that** プロファイル設定を重複させることなく、複数の接続パターンを管理できる

**受け入れ条件:**
- [ ] 1つのYAMLファイルで複数のルートを定義できる
- [ ] 各ルートに一意な名前を付けられる
- [ ] 複数のルート間でプロファイルを再利用できる
- [ ] 各ルートに対して個別のTTLファイルが生成される

### US-2: 複数TTLファイルの自動生成

**As a** ttlxユーザー
**I want to** `ttlx build`コマンドを実行すると、定義された全ルートに対応するTTLファイルが自動生成される
**So that** 手動で複数のYAMLファイルを作成・管理する手間を省ける

**受け入れ条件:**
- [ ] `ttlx build config.yml`で複数のTTLファイルが生成される
- [ ] 生成されるTTLファイル名はルート名に基づく（例: `target1-connection.ttl`）
- [ ] 出力先ディレクトリを指定できる（`-o`オプション）
- [ ] 生成されたファイル一覧が表示される

## 機能要件

### YAML仕様

```yaml
version: "1.0"

profiles:
  bastion:
    host: bastion.example.com
    user: user1
    prompt_marker: "$ "
    auth:
      type: password
      prompt: true

  target1:
    host: target1.example.com
    user: user2
    prompt_marker: "$ "
    auth:
      type: password
      env: TARGET1_PASSWORD
      password_prompt: "password:"

  target2:
    host: target2.example.com
    user: user3
    prompt_marker: "$ "
    auth:
      type: password
      env: TARGET2_PASSWORD
      password_prompt: "password:"

routes:
  target1-connection:
    - profile: bastion
      commands:
        - echo "Connected to bastion"
    - profile: target1
      commands:
        - hostname
        - df -h

  target2-connection:
    - profile: bastion
      commands:
        - echo "Connected to bastion"
    - profile: target2
      commands:
        - hostname
        - ps aux

options:
  timeout: 30
  auto_disconnect: false
```

### CLIコマンドの動作

#### 基本的な使い方

```bash
$ ttlx build config.yml

Generated TTL files:
  - target1-connection.ttl
  - target2-connection.ttl
```

#### 出力先ディレクトリ指定

```bash
$ ttlx build config.yml -o output/

Generated TTL files in output/:
  - target1-connection.ttl
  - target2-connection.ttl
```

### バリデーション要件

- [ ] `routes`は必須フィールド（定義されていない場合はエラー）
- [ ] `routes`には最低1つのルートが必要
- [ ] 各ルート名は一意である必要がある
- [ ] ルート名は有効なファイル名として使用できる文字のみ（英数字、ハイフン、アンダースコア）
- [ ] 旧仕様の`route`（単数形）が使われている場合は分かりやすいエラーメッセージを表示
- [ ] 各ルートは最低1つのステップを持つ必要がある
- [ ] 各ルート内のプロファイル参照は有効である必要がある

## 非機能要件

### パフォーマンス

- [ ] 複数TTL生成時も実行時間は許容範囲内（10ルート以下で1秒以内）

### 保守性

- [ ] 既存のコード構造を大きく変更しない
- [ ] テストカバレッジを維持（95%以上）

### ユーザビリティ

- [ ] エラーメッセージは具体的で分かりやすい
- [ ] 生成されたファイル一覧が明確に表示される

## 制約事項

### 技術的制約

- Go 1.21以上
- gopkg.in/yaml.v3でのパース互換性維持

### ビジネス制約

- リリース前のため、破壊的変更を許容
- 旧仕様からの移行エラーメッセージを分かりやすくする

### スコープ外

- ルート間の依存関係管理（例: route1が成功したらroute2を実行）
- 並列TTL実行機能
- ルートの条件分岐（if/else）

## 成功基準

- [ ] 1つのYAMLファイルで複数ルートを定義できる
- [ ] プロファイルが複数ルート間で再利用できる
- [ ] 旧仕様の`route`使用時に分かりやすいエラーメッセージが表示される
- [ ] すべてのユニットテストが成功する
- [ ] テストカバレッジが95%以上を維持する
- [ ] ドキュメント（README、functional-design.md）が更新されている
- [ ] 既存のフィクスチャファイルを新仕様に移行する
