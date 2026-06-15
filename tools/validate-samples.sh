#!/usr/bin/env bash
#
# validate-samples.sh — 各サンプルが規約を満たしているか検証する
#
# チェック内容:
#   - samples/<category>/<name>/ に README.md と metadata.yaml が存在する
#   - metadata.yaml に必須キー(name/category/status/languages/summary)がある
#   - category が定義済みカテゴリのいずれか
#   - status が許可された値のいずれか
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SAMPLES_DIR="${REPO_ROOT}/samples"
VALID_CATEGORIES="monolithic microservices event-driven serverless frontend"
VALID_STATUS="draft active maintenance deprecated"
REQUIRED_KEYS="name category status languages summary"

errors=0
checked=0

fail() { echo "  ✗ $*"; errors=$((errors + 1)); }

for category in $VALID_CATEGORIES; do
  cat_dir="${SAMPLES_DIR}/${category}"
  [ -d "$cat_dir" ] || continue
  for sample in "$cat_dir"/*/; do
    [ -d "$sample" ] || continue
    name="$(basename "$sample")"
    [ "$name" = "_template" ] && continue
    checked=$((checked + 1))
    rel="samples/${category}/${name}"
    echo "• ${rel}"

    [ -f "${sample}README.md" ] || fail "README.md がありません"

    meta="${sample}metadata.yaml"
    if [ ! -f "$meta" ]; then
      fail "metadata.yaml がありません"
      continue
    fi

    for key in $REQUIRED_KEYS; do
      grep -qE "^${key}:" "$meta" || fail "metadata.yaml に '${key}:' がありません"
    done

    meta_category="$(grep -E '^category:' "$meta" | head -1 | sed 's/^category:[[:space:]]*//')"
    echo " $VALID_CATEGORIES " | grep -q " ${meta_category} " \
      || fail "category '${meta_category}' が未定義、または '${category}' と不一致"

    meta_status="$(grep -E '^status:' "$meta" | head -1 | sed 's/^status:[[:space:]]*//')"
    echo " $VALID_STATUS " | grep -q " ${meta_status} " \
      || fail "status '${meta_status}' が不正です（許可: ${VALID_STATUS}）"
  done
done

echo "---"
if [ "$checked" -eq 0 ]; then
  echo "ℹ️  検証対象のサンプルはまだありません。"
fi
if [ "$errors" -gt 0 ]; then
  echo "❌ ${errors} 件の問題が見つかりました。"
  exit 1
fi
echo "✅ 検証完了（${checked} サンプル、問題なし）。"
