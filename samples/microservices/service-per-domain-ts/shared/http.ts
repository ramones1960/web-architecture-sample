// 汎用的な JSON HTTP ヘルパー。
//
// ここにはドメインロジックを一切置かない（「サービスごとにドメインを分ける」原則を守るため）。
// 各サービスはリクエスト/レスポンスの組み立てとサービス間呼び出しにのみこれを使う。

import http from "node:http";

/** HTTP レスポンスを JSON で返す。 */
export function sendJson(
  res: http.ServerResponse,
  status: number,
  body: unknown,
): void {
  const payload = JSON.stringify(body);
  res.writeHead(status, {
    "Content-Type": "application/json; charset=utf-8",
    "Content-Length": Buffer.byteLength(payload),
  });
  res.end(payload);
}

/** リクエストボディを読み取り JSON としてパースする。空ボディなら {} を返す。 */
export async function readJson<T = unknown>(
  req: http.IncomingMessage,
): Promise<T> {
  const chunks: Buffer[] = [];
  for await (const chunk of req) {
    chunks.push(chunk as Buffer);
  }
  const raw = Buffer.concat(chunks).toString("utf-8").trim();
  if (raw === "") return {} as T;
  return JSON.parse(raw) as T;
}

/** サービス間の同期呼び出し（GET / POST 等）に使う汎用 JSON クライアント。 */
export function requestJson<T = unknown>(
  url: string,
  options: { method?: string; body?: unknown } = {},
): Promise<{ status: number; body: T }> {
  const { method = "GET", body } = options;
  const payload = body === undefined ? undefined : JSON.stringify(body);

  return new Promise((resolve, reject) => {
    const u = new URL(url);
    const req = http.request(
      {
        hostname: u.hostname,
        port: u.port,
        path: u.pathname + u.search,
        method,
        headers: {
          "Content-Type": "application/json; charset=utf-8",
          ...(payload
            ? { "Content-Length": Buffer.byteLength(payload) }
            : {}),
        },
      },
      (res) => {
        const chunks: Buffer[] = [];
        res.on("data", (c) => chunks.push(c as Buffer));
        res.on("end", () => {
          const raw = Buffer.concat(chunks).toString("utf-8").trim();
          let parsed: unknown = undefined;
          try {
            parsed = raw === "" ? undefined : JSON.parse(raw);
          } catch {
            parsed = raw;
          }
          resolve({ status: res.statusCode ?? 0, body: parsed as T });
        });
      },
    );
    req.on("error", reject);
    if (payload) req.write(payload);
    req.end();
  });
}
