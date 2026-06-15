"""メソッド + パス → ハンドラ関数の対応付け（純粋・テスト容易）。

API Gateway のルーティング相当を、外部依存なしの純粋な関数で表現する。
``{id}`` のようなパスパラメータを抽出し、event["pathParameters"] に詰める。

ルータ自体は HTTP を知らない。``match()`` は
「どのハンドラを、どんな pathParameters で呼ぶか」を返すだけなので、
ハンドラと同様に単体テストできる。
"""

from __future__ import annotations

from typing import Callable, Optional

from . import handlers

# レスポンス dict を返す純粋ハンドラの型。
Handler = Callable[[dict, object], dict]

# (httpMethod, パステンプレート) -> ハンドラ
# パステンプレートの "{id}" 部分は可変セグメントを表す。
ROUTES: dict[tuple[str, str], Handler] = {
    ("POST", "/notes"): handlers.create_note,
    ("GET", "/notes"): handlers.list_notes,
    ("GET", "/notes/{id}"): handlers.get_note,
}


def _segments(path: str) -> list[str]:
    """パスを空要素を除いたセグメント列に分解する。"""
    return [seg for seg in path.split("/") if seg != ""]


def _match_template(template: str, path: str) -> Optional[dict]:
    """パステンプレートと実パスを照合し、一致すれば pathParameters を返す。

    一致しなければ None。"{name}" は1セグメントの可変部として扱う。
    """
    t_segs = _segments(template)
    p_segs = _segments(path)
    if len(t_segs) != len(p_segs):
        return None

    params: dict[str, str] = {}
    for t, p in zip(t_segs, p_segs):
        if t.startswith("{") and t.endswith("}"):
            params[t[1:-1]] = p
        elif t != p:
            return None
    return params


def match(method: str, path: str) -> tuple[Optional[Handler], dict]:
    """method + path に対応するハンドラと pathParameters を返す。

    対応が無ければ (None, {}) を返す（呼び出し側で 404 を組み立てる）。
    """
    for (route_method, template), handler in ROUTES.items():
        if route_method != method:
            continue
        params = _match_template(template, path)
        if params is not None:
            return handler, params
    return None, {}
