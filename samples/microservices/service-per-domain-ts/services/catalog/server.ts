// catalog-service（カタログサービス / port 4001）
//
// 担当チーム: カタログチーム。
// 責務: 商品（id, name, priceYen）の管理。自身のインメモリデータストアのみを持ち、
//       他サービスのデータには一切アクセスしない（service per domain）。
//
// エンドポイント:
//   GET /products       -> 全商品一覧
//   GET /products/:id   -> 単一商品（無ければ 404）

import http from "node:http";
import { sendJson } from "../../shared/http.ts";

const PORT = Number(process.env.CATALOG_PORT ?? 4001);

export interface Product {
  id: string;
  name: string;
  priceYen: number;
}

// このサービス専用のインメモリデータストア（シードデータ）。
const products: ReadonlyArray<Product> = [
  { id: "p1", name: "コーヒー豆 200g", priceYen: 1200 },
  { id: "p2", name: "ステンレスマグ", priceYen: 2800 },
  { id: "p3", name: "ドリッパー", priceYen: 1500 },
];

const server = http.createServer((req, res) => {
  const url = new URL(req.url ?? "/", `http://localhost:${PORT}`);
  const path = url.pathname;

  if (req.method === "GET" && path === "/products") {
    sendJson(res, 200, { products });
    return;
  }

  const match = path.match(/^\/products\/([^/]+)$/);
  if (req.method === "GET" && match) {
    const id = decodeURIComponent(match[1]);
    const product = products.find((p) => p.id === id);
    if (!product) {
      sendJson(res, 404, { error: "product_not_found", id });
      return;
    }
    sendJson(res, 200, { product });
    return;
  }

  if (req.method === "GET" && path === "/health") {
    sendJson(res, 200, { service: "catalog", status: "ok" });
    return;
  }

  sendJson(res, 404, { error: "not_found", path });
});

server.listen(PORT, () => {
  console.log(`[catalog] listening on http://localhost:${PORT}`);
});
