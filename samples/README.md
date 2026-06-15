# samples — アーキテクチャサンプル索引

各アーキテクチャサンプルは `samples/<category>/<architecture-name>/` に配置されます。
カテゴリの定義は [docs/architecture-catalog.md](../docs/architecture-catalog.md) を参照してください。

## カテゴリ

| カテゴリ | 説明 | ディレクトリ |
| --- | --- | --- |
| モノリス系 | 単一デプロイ単位に集約する構成 | [`monolithic/`](monolithic/) |
| マイクロサービス系 | ドメイン単位でサービスを分割 | [`microservices/`](microservices/) |
| イベント駆動系 | イベントで疎結合に連携 | [`event-driven/`](event-driven/) |
| サーバーレス系 | FaaS / BaaS / エッジ | [`serverless/`](serverless/) |
| フロントエンド系 | SPA / SSR・SSG / マイクロフロントエンド / BFF | [`frontend/`](frontend/) |

## 収録サンプル一覧

> サンプルを追加したら、この表に 1 行追加してください。

| サンプル | カテゴリ | 言語 | 役割 | ステータス |
| --- | --- | --- | --- | --- |
| [`layered-architecture-go`](monolithic/layered-architecture-go/) | monolithic | Go | 業務アプリのバックエンド | active |
| [`service-per-domain-ts`](microservices/service-per-domain-ts/) | microservices | TypeScript | 大規模業務システム（複数チーム） | active |
| [`faas-rest-api-python`](serverless/faas-rest-api-python/) | serverless | Python | API サーバ（FaaS） | active |

## 新しいサンプルを追加するには

```bash
./tools/new-sample.sh <category> <architecture-name>
```

詳細は [新しいアーキテクチャの追加ガイド](../docs/guides/adding-a-new-architecture.md) を参照してください。
ひな型は [`_template/`](_template/) にあります。
