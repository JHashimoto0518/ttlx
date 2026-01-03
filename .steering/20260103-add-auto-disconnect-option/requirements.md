# 自動切断オプションの追加

## 概要

SSH接続完了後の動作モードを選択できる機能を追加します。ユーザーが「完全自動化」と「手動操作」のどちらのワークフローを使用するかを、YAML設定で制御できるようにします。

---

## 背景

現在のttlxの設計では、SSH接続後にコマンドを実行することは想定されていますが、その後の動作（自動切断 vs 接続保持）が明確に定義されていません。

ユーザーのワークフローには以下の2つのパターンがあります：

1. **完全自動化ワークフロー**
   - 定型的なコマンドを実行して結果を取得する
   - 例: ログ収集、ステータス確認、バッチ処理
   - コマンド実行後、自動的に切断してスクリプト終了

2. **半自動化ワークフロー**
   - サーバーに接続するまでの手順を自動化
   - 接続後は手動で作業を行う
   - 例: トラブルシューティング、設定変更、調査作業
   - ユーザーが `exit` で終了するまで接続を保持

両方のワークフローをサポートすることで、ttlxの適用範囲を広げます。

---

## 要求事項

### 機能要件

#### FR-AD-001: auto_disconnect オプションの追加

**Optionsレベル**（グローバルオプション）で `auto_disconnect` オプションを提供します。

**YAML設定例**:
```yaml
route:
  - profile: bastion
    commands:
      - echo "Connected to bastion"

  - profile: target
    commands:
      - ls -la
      - pwd

options:
  auto_disconnect: true  # 最終ステップ完了後、自動的に切断してスクリプト終了
```

**オプション仕様**:
- `options.auto_disconnect: true` - 最終ステップのコマンド実行後、自動的にSSH接続を切断し、TTLスクリプトを終了
- `options.auto_disconnect: false` - 最終ステップのコマンド実行後も接続を保持し、ユーザーの手動操作を待つ
- 未指定の場合のデフォルト値: `false`（接続保持）

---

#### FR-AD-002: デフォルト動作の定義

`options.auto_disconnect` が未指定の場合のデフォルト動作: `false`（接続保持）

**理由**:
- 半自動化ワークフロー（手動操作）がより一般的なユースケース
- デフォルトで接続を保持することで、初心者が混乱しない
- 完全自動化したいユーザーは明示的に `options.auto_disconnect: true` を指定

**デフォルト動作の例**:
```yaml
route:
  - profile: bastion
  - profile: target
    commands:
      - ls -la

# options.auto_disconnect 未指定 → デフォルト false（接続保持）
```

---

#### FR-AD-003: TTL生成ロジックの変更

**auto_disconnect: true の場合**:
- コマンド実行後、`exit` コマンドを送信してSSH接続を切断
- 多段接続の場合、全ての段階を順次切断
- 最終的に Tera Term を終了（`closett` コマンド）

**auto_disconnect: false の場合**:
- コマンド実行後、何もせず接続を保持
- ユーザーが手動で `exit` を入力するまで待機
- TTLスクリプトは終了するが、Tera Termのセッションは維持

**生成されるTTLの例（auto_disconnect: true の場合）**:
```ttl
; 最終ステップのコマンド実行
sendln 'ls -la'
wait prompt_marker
sendln 'pwd'
wait prompt_marker

; 自動切断処理
sendln 'exit'
wait 'Connection closed'

; スクリプト終了
end
```

**生成されるTTLの例（auto_disconnect: false の場合）**:
```ttl
; 最終ステップのコマンド実行
sendln 'ls -la'
wait prompt_marker
sendln 'pwd'
wait prompt_marker

; 接続保持（何もしない）
; ユーザーが手動で操作を続けられる

; スクリプト終了（Tera Termセッションは維持）
end
```

---

### 非機能要件

#### NFR-AD-001: 後方互換性

- 既存のYAML設定ファイル（`auto_disconnect` 未指定）は、引き続き動作する
- デフォルト動作（接続保持）により、既存ユーザーの期待動作を維持

#### NFR-AD-002: バリデーション

- `auto_disconnect` フィールドはboolean型のみを受け入れる
- 不正な値（文字列、数値など）が指定された場合、バリデーションエラーを表示

---

## ユーザーストーリー

### US-AD-001: 完全自動化ワークフロー

**As a** インフラエンジニア

**I want to** サーバーに接続してコマンドを実行した後、自動的に切断したい

**So that** スクリプトを無人実行できる

**受け入れ条件**:
- `options.auto_disconnect: true` を指定できる
- 最終ステップのコマンド実行後、自動的にSSH接続が切断される
- TTLスクリプトが終了する

**例**:
```yaml
route:
  - profile: bastion
  - profile: target
    commands:
      - /opt/backup.sh

options:
  auto_disconnect: true  # バックアップ実行後、自動切断
```

---

### US-AD-002: 半自動化ワークフロー

**As a** サーバー管理者

**I want to** サーバーに接続した後、手動で作業を続けたい

**So that** トラブルシューティングや調査作業ができる

**受け入れ条件**:
- `options.auto_disconnect: false` を指定できる（またはデフォルトで省略可）
- 最終ステップのコマンド実行後も接続が保持される
- ユーザーが手動で `exit` を入力するまでセッションが維持される

**例**:
```yaml
route:
  - profile: bastion
  - profile: target
    commands:
      - echo "Ready for manual work"

options:
  auto_disconnect: false  # 手動作業のため接続保持（デフォルトなので省略可）
```

---

## 影響範囲

### 変更が必要なコンポーネント

1. **internal/config/model.go**
   - `Options` 構造体に `AutoDisconnect *bool` フィールドを追加
   - `SetDefaults()` メソッドでデフォルト値（`false`）を設定

2. **internal/config/validator.go**
   - `options.auto_disconnect` フィールドのバリデーション追加（型チェック）

3. **internal/generator/ttl_generator.go**
   - 最終ステップ完了後の切断処理生成ロジック追加
   - `options.auto_disconnect` フラグに基づく条件分岐
   - `auto_disconnect: true` の場合: `exit` コマンド送信と接続切断処理
   - `auto_disconnect: false` の場合: 何もしない（接続保持）

4. **test/fixtures/**
   - テスト用YAMLファイルの追加
   - `options.auto_disconnect: true` のケース
   - `options.auto_disconnect: false` のケース
   - デフォルト動作のケース（未指定）

5. **ドキュメント**
   - `docs/product-requirements.md` - 機能要件に追加
   - `docs/functional-design.md` - YAML仕様に追加
   - `README.md` / `README.en.md` - 使用例の更新

---

## 成功基準

- [ ] `options.auto_disconnect: true` を指定した場合、最終ステップ完了後に自動切断される
- [ ] `options.auto_disconnect: false` を指定した場合、接続が保持される
- [ ] デフォルト（未指定）の場合、接続が保持される（デフォルト: `false`）
- [ ] 生成されたTTLが Tera Term で正しく動作する
- [ ] バリデーションが不正な値を検出する
- [ ] ユニットテストが追加され、カバレッジが維持される（80%以上）

---

## リスクと対策

### リスク1: Tera Term の挙動の不確実性

**発生確率**: 中

**影響度**: 高

**対策**:
- 実際に Tera Term でテストして動作確認
- `closett` コマンドの動作を検証
- 切断処理のタイミングを調整（wait コマンドの使用）

---

### リスク2: デフォルト動作の選択ミス

**発生確率**: 低

**影響度**: 中

**対策**:
- ユーザーにデフォルト動作（false）を明確に説明
- ドキュメントに使用例を豊富に記載
- バリデーション時に、デフォルト値を使用していることを通知（オプション）

---

## まとめ

`auto_disconnect` オプションを追加することで、ttlxは「完全自動化」と「半自動化」の両方のワークフローをサポートできるようになります。デフォルトで接続を保持することで、既存ユーザーの期待動作を維持しつつ、自動化したいユーザーには明示的なオプションを提供します。
