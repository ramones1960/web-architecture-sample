# microservices — マイクロサービス系

業務領域ごとにサービスを分割し、独立デプロイ・独立スケールを狙う構成のサンプル群。
組織・トラフィックの拡大に応じた分割のトレードオフを示す。

## 想定サンプル

- `service-per-domain` — ドメイン単位のサービス分割 + API Gateway
- `saga-orchestration` — 分散トランザクションを Saga で調整

## 収録サンプル

| サンプル | 言語 | 役割 | ステータス |
| --- | --- | --- | --- |
| [`service-per-domain-ts`](service-per-domain-ts/) | TypeScript (Node) | 大規模業務システムのバックエンド（API Gateway + ドメイン別サービス） | active |

> 追加方法は [追加ガイド](../../docs/guides/adding-a-new-architecture.md) を参照。
