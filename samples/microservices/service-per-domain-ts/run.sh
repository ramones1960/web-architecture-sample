#!/usr/bin/env bash
# ローカル開発用: catalog / order / gateway の 3 サービスをまとめて起動する。
# Ctrl+C（SIGINT）で全プロセスを停止する。
#
# 使い方: bash run.sh

set -euo pipefail
cd "$(dirname "$0")"

pids=()

cleanup() {
  echo
  echo "[run] stopping services..."
  for pid in "${pids[@]}"; do
    kill "$pid" 2>/dev/null || true
  done
  wait 2>/dev/null || true
  echo "[run] stopped."
}
trap cleanup INT TERM

echo "[run] starting catalog-service (4001)..."
node services/catalog/server.ts &
pids+=("$!")

echo "[run] starting order-service (4002)..."
node services/order/server.ts &
pids+=("$!")

echo "[run] starting api-gateway (4000)..."
node gateway/server.ts &
pids+=("$!")

echo
echo "[run] all services started. 公開窓口は API Gateway のみ:"
echo "[run]   Gateway : http://localhost:4000"
echo "[run]     例) curl localhost:4000/catalog/products"
echo "[run]     例) curl -X POST localhost:4000/orders -d '{\"productId\":\"p1\",\"qty\":2}'"
echo "[run]   (内部) catalog: http://localhost:4001 / order: http://localhost:4002"
echo "[run] Ctrl+C で停止します。"
echo

# いずれかのプロセスが終了するまで待つ。
wait
