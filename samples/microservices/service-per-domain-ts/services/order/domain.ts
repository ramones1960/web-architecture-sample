// order-service の純粋なドメインロジック。
//
// ここにはネットワーク・HTTP・I/O を一切含めない。
// そのため node:test で他サービスを起動せずに単体テストできる（offline 可能）。

export interface OrderRequest {
  productId: string;
  qty: number;
}

export interface Order {
  id: string;
  productId: string;
  productName: string;
  unitPriceYen: number;
  qty: number;
  totalYen: number;
}

/** 注文リクエストのバリデーション。問題があればエラーメッセージ、無ければ null。 */
export function validateOrderRequest(input: unknown): string | null {
  if (typeof input !== "object" || input === null) {
    return "request body must be an object";
  }
  const { productId, qty } = input as Record<string, unknown>;
  if (typeof productId !== "string" || productId.trim() === "") {
    return "productId is required";
  }
  if (typeof qty !== "number" || !Number.isInteger(qty) || qty <= 0) {
    return "qty must be a positive integer";
  }
  return null;
}

/** 注文合計金額の計算（単価 × 数量）。 */
export function calcOrderTotal(unitPriceYen: number, qty: number): number {
  if (!Number.isFinite(unitPriceYen) || unitPriceYen < 0) {
    throw new Error("unitPriceYen must be a non-negative number");
  }
  if (!Number.isInteger(qty) || qty <= 0) {
    throw new Error("qty must be a positive integer");
  }
  return unitPriceYen * qty;
}

/** 商品情報と注文リクエストから確定した注文を組み立てる。 */
export function buildOrder(
  id: string,
  req: OrderRequest,
  product: { id: string; name: string; priceYen: number },
): Order {
  return {
    id,
    productId: product.id,
    productName: product.name,
    unitPriceYen: product.priceYen,
    qty: req.qty,
    totalYen: calcOrderTotal(product.priceYen, req.qty),
  };
}
