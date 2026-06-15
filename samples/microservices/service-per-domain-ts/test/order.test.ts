// order-service の純粋なドメインロジックの単体テスト。
//
// 他サービスの起動・ネットワークを必要としない（完全オフライン）。
// 実行: node --test test/order.test.ts

import { test } from "node:test";
import assert from "node:assert/strict";
import {
  buildOrder,
  calcOrderTotal,
  validateOrderRequest,
} from "../services/order/domain.ts";

test("calcOrderTotal: 単価 × 数量を計算する", () => {
  assert.equal(calcOrderTotal(1200, 2), 2400);
  assert.equal(calcOrderTotal(2800, 1), 2800);
});

test("calcOrderTotal: 不正な数量で例外", () => {
  assert.throws(() => calcOrderTotal(1200, 0));
  assert.throws(() => calcOrderTotal(1200, -1));
  assert.throws(() => calcOrderTotal(1200, 1.5));
});

test("calcOrderTotal: 不正な単価で例外", () => {
  assert.throws(() => calcOrderTotal(-1, 2));
  assert.throws(() => calcOrderTotal(Number.NaN, 2));
});

test("validateOrderRequest: 正常系は null", () => {
  assert.equal(validateOrderRequest({ productId: "p1", qty: 2 }), null);
});

test("validateOrderRequest: productId 欠落でエラー", () => {
  assert.equal(
    validateOrderRequest({ qty: 2 }),
    "productId is required",
  );
  assert.equal(
    validateOrderRequest({ productId: "  ", qty: 2 }),
    "productId is required",
  );
});

test("validateOrderRequest: qty が不正でエラー", () => {
  assert.equal(
    validateOrderRequest({ productId: "p1", qty: 0 }),
    "qty must be a positive integer",
  );
  assert.equal(
    validateOrderRequest({ productId: "p1", qty: -3 }),
    "qty must be a positive integer",
  );
  assert.equal(
    validateOrderRequest({ productId: "p1", qty: 1.5 }),
    "qty must be a positive integer",
  );
});

test("validateOrderRequest: オブジェクト以外でエラー", () => {
  assert.equal(validateOrderRequest(null), "request body must be an object");
  assert.equal(validateOrderRequest("x"), "request body must be an object");
});

test("buildOrder: 商品情報から確定注文を組み立てる", () => {
  const order = buildOrder(
    "o1",
    { productId: "p1", qty: 3 },
    { id: "p1", name: "コーヒー豆 200g", priceYen: 1200 },
  );
  assert.deepEqual(order, {
    id: "o1",
    productId: "p1",
    productName: "コーヒー豆 200g",
    unitPriceYen: 1200,
    qty: 3,
    totalYen: 3600,
  });
});
