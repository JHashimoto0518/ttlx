# 要求定義: 環境変数認証からgetpasswordへの移行

## 概要

環境変数を使用したパスワード認証（`auth.env`）を廃止し、Tera Term標準の`getpassword`コマンドを使用した実行時パスワード入力方式に統一する。

## 背景・目的

### 現状の課題

**環境変数認証（`auth.env`）の問題点:**
1. **セキュリティリスク**
   - 環境変数はプロセスリストに露出する可能性がある
   - ログに誤って記録されるリスク
   - シェル履歴に残る可能性

2. **一貫性の欠如**
   - TTL標準のパスワード管理方法ではない
   - `expandenv`という特殊なアプローチを使用

3. **運用の複雑さ**
   - 実行前に環境変数を設定する必要がある
   - 設定方法がTera Termの標準的な使い方ではない

**`passwordbox`の問題点:**
1. **UXの低さ**
   - `getpassword`と比較してユーザー体験が劣る

### 改善の方針

**`getpassword`を使用する利点:**
1. **セキュリティ向上**
   - パスワードファイルで一元管理
   - 環境変数への露出を回避
   - TTL標準のセキュアなパスワード管理

2. **TTL標準準拠**
   - Tera Term公式のパスワード管理方法
   - パスワードファイル方式を採用

3. **シンプルな運用**
   - パスワードファイルで複数のパスワードを管理
   - プロファイル名で自動的にパスワードを識別
   - 環境変数の事前設定が不要

## ユーザーストーリー

### US-1: セキュアなパスワード管理
**As a** セキュリティを重視するユーザー

**I want** パスワードをTTLファイルや環境変数に保存せず、パスワードファイルで管理したい

**So that** パスワードが漏洩するリスクを最小化できる

**受け入れ条件:**
- [ ] TTLファイルにパスワードが含まれない
- [ ] 環境変数にパスワードを設定する必要がない
- [ ] パスワードファイルで一元管理できる

### US-2: シンプルな設定
**As a** ttlxのユーザー

**I want** パスワード認証の設定をシンプルにしたい

**So that** 設定ファイルの記述が簡潔になり、理解しやすくなる

**受け入れ条件:**
- [ ] パスワード認証方法が2種類に整理される（直接指定/プロンプト）
- [ ] 環境変数関連の設定が不要になる
- [ ] ドキュメントがシンプルになる

### US-3: TTL標準準拠
**As a** Tera Termユーザー

**I want** TTL標準の方法でパスワードを管理したい

**So that** Tera Termの公式ドキュメントと整合性が取れる

**受け入れ条件:**
- [ ] `getpassword`コマンドを使用している
- [ ] Tera Term公式の推奨方法に従っている

## 機能要件

### FR-1: auth.envの削除
- `auth.env`フィールドをconfigスキーマから削除
- 環境変数認証関連のコードを削除
- `expandenv`を使用したパスワード認証を削除

### FR-2: passwordboxの削除
- `passwordbox`を使用したパスワード認証を削除
- `auth.prompt`フィールドを削除

### FR-3: getpasswordの採用
- パスワードファイル方式の`getpassword`を使用
- `password_file`フィールドを追加（オプション）
- password nameはプロファイル名を自動使用

### FR-4: パスワード認証方法の整理
変更後のパスワード認証方法を2つに統一:

#### 1. 直接指定（auth.value）
```yaml
auth:
  type: password
  value: "password123"
```
- **用途**: テスト、自動化、非対話的な実行
- **セキュリティ**: 低（パスワードが平文で保存される）
- **生成されるTTL**:
  ```ttl
  sendln 'password123'
  ```

#### 2. パスワードファイル（getpassword）
```yaml
auth:
  type: password
  password_file: "C:\\passwords.dat"  # optional
```
- **用途**: 本番環境、セキュアなパスワード管理
- **セキュリティ**: 高（パスワードファイルで一元管理）
- **デフォルト値**: `password_file`が省略された場合は`"passwords.dat"`
- **password name**: プロファイル名を自動使用（例: `bastion`）
- **生成されるTTL（第1ステップ - connect）**:
  ```ttl
  ; Password authentication (from password file)
  getpassword 'C:\passwords.dat' 'bastion' password

  ; Build connect command with password
  strconcat connectcmd 'localhost:22220 /ssh /auth=password /user=testuser /passwd='
  strconcat connectcmd password
  connect connectcmd
  if result <> 2 then
      goto ERROR_CONNECT_BASTION
  endif
  wait '$ '
  if result = 0 then
      goto TIMEOUT_BASTION
  endif
  ```
- **生成されるTTL（第2ステップ以降 - ssh）**:
  ```ttl
  sendln 'ssh testuser@target -p 2222'
  wait 'password:'
  if result = 0 then
      goto TIMEOUT_TARGET
  endif

  ; Password authentication (from password file)
  getpassword 'C:\passwords.dat' 'target' password
  sendln password
  ```

### FR-5: password_promptフィールドの扱い
- `password_prompt`フィールドは2段階目以降のSSH接続で引き続き使用
- パスワードプロンプトの文字列を指定するために必要

## 非機能要件

### NFR-1: 後方互換性
- **破壊的変更**: `auth.env`を使用している既存の設定ファイルは動作しなくなる
- **対応**: ベータ版（0.1.0-beta）なので破壊的変更は許容される
- **マイグレーション**: 考慮不要

### NFR-2: テストカバレッジ
- 既存のテストを更新
- 新しい実装のテストを追加
- E2Eテストケースを更新

### NFR-3: ドキュメント
- README.mdの更新
- CHANGELOG.mdへの破壊的変更の記載
- E2Eテストドキュメントの更新

### NFR-4: Tera Termバージョン互換性
- `getpassword`コマンドを使用（Tera Term 4.x以降で利用可能）
- `getpassword2`は使用しない（バージョン5.3以降が必須のため）

## 制約事項

### 技術的制約
- Tera Termの`getpassword`コマンド仕様に準拠する必要がある
- 既存のconfigスキーマバリデーションロジックを更新

### スコープ外
- パスワードの暗号化保存機能（将来的な検討事項）
- パスワードマネージャーとの連携

## 影響範囲

### 削除されるファイル/機能
- 環境変数認証関連のコード:
  - `connectWithEnvPasswordTemplate`
  - `passwordEnvTemplate`（環境変数版）
  - `auth.env`のサポート
- プロンプト認証関連のコード:
  - `passwordPromptTemplate`
  - `auth.prompt`のサポート
- テストファイル:
  - `test/e2e/configs/02-password-env.yml`
  - `test/e2e/configs/05-multi-hop-env.yml`

### 修正が必要なファイル
- `internal/config/config.go`
  - `auth.env`フィールド削除
  - `auth.prompt`フィールド削除
  - `auth.password_file`フィールド追加
- `internal/config/validator.go`
  - バリデーションロジック更新
  - `password_file`のデフォルト値設定
- `internal/generator/template.go`
  - `connectWithEnvPasswordTemplate`削除
  - `passwordEnvTemplate`削除
  - `passwordPromptTemplate`削除
  - `getpassword`用の新しいテンプレート追加
- `internal/generator/generator.go`
  - 環境変数認証ロジック削除
  - プロンプト認証ロジック削除
  - `getpassword`生成ロジック追加
- テストファイル群
  - 環境変数認証テストの削除
  - プロンプト認証テストの削除
  - `getpassword`テストの追加
- README.md - パスワード認証の説明更新
- CHANGELOG.md - 破壊的変更の記載

## 成功の定義

1. **機能面**
   - [ ] `auth.env`を使用した設定がエラーになる（バリデーション）
   - [ ] `auth.prompt`を使用した設定がエラーになる（バリデーション）
   - [ ] `password_file`が省略された場合、デフォルト値`"passwords.dat"`が使用される
   - [ ] password nameとしてプロファイル名が自動使用される
   - [ ] `getpassword`コマンドが正しく生成される
   - [ ] 既存の`auth.value`は引き続き動作する

2. **品質面**
   - [ ] 全テストがパスする
   - [ ] Linterチェックがパスする
   - [ ] E2Eテストが更新され、動作確認できる

3. **ドキュメント面**
   - [ ] CHANGELOG.mdに破壊的変更が記載される
   - [ ] README.mdが更新される
   - [ ] E2Eテストドキュメントが更新される

## 備考

- この変更はベータ版での実施のため、影響は限定的
- よりセキュアで標準的な実装に移行できる
- 将来的なパスワード管理機能拡張の基盤となる
- `getpassword`を使用（Tera Term 4.x以降で利用可能）
- `getpassword2`は採用しない（バージョン5.3以降が必須のため）
