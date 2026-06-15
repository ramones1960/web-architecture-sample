#!/usr/bin/env bash
#
# new-sample.sh — テンプレートから新しいアーキテクチャサンプルを生成する
#
# 使い方:
#   ./tools/new-sample.sh <category> <architecture-name>
#
# 例:
#   ./tools/new-sample.sh monolithic clean-architecture-ts
#
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE_DIR="${REPO_ROOT}/samples/_template"
VALID_CATEGORIES=(monolithic microservices event-driven serverless frontend)

err() { echo "Error: $*" >&2; exit 1; }

[ $# -eq 2 ] || err "使い方: $0 <category> <architecture-name>"

category="$1"
name="$2"

# カテゴリの妥当性チェック
found=false
for c in "${VALID_CATEGORIES[@]}"; do
  [ "$c" = "$category" ] && found=true && break
done
$found || err "カテゴリ '$category' は未定義です。有効: ${VALID_CATEGORIES[*]}
新カテゴリは docs/architecture-catalog.md に追記してから使用してください。"

# 名前は kebab-case のみ許可
[[ "$name" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]] || err "サンプル名は kebab-case で指定してください（例: clean-architecture-ts）"

dest="${REPO_ROOT}/samples/${category}/${name}"
[ -e "$dest" ] && err "既に存在します: ${dest#$REPO_ROOT/}"

cp -R "$TEMPLATE_DIR" "$dest"

# プレースホルダを置換
if command -v sed >/dev/null 2>&1; then
  grep -rl --exclude-dir=.git "<architecture-name>\|<category>" "$dest" 2>/dev/null | while read -r f; do
    sed -i.bak "s/<architecture-name>/${name}/g; s/<category>/${category}/g" "$f" && rm -f "${f}.bak"
  done
fi

echo "✅ 生成しました: ${dest#$REPO_ROOT/}"
echo "次のステップ:"
echo "  1. ${dest#$REPO_ROOT/}/metadata.yaml を記入"
echo "  2. README.md / docs/architecture.md を実装に合わせて更新"
echo "  3. samples/README.md と docs/architecture-catalog.md の索引を更新"
echo "  4. ./tools/validate-samples.sh で検証"
