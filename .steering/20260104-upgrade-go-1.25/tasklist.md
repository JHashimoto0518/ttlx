# タスクリスト: Go バージョンアップグレード (1.21.13 → 1.25.5)

## 概要

Go のバージョンを 1.21.13 から 1.25.5 にアップグレードするタスクリスト。
各タスクは順番に実行し、問題があれば即座に対応する。

---

## Phase 1: Go バージョン更新

### Task 1.1: go.mod の更新
- **説明**: `go.mod` の Go バージョンを 1.25.5 に更新
- **ファイル**: `go.mod`
- **変更内容**: `go 1.21.13` → `go 1.25.5`
- **検証**: `go version` で Go 1.25.5 であることを確認

### Task 1.2: 依存関係の整理
- **説明**: `go mod tidy` で依存関係を整理
- **コマンド**: `go mod tidy`
- **検証**: `go mod verify` でモジュールの整合性を確認

---

## Phase 2: テスト実行

### Task 2.1: ユニットテストの実行
- **説明**: すべてのユニットテストを実行
- **コマンド**: `go test ./...`
- **検証**: すべてのテストが PASS すること

### Task 2.2: 統合テストの実行
- **説明**: 統合テストを実行
- **コマンド**: `go test ./test/integration/...`
- **検証**: すべてのテストが PASS すること

### Task 2.3: テストカバレッジの確認
- **説明**: テストカバレッジを計測
- **コマンド**:
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -func=coverage.out
  ```
- **検証**: カバレッジが 97.8% 以上であること

### Task 2.4: Linter の実行
- **説明**: golangci-lint を実行
- **コマンド**: `golangci-lint run`
- **検証**: エラーがないこと、警告が増加していないこと

---

## Phase 3: ビルドテスト

### Task 3.1: ビルドの実行
- **説明**: ttlx バイナリをビルド
- **コマンド**: `go build -o ttlx ./cmd/ttlx`
- **検証**: ビルドが成功すること

### Task 3.2: バージョン確認
- **説明**: バイナリのバージョンを確認
- **コマンド**: `./ttlx version`
- **検証**: バージョン情報が正しく表示されること

### Task 3.3: TTL 生成テスト
- **説明**: テスト用 YAML から TTL を生成
- **コマンド**: `./ttlx build test/fixtures/valid/simple.yml --dry-run`
- **検証**: TTL が正常に生成され、内容が正しいこと

---

## Phase 4: 開発環境更新

### Task 4.1: devcontainer の確認
- **説明**: `.devcontainer/devcontainer.json` が存在するか確認
- **確認事項**:
  - ファイルの存在確認
  - Go バージョンの指定方法を確認

### Task 4.2: devcontainer の更新（存在する場合）
- **説明**: devcontainer の Go バージョンを 1.25.5 に更新
- **ファイル**: `.devcontainer/devcontainer.json`
- **検証**: devcontainer のリビルド後、Go 1.25.5 が使用されること

### Task 4.3: CI/CD 設定の確認
- **説明**: CI/CD で使用している Go バージョンを確認
- **確認対象**:
  - GitHub Actions ワークフロー
  - その他の CI/CD 設定

### Task 4.4: CI/CD 設定の更新（存在する場合）
- **説明**: CI/CD の Go バージョンを 1.25.5 に更新
- **検証**: CI/CD パイプラインが成功すること

---

## Phase 5: ドキュメント更新

### Task 5.1: README.md の更新
- **説明**: 日本語版 README の必要要件を更新
- **ファイル**: `README.md`
- **変更内容**:
  - 「Go 1.21以降」 → 「Go 1.25以降」
- **セクション**: 開発 > 必要要件

### Task 5.2: README.en.md の更新
- **説明**: 英語版 README の必要要件を更新
- **ファイル**: `README.en.md`
- **変更内容**:
  - 「Go 1.21 or later」 → 「Go 1.25 or later」
- **セクション**: Development > Prerequisites

### Task 5.3: docs/architecture.md の更新
- **説明**: 技術仕様書の Go バージョンを更新
- **ファイル**: `docs/architecture.md`
- **変更内容**:
  - 「Go 1.21以降」 → 「Go 1.25以降」
- **セクション**: テクノロジースタック > 実装言語 > バージョン

---

## Phase 6: 最終検証

### Task 6.1: 全テストの再実行
- **説明**: すべてのテストを再度実行して最終確認
- **コマンド**:
  ```bash
  go test ./...
  golangci-lint run
  ```
- **検証**: すべてのテストと Linter が成功すること

### Task 6.2: ビルドの再実行
- **説明**: 最終的なビルドを実行
- **コマンド**: `go build -o ttlx ./cmd/ttlx`
- **検証**: ビルドが成功すること

### Task 6.3: E2E TTL 生成の確認
- **説明**: E2E テスト用の TTL ファイルを生成
- **コマンド**:
  ```bash
  ./ttlx build test/e2e/configs/02-password-file.yml -o test/e2e/output/
  ./ttlx build test/e2e/configs/05-multi-hop-password-file.yml -o test/e2e/output/
  ```
- **検証**: TTL ファイルが正常に生成され、内容に変更がないこと

### Task 6.4: Git 変更内容の確認
- **説明**: 変更されたファイルを確認
- **コマンド**: `git status`
- **検証**: 想定通りのファイルのみが変更されていること

---

## 完了条件チェックリスト

### 必須項目

- [ ] `go.mod` の Go バージョンが 1.25.5 に更新されている
- [ ] `go mod tidy` が実行され、依存関係が整理されている
- [ ] すべてのユニットテストが成功している
- [ ] すべての統合テストが成功している
- [ ] テストカバレッジが 97.8% 以上維持されている
- [ ] golangci-lint がエラーなしで成功している
- [ ] `go build` が成功している
- [ ] 生成される TTL に変更がないことを確認している
- [ ] README.md (日本語) が更新されている
- [ ] README.en.md (英語) が更新されている
- [ ] docs/architecture.md が更新されている

### オプション項目

- [ ] devcontainer が更新されている（該当する場合）
- [ ] CI/CD 設定が更新されている（該当する場合）
- [ ] CI/CD パイプラインが成功している（該当する場合）

---

## トラブルシューティング

### テストが失敗した場合

1. エラーログを確認
2. Go 1.25 で非推奨になった API がないか確認
3. 依存ライブラリのバージョンを確認
4. 必要に応じて中間バージョン（1.23, 1.24）で検証

### ビルドが失敗した場合

1. `go mod verify` でモジュールの整合性を確認
2. `go clean -cache` でビルドキャッシュをクリア
3. `go mod tidy` を再実行

### Linter が失敗した場合

1. 新しい Lint ルールが追加されていないか確認
2. Go 1.25 で推奨される書き方に変更
3. 必要に応じて `.golangci.yml` を更新

---

## 参考コマンド

### 開発環境の確認

```bash
# Go バージョンの確認
go version

# 依存関係の確認
go list -m all

# モジュールの検証
go mod verify
```

### テストとビルド

```bash
# すべてのテスト実行
go test ./...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Linter 実行
golangci-lint run

# ビルド
go build -o ttlx ./cmd/ttlx

# クリーンビルド
go clean -cache
go build -o ttlx ./cmd/ttlx
```

### デバッグ

```bash
# ビルド時の詳細ログ
go build -v -x -o ttlx ./cmd/ttlx

# テスト時の詳細ログ
go test -v ./...

# 依存関係のグラフ表示
go mod graph
```

---

## タスク実行ログ

各タスクの実行結果を記録：

| Phase | Task | 状態 | 実行日時 | 備考 |
|-------|------|------|----------|------|
| 1 | 1.1 | ⏳ 未実施 | - | - |
| 1 | 1.2 | ⏳ 未実施 | - | - |
| 2 | 2.1 | ⏳ 未実施 | - | - |
| 2 | 2.2 | ⏳ 未実施 | - | - |
| 2 | 2.3 | ⏳ 未実施 | - | - |
| 2 | 2.4 | ⏳ 未実施 | - | - |
| 3 | 3.1 | ⏳ 未実施 | - | - |
| 3 | 3.2 | ⏳ 未実施 | - | - |
| 3 | 3.3 | ⏳ 未実施 | - | - |
| 4 | 4.1 | ⏳ 未実施 | - | - |
| 4 | 4.2 | ⏳ 未実施 | - | - |
| 4 | 4.3 | ⏳ 未実施 | - | - |
| 4 | 4.4 | ⏳ 未実施 | - | - |
| 5 | 5.1 | ⏳ 未実施 | - | - |
| 5 | 5.2 | ⏳ 未実施 | - | - |
| 5 | 5.3 | ⏳ 未実施 | - | - |
| 6 | 6.1 | ⏳ 未実施 | - | - |
| 6 | 6.2 | ⏳ 未実施 | - | - |
| 6 | 6.3 | ⏳ 未実施 | - | - |
| 6 | 6.4 | ⏳ 未実施 | - | - |
