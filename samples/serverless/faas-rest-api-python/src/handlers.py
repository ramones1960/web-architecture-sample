"""純粋なハンドラ関数群（AWS Lambda + API Gateway プロキシ統合スタイル）。

各エンドポイントは ``handler(event, context) -> response`` という
1つの関数として表現される。

- event : API Gateway が組み立てる dict。
  主なキー: ``httpMethod`` / ``path`` / ``pathParameters`` / ``body``（JSON 文字列）
- context: 実行コンテキスト（リクエスト ID など）。このサンプルでは未使用。
- response: ``{"statusCode": int, "headers": {...}, "body": "<JSON 文字列>"}``

これらの関数は HTTP サーバを一切必要としない。
event dict を渡すだけで単体テストできる（tests/test_handlers.py 参照）。
これが FaaS のプログラミングモデルの肝であり、
テスト容易性・移植性（別クラウドや別ゲートウェイへの載せ替え）を生む。
"""

from __future__ import annotations

import json
from typing import Any

from .store import store

JSON_HEADERS = {"Content-Type": "application/json"}


def _response(status_code: int, payload: Any) -> dict:
    """body を JSON 文字列にした API Gateway 形式のレスポンスを作る。"""
    return {
        "statusCode": status_code,
        "headers": JSON_HEADERS,
        "body": json.dumps(payload, ensure_ascii=False),
    }


def _parse_body(event: dict) -> dict:
    """event["body"]（JSON 文字列 or None）を dict にして返す。

    不正な JSON の場合は ValueError を送出する。
    """
    raw = event.get("body")
    if raw is None or raw == "":
        return {}
    if isinstance(raw, dict):
        # 念のため: 既に dict が渡ってきた場合も許容する。
        return raw
    return json.loads(raw)


def create_note(event: dict, context: Any = None) -> dict:
    """POST /notes — メモを作成する。

    body: {"text": "..."}
    成功時: 201 と {id, text, createdAt}
    text が無い/空/型不正: 400
    body が不正な JSON: 400
    """
    try:
        body = _parse_body(event)
    except ValueError:
        return _response(400, {"message": "リクエストボディが不正な JSON です"})

    text = body.get("text")
    if not isinstance(text, str) or text.strip() == "":
        return _response(400, {"message": "'text' は必須の文字列です"})

    note = store.add(text)
    return _response(201, note)


def list_notes(event: dict, context: Any = None) -> dict:
    """GET /notes — メモ一覧を返す。常に 200。"""
    return _response(200, store.list_all())


def get_note(event: dict, context: Any = None) -> dict:
    """GET /notes/{id} — メモを1件返す。

    見つかれば 200、無ければ 404。
    """
    path_params = event.get("pathParameters") or {}
    note_id = path_params.get("id")

    note = store.get(note_id) if note_id else None
    if note is None:
        return _response(404, {"message": "指定された ID のメモは存在しません"})

    return _response(200, note)
