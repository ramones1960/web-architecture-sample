// order-service（注文サービス / port 4002）
//
// 担当チーム: 注文チーム。
// 責務: 注文の作成・取得。注文作成時に catalog-service を HTTP で同期呼び出しして
//       商品の存在・価格を検証する（サービス間通信）。
//       注文データは自身のインメモリストアのみに保持し、商品データは保持しない。
//
// エンドポイント:
//   POST /orders      -> {productId, qty} を受け取り、catalog で検証後に注文を作成
//   GET  /orders/:id  -> 単一注文（無ければ 404）

import http from "node:http";
import { readJson, sendJson, requestJson } from "../../shared/http.ts";
import {
  buildOrder,
  validateOrderRequest,
  type Order,
  type OrderRequest,
} from "./domain.ts";

const PORT = Number(process.env.ORDER_PORT ?? 4002);
// サービス間呼び出し先。設定で差し替え可能（環境ごとのサービスディスカバリを想定）。
const CATALOG_URL = process.env.CATALOG_URL ?? "http://localhost:4001";

// このサービス専用のインメモリデータストア。
const orders = new Map<string, Order>();
let seq = 0;
const nextId = () => `o${++seq}`;

interface CatalogProductResponse {
  product?: { id: string; name: string; priceYen: number };
  error?: string;
}

const server = http.createServer(async (req, res) => {
  const url = new URL(req.url ?? "/", `http://localhost:${PORT}`);
  const path = url.pathname;

  try {
    if (req.method === "POST" && path === "/orders") {
      const body = await readJson(req);
      const validationError = validateOrderRequest(body);
      if (validationError) {
        sendJson(res, 400, { error: "invalid_request", detail: validationError });
        return;
      }
      const orderReq = body as OrderRequest;

      // --- サービス間同期呼び出し: catalog-service で商品を検証 ---
      let catalogRes;
      try {
        catalogRes = await requestJson<CatalogProductResponse>(
          `${CATALOG_URL}/products/${encodeURIComponent(orderReq.productId)}`,
        );
      } catch {
        // ネットワーク障害（分散システム特有の失敗）。502 で表現する。
        sendJson(res, 502, {
          error: "catalog_unavailable",
          detail: "failed to reach catalog-service",
        });
        return;
      }

      if (catalogRes.status === 404 || !catalogRes.body.product) {
        sendJson(res, 400, {
          error: "unknown_product",
          productId: orderReq.productId,
        });
        return;
      }

      const order = buildOrder(nextId(), orderReq, catalogRes.body.product);
      orders.set(order.id, order);
      sendJson(res, 201, { order });
      return;
    }

    const match = path.match(/^\/orders\/([^/]+)$/);
    if (req.method === "GET" && match) {
      const id = decodeURIComponent(match[1]);
      const order = orders.get(id);
      if (!order) {
        sendJson(res, 404, { error: "order_not_found", id });
        return;
      }
      sendJson(res, 200, { order });
      return;
    }

    if (req.method === "GET" && path === "/health") {
      sendJson(res, 200, { service: "order", status: "ok" });
      return;
    }

    sendJson(res, 404, { error: "not_found", path });
  } catch (err) {
    sendJson(res, 500, {
      error: "internal_error",
      detail: err instanceof Error ? err.message : String(err),
    });
  }
});

server.listen(PORT, () => {
  console.log(`[order] listening on http://localhost:${PORT} (catalog=${CATALOG_URL})`);
});
