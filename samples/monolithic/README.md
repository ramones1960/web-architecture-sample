# monolithic — モノリス系

機能を 1 つのデプロイ単位に集約する構成のサンプル群。
小〜中規模、プロダクト立ち上げ期、チームが小さい場合に有効。

## 想定サンプル

- `layered-architecture` — プレゼン / アプリ / ドメイン / インフラの層分割
- `modular-monolith` — モジュール境界を持つ単一デプロイ
- `clean-architecture` — 依存方向を内向きに統一（ヘキサゴナル / オニオン）

## 収録サンプル

| サンプル | 言語 | 役割 | ステータス |
| --- | --- | --- | --- |
| [`layered-architecture-go`](layered-architecture-go/) | Go | 業務アプリのバックエンド（タスク管理 API） | active |

> 追加方法は [追加ガイド](../../docs/guides/adding-a-new-architecture.md) を参照。
