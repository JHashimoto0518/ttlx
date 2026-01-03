# ttlx E2E テスト

## 概要

このディレクトリには、ttlxが生成したTTLマクロをTera Termで実際に実行し、SSH接続が正しく動作することを確認するE2Eテスト環境が含まれています。

## テスト環境の構成

### Dockerコンテナ

- **bastion** (localhost:2222)
  - 踏み台サーバー
  - ユーザー: testuser
  - パスワード: testpass123
  - 公開鍵認証: 対応

- **target** (内部ネットワーク)
  - 内部サーバー
  - ユーザー: testuser
  - パスワード: testpass123
  - bastionからのみ接続可能

### テストシナリオ

| ファイル | テスト内容 | 検証項目 |
|---------|----------|---------|
| `01-keyfile-auth.yml` | 公開鍵認証 | keyfile認証でSSH接続できるか |
| `02-password-env.yml` | 環境変数パスワード認証 | 環境変数からパスワードを取得して接続できるか |
| `03-auto-disconnect.yml` | 自動切断 | auto_disconnect=true で接続が自動切断されるか |
| `04-multiple-routes.yml` | 複数ルート | 1つのYAMLから複数のTTLが生成され、それぞれが動作するか |

## セットアップと実行

### 1. テスト環境のセットアップ

```bash
cd test/e2e
./setup.sh
```

**実行内容:**
- SSH鍵ペアの生成（存在しない場合）
- Dockerコンテナの起動（bastion, target）
- TTLファイルの生成（configs/*.yml → output/*.ttl）

### 2. Tera Termでの手動テスト

#### Windows環境での実行手順

1. **Tera Termを起動**

2. **環境変数の設定**（02-password-env.ttl をテストする場合）
   ```cmd
   set TEST_SSH_PASSWORD=testpass123
   ```

3. **TTLマクロの実行**
   ```
   Tera Term メニュー:
   Control → Macro
   → test/e2e/output/<テストしたいTTL>.ttl を選択
   ```

4. **動作確認**
   - SSH接続が成功するか
   - コマンドが実行されるか
   - 期待通りの動作をするか
   - エラーが発生しないか

#### テストシナリオ別の確認ポイント

##### 01-keyfile-auth.yml
- [ ] 公開鍵認証でbastionに接続できる
- [ ] パスワード入力なしで接続完了
- [ ] `whoami`, `hostname` コマンドが実行される

##### 02-password-env.yml
- [ ] 環境変数 `TEST_SSH_PASSWORD` が読み込まれる
- [ ] パスワード認証でbastionに接続できる
- [ ] コマンドが正常に実行される

##### 03-auto-disconnect.yml
- [ ] bastionに接続できる
- [ ] targetに多段SSH接続できる
- [ ] スクリプト終了時に自動的に切断される
- [ ] Tera Termが自動的に閉じる

##### 04-multiple-routes.yml
- **route-single-hop.ttl**:
  - [ ] bastionに接続できる
  - [ ] コマンドが実行される
  - [ ] 接続が維持される（auto_disconnect=false）

- **route-multi-hop.ttl**:
  - [ ] bastionに接続できる
  - [ ] targetに多段接続できる
  - [ ] 各ステップのコマンドが実行される
  - [ ] 接続が維持される

### 3. テスト環境のクリーンアップ

```bash
cd test/e2e
./teardown.sh
```

**実行内容:**
- Dockerコンテナの停止と削除
- 生成されたTTLファイルの削除

## ディレクトリ構造

```
test/e2e/
├── README.md                  # このファイル
├── docker-compose.yml         # Docker環境定義
├── setup.sh                   # セットアップスクリプト
├── teardown.sh                # クリーンアップスクリプト
├── configs/                   # テスト用YAMLファイル
│   ├── 01-keyfile-auth.yml
│   ├── 02-password-env.yml
│   ├── 03-auto-disconnect.yml
│   └── 04-multiple-routes.yml
├── output/                    # 生成されたTTLファイル（自動生成）
│   ├── keyfile-test.ttl
│   ├── env-password-test.ttl
│   ├── auto-disconnect-test.ttl
│   ├── route-single-hop.ttl
│   └── route-multi-hop.ttl
└── ssh-keys/                  # テスト用SSH鍵（自動生成）
    ├── id_rsa
    ├── id_rsa.pub
    └── authorized_keys
```

## トラブルシューティング

### Dockerコンテナが起動しない

```bash
# コンテナのログを確認
docker-compose logs bastion
docker-compose logs target

# コンテナを再起動
docker-compose restart
```

### SSH接続が失敗する

```bash
# bastionコンテナに手動接続して確認
docker exec -it ttlx-e2e-bastion sh

# SSHサービスの状態確認
docker exec ttlx-e2e-bastion ps aux | grep ssh
```

### TTLファイルが生成されない

```bash
# ttlxを再ビルド
cd ../..  # プロジェクトルート
go build -o ttlx ./cmd/ttlx

# 手動でTTL生成
./ttlx build test/e2e/configs/01-keyfile-auth.yml -o test/e2e/output/
```

## 注意事項

- **セキュリティ**: このテスト環境は開発・テスト専用です。本番環境では使用しないでください
- **パスワード**: テスト用のパスワード（testpass123）はハードコードされています
- **SSH鍵**: test/e2e/ssh-keys/ 内の鍵はテスト専用です。本番環境では使用しないでください
- **ポート競合**: ホストの2222ポートが使用されている場合、docker-compose.ymlでポート番号を変更してください

## 将来の拡張

現在は手動テストですが、将来的には以下の自動化を検討できます：

- TTL実行の自動化（Tera Term CLI モードまたは互換ツール）
- 接続成功/失敗の自動判定
- CI/CDパイプラインへの統合
