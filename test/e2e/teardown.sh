#!/bin/bash
set -e

echo "========================================="
echo "ttlx E2E Test Environment Teardown"
echo "========================================="
echo ""

# 現在のディレクトリを保存
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 1. Dockerコンテナの停止と削除
echo "[1/2] Stopping and removing Docker containers..."
docker-compose down -v
echo "  ✓ Containers removed"

# 2. 生成されたTTLファイルのクリーンアップ
echo ""
echo "[2/2] Cleaning up generated TTL files..."
if [ -d "output" ]; then
    rm -rf output/*
    echo "  ✓ Output files cleaned"
else
    echo "  ✓ No output files to clean"
fi

echo ""
echo "========================================="
echo "✓ E2E Test Environment Cleaned Up!"
echo "========================================="
echo ""
