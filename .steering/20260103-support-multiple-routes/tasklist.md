# タスクリスト: 複数ルート対応

## フェーズ1: データ構造の変更

### 1.1 Config構造体の変更
- [ ] `internal/config/model.go` の `Config` 構造体を変更
  - `Route []RouteStep` → `Routes map[string][]RouteStep`
- [ ] `SetDefaults()` メソッドを修正
  - Routes の各ルートに対してデフォルト設定を適用

### 1.2 Config構造体のテスト
- [ ] `internal/config/model_test.go` にテストケースを追加
  - `TestConfig_Parse_Routes`: 複数ルートのパーステスト
  - `TestConfig_SetDefaults_Routes`: Routes のデフォルト設定テスト
- [ ] 既存のテストを修正（route → routes）

## フェーズ2: バリデーションの変更

### 2.1 Validator の変更
- [ ] `internal/config/validator.go` の `Validate()` 関数を変更
  - routesが空の場合のチェック
  - ルート名のバリデーション（空文字、無効文字）
  - 各ルートが空でないかチェック
  - 各ルート内のプロファイル参照チェック
  - 各ルート内の password_prompt チェック
- [ ] `isValidFileName()` ヘルパー関数を追加
  - 英数字、ハイフン、アンダースコアのみ許可

### 2.2 Validator のテスト
- [ ] `internal/config/validator_test.go` にテストケースを追加
  - `TestValidate_Routes_Empty`: routesが空
  - `TestValidate_Routes_InvalidRouteName`: 無効なルート名
  - `TestValidate_Routes_EmptyRoute`: ルートが空
  - `TestValidate_Routes_ProfileNotFound`: プロファイル参照エラー
  - `TestValidate_Routes_PasswordPrompt`: password_prompt 必須チェック

## フェーズ3: Generator の変更

### 3.1 Generator の変更
- [ ] `internal/generator/generator.go` をリファクタリング
  - `Generate()` → `GenerateAll()` に変更
  - 返り値を `string` → `map[string]string` に変更
  - `generateRoute()` 内部関数を追加（単一ルート生成）
  - `generateHeader()` にルート名パラメータを追加
- [ ] `generateRoute()` の実装
  - 既存の `Generate()` ロジックを流用
  - `cfg.Route` → `route` 引数に変更
- [ ] `GenerateAll()` の実装
  - 各ルートをループして `generateRoute()` を呼び出し
  - 結果を map に格納

### 3.2 Generator のテスト
- [ ] `internal/generator/generator_test.go` を修正
  - `TestGenerate` → `TestGenerateAll` に変更
  - 複数ルート生成のテストケースを追加
  - 各ルートのTTL内容を検証
  - エラーケースのテスト

## フェーズ4: CLI の変更

### 4.1 build コマンドの変更
- [ ] `internal/cli/build.go` を変更
  - `generator.Generate()` → `generator.GenerateAll()` に変更
  - 返り値を `string` → `map[string]string` に変更
  - 複数ファイル書き込みロジックを追加
  - `-o` オプションをディレクトリパスとして扱う
  - 出力ディレクトリが存在しない場合は作成
  - 生成されたファイル一覧を表示（ソート済み）
- [ ] dry-run モードの対応
  - 各ルートのTTLを区切って表示

### 4.2 validate コマンドの確認
- [ ] `internal/cli/validate.go` の動作確認
  - エラーメッセージが正しく表示されるか

## フェーズ5: テストフィクスチャの移行

### 5.1 valid フィクスチャの移行
- [ ] `test/fixtures/valid/simple.yml` を変更
  - `route:` → `routes:` + ルート名追加
- [ ] `test/fixtures/valid/full.yml` を変更
  - `route:` → `routes:` + ルート名追加
- [ ] `test/fixtures/valid/auto_disconnect_true.yml` を変更
  - `route:` → `routes:` + ルート名追加
- [ ] `test/fixtures/valid/auto_disconnect_false.yml` を変更
  - `route:` → `routes:` + ルート名追加
- [ ] `test/fixtures/valid/auto_disconnect_default.yml` を変更
  - `route:` → `routes:` + ルート名追加
- [ ] 新規フィクスチャを追加
  - `test/fixtures/valid/multiple-routes.yml`: 複数ルート定義の例

### 5.2 invalid フィクスチャの追加
- [ ] `test/fixtures/invalid/empty-routes.yml` を作成
  - routesが空の場合
- [ ] `test/fixtures/invalid/invalid-route-name.yml` を作成
  - 無効なルート名（スペース、スラッシュなど）
- [ ] `test/fixtures/invalid/empty-route.yml` を作成
  - ルートのステップが空

### 5.3 test-config.yml の移行
- [ ] `test/test-config.yml` を変更
  - `route:` → `routes:` + ルート名追加

## フェーズ6: ドキュメント更新

### 6.1 functional-design.md の更新
- [ ] `docs/functional-design.md` を更新
  - YAML仕様を `routes` に変更
  - 複数ルートの例を追加
  - スキーマ定義を更新

### 6.2 README の更新
- [ ] `README.md` を更新
  - クイックスタートの例を `routes` に変更
  - 複数ルートの使用例を追加
  - CLI出力例を更新
- [ ] `README.en.md` を更新
  - 日本語版と同じ内容を英語で

## フェーズ7: 統合テスト

### 7.1 全テストの実行
- [ ] すべてのユニットテストが成功することを確認
  ```bash
  go test ./...
  ```
- [ ] テストカバレッジを確認（95%以上）
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

### 7.2 Lint の実行
- [ ] golangci-lint でエラーがないことを確認
  ```bash
  golangci-lint run
  ```

### 7.3 実機テスト
- [ ] `simple.yml` をビルドして動作確認
  ```bash
  ./ttlx build test/fixtures/valid/simple.yml
  ```
- [ ] `multiple-routes.yml` をビルドして複数ファイル生成を確認
  ```bash
  ./ttlx build test/fixtures/valid/multiple-routes.yml
  ```
- [ ] `-o` オプションでディレクトリ指定を確認
  ```bash
  ./ttlx build test/fixtures/valid/multiple-routes.yml -o output/
  ```
- [ ] dry-run モードの確認
  ```bash
  ./ttlx build test/fixtures/valid/multiple-routes.yml --dry-run
  ```

### 7.4 validate コマンドのテスト
- [ ] valid フィクスチャのバリデーション成功を確認
  ```bash
  ./ttlx validate test/fixtures/valid/simple.yml
  ```
- [ ] invalid フィクスチャのバリデーション失敗を確認
  ```bash
  ./ttlx validate test/fixtures/invalid/empty-routes.yml
  ```

## 完了条件

- [ ] すべてのユニットテストが成功
- [ ] テストカバレッジが95%以上
- [ ] golangci-lint でエラーなし
- [ ] 実機テストですべての機能が動作確認済み
- [ ] ドキュメントが更新済み
- [ ] 既存フィクスチャが新仕様に移行済み

## 備考

### 実装時の注意点
- 既存のテストケースは route → routes に修正が必要
- generator の既存関数は引数調整が必要（cfg.Route への直接参照を避ける）
- ルート名のバリデーションは厳密に（ファイル名として使用するため）
- エラーメッセージにはルート名を含める（デバッグしやすくする）

### 想定される課題
- テストフィクスチャの移行漏れ
- generator の既存ロジックの依存関係
- CLI の出力フォーマット（見やすさ）
