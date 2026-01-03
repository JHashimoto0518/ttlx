#!/bin/bash
set -e

echo "========================================="
echo "ttlx E2E Test Environment Setup"
echo "========================================="
echo ""

# 現在のディレクトリを保存
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 1. SSH鍵の確認
echo "[1/5] Checking SSH keys..."
if [ ! -f "ssh-keys/id_rsa" ]; then
    echo "  Generating SSH key pair..."
    ssh-keygen -t rsa -b 2048 -f ssh-keys/id_rsa -N "" -C "ttlx-e2e-test"
    cp ssh-keys/id_rsa.pub ssh-keys/authorized_keys
    chmod 600 ssh-keys/id_rsa
    chmod 644 ssh-keys/authorized_keys
    echo "  ✓ SSH keys generated"
else
    echo "  ✓ SSH keys already exist"
fi

# 2. Dockerコンテナの起動
echo ""
echo "[2/5] Starting Docker containers..."
docker-compose down 2>/dev/null || true
docker-compose up -d

# 3. コンテナの起動待機
echo ""
echo "[3/5] Waiting for containers to be ready..."
echo "  Waiting for bastion..."
timeout 30 bash -c 'until docker exec ttlx-e2e-bastion nc -z localhost 2222 2>/dev/null; do sleep 1; done'
echo "  ✓ Bastion is ready"

echo "  Waiting for target..."
timeout 30 bash -c 'until docker exec ttlx-e2e-target nc -z localhost 2222 2>/dev/null; do sleep 1; done'
echo "  ✓ Target is ready"

# 4. TTLファイルの生成
echo ""
echo "[4/5] Generating TTL files from test configs..."
cd ../..  # プロジェクトルートへ移動

# ttlxをビルド
if [ ! -f "ttlx" ]; then
    echo "  Building ttlx..."
    go build -o ttlx ./cmd/ttlx
fi

# 各YAMLからTTLを生成
for config in test/e2e/configs/*.yml; do
    echo "  Generating TTL from $(basename "$config")..."
    ./ttlx build "$config" -o test/e2e/output/
done

cd test/e2e

# 5. 環境変数の確認
echo ""
echo "[5/5] Environment variable check..."
if [ -z "$TEST_SSH_PASSWORD" ]; then
    echo "  ⚠️  Warning: TEST_SSH_PASSWORD is not set"
    echo "  To test 02-password-env.yml, run:"
    echo "    export TEST_SSH_PASSWORD=testpass123"
else
    echo "  ✓ TEST_SSH_PASSWORD is set"
fi

echo ""
echo "========================================="
echo "✓ E2E Test Environment Ready!"
echo "========================================="
echo ""
echo "Next steps:"
echo "1. Open Tera Term on Windows"
echo "2. Load and execute TTL files from: test/e2e/output/"
echo "3. Verify SSH connections work correctly"
echo ""
echo "Available test scenarios:"
echo "  - output/keyfile-test.ttl       : Keyfile authentication test"
echo "  - output/env-password-test.ttl  : Environment variable password test"
echo "  - output/auto-disconnect-test.ttl : Auto-disconnect test"
echo "  - output/route-single-hop.ttl   : Multiple routes test (single hop)"
echo "  - output/route-multi-hop.ttl    : Multiple routes test (multi hop)"
echo ""
echo "To teardown: ./teardown.sh"
echo ""
