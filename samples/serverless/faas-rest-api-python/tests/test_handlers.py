"""ハンドラの単体テスト。

ポイント: HTTP サーバを一切立てず、event dict を直接渡してハンドラを呼ぶ。
これが FaaS のテスト容易性であり、純粋関数モデルの利点。

実行:
    python3 -m unittest discover -s tests
"""

import json
import os
import sys
import unittest

# tests/ から実行されるため、サンプルルートを import パスに追加する。
SAMPLE_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if SAMPLE_ROOT not in sys.path:
    sys.path.insert(0, SAMPLE_ROOT)

from src import handlers  # noqa: E402
from src.store import store  # noqa: E402


def make_event(method, path, body=None, path_params=None):
    """テスト用に API Gateway 形式の event を作る小さなヘルパ。"""
    return {
        "httpMethod": method,
        "path": path,
        "pathParameters": path_params or {},
        "body": json.dumps(body, ensure_ascii=False) if body is not None else None,
    }


def body_of(response):
    """response["body"]（JSON 文字列）を Python オブジェクトに戻す。"""
    return json.loads(response["body"])


class NoteHandlerTests(unittest.TestCase):
    def setUp(self):
        # テスト間で状態が漏れないようストアをリセットする。
        store.clear()

    def tearDown(self):
        store.clear()

    def test_create_returns_201_with_id(self):
        resp = handlers.create_note(make_event("POST", "/notes", {"text": "買い物"}))
        self.assertEqual(resp["statusCode"], 201)
        payload = body_of(resp)
        self.assertIn("id", payload)
        self.assertTrue(payload["id"])
        self.assertEqual(payload["text"], "買い物")
        self.assertIn("createdAt", payload)

    def test_create_with_missing_text_returns_400(self):
        resp = handlers.create_note(make_event("POST", "/notes", {}))
        self.assertEqual(resp["statusCode"], 400)

    def test_create_with_empty_text_returns_400(self):
        resp = handlers.create_note(make_event("POST", "/notes", {"text": "   "}))
        self.assertEqual(resp["statusCode"], 400)

    def test_create_with_invalid_json_returns_400(self):
        event = {"httpMethod": "POST", "path": "/notes",
                 "pathParameters": {}, "body": "{not-json"}
        resp = handlers.create_note(event)
        self.assertEqual(resp["statusCode"], 400)

    def test_list_returns_created_notes(self):
        handlers.create_note(make_event("POST", "/notes", {"text": "1つ目"}))
        handlers.create_note(make_event("POST", "/notes", {"text": "2つ目"}))

        resp = handlers.list_notes(make_event("GET", "/notes"))
        self.assertEqual(resp["statusCode"], 200)
        notes = body_of(resp)
        self.assertEqual(len(notes), 2)
        texts = {n["text"] for n in notes}
        self.assertEqual(texts, {"1つ目", "2つ目"})

    def test_list_is_empty_initially(self):
        resp = handlers.list_notes(make_event("GET", "/notes"))
        self.assertEqual(resp["statusCode"], 200)
        self.assertEqual(body_of(resp), [])

    def test_get_returns_200_for_existing(self):
        created = body_of(
            handlers.create_note(make_event("POST", "/notes", {"text": "memo"}))
        )
        resp = handlers.get_note(
            make_event("GET", f"/notes/{created['id']}", path_params={"id": created["id"]})
        )
        self.assertEqual(resp["statusCode"], 200)
        self.assertEqual(body_of(resp)["id"], created["id"])

    def test_get_returns_404_for_missing(self):
        resp = handlers.get_note(
            make_event("GET", "/notes/does-not-exist",
                       path_params={"id": "does-not-exist"})
        )
        self.assertEqual(resp["statusCode"], 404)


if __name__ == "__main__":
    unittest.main()
