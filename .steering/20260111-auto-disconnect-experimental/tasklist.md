# タスクリスト: auto_disconnect 実験的機能への変更

## 実装タスク

- [x] `generateAutoDisconnect`関数をループベースに変更
- [x] `collectPromptMarkers`関数を追加
- [x] 呼び出し元を更新（prompt_markersを渡す）
- [x] テストを更新
- [x] 単体テスト実行・確認
- [x] Lint実行・確認

## ドキュメントタスク

- [x] CHANGELOG.md更新
- [x] README.md更新
  - [x] 実験的機能であることを明記
  - [x] 制限事項を追加
- [x] README.en.md更新
  - [x] 実験的機能であることを明記
  - [x] 制限事項を追加
- [x] docs/functional-design.md更新
  - [x] 実験的機能であることを追記

## 完了条件

- [x] 全テストがパス
- [x] Lintエラーなし
- [x] ドキュメント更新完了
- [ ] PRマージ

## 備考

- 現在のブランチ: `fix/auto-disconnect-loop`
- PR: #33（ドキュメント更新後に再プッシュが必要）
