"""ルータの単体テスト（HTTP 不要・純粋関数）。"""

import os
import sys
import unittest

SAMPLE_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
if SAMPLE_ROOT not in sys.path:
    sys.path.insert(0, SAMPLE_ROOT)

from src import handlers, router  # noqa: E402


class RouterTests(unittest.TestCase):
    def test_static_routes(self):
        handler, params = router.match("POST", "/notes")
        self.assertIs(handler, handlers.create_note)
        self.assertEqual(params, {})

        handler, params = router.match("GET", "/notes")
        self.assertIs(handler, handlers.list_notes)

    def test_path_parameter_extraction(self):
        handler, params = router.match("GET", "/notes/abc123")
        self.assertIs(handler, handlers.get_note)
        self.assertEqual(params, {"id": "abc123"})

    def test_unknown_route_returns_none(self):
        handler, params = router.match("DELETE", "/notes/abc123")
        self.assertIsNone(handler)
        self.assertEqual(params, {})

        handler, _ = router.match("GET", "/unknown")
        self.assertIsNone(handler)


if __name__ == "__main__":
    unittest.main()
