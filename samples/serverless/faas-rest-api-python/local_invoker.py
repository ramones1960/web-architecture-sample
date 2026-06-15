"""ローカルインボーカ — API Gateway を模した最小の HTTP サーバ。

標準ライブラリの http.server だけで動く。役割は次の3つ。

1. 受け取った HTTP リクエストを API Gateway 形式の ``event`` dict に変換する
   （httpMethod / path / pathParameters / body[JSON 文字列] など）。
2. router で該当ハンドラを引き当て、``handler(event, context)`` を呼ぶ。
3. ハンドラが返した ``response`` dict（statusCode / headers / body）を
   実際の HTTP レスポンスへ変換して返す。

これにより、クラウドに一切繋がずローカルだけで FaaS の挙動を再現できる。
ハンドラ自体は HTTP を知らない純粋関数のままなので、
このインボーカは「ゲートウェイの差し替え可能な一実装」にすぎない。

起動:
    python3 local_invoker.py
ポート（既定 8000）は環境変数 PORT で変更可能。
"""

from __future__ import annotations

import json
import os
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from urllib.parse import urlsplit

from src import router

DEFAULT_PORT = 8000


def build_event(method: str, raw_path: str, body: str | None) -> dict:
    """HTTP リクエストの素材から API Gateway 形式の event を作る。

    pathParameters はルータ照合後に詰めるため、ここでは空にしておく。
    """
    parsed = urlsplit(raw_path)
    return {
        "httpMethod": method,
        "path": parsed.path,
        "pathParameters": {},
        "queryStringParameters": parsed.query or None,
        "headers": {},
        "body": body,
    }


def invoke(method: str, raw_path: str, body: str | None) -> dict:
    """event を組み立て、ルータでハンドラを引き当てて実行する。

    一致するルートが無ければ 404 のレスポンス dict を返す。
    戻り値は API Gateway 形式の response dict。
    """
    event = build_event(method, raw_path, body)
    handler, path_params = router.match(method, event["path"])

    if handler is None:
        return {
            "statusCode": 404,
            "headers": {"Content-Type": "application/json"},
            "body": json.dumps(
                {"message": "ルートが見つかりません"}, ensure_ascii=False
            ),
        }

    event["pathParameters"] = path_params
    # context はこのサンプルでは未使用（本番では request id 等が入る）。
    return handler(event, None)


class GatewayRequestHandler(BaseHTTPRequestHandler):
    """API Gateway の代役。HTTP <-> event/response の変換だけを担う。"""

    server_version = "FaaSLocalInvoker/1.0"

    def _read_body(self) -> str | None:
        length = int(self.headers.get("Content-Length") or 0)
        if length <= 0:
            return None
        return self.rfile.read(length).decode("utf-8")

    def _dispatch(self, method: str) -> None:
        body = self._read_body()
        response = invoke(method, self.path, body)

        payload = (response.get("body") or "").encode("utf-8")
        self.send_response(response.get("statusCode", 200))
        for key, value in (response.get("headers") or {}).items():
            self.send_header(key, value)
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def do_GET(self) -> None:  # noqa: N802 (http.server の規約)
        self._dispatch("GET")

    def do_POST(self) -> None:  # noqa: N802
        self._dispatch("POST")

    def log_message(self, fmt: str, *args) -> None:
        # 既定のログをやや簡素化（標準エラーへ）。
        print("[invoker] " + (fmt % args))


def main() -> None:
    port = int(os.environ.get("PORT", DEFAULT_PORT))
    server = ThreadingHTTPServer(("127.0.0.1", port), GatewayRequestHandler)
    print(f"FaaS local invoker (API Gateway 相当) を起動しました: "
          f"http://127.0.0.1:{port}")
    print("停止するには Ctrl+C を押してください。")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n停止します。")
    finally:
        server.server_close()


if __name__ == "__main__":
    main()
