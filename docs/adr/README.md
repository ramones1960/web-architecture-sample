# Architecture Decision Records (ADR)

ADR は「なぜその設計・技術を選んだのか」を記録する軽量ドキュメントです。
[Michael Nygard 形式](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) に準拠します。

## ルール

- 1 つの決定 = 1 ファイル。`NNNN-<kebab-title>.md`（連番 4 桁ゼロ埋め）。
- 一度 `Accepted` になった ADR は **書き換えず**、覆す場合は新しい ADR を作り、
  旧 ADR の Status を `Superseded by ADR-NNNN` に更新する。
- 雛形は [template.md](template.md) をコピーして使う。

## 索引

| # | タイトル | Status |
| --- | --- | --- |
| [0001](0001-record-architecture-decisions.md) | ADR を導入して意思決定を記録する | Accepted |
| [0002](0002-monorepo-with-independent-samples.md) | モノレポ + サンプル自己完結方針 | Accepted |
