# 設計書: 環境変数認証からgetpasswordへの移行

## 概要

環境変数認証（`auth.env`）とプロンプト認証（`auth.prompt`）を廃止し、Tera Term標準の`getpassword`コマンドを使用したパスワードファイル認証方式に統一する。

## アーキテクチャ

### 変更前のフロー

```
YAML設定
  ↓
Config解析
  ↓
┌─────────────────────────┐
│ パスワード認証方法      │
│ 1. auth.value (直接)    │
│ 2. auth.env (環境変数)  │
│ 3. auth.prompt (prompt) │
└─────────────────────────┘
  ↓
TTL生成
  ↓
┌─────────────────────────┐
│ 1. sendln 'password'    │
│ 2. expandenv + sendln   │
│ 3. passwordbox + sendln │
└─────────────────────────┘
```

### 変更後のフロー

```
YAML設定
  ↓
Config解析
  ↓
┌─────────────────────────────┐
│ パスワード認証方法          │
│ 1. auth.value (直接)        │
│ 2. password_file (ファイル) │
└─────────────────────────────┘
  ↓
デフォルト値適用
(password_file未指定 → "passwords.dat")
  ↓
TTL生成
  ↓
┌─────────────────────────────┐
│ 1. sendln 'password'        │
│ 2. getpassword + sendln     │
└─────────────────────────────┘
```

## データ構造の変更

### internal/config/config.go

#### 変更前
```go
type Auth struct {
    Type           string  `yaml:"type"`
    Value          string  `yaml:"value,omitempty"`
    Path           string  `yaml:"path,omitempty"`
    Env            string  `yaml:"env,omitempty"`           // 削除
    Prompt         bool    `yaml:"prompt,omitempty"`        // 削除
    PasswordPrompt string  `yaml:"password_prompt,omitempty"`
}
```

#### 変更後
```go
type Auth struct {
    Type           string  `yaml:"type"`
    Value          string  `yaml:"value,omitempty"`
    Path           string  `yaml:"path,omitempty"`
    PasswordFile   string  `yaml:"password_file,omitempty"` // 追加
    PasswordPrompt string  `yaml:"password_prompt,omitempty"`
}
```

### バリデーションロジック

#### internal/config/validator.go

**追加するバリデーション:**

1. **password_fileのデフォルト値設定**
   ```go
   // パスワード認証でvalueもpassword_fileも指定されていない場合、デフォルト値を設定
   if auth.Type == "password" && auth.Value == "" && auth.PasswordFile == "" {
       auth.PasswordFile = "passwords.dat"
   }
   ```

2. **相互排他性の検証**
   ```go
   // password認証の場合、valueとpassword_fileは排他的
   if auth.Type == "password" {
       if auth.Value != "" && auth.PasswordFile != "" {
           return fmt.Errorf("auth.value and auth.password_file are mutually exclusive")
       }
   }
   ```

## TTL生成ロジック

### internal/generator/template.go

#### 削除するテンプレート

```go
// 削除
connectWithEnvPasswordTemplate = `...`
passwordEnvTemplate = `...`
passwordPromptTemplate = `...`
```

#### 追加するテンプレート

**第1ステップ用（connect前）:**
```go
// パスワード認証テンプレート（getpassword - connect用）
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

**パラメータ:**
1. `%s` - password_file
2. `%s` - password name (プロファイル名)
3. `%s` - host
4. `%d` - port
5. `%s` - auth type
6. `%s` - user
7. `%s` - keyfile option (空文字列)
8. `%s` - upper profile name (ERROR label)
9. `%s` - prompt marker
10. `%s` - upper profile name (TIMEOUT label)

**第2ステップ以降用（ssh後）:**
```go
// パスワード認証テンプレート（getpassword）
passwordFileTemplate = `; Password authentication (from password file)
getpassword '%s' '%s' password
sendln password

`
```

**パラメータ:**
1. `%s` - password_file
2. `%s` - password name (プロファイル名)

#### 既存のテンプレートは維持

```go
// パスワード認証テンプレート（直接指定）
passwordValueTemplate = `; Password authentication
sendln '%s'

`
```

### internal/generator/generator.go

#### generateConnect() の変更

**変更前:**
```go
func generateConnect(stepNum int, profileName, upperProfileName string, profile *config.Profile) string {
    authType := profile.Auth.Type
    keyfileOption := ""
    passwordOption := ""

    // 環境変数からパスワードを取得する場合は専用テンプレートを使用
    if authType == "password" && profile.Auth.Env != "" {
        return fmt.Sprintf(connectWithEnvPasswordTemplate, ...)
    }

    if authType == "keyfile" {
        keyfileOption = fmt.Sprintf(" /keyfile=%s", profile.Auth.Path)
    } else if authType == "password" && profile.Auth.Value != "" {
        passwordOption = fmt.Sprintf(" /passwd=%s", profile.Auth.Value)
    }

    return fmt.Sprintf(connectTemplate, ...)
}
```

**変更後:**
```go
func generateConnect(stepNum int, profileName, upperProfileName string, profile *config.Profile) string {
    authType := profile.Auth.Type
    keyfileOption := ""
    passwordOption := ""

    // パスワードファイル認証の場合は専用テンプレートを使用
    // getpassword + strconcat + connectを含む完全なシーケンスを生成
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

    // keyfile認証
    if authType == "keyfile" {
        keyfileOption = fmt.Sprintf(" /keyfile=%s", profile.Auth.Path)
    }

    // パスワード直接指定
    if authType == "password" && profile.Auth.Value != "" {
        passwordOption = fmt.Sprintf(" /passwd=%s", profile.Auth.Value)
    }

    return fmt.Sprintf(
        connectTemplate,
        stepNum,
        profileName,
        upperProfileName,
        profile.Host,
        profile.Port,
        authType,
        profile.User,
        keyfileOption,
        passwordOption,
        upperProfileName,
        profile.PromptMarker,
        upperProfileName,
    )
}
```

**重要な変更点:**
- `profile.Auth.Env`のチェックを削除
- `profile.Auth.PasswordFile`のチェックを追加
- パスワードファイル認証の場合、`passwordFileConnectTemplate`を使用
- このテンプレートは`getpassword` + `strconcat` + `connect`の完全なシーケンスを含む
- そのため、第1ステップの生成ロジックで追加のパスワード認証コードを生成する必要がない

#### generatePasswordAuth() の変更

**変更前:**
```go
func generatePasswordAuth(profileName string, auth *config.Auth) string {
    if auth.Env != "" {
        return fmt.Sprintf(passwordEnvTemplate, auth.Env)
    }
    if auth.Prompt {
        return fmt.Sprintf(passwordPromptTemplate, profileName)
    }
    if auth.Value != "" {
        return fmt.Sprintf(passwordValueTemplate, auth.Value)
    }
    return ""
}
```

**変更後:**
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

#### 第1ステップの生成ロジック変更

**変更前:**
```go
if i == 0 {
    // 最初のステップ: connect コマンド
    sb.WriteString(generateConnect(i+1, step.Profile, upperProfileName, profile))

    // パスワード認証処理（connect コマンドに含まれていない場合のみ）
    // 環境変数認証の場合はconnectコマンド内でexpandenvを使用するため、ここでは処理しない
    if profile.Auth.Type == "password" && profile.Auth.Value == "" && profile.Auth.Env == "" {
        sb.WriteString(generatePasswordAuth(step.Profile, profile.Auth))
    }
}
```

**変更後:**
```go
if i == 0 {
    // 最初のステップ: connect コマンド
    // password_fileが設定されている場合、generateConnect()が完全なシーケンスを生成するため
    // 追加のパスワード認証コードは不要
    sb.WriteString(generateConnect(i+1, step.Profile, upperProfileName, profile))

    // パスワード認証処理は不要
    // - password_fileの場合: generateConnect()内で完全なシーケンスを生成済み
    // - auth.valueの場合: connectコマンドに/passwd=が含まれている
}
```

**変更点:**
- パスワード認証処理のコードを完全に削除
- `password_file`の場合、`generateConnect()`が`passwordFileConnectTemplate`を使用し、完全なシーケンス（getpassword + strconcat + connect）を生成
- `auth.value`の場合、connectコマンドに`/passwd=`が既に含まれている
- どちらの場合も追加のコードは不要

## テスト戦略

### 削除するテスト

1. **環境変数認証のテスト**
   - `test/e2e/configs/02-password-env.yml`
   - `test/e2e/configs/05-multi-hop-env.yml`
   - 環境変数認証に関するユニットテスト

2. **プロンプト認証のテスト**
   - `passwordbox`に関するテスト
   - `auth.prompt`に関するテスト

### 追加するテスト

#### 1. getpassword生成のユニットテスト

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
    assert.Contains(t, ttl, "sendln password")
}
```

#### 2. デフォルト値のテスト

```go
func TestGenerate_PasswordFileDefault(t *testing.T) {
    cfg := &config.Config{
        Profiles: map[string]*config.Profile{
            "server": {
                Host: "example.com",
                Port: 22,
                User: "user",
                Auth: &config.Auth{
                    Type: "password",
                    // PasswordFile未指定 → デフォルト値が適用されるはず
                },
                PromptMarker: "$ ",
            },
        },
    }

    // バリデーション後にデフォルト値が設定される
    // テストではバリデーション後のconfigを使用

    route := []*config.RouteStep{
        {Profile: "server", Commands: []string{"ls"}},
    }

    ttl, err := generateRoute(cfg, "test", route, "test.yml")
    require.NoError(t, err)

    // デフォルトのパスワードファイルが使用される
    assert.Contains(t, ttl, "getpassword 'passwords.dat' 'server' password")
}
```

#### 3. バリデーションエラーのテスト

```go
func TestValidate_MutuallyExclusivePasswordFields(t *testing.T) {
    // valueとpassword_fileの両方を指定した場合はエラー
    cfg := &config.Config{
        Profiles: map[string]*config.Profile{
            "server": {
                Auth: &config.Auth{
                    Type:         "password",
                    Value:        "mypassword",
                    PasswordFile: "passwords.dat",
                },
            },
        },
    }

    err := Validate(cfg)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "mutually exclusive")
}
```

#### 4. E2Eテストケース

**test/e2e/configs/06-password-file.yml** (新規作成)
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

**期待されるTTL出力:**
```ttl
; === Step 1: bastion ===
:CONNECT_BASTION
; Password authentication (from password file)
getpassword 'test_passwords.dat' 'bastion' password

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

; Command: whoami
sendln 'whoami'
...
```

### テスト実行順序

1. **ユニットテスト**: config, validator, generator
2. **統合テスト**: build_test.go
3. **E2Eテスト**: 実際のTera Termでの動作確認（手動）

## マイグレーション戦略

マイグレーションは考慮不要（ベータ版のため）。

既存ユーザーへの影響:
- `auth.env`または`auth.prompt`を使用している設定はバリデーションエラーになる
- エラーメッセージで`password_file`の使用を案内

## リスク分析

### 高リスク
なし

### 中リスク

1. **パスワードファイルの作成方法が不明瞭**
   - **対策**: READMEにパスワードファイルの作成方法を記載
   - **対策**: サンプルファイルを提供

2. **デフォルト値が適切でない可能性**
   - **対策**: デフォルト値`"passwords.dat"`を相対パスで設定
   - **対策**: ユーザーが絶対パスで上書き可能

### 低リスク

1. **既存テストの更新漏れ**
   - **対策**: 全テストを実行してエラーがないことを確認

## パフォーマンス影響

- TTL生成ロジックの変更による影響: なし（軽微な変更のため）
- ファイルI/O: Tera Term実行時のみ（ttlx自体には影響なし）

## セキュリティ考慮事項

### 改善点
1. **環境変数への露出を回避**
   - パスワードファイル方式により、環境変数にパスワードを設定不要

2. **TTL標準準拠**
   - Tera Term公式の方法でパスワードを管理

### 注意点
1. **パスワードファイルの保護**
   - ユーザーはパスワードファイルのアクセス権限を適切に設定する必要がある
   - READMEに注意事項を記載

2. **直接パスワード指定（auth.value）**
   - 引き続きサポートするが、セキュリティリスクがあることを明示
   - テスト/自動化用途に限定することを推奨

## ドキュメント更新

### README.md

**追加する内容:**

#### パスワード認証方式

**1. パスワードファイル（推奨）**
```yaml
profiles:
  server:
    auth:
      type: password
      password_file: "C:\\my_passwords.dat"  # optional
```

パスワードファイルの作成方法:
```
Tera Termメニュー: Setup → Additional settings → Copy
"Crypt" タブでパスワードを登録し、ファイルに保存
```

**2. 直接指定（テスト用）**
```yaml
profiles:
  server:
    auth:
      type: password
      value: "mypassword"  # セキュリティリスクあり
```

### CHANGELOG.md

```markdown
## [Unreleased]

### Changed
- **BREAKING**: パスワード認証方式を変更
  - `auth.env`（環境変数）を削除 → `password_file`（パスワードファイル）を使用
  - `auth.prompt`（プロンプト）を削除 → `password_file`（パスワードファイル）を使用
  - `getpassword`コマンドを使用したTTL標準のパスワード管理に統一

### Migration
マイグレーション考慮不要（ベータ版のため）
```

## 実装順序

1. **Config変更** (`internal/config/`)
   - Auth構造体の更新
   - バリデーションロジックの追加

2. **Template変更** (`internal/generator/template.go`)
   - 不要なテンプレートの削除
   - `passwordFileTemplate`の追加

3. **Generator変更** (`internal/generator/generator.go`)
   - `generatePasswordAuth()`の更新
   - 第1ステップ生成ロジックの修正

4. **テスト更新**
   - 既存テストの削除
   - 新しいテストの追加

5. **E2Eテスト**
   - テストファイルの削除/追加
   - READMEの更新

6. **ドキュメント更新**
   - README.md
   - CHANGELOG.md

## 品質保証

### レビューポイント
- [ ] Auth構造体の変更が正しいか
- [ ] バリデーションロジックが適切か
- [ ] TTL生成ロジックが正しいか
- [ ] テストカバレッジが十分か
- [ ] ドキュメントが明確か

### テスト項目
- [ ] ユニットテストがすべてパス
- [ ] 統合テストがすべてパス
- [ ] Linterチェックがパス
- [ ] E2Eテストで実際の動作確認

## 完了条件

- [ ] すべてのコード変更が完了
- [ ] すべてのテストがパス
- [ ] ドキュメントが更新される
- [ ] PRがマージされる
