# タスクリスト: 環境変数認証からgetpasswordへの移行

## 進捗サマリー

- **未着手**: 24タスク
- **進行中**: 0タスク
- **完了**: 0タスク

## フェーズ1: Config変更

### Task 1.1: Auth構造体の更新
**ファイル**: `internal/config/config.go`

**作業内容:**
- [ ] `Env string` フィールドを削除
- [ ] `Prompt bool` フィールドを削除
- [ ] `PasswordFile string` フィールドを追加
- [ ] yamlタグを正しく設定: `yaml:"password_file,omitempty"`

**完了条件:**
- Auth構造体が正しく更新されている
- ビルドエラーがない

---

### Task 1.2: バリデーションロジックの追加
**ファイル**: `internal/config/validator.go`

**作業内容:**
- [ ] password_fileのデフォルト値設定ロジックを追加
  - パスワード認証でvalueもpassword_fileも未指定の場合
  - デフォルト値: `"passwords.dat"`
- [ ] valueとpassword_fileの相互排他性チェックを追加
  - 両方指定されている場合はエラー

**完了条件:**
- デフォルト値が正しく設定される
- 相互排他性のバリデーションが機能する
- ビルドエラーがない

---

## フェーズ2: Template変更

### Task 2.1: 不要なテンプレートの削除
**ファイル**: `internal/generator/template.go`

**作業内容:**
- [ ] `connectWithEnvPasswordTemplate` を削除
- [ ] `passwordEnvTemplate` を削除
- [ ] `passwordPromptTemplate` を削除

**完了条件:**
- 削除したテンプレートへの参照がないことを確認
- ビルドエラーを確認（この時点では発生する）

---

### Task 2.2: 新しいテンプレートの追加
**ファイル**: `internal/generator/template.go`

**作業内容:**
- [ ] `passwordFileConnectTemplate` を追加（第1ステップ用）
  ```go
  passwordFileConnectTemplate = `; Password authentication (from password file)
  getpassword '%s' '%s' password

  ; Build connect command with password
  strconcat connectcmd '%s:%d /ssh /auth=%s /user=%s%s /passwd='
  strconcat connectcmd password
  connect connectcmd
  if result <> 2 then
      goto ERROR_CONNECT_%s
  endif
  wait '%s'
  if result = 0 then
      goto TIMEOUT_%s
  endif

  `
  ```
- [ ] `passwordFileTemplate` を追加（第2ステップ以降用）
  ```go
  passwordFileTemplate = `; Password authentication (from password file)
  getpassword '%s' '%s' password
  sendln password

  `
  ```

**完了条件:**
- テンプレートが正しく定義されている
- パラメータの数と順序が設計書と一致

---

## フェーズ3: Generator変更

### Task 3.1: generateConnect()の更新
**ファイル**: `internal/generator/generator.go`

**作業内容:**
- [ ] `profile.Auth.Env`のチェックを削除
- [ ] `profile.Auth.PasswordFile`のチェックを追加
- [ ] password_fileが設定されている場合、`passwordFileConnectTemplate`を使用するロジックを追加
- [ ] テンプレートパラメータを正しく渡す

**実装例:**
```go
// パスワードファイル認証の場合は専用テンプレートを使用
if authType == "password" && profile.Auth.PasswordFile != "" {
    return fmt.Sprintf(
        passwordFileConnectTemplate,
        profile.Auth.PasswordFile,
        profileName,
        profile.Host,
        profile.Port,
        authType,
        profile.User,
        keyfileOption,
        upperProfileName,
        profile.PromptMarker,
        upperProfileName,
    )
}
```

**完了条件:**
- password_file認証で正しいテンプレートが使用される
- ビルドエラーがない

---

### Task 3.2: generatePasswordAuth()の更新
**ファイル**: `internal/generator/generator.go`

**作業内容:**
- [ ] `auth.Env`のチェックを削除
- [ ] `auth.Prompt`のチェックを削除
- [ ] `auth.PasswordFile`のチェックを追加
- [ ] password_fileが設定されている場合、`passwordFileTemplate`を使用

**実装例:**
```go
func generatePasswordAuth(profileName string, auth *config.Auth) string {
    // password_fileが設定されている場合（デフォルト値含む）
    if auth.PasswordFile != "" {
        return fmt.Sprintf(passwordFileTemplate, auth.PasswordFile, profileName)
    }

    // 直接パスワード指定
    if auth.Value != "" {
        return fmt.Sprintf(passwordValueTemplate, auth.Value)
    }

    return ""
}
```

**完了条件:**
- password_file認証で正しいコードが生成される
- ビルドエラーがない

---

### Task 3.3: 第1ステップ生成ロジックの更新
**ファイル**: `internal/generator/generator.go`

**作業内容:**
- [ ] 第1ステップでのパスワード認証処理コードを削除
  - `profile.Auth.Type == "password" && profile.Auth.Value == "" && profile.Auth.Env == ""`の条件分岐を削除
  - `generatePasswordAuth()`の呼び出しを削除

**理由:**
- password_fileの場合、`generateConnect()`が完全なシーケンスを生成
- auth.valueの場合、connectコマンドに`/passwd=`が含まれている
- どちらの場合も追加のコードは不要

**完了条件:**
- 第1ステップでパスワード認証コードが重複生成されない
- ビルドエラーがない

---

## フェーズ4: テスト更新

### Task 4.1: 環境変数認証テストの削除
**ファイル**: `internal/generator/generator_test.go`, `test/integration/build_test.go`

**作業内容:**
- [ ] 環境変数認証に関するテストケースを削除
- [ ] `expandenv`に関するアサーションを削除

**完了条件:**
- 環境変数認証のテストが完全に削除されている

---

### Task 4.2: プロンプト認証テストの削除
**ファイル**: `internal/generator/generator_test.go`

**作業内容:**
- [ ] プロンプト認証（passwordbox）に関するテストケースを削除

**完了条件:**
- passwordboxのテストが完全に削除されている

---

### Task 4.3: getpassword生成のユニットテスト追加
**ファイル**: `internal/generator/generator_test.go`

**作業内容:**
- [ ] `TestGenerate_PasswordFile`を追加
  - password_file指定時の動作確認
  - `getpassword`コマンドの生成確認
  - `strconcat`コマンドの生成確認（第1ステップ）
  - `sendln password`の生成確認

**テストケース:**
```go
func TestGenerate_PasswordFile(t *testing.T) {
    cfg := &config.Config{
        Profiles: map[string]*config.Profile{
            "server": {
                Host: "example.com",
                Port: 22,
                User: "user",
                Auth: &config.Auth{
                    Type:         "password",
                    PasswordFile: "C:\\passwords.dat",
                },
                PromptMarker: "$ ",
            },
        },
    }

    route := []*config.RouteStep{
        {Profile: "server", Commands: []string{"ls"}},
    }

    ttl, err := generateRoute(cfg, "test", route, "test.yml")
    require.NoError(t, err)

    // getpasswordコマンドの確認
    assert.Contains(t, ttl, "getpassword 'C:\\passwords.dat' 'server' password")
    // strconcatの確認
    assert.Contains(t, ttl, "strconcat connectcmd")
    // sendln passwordの確認は不要（第1ステップではsendlnは使用しない）
}
```

**完了条件:**
- テストがパスする
- 必要なアサーションが含まれている

---

### Task 4.4: デフォルト値のユニットテスト追加
**ファイル**: `internal/generator/generator_test.go`

**作業内容:**
- [ ] `TestGenerate_PasswordFileDefault`を追加
  - password_file未指定時のデフォルト値確認
  - `getpassword 'passwords.dat'`が生成されることを確認

**完了条件:**
- テストがパスする
- デフォルト値が正しく適用される

---

### Task 4.5: 相互排他性バリデーションのテスト追加
**ファイル**: `internal/config/validator_test.go`

**作業内容:**
- [ ] `TestValidate_MutuallyExclusivePasswordFields`を追加
  - valueとpassword_fileの両方指定時にエラーになることを確認

**完了条件:**
- テストがパスする
- エラーメッセージが適切

---

### Task 4.6: 統合テストの更新
**ファイル**: `test/integration/build_test.go`

**作業内容:**
- [ ] 環境変数認証に関するアサーションを削除
- [ ] password_file認証のアサーションを追加（必要に応じて）

**完了条件:**
- 全統合テストがパスする

---

## フェーズ5: E2Eテスト

### Task 5.1: 環境変数認証テストファイルの削除
**ファイル**: `test/e2e/configs/`

**作業内容:**
- [ ] `02-password-env.yml`を削除
- [ ] `05-multi-hop-env.yml`を削除
- [ ] 対応するTTL出力ファイルを削除

**完了条件:**
- 環境変数認証のテストファイルが完全に削除されている

---

### Task 5.2: password_fileテストケースの追加
**ファイル**: `test/e2e/configs/06-password-file.yml`（新規作成）

**作業内容:**
- [ ] password_file認証のテストYAMLを作成
  ```yaml
  version: "1.0"

  profiles:
    bastion:
      host: localhost
      port: 22220
      user: testuser
      prompt_marker: "$ "
      auth:
        type: password
        password_file: "test_passwords.dat"

  routes:
    password-file-test:
      - profile: bastion
        commands:
          - whoami
          - hostname

  options:
    timeout: 30
  ```

**完了条件:**
- テストYAMLファイルが作成されている
- TTLが正しく生成される

---

### Task 5.3: E2E READMEの更新
**ファイル**: `test/e2e/README.md`

**作業内容:**
- [ ] 02-password-env.ymlのテストシナリオを削除
- [ ] 05-multi-hop-env.ymlのテストシナリオを削除
- [ ] 06-password-file.ymlのテストシナリオを追加
- [ ] パスワードファイルの作成方法を記載

**完了条件:**
- READMEが最新の状態に更新されている
- パスワードファイル作成方法が明記されている

---

## フェーズ6: ドキュメント更新

### Task 6.1: README.md（日本語）の更新
**ファイル**: `README.md`（デフォルト - 日本語）

**作業内容:**
- [ ] パスワード認証セクションを更新
- [ ] auth.envとauth.promptの記述を削除
- [ ] password_fileの説明を追加
- [ ] パスワードファイルの作成方法を追加
  ```markdown
  #### パスワード認証

  **1. パスワードファイル（推奨）**
  ```yaml
  profiles:
    server:
      auth:
        type: password
        password_file: "C:\\my_passwords.dat"  # optional (default: "passwords.dat")
  ```

  パスワードファイルの作成方法:
  - Tera Termメニュー: Setup → Additional settings → Copy
  - "Crypt" タブでパスワードを登録し、ファイルに保存

  **2. 直接指定（テスト用）**
  ```yaml
  profiles:
    server:
      auth:
        type: password
        value: "mypassword"  # セキュリティリスクあり
  ```
  ```

**完了条件:**
- auth.env/auth.promptの記述が削除されている
- password_fileの説明が追加されている
- パスワードファイル作成方法が明記されている

---

### Task 6.2: README.en.md（英語）の更新
**ファイル**: `README.en.md`

**作業内容:**
- [ ] README.mdと同じ内容を英語で更新
- [ ] パスワード認証セクションを更新
- [ ] auth.envとauth.promptの記述を削除
- [ ] password_fileの説明を追加
- [ ] パスワードファイルの作成方法を追加（英語）

**完了条件:**
- README.mdとREADME.en.mdの内容が同期している
- 両言語で同じ情報が提供されている

---

### Task 6.3: docs/の確認と更新
**ファイル**: `docs/development-guidelines.md`, `docs/architecture.md`等

**作業内容:**
- [ ] docs/配下のドキュメントでパスワード認証に言及している箇所を検索
- [ ] auth.env/auth.promptの記述があれば削除
- [ ] password_fileの記述に更新（必要に応じて）

**検索コマンド例:**
```bash
grep -r "auth.env\|auth.prompt\|passwordbox\|expandenv" docs/
```

**完了条件:**
- docs/配下のドキュメントが最新の仕様に更新されている
- 古い認証方式への言及が削除されている

---

### Task 6.4: CHANGELOG.mdの更新
**ファイル**: `CHANGELOG.md`

**作業内容:**
- [ ] `## [Unreleased]`セクションに破壊的変更を記載
  ```markdown
  ## [Unreleased]

  ### Changed
  - **BREAKING**: パスワード認証方式を変更
    - `auth.env`（環境変数）を削除 → `password_file`（パスワードファイル）を使用
    - `auth.prompt`（プロンプト）を削除 → `password_file`（パスワードファイル）を使用
    - `getpassword`コマンドを使用したTTL標準のパスワード管理に統一
    - デフォルトパスワードファイル: `passwords.dat`
  ```

**完了条件:**
- 破壊的変更が明記されている
- 変更内容が明確

---

## フェーズ7: 最終確認

### Task 7.1: 全テストの実行
**作業内容:**
- [ ] `go test ./...`を実行
- [ ] 全テストがパスすることを確認

**完了条件:**
- 全ユニットテストがパス
- 全統合テストがパス

---

### Task 7.2: Linterチェック
**作業内容:**
- [ ] `golangci-lint run`を実行
- [ ] Lintエラーがないことを確認

**完了条件:**
- Lintエラーが0件

---

### Task 7.3: TTL生成の動作確認
**作業内容:**
- [ ] サンプルYAMLでTTLを生成
- [ ] password_file認証のTTLが正しく生成されることを確認
  - `getpassword`コマンドが含まれている
  - `strconcat`コマンドが含まれている（第1ステップ）
  - password nameがプロファイル名と一致している

**完了条件:**
- 生成されるTTLが期待通り

---

### Task 7.4: E2Eテストの手動実行
**作業内容:**
- [ ] Dockerコンテナを起動
- [ ] パスワードファイルを作成（Tera Term）
- [ ] 06-password-file.ttlをTera Termで実行
- [ ] SSH接続が成功することを確認

**完了条件:**
- Tera Termでの実行が成功
- パスワードファイルからパスワードが正しく読み込まれる

---

## 完了基準

### 必須項目
- [ ] すべてのタスクが完了している
- [ ] 全テストがパスしている
- [ ] Linterチェックがパスしている
- [ ] ドキュメントが更新されている
- [ ] E2Eテストで動作確認できている

### 品質基準
- [ ] コードレビューを受けている
- [ ] 設計書との整合性が取れている
- [ ] 破壊的変更がCHANGELOG.mdに記載されている

---

## 備考

- 各タスク完了後、都度コミットすることを推奨
- フェーズごとにテストを実行して早期に問題を発見
- 不明点があれば設計書（design.md）を参照
