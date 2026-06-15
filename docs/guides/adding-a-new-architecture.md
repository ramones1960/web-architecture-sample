# 新しいアーキテクチャサンプルを追加する

このガイドは、`samples/` 配下に新しいアーキテクチャサンプルを追加する
詳細手順を説明します。概要は [ルート README](../../README.md#-新しいアーキテクチャサンプルの追加手順) を参照してください。

---

## 0. 事前確認

- 追加したいアーキテクチャが [カタログ](../architecture-catalog.md) のどの **カテゴリ** に属するか決める。
  - 既存カテゴリに収まらない場合のみ、新カテゴリ追加を検討（カタログ末尾の手順参照）。
- 既に同等のサンプルが無いか `samples/` を確認する。
  - 同じアーキテクチャでも **言語が異なる** 比較サンプルは歓迎。
    その場合は名前に言語を含める（例: `clean-architecture-go`, `clean-architecture-ts`）。

## 1. 提案（Issue）

Issue テンプレート **「New Architecture Sample」** を使い、以下を明記します。

- 対象アーキテクチャ名 / カテゴリ
- 使用言語・主要フレームワーク
- 題材（何を作るか。例: ToDo API、ミニ EC）
- このサンプルで示したい設計上のポイント

> 小さな追加であれば Issue を省略しても構いませんが、
> 設計議論が必要なものは Issue で合意してから着手すると手戻りが減ります。

## 2. 雛形を生成する

```bash
./tools/new-sample.sh <category> <architecture-name>

# 例
./tools/new-sample.sh monolithic clean-architecture-ts
```

これにより `samples/<category>/<architecture-name>/` が
`samples/_template/` の内容をもとに作成されます。

## 3. 命名規約

| 対象 | 規約 | 例 |
| --- | --- | --- |
| ディレクトリ名 | 小文字 + ハイフン（kebab-case） | `event-sourcing-ts` |
| 言語比較する場合 | 末尾に言語サフィックス | `-go` / `-ts` / `-py` / `-java` |
| カテゴリ | カタログで定義済みのものを使用 | `microservices` |

## 4. 実装

各サンプルは **自己完結** が原則です。次を満たしてください。

- 依存・ロックファイル・ビルド設定をサンプル内に閉じる。
- 他サンプルのコードを import しない。
- `README.md` の「セットアップ / 起動 / テスト」手順どおりに、
  クリーン環境で動作すること。
- 可能なら `docker compose` などで **ワンコマンド起動** を提供する。

### サンプル内の推奨レイアウト（言語により調整可）

```text
<architecture-name>/
├── README.md            # 必須: 目的・構成図・起動手順・学べること
├── metadata.yaml        # 必須: カテゴリ/言語/ステータス/概要
├── docs/
│   └── architecture.md  # 構成図や設計判断の詳細
├── src/ (or app/, cmd/) # 実装
└── (Makefile / docker-compose.yml / package.json 等)
```

## 5. メタデータを記入する

`metadata.yaml` を埋めます（スキーマは [_template/metadata.yaml](../../samples/_template/metadata.yaml) 参照）。

```yaml
name: clean-architecture-ts
category: monolithic
status: draft        # draft -> active へ更新
languages: [typescript]
summary: TypeScript によるクリーンアーキテクチャの最小実装
tags: [clean-architecture, dependency-inversion]
```

## 6. ドキュメントを更新する

- `metadata.yaml` の `status` を実態に合わせる。
- [カタログ](../architecture-catalog.md) の該当行の状態（📝/🚧/✅）を更新する。
- 重要な設計判断があれば [ADR](../adr/README.md) を追加する。

## 7. 検証して PR を出す

```bash
./tools/validate-samples.sh   # 構造・必須項目をチェック
```

- ブランチ名: `feat/<category>-<architecture-name>`
- コミット: Conventional Commits（例: `feat(monolithic): add clean-architecture-ts sample`）
- PR テンプレートのチェックリストを満たす。

---

## チェックリスト（PR 前の最終確認）

- [ ] `samples/<category>/<name>/README.md` がある
- [ ] `metadata.yaml` の必須項目（name/category/status/languages/summary）が埋まっている
- [ ] README の手順でクリーン環境から起動・テストできる
- [ ] 他サンプルへの依存がない（自己完結している）
- [ ] カタログの状態表を更新した
- [ ] `./tools/validate-samples.sh` が通る
