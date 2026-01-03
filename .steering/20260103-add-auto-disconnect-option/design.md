# 自動切断オプションの実装設計

## 概要

本ドキュメントでは、`options.auto_disconnect` オプションの実装設計を定義します。最終ステップ完了後の動作（自動切断 vs 接続保持）をグローバルオプションで制御できるようにします。

---

## 実装アプローチ

### 基本方針

1. **グローバルオプション**: `Options` 構造体に `AutoDisconnect *bool` フィールドを追加
2. **デフォルト値**: 未指定時は `false`（接続保持）
3. **TTL生成ロジックの変更**: 成功終了処理を `auto_disconnect` の値に基づいて切り替え
4. **バリデーション**: `auto_disconnect` の型チェック（boolean以外はエラー）

---

## データ構造の変更

### 1. Options 構造体の変更

**ファイル**: `internal/config/model.go`

**変更前**:
```go
// Options represents global options.
type Options struct {
	Timeout int    `yaml:"timeout,omitempty"`
	Retry   int    `yaml:"retry,omitempty"`
	Log     bool   `yaml:"log,omitempty"`
	LogFile string `yaml:"log_file,omitempty"`
}
```

**変更後**:
```go
// Options represents global options.
type Options struct {
	Timeout        int    `yaml:"timeout,omitempty"`
	Retry          int    `yaml:"retry,omitempty"`
	Log            bool   `yaml:"log,omitempty"`
	LogFile        string `yaml:"log_file,omitempty"`
	AutoDisconnect *bool  `yaml:"auto_disconnect,omitempty"` // 最終ステップ完了後に自動切断するか
}
```

**設計上の注意点**:
- `*bool` (ポインタ型) を使用する理由: 未指定（nil）とfalseを区別するため
- YAML タグに `omitempty` を使用: 未指定時にYAMLに出力されないようにする

---

### 2. SetDefaults メソッドの変更

**ファイル**: `internal/config/model.go`

**変更前**:
```go
// SetDefaults sets default values for the config.
func (c *Config) SetDefaults() {
	for _, profile := range c.Profiles {
		if profile.Port == 0 {
			profile.Port = 22
		}
	}

	if c.Options == nil {
		c.Options = &Options{}
	}
	if c.Options.Timeout == 0 {
		c.Options.Timeout = 30
	}
}
```

**変更後**:
```go
// SetDefaults sets default values for the config.
func (c *Config) SetDefaults() {
	for _, profile := range c.Profiles {
		if profile.Port == 0 {
			profile.Port = 22
		}
	}

	if c.Options == nil {
		c.Options = &Options{}
	}
	if c.Options.Timeout == 0 {
		c.Options.Timeout = 30
	}
	if c.Options.AutoDisconnect == nil {
		defaultAutoDisconnect := false
		c.Options.AutoDisconnect = &defaultAutoDisconnect
	}
}
```

**設計上の注意点**:
- デフォルト値は `false`（接続保持）
- `nil` の場合のみデフォルト値を設定（明示的に指定された値は上書きしない）

---

## TTL生成ロジックの変更

### 3. Generate 関数の変更

**ファイル**: `internal/generator/generator.go`

**変更箇所**: 成功終了処理の生成部分（現在の57-58行目）

**変更前**:
```go
// 成功終了
sb.WriteString(successTemplate)
```

**変更後**:
```go
// 成功終了（auto_disconnect に基づいて処理を切り替え）
autoDisconnect := false
if cfg.Options != nil && cfg.Options.AutoDisconnect != nil {
	autoDisconnect = *cfg.Options.AutoDisconnect
}

if autoDisconnect {
	// 自動切断: 多段接続を順次exit、最後にclosett
	sb.WriteString(generateAutoDisconnect(len(cfg.Route)))
} else {
	// 接続保持: セッションを維持したまま終了
	sb.WriteString(successKeepAliveTemplate)
}
```

---

### 4. 自動切断処理の生成関数

**ファイル**: `internal/generator/generator.go`

**新規追加**:
```go
// generateAutoDisconnect generates disconnect sequence for all route steps.
func generateAutoDisconnect(routeSteps int) string {
	var sb strings.Builder

	sb.WriteString("; === Auto Disconnect ===\n")

	// 多段接続の場合、すべての接続を順次exit
	if routeSteps > 1 {
		for i := routeSteps - 1; i > 0; i-- {
			sb.WriteString(fmt.Sprintf("; Disconnect from step %d\n", i+1))
			sb.WriteString("sendln 'exit'\n")
			sb.WriteString("pause 1\n") // 切断処理の完了を待つ
		}
	}

	// 成功終了（Tera Term終了）
	sb.WriteString("\n:SUCCESS\n")
	sb.WriteString("closett\n")
	sb.WriteString("end\n\n")

	return sb.String()
}
```

**設計上の注意点**:
- 多段接続の場合、最後のステップから順に `exit` で戻る
- 1段接続の場合、直接 `closett` で終了
- `pause 1` で切断処理の完了を待つ（Tera Termが安定して切断処理を行うため）

---

### 5. 接続保持テンプレートの追加

**ファイル**: `internal/generator/template.go`

**新規追加**:
```go
// 成功終了テンプレート（接続保持）
successKeepAliveTemplate = `; === Success (Keep connection alive) ===
:SUCCESS
end

`
```

**既存の successTemplate は削除せず、auto_disconnect: true の場合に使用**:
```go
// 成功終了テンプレート（自動切断）
successTemplate = `:SUCCESS
closett
end

`
```

**設計上の注意点**:
- `successKeepAliveTemplate` は `end` のみで終了（`closett` しない）
- TTLスクリプトは終了するが、Tera Termのセッションは維持される
- ユーザーが手動で `exit` を入力するまで接続が保持される

---

## バリデーションの変更

### 6. Validate 関数の拡張

**ファイル**: `internal/config/validator.go`

**追加箇所**: `Validate()` 関数の末尾（return の前）

**追加コード**:
```go
// auto_disconnect の型チェック（YAMLパース時に自動的にチェックされるが、念のため）
if config.Options != nil && config.Options.AutoDisconnect != nil {
	// ポインタ型なので、nilでなければboolean型として有効
	// YAMLパース時に型が合わない場合はエラーになるため、ここでは追加のチェック不要
}
```

**設計上の注意点**:
- YAMLパース時に `gopkg.in/yaml.v3` が型チェックを行うため、明示的なバリデーションは不要
- `auto_disconnect: "true"` のような文字列が指定された場合、YAMLパース時にエラーになる

---

## テストケースの追加

### 7. model_test.go の拡張

**ファイル**: `internal/config/model_test.go`

**追加テストケース**:
```go
func TestSetDefaults_AutoDisconnect(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name: "auto_disconnect not specified (default to false)",
			config: &Config{
				Options: &Options{},
			},
			expected: false,
		},
		{
			name: "auto_disconnect explicitly set to true",
			config: &Config{
				Options: &Options{
					AutoDisconnect: boolPtr(true),
				},
			},
			expected: true,
		},
		{
			name: "auto_disconnect explicitly set to false",
			config: &Config{
				Options: &Options{
					AutoDisconnect: boolPtr(false),
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.SetDefaults()
			if *tt.config.Options.AutoDisconnect != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, *tt.config.Options.AutoDisconnect)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
```

---

### 8. generator_test.go の拡張

**ファイル**: `internal/generator/generator_test.go`

**追加テストケース**:
```go
func TestGenerate_AutoDisconnect(t *testing.T) {
	tests := []struct {
		name           string
		autoDisconnect *bool
		routeSteps     int
		expectClosett  bool
		expectExit     int // 期待される exit コマンドの数
	}{
		{
			name:           "auto_disconnect: true (1 step)",
			autoDisconnect: boolPtr(true),
			routeSteps:     1,
			expectClosett:  true,
			expectExit:     0,
		},
		{
			name:           "auto_disconnect: true (2 steps)",
			autoDisconnect: boolPtr(true),
			routeSteps:     2,
			expectClosett:  true,
			expectExit:     1,
		},
		{
			name:           "auto_disconnect: false",
			autoDisconnect: boolPtr(false),
			routeSteps:     2,
			expectClosett:  false,
			expectExit:     0,
		},
		{
			name:           "auto_disconnect: not specified (default)",
			autoDisconnect: nil,
			routeSteps:     2,
			expectClosett:  false,
			expectExit:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テストケースに応じた Config を構築
			cfg := buildTestConfig(tt.autoDisconnect, tt.routeSteps)

			// TTL生成
			result, err := Generate(cfg, "test.yml")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// closett の存在チェック
			if tt.expectClosett {
				if !strings.Contains(result, "closett") {
					t.Error("expected 'closett' in generated TTL")
				}
			} else {
				if strings.Contains(result, "closett") {
					t.Error("unexpected 'closett' in generated TTL")
				}
			}

			// exit コマンドの数チェック
			exitCount := strings.Count(result, "sendln 'exit'")
			if exitCount != tt.expectExit {
				t.Errorf("expected %d 'exit' commands, got %d", tt.expectExit, exitCount)
			}
		})
	}
}
```

---

### 9. フィクスチャファイルの追加

**ディレクトリ**: `test/fixtures/`

**新規作成**:
1. `test/fixtures/valid/auto_disconnect_true.yml` - `auto_disconnect: true` のケース
2. `test/fixtures/valid/auto_disconnect_false.yml` - `auto_disconnect: false` のケース
3. `test/fixtures/valid/auto_disconnect_default.yml` - 未指定のケース

**auto_disconnect_true.yml の例**:
```yaml
version: "1.0"

profiles:
  bastion:
    host: bastion.example.com
    port: 22
    user: user1
    prompt_marker: "$ "
    auth:
      type: password
      env: BASTION_PASS

  target:
    host: 10.0.0.50
    port: 22
    user: user2
    prompt_marker: "$ "
    auth:
      type: password
      password_prompt: "password:"
      env: TARGET_PASS

route:
  - profile: bastion
  - profile: target
    commands:
      - ls -la
      - pwd

options:
  auto_disconnect: true
```

---

## 実装順序

1. **Phase 1: データ構造の変更**
   - [ ] `Options` 構造体に `AutoDisconnect *bool` フィールドを追加
   - [ ] `SetDefaults()` メソッドにデフォルト値設定を追加
   - [ ] ユニットテスト追加（`model_test.go`）

2. **Phase 2: TTL生成ロジックの変更**
   - [ ] `successKeepAliveTemplate` を `template.go` に追加
   - [ ] `generateAutoDisconnect()` 関数を `generator.go` に追加
   - [ ] `Generate()` 関数で `auto_disconnect` に基づく処理分岐を実装
   - [ ] ユニットテスト追加（`generator_test.go`）

3. **Phase 3: フィクスチャとテスト**
   - [ ] フィクスチャファイルを追加
   - [ ] 統合テストの実行
   - [ ] 生成されたTTLの手動確認

4. **Phase 4: ドキュメント更新**
   - [ ] `docs/product-requirements.md` 更新
   - [ ] `docs/functional-design.md` 更新
   - [ ] `README.md` / `README.en.md` 更新

---

## リスクと対策

### リスク1: 多段接続の切断処理の不安定性

**リスク**: `exit` コマンドの送信タイミングによって、切断が失敗する可能性がある

**対策**:
- `pause 1` で切断処理の完了を待つ
- 実際に Tera Term で動作確認を行う
- 必要に応じて `pause` の時間を調整

---

### リスク2: Tera Term のバージョンによる挙動の違い

**リスク**: Tera Term のバージョンによって `closett` や `end` の挙動が異なる可能性がある

**対策**:
- Tera Term 4.106 および 5.x で動作確認
- ドキュメントに動作確認済みバージョンを記載

---

### リスク3: デフォルト値の選択による混乱

**リスク**: デフォルトが `false`（接続保持）のため、自動化を期待するユーザーが混乱する可能性がある

**対策**:
- README に明確に記載
- バリデーション時に、デフォルト値を使用していることを通知（オプション、実装するか検討）

---

## まとめ

本設計では、`options.auto_disconnect` オプションをグローバルオプションとして実装します。デフォルト値を `false`（接続保持）とすることで、既存ユーザーの期待動作を維持しつつ、完全自動化したいユーザーには明示的な指定を求めます。

実装は4つのフェーズに分けて段階的に進め、各フェーズでテストを実施することで品質を確保します。
