// api-gateway（API ゲートウェイ / port 4000）
//
// 責務: クライアントに対する唯一の公開窓口（single public surface）。
//       パスに応じて適切なドメインサービスへリクエストをルーティングする。
//       ドメインロジックは持たず、ルーティングと転送のみを担う。
//
// ルーティング:
//   /catalog/*  -> catalog-service (4001)  例: /catalog/products -> /products
//   /orders/*   -> order-service   (4002)  例: /orders          -> /orders
//   /orders     -> order-service   (4002)

import http from "node:http";

const PORT = Number(process.env.GATEWAY_PORT ?? 4000);
const CATALOG_URL = process.env.CATALOG_URL ?? "http://localhost:4001";
const ORDER_URL = process.env.ORDER_URL ?? "http://localhost:4002";

interface Route {
  // 一致するパスの接頭辞
  prefix: string;
  // 転送先のベース URL
  target: string;
  // ゲートウェイのパスを下流サービスのパスへ書き換える
  rewrite: (path: string) => string;
}

const routes: Route[] = [
  {
    prefix: "/catalog",
    target: CATALOG_URL,
    // /catalog/products -> /products
    rewrite: (p) => p.replace(/^\/catalog/, "") || "/",
  },
  {
    prefix: "/orders",
    target: ORDER_URL,
    // /orders, /orders/:id はそのまま order-service の /orders, /orders/:id へ
    rewrite: (p) => p,
  },
];

function matchRoute(path: string): Route | undefined {
  return routes.find(
    (r) => path === r.prefix || path.startsWith(r.prefix + "/"),
  );
}

const server = http.createServer((req, res) => {
  const url = new URL(req.url ?? "/", `http://localhost:${PORT}`);
  const path = url.pathname;

  if (req.method === "GET" && path === "/health") {
    res.writeHead(200, { "Content-Type": "application/json; charset=utf-8" });
    res.end(JSON.stringify({ service: "gateway", status: "ok" }));
    return;
  }

  const route = matchRoute(path);
  if (!route) {
    res.writeHead(404, { "Content-Type": "application/json; charset=utf-8" });
    res.end(JSON.stringify({ error: "no_route", path }));
    return;
  }

  const target = new URL(route.target);
  const downstreamPath = route.rewrite(path) + url.search;

  // リクエストボディを含めてストリームでそのまま転送する（透過プロキシ）。
  const proxyReq = http.request(
    {
      hostname: target.hostname,
      port: target.port,
      path: downstreamPath,
      method: req.method,
      headers: req.headers,
    },
    (proxyRes) => {
      res.writeHead(proxyRes.statusCode ?? 502, proxyRes.headers);
      proxyRes.pipe(res);
    },
  );

  proxyReq.on("error", () => {
    res.writeHead(502, { "Content-Type": "application/json; charset=utf-8" });
    res.end(
      JSON.stringify({
        error: "bad_gateway",
        detail: `failed to reach ${route.target}`,
      }),
    );
  });

  req.pipe(proxyReq);
});

server.listen(PORT, () => {
  console.log(`[gateway] listening on http://localhost:${PORT}`);
  console.log(`[gateway]   /catalog/* -> ${CATALOG_URL}`);
  console.log(`[gateway]   /orders/*  -> ${ORDER_URL}`);
});
