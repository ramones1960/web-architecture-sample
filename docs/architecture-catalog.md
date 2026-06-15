# アーキテクチャカタログ

このリポジトリで扱うアーキテクチャの **分類（カテゴリ）** と、
各カテゴリに含める代表的なアーキテクチャの一覧です。
カテゴリは `samples/<category>/` のディレクトリに対応します。

> ✅ = 収録済み / 🚧 = 作成中(draft) / 📝 = 収録予定(未着手)
> ステータスは各サンプルの `metadata.yaml` を正とし、本表は索引として更新します。

---

## カテゴリ定義

### `monolithic` — モノリス系
1 つのデプロイ単位に機能を集約する構成。小〜中規模や立ち上げ期に有効。

| アーキテクチャ | 概要 | 状態 |
| --- | --- | --- |
| layered-architecture | プレゼン / アプリ / ドメイン / インフラの層分割 | ✅ [`layered-architecture-go`](../samples/monolithic/layered-architecture-go/) |
| modular-monolith | モジュール境界を持つ単一デプロイ | 📝 |
| clean-architecture | 依存方向を内向きに統一（ヘキサゴナル / オニオン） | 📝 |

### `microservices` — マイクロサービス系
業務領域ごとにサービスを分割し、独立デプロイ・独立スケールを狙う構成。

| アーキテクチャ | 概要 | 状態 |
| --- | --- | --- |
| service-per-domain | ドメイン単位のサービス分割 + API Gateway | ✅ [`service-per-domain-ts`](../samples/microservices/service-per-domain-ts/) |
| saga-orchestration | 分散トランザクションを Saga で調整 | 📝 |

### `event-driven` — イベント駆動系
状態変化をイベントとして伝播し、疎結合・非同期で連携する構成。

| アーキテクチャ | 概要 | 状態 |
| --- | --- | --- |
| pub-sub | メッセージブローカ経由の発行/購読 | 📝 |
| cqrs-event-sourcing | コマンド/クエリ分離 + イベントソーシング | 📝 |

### `serverless` — サーバーレス系
FaaS / BaaS / エッジ実行を中心とした、運用負荷とコストを抑える構成。

| アーキテクチャ | 概要 | 状態 |
| --- | --- | --- |
| faas-rest-api | 関数ベースの REST API | ✅ [`faas-rest-api-python`](../samples/serverless/faas-rest-api-python/) |
| edge-rendering | エッジでのレンダリング/配信 | 📝 |

### `frontend` — フロントエンド系
クライアント側／配信側のアーキテクチャ。

| アーキテクチャ | 概要 | 状態 |
| --- | --- | --- |
| spa | シングルページアプリケーション | 📝 |
| ssr-ssg | サーバーサイドレンダリング / 静的生成 | 📝 |
| micro-frontends | 画面単位での独立デプロイ | 📝 |
| bff | Backend for Frontend | 📝 |

---

## カテゴリの追加について

新しいカテゴリ（例: `data-streaming`、`p2p` など）が必要になった場合は、

1. このカタログにカテゴリ定義を追記し、
2. `samples/<new-category>/README.md` を作成し、
3. なぜ新カテゴリが必要かを ADR に残す（推奨）

という手順で追加してください。既存カテゴリに収まるものは新設しないでください。
