# TTLテスト環境

**Read this in other languages: 日本語 | [English](README.en.md)**

Docker Composeを使用したTera Term マクロ（TTL）のテスト環境です。

## 前提条件

- Docker Desktop for Windows（WSL2バックエンド）
- または Docker Engine（Linux/macOS）

## テスト手順

### 1. テスト用SSHサーバーの起動

```bash
# WSL2またはLinux環境で実行
docker compose -f docker-compose.test.yml up -d

# サーバーが起動するまで少し待つ（初回は時間がかかります）
sleep 10

# 起動確認
docker compose -f docker-compose.test.yml ps
```

### 2. 接続テスト（Linux側）

```bash
# bastionへの接続テスト
ssh -p 2222 testuser@localhost
# パスワード: testpass123

# targetへの接続テスト
ssh -p 2223 testuser@localhost
# パスワード: testpass123
```

### 3. TTLスクリプトの生成

```bash
# テスト用TTLを生成（-o でディレクトリパスを指定）
./ttlx build test/test-config.yml -o test/

# 生成内容を確認
cat test/test-connection.ttl
```

### 4. Tera Termでの実行（Windows側）

1. **Tera Termを起動**
2. **マクロを実行**: メニュー → Control → Macro
3. **ファイルを選択**: `\\wsl$\Ubuntu\workspaces\ttlx\test\test-connection.ttl`
   - パスは環境に応じて調整してください
4. **動作確認**:
   - bastionに接続できるか
   - targetに多段接続できるか
   - コマンドが実行されるか

### 5. テスト環境のクリーンアップ

```bash
# コンテナを停止・削除
docker compose -f docker-compose.test.yml down

# イメージも削除する場合
docker compose -f docker-compose.test.yml down --rmi all
```

## テスト設定

### 接続情報

| サーバー | ホスト | ポート | ユーザー | パスワード |
|---------|--------|--------|---------|-----------|
| bastion | localhost | 2222 | testuser | testpass123 |
| target  | localhost | 2223 | testuser | testpass123 |

**注意**: これはテスト環境専用の認証情報です。本番環境では絶対に使用しないでください。

### テスト項目

- [ ] 単段SSH接続（bastion）
- [ ] 多段SSH接続（bastion → target）
- [ ] パスワード認証（直接指定）
- [ ] コマンド実行
- [ ] 接続エラーハンドリング
- [ ] タイムアウト処理

## トラブルシューティング

### ポートが既に使用されている

```bash
# ポート使用状況を確認
netstat -an | grep -E "2222|2223"

# 別のポートを使用する場合は docker-compose.test.yml を編集
```

### コンテナが起動しない

```bash
# ログを確認
docker compose -f docker-compose.test.yml logs

# コンテナを再起動
docker compose -f docker-compose.test.yml restart
```

### SSH接続が失敗する

```bash
# コンテナ内でSSHサーバーの状態を確認
docker exec ttlx-test-bastion ps aux | grep ssh

# SSHの詳細ログを確認
ssh -vvv -p 2222 testuser@localhost
```

## 自動化スクリプト（オプション）

テスト環境の起動からTTL生成までを自動化：

```bash
#!/bin/bash
# test-setup.sh

echo "🚀 テスト環境をセットアップ中..."

# 1. コンテナ起動
docker compose -f docker-compose.test.yml up -d

# 2. サーバーが起動するまで待機
echo "⏳ SSHサーバーの起動を待機中..."
sleep 15

# 3. 接続テスト
echo "🔍 接続テスト中..."
ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2222 testuser@localhost "echo 'bastion OK'" 2>/dev/null
ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2223 testuser@localhost "echo 'target OK'" 2>/dev/null

# 4. TTL生成
echo "📝 TTLスクリプトを生成中..."
./ttlx build test/test-config.yml -o test/

echo "✅ セットアップ完了！"
echo ""
echo "次のステップ:"
echo "1. Windowsでエクスプローラーを開く"
echo "2. パス: \\\\wsl\$\\Ubuntu\\workspaces\\ttlx\\test\\test-connection.ttl"
echo "3. Tera Termでマクロを実行"
```

使用方法：

```bash
chmod +x test-setup.sh
./test-setup.sh
```
