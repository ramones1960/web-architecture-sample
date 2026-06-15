"""インメモリのメモ（note）ストア。

このストアは「外部データストア（DynamoDB 等）の代役」です。
ローカルデモを動かすためだけにメモリ上へ保持しており、
本番の FaaS では以下の理由でこの実装は使えません。

- 関数はステートレス: 各呼び出しは別インスタンスで実行され得る。
- メモリはインスタンス間で共有されない（プロセスが消えれば消える）。
- スケールアウト時に複数インスタンスが立つと状態が分散・不整合になる。

本番では DynamoDB / Aurora Serverless / Firestore などの
外部ストアに置き換え、ハンドラからはこのモジュールと同じ
インターフェース（add / list_all / get / clear）越しにアクセスする。
"""

from __future__ import annotations

import threading
import uuid
from datetime import datetime, timezone
from typing import Optional


class InMemoryNoteStore:
    """メモを保持する最小のストア実装（ローカルデモ専用）。"""

    def __init__(self) -> None:
        # 同一プロセス内で local_invoker が複数スレッドから触るため、
        # 念のためロックで保護する。
        self._lock = threading.Lock()
        self._notes: dict[str, dict] = {}

    def add(self, text: str) -> dict:
        """メモを1件作成して保存し、保存したレコードを返す。"""
        note = {
            "id": uuid.uuid4().hex,
            "text": text,
            "createdAt": datetime.now(timezone.utc).isoformat(),
        }
        with self._lock:
            self._notes[note["id"]] = note
        return note

    def list_all(self) -> list[dict]:
        """全メモを作成日時の昇順で返す。"""
        with self._lock:
            return sorted(self._notes.values(), key=lambda n: n["createdAt"])

    def get(self, note_id: str) -> Optional[dict]:
        """ID 指定で1件取得。存在しなければ None。"""
        with self._lock:
            return self._notes.get(note_id)

    def clear(self) -> None:
        """全件削除（主にテストのリセット用）。"""
        with self._lock:
            self._notes.clear()


# プロセス内で共有する既定インスタンス。
# 本番ではこのシングルトンを外部ストアのクライアントに差し替える。
store = InMemoryNoteStore()
