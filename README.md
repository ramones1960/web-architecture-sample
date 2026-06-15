# web-architecture-sample

さまざまな **Web アーキテクチャ** の実装サンプルを 1 つのモノレポに集約し、
比較・学習・参照ができるようにするためのリポジトリです。

レイヤードアーキテクチャ、クリーンアーキテクチャ、マイクロサービス、
イベント駆動、CQRS / Event Sourcing、サーバーレス、マイクロフロントエンドなど、
代表的なアーキテクチャを **それぞれ独立したサンプル** として育てていきます。

---

## 🎯 目的（このリポジトリで実現したいこと）

1. **比較学習** — 同じ題材（例: ToDo / EC / ブログ）を異なるアーキテクチャで実装し、
   設計判断のトレードオフを横並びで見られるようにする。
2. **参照実装** — 実プロジェクトを始めるときの「型」として、すぐにコピーして使える
   ミニマルかつ実用的なサンプルを提供する。
3. **意思決定の記録** — なぜそのアーキテクチャ・技術を選んだのかを
   ADR（Architecture Decision Record）として残し、判断の背景を追えるようにする。
4. **将来の拡張** — 言語・フレームワークに依存しない運用ルールを定め、
   サンプルを継続的に追加できる土台を保つ。

### 設計原則

| 原則 | 内容 |
| --- | --- |
| **サンプルは自己完結** | 各サンプルは自前の依存・ビルド・テスト・実行手順を持ち、単体で起動できる |
| **疎結合** | サンプル間で直接コードを共有しない（共有したい知見はドキュメント化する） |
| **ポリグロット前提** | TypeScript / Go / Python / Java など、アーキテクチャに適した言語を自由に選べる |
| **ドキュメント駆動** | コードと同じ場所に「何を・なぜ」を記した README を必ず置く |
| **再現可能** | `README` の手順どおりに誰でも同じ結果を再現できる |

---

## 📁 ディレクトリ構造

```text
web-architecture-sample/
├── README.md                  # このファイル（全体の目的・運用ルール）
├── CONTRIBUTING.md            # コントリビューション / サンプル追加の手引き
├── .editorconfig             # エディタ共通設定
├── .gitignore
├── .github/                   # CI・Issue/PR テンプレート
│   ├── workflows/
│   │   └── ci.yml             # 構造検証 + 変更サンプルごとのビルド
│   ├── ISSUE_TEMPLATE/
│   │   └── new-architecture-sample.md
│   └── pull_request_template.md
├── docs/                      # リポジトリ横断ドキュメント
│   ├── architecture-catalog.md   # 収録アーキテクチャの一覧・分類
│   ├── guides/
│   │   └── adding-a-new-architecture.md  # サンプル追加の詳細手順
│   └── adr/                   # Architecture Decision Records
│       ├── README.md
│       ├── template.md
│       ├── 0001-record-architecture-decisions.md
│       └── 0002-monorepo-with-independent-samples.md
├── tools/                     # リポジトリ運用スクリプト
│   ├── new-sample.sh         # テンプレートから新サンプルを生成
│   └── validate-samples.sh   # 各サンプルの規約準拠を検証
└── samples/                   # ★ アーキテクチャサンプル本体
    ├── README.md             # カテゴリ説明 + サンプル索引
    ├── _template/            # 新サンプルのひな型（コピー元）
    │   ├── README.md
    │   ├── metadata.yaml
    │   ├── .gitignore
    │   └── docs/architecture.md
    ├── monolithic/           # モノリス系（レイヤード / モジュラーモノリス / クリーン）
    ├── microservices/        # マイクロサービス系
    ├── event-driven/         # イベント駆動 / Pub-Sub / CQRS / Event Sourcing
    ├── serverless/           # FaaS / BaaS / エッジ
    └── frontend/             # SPA / SSR・SSG / マイクロフロントエンド / BFF
```

### なぜこの構造か

- **`samples/<category>/<architecture-name>/`** という 2 階層で分類することで、
  アーキテクチャが増えても一覧性を保てる（カテゴリは [カタログ](docs/architecture-catalog.md) で定義）。
- **各サンプルディレクトリが管理単位**。依存ファイル・ロックファイル・CI 設定・
  README をサンプル内に閉じ込め、他サンプルへ影響させない。
- ルート直下は **「規約・ドキュメント・運用ツール」だけ** を置き、
  特定言語のツールチェーン（package.json 等）でルートを汚さない。

---

## 🗂 アーキテクチャカタログ

収録済み・収録予定のアーキテクチャ一覧は **[docs/architecture-catalog.md](docs/architecture-catalog.md)** を参照してください。
各サンプルの索引は **[samples/README.md](samples/README.md)** にもまとまっています。

---

## ➕ 新しいアーキテクチャサンプルの追加手順

詳細は **[docs/guides/adding-a-new-architecture.md](docs/guides/adding-a-new-architecture.md)** にあります。要約すると次の 6 ステップです。

1. **提案** — Issue テンプレート「New Architecture Sample」で目的・対象アーキテクチャ・言語を提案する。
2. **生成** — ひな型から雛形を作る。

   ```bash
   ./tools/new-sample.sh <category> <architecture-name>
   # 例: ./tools/new-sample.sh monolithic clean-architecture-go
   ```

3. **実装** — サンプルを実装する。`README.md` の起動手順を必ず動く状態に保つ。
4. **メタデータ記入** — `metadata.yaml`（カテゴリ・言語・ステータス・概要）を埋める。
5. **意思決定の記録** — 設計上の重要な判断があれば ADR を追加する（任意）。
6. **検証 & PR** — `./tools/validate-samples.sh` が通ることを確認し、PR を出す。

---

## 🛠 リポジトリ運用方法

### ブランチ戦略

- `main` … 常にグリーン（CI 通過）を保つ保護ブランチ。直接 push しない。
- 作業ブランチ … `<type>/<scope>-<short-desc>` 形式。
  - 例: `feat/microservices-order-saga`、`docs/update-catalog`、`fix/serverless-cold-start`
- 1 PR = 1 つの目的（1 サンプル追加 / 1 改善）に絞る。

### コミットメッセージ

[Conventional Commits](https://www.conventionalcommits.org/) に準拠します。

```text
<type>(<scope>): <subject>

例:
feat(monolithic): add layered architecture sample in TypeScript
docs(adr): record decision to keep samples self-contained
chore(ci): validate metadata.yaml schema
```

`type`: `feat` / `fix` / `docs` / `refactor` / `test` / `chore` / `ci`
`scope`: カテゴリ名やサンプル名（例: `microservices`, `frontend/bff`）

### Pull Request

- PR テンプレートのチェックリストを満たすこと。
- CI（構造検証 + 変更サンプルのビルド/テスト）がグリーンであること。
- サンプル追加時は README どおりに起動できることをレビュアーが確認できる状態にする。

### CI / 品質ゲート

- **構造検証**: 各サンプルに `README.md` と `metadata.yaml` が存在し、
  必須項目が埋まっているかを `tools/validate-samples.sh` でチェック。
- **サンプル個別ビルド**: 変更されたサンプルのみを対象に、各サンプルが提供する
  ビルド/テストタスクを実行（言語非依存にするため、各サンプルが手順を自前で定義）。

### バージョニング / ライフサイクル

各サンプルは `metadata.yaml` の `status` でライフサイクルを管理します。

| status | 意味 |
| --- | --- |
| `draft` | 作成中。動かない部分があってよい |
| `active` | 完成。README どおりに動作する |
| `maintenance` | メンテのみ。新機能追加はしない |
| `deprecated` | 非推奨。参考用に残すが更新しない |

---

## 🚀 はじめかた

```bash
# リポジトリを取得
git clone <this-repo-url>
cd web-architecture-sample

# サンプル一覧を見る
cat samples/README.md

# 個別サンプルを動かす（各サンプルの README に従う）
cd samples/<category>/<architecture-name>
cat README.md
```

---

## 🤝 コントリビューション

[CONTRIBUTING.md](CONTRIBUTING.md) を参照してください。

## 📄 ライセンス

このリポジトリのライセンスは別途定義します（未設定の場合は追加してください）。
