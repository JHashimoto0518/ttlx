# 自動切断オプション実装タスクリスト

## Phase 1: データ構造の変更

### Task 1.1: Options 構造体の変更
- [ ] `internal/config/model.go` の `Options` 構造体に `AutoDisconnect *bool` フィールドを追加
- [ ] YAML タグ `yaml:"auto_disconnect,omitempty"` を設定
- [ ] コメントを追加

**完了条件**:
- Options 構造体に AutoDisconnect フィールドが追加されている
- コンパイルエラーがない

---

### Task 1.2: SetDefaults メソッドの変更
- [ ] `internal/config/model.go` の `SetDefaults()` メソッドにデフォルト値設定を追加
- [ ] `AutoDisconnect` が `nil` の場合、`false` を設定

**完了条件**:
- SetDefaults() メソッドで AutoDisconnect のデフォルト値が設定される
- コンパイルエラーがない

---

### Task 1.3: ユニットテストの追加
- [ ] `internal/config/model_test.go` に `TestSetDefaults_AutoDisconnect` テストを追加
- [ ] 以下のケースをテスト:
  - 未指定の場合（デフォルト: false）
  - 明示的に true を指定した場合
  - 明示的に false を指定した場合

**完了条件**:
- テストが実装されている
- `go test ./internal/config` が成功する

---

## Phase 2: TTL生成ロジックの変更

### Task 2.1: テンプレートの追加
- [ ] `internal/generator/template.go` に `successKeepAliveTemplate` を追加

**完了条件**:
- テンプレート定数が追加されている
- コンパイルエラーがない

---

### Task 2.2: 自動切断処理の生成関数を追加
- [ ] `internal/generator/generator.go` に `generateAutoDisconnect()` 関数を追加
- [ ] 多段接続の場合、すべての接続を順次 `exit` で戻る処理を実装
- [ ] 1段接続の場合、直接 `closett` で終了する処理を実装

**完了条件**:
- generateAutoDisconnect() 関数が実装されている
- コンパイルエラーがない

---

### Task 2.3: Generate 関数の変更
- [ ] `internal/generator/generator.go` の `Generate()` 関数を変更
- [ ] `auto_disconnect` の値に基づいて処理を分岐
- [ ] `auto_disconnect: true` の場合: `generateAutoDisconnect()` を呼び出し
- [ ] `auto_disconnect: false` の場合: `successKeepAliveTemplate` を使用

**完了条件**:
- Generate() 関数で auto_disconnect に基づく処理分岐が実装されている
- コンパイルエラーがない

---

### Task 2.4: ユニットテストの追加
- [ ] `internal/generator/generator_test.go` に `TestGenerate_AutoDisconnect` テストを追加
- [ ] 以下のケースをテスト:
  - auto_disconnect: true (1段接続)
  - auto_disconnect: true (2段接続)
  - auto_disconnect: false
  - 未指定（デフォルト）

**完了条件**:
- テストが実装されている
- `go test ./internal/generator` が成功する

---

## Phase 3: フィクスチャとテスト

### Task 3.1: フィクスチャファイルの追加
- [ ] `test/fixtures/valid/auto_disconnect_true.yml` を作成
- [ ] `test/fixtures/valid/auto_disconnect_false.yml` を作成
- [ ] `test/fixtures/valid/auto_disconnect_default.yml` を作成

**完了条件**:
- 3つのフィクスチャファイルが作成されている
- YAML構文が正しい

---

### Task 3.2: 統合テストの実行
- [ ] `go test ./...` を実行してすべてのテストが成功することを確認
- [ ] テストカバレッジを確認（80%以上）

**完了条件**:
- すべてのユニットテストが成功する
- テストカバレッジが80%以上

---

### Task 3.3: 生成されたTTLの手動確認
- [ ] フィクスチャファイルから TTL を生成
- [ ] `auto_disconnect: true` の場合の TTL が正しいか確認
  - `closett` が含まれる
  - 多段接続の場合、`exit` コマンドが含まれる
- [ ] `auto_disconnect: false` の場合の TTL が正しいか確認
  - `closett` が含まれない
  - `end` のみで終了

**完了条件**:
- 生成されたTTLが期待通りの内容である

---

## Phase 4: ドキュメント更新

### Task 4.1: product-requirements.md の更新
- [ ] `docs/product-requirements.md` に `auto_disconnect` オプションの機能要件を追加
- [ ] ユーザーストーリーを追加（必要に応じて）

**完了条件**:
- product-requirements.md が更新されている

---

### Task 4.2: functional-design.md の更新
- [ ] `docs/functional-design.md` に YAML仕様として `auto_disconnect` を追加
- [ ] 使用例を追加

**完了条件**:
- functional-design.md が更新されている

---

### Task 4.3: README の更新
- [ ] `README.md` に `auto_disconnect` オプションの使用例を追加
- [ ] `README.en.md` に `auto_disconnect` オプションの使用例を追加（英語）

**完了条件**:
- README.md と README.en.md が更新されている

---

## 最終確認

### Task 5.1: 全体テスト
- [ ] `go test ./...` を実行してすべてのテストが成功することを確認
- [ ] `go build` を実行してビルドが成功することを確認
- [ ] `golangci-lint run` を実行してリントエラーがないことを確認

**完了条件**:
- すべてのテストが成功する
- ビルドが成功する
- リントエラーがない

---

### Task 5.2: ドキュメントの最終確認
- [ ] すべてのドキュメントが更新されているか確認
- [ ] コード例が正しいか確認
- [ ] 英語・日本語版の内容が一致しているか確認

**完了条件**:
- ドキュメントが整合性を持っている

---

## タスクサマリー

- **Phase 1**: データ構造の変更（3タスク）
- **Phase 2**: TTL生成ロジックの変更（4タスク）
- **Phase 3**: フィクスチャとテスト（3タスク）
- **Phase 4**: ドキュメント更新（3タスク）
- **最終確認**: 全体テスト（2タスク）

**合計**: 15タスク
