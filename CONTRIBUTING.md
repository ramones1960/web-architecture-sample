# コントリビューションガイド

このリポジトリへの貢献を歓迎します。多様な Web アーキテクチャのサンプルを
育てるための運用ルールをまとめます。

## 貢献の種類

- 🆕 新しいアーキテクチャサンプルの追加
- 🛠 既存サンプルの改善・バグ修正
- 📖 ドキュメント・カタログ・運用ルールの改善

## はじめる前に

1. [ルート README](README.md) で目的と設計原則を確認する。
2. サンプルを追加する場合は [追加ガイド](docs/guides/adding-a-new-architecture.md) を読む。
3. 大きめの変更や設計議論が必要なものは、まず Issue で合意する。

## 開発フロー

```text
1. ブランチを切る        feat/<category>-<name>  など
2. 変更を実装する        サンプルは自己完結を維持する
3. 検証する             ./tools/validate-samples.sh
4. コミットする          Conventional Commits
5. PR を出す            PR テンプレートのチェックリストを満たす
```

### ブランチ名

`<type>/<scope>-<short-desc>`（例: `feat/microservices-saga`, `docs/update-catalog`）

### コミットメッセージ

[Conventional Commits](https://www.conventionalcommits.org/) に準拠。

```text
feat(monolithic): add layered architecture sample in TypeScript
docs(adr): record decision to keep samples self-contained
```

`type`: `feat` / `fix` / `docs` / `refactor` / `test` / `chore` / `ci`

## サンプルの品質基準

- **自己完結**: 依存・ビルド・テストをサンプル内に閉じる。他サンプルを import しない。
- **再現可能**: README の手順どおりにクリーン環境で動く。
- **説明責任**: README に「目的・構成・起動手順・学べること・トレードオフ」を書く。
- **メタデータ**: `metadata.yaml` の必須項目を埋める。

## レビュー

- CI（構造検証）がグリーンであること。
- レビュアーは README どおりに起動・理解できるかを確認する。
- 設計上の重要な判断は ADR への記録を推奨する。

## 行動規範

建設的・敬意あるコミュニケーションを心がけてください。
