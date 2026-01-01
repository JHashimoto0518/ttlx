#!/bin/bash
set -e

echo "🚀 ttlx テスト環境をセットアップ中..."
echo ""

# 1. コンテナ起動
echo "📦 Docker コンテナを起動中..."
docker compose -f docker-compose.test.yml up -d

# 2. サーバーが起動するまで待機
echo "⏳ SSHサーバーの起動を待機中（15秒）..."
sleep 15

# 3. 接続テスト
echo "🔍 接続テスト中..."
if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2222 testuser@localhost "echo 'bastion OK'" 2>/dev/null; then
    echo "  ✅ bastion: OK"
else
    echo "  ❌ bastion: 接続失敗"
fi

if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -p 2223 testuser@localhost "echo 'target OK'" 2>/dev/null; then
    echo "  ✅ target: OK"
else
    echo "  ❌ target: 接続失敗"
fi

# 4. TTL生成
echo ""
echo "📝 TTLスクリプトを生成中..."
./ttlx build test/test-config.yml -o test/test-connection.ttl

echo ""
echo "✅ セットアップ完了！"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 次のステップ:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Windowsでエクスプローラーを開く"
echo "2. アドレスバーに以下を入力:"
echo "   \\\\wsl\$\\Ubuntu\\workspaces\\ttlx\\test"
echo ""
echo "3. test-connection.ttl をダブルクリック"
echo "   （Tera Termがインストールされている場合）"
echo ""
echo "   または Tera Term を起動して:"
echo "   Control → Macro → test-connection.ttl を選択"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔧 テスト環境情報:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  bastion: localhost:2222 (testuser/testpass123)"
echo "  target:  localhost:2223 (testuser/testpass123)"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧹 クリーンアップ:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  docker-compose -f docker-compose.test.yml down"
echo ""
