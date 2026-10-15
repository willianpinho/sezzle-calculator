import { expect, test } from "@playwright/test";

const headers = { "Content-Type": "application/json" };

test.describe("live API contract", () => {
  test("reports health and rejects malformed or invalid calculations", async ({
    request,
  }) => {
    const health = await request.get("/api/health");
    expect(health.ok()).toBe(true);
    expect(await health.json()).toEqual({ status: "ok" });

    const cases = [
      ["malformed JSON", "{", 400, "invalid_json"],
      [
        "unknown operation",
        '{"operation":"mod","a":1,"b":2}',
        422,
        "validation_error",
      ],
      [
        "extra field",
        '{"operation":"add","a":1,"b":2,"extra":true}',
        422,
        "validation_error",
      ],
      [
        "sqrt with b",
        '{"operation":"sqrt","a":9,"b":1}',
        422,
        "validation_error",
      ],
      ["negative sqrt", '{"operation":"sqrt","a":-1}', 422, "domain_error"],
      [
        "division by zero",
        '{"operation":"divide","a":1,"b":0}',
        422,
        "domain_error",
      ],
      [
        "overflow",
        '{"operation":"multiply","a":1e308,"b":1e308}',
        422,
        "domain_error",
      ],
    ] as const;

    for (const [, data, status, code] of cases) {
      const response = await request.post("/api/calculate", { data, headers });
      expect(response.status()).toBe(status);
      expect((await response.json()).error.code).toBe(code);
    }
  });
});
