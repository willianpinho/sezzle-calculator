import { expect, test } from "@playwright/test";

const headers = { "Content-Type": "application/json" };

test.describe("documented binary64 boundaries", () => {
  test("returns defined finite and zero results", async ({ request }) => {
    const cases = [
      ['{"operation":"divide","a":0,"b":2}', 0],
      ['{"operation":"power","a":0,"b":0}', 1],
      ['{"operation":"power","a":-2,"b":3}', -8],
      ['{"operation":"power","a":9,"b":0.5}', 3],
      ['{"operation":"sqrt","a":5e-324}', 2.2227587494850775e-162],
      ['{"operation":"multiply","a":5e-324,"b":0.5}', 0],
      ['{"operation":"percentage","a":5e-324,"b":50}', 0],
      ['{"operation":"add","a":9007199254740992,"b":1}', 9007199254740992],
      [
        '{"operation":"percentage","a":100,"b":1.7976931348623157e308}',
        1.7976931348623157e308,
      ],
    ] as const;

    for (const [data, result] of cases) {
      const response = await request.post("/api/calculate", { data, headers });
      expect(response.status()).toBe(200);
      expect(await response.json()).toEqual({ result });
    }

    const signedZero = await request.post("/api/calculate", {
      data: '{"operation":"sqrt","a":-0}',
      headers,
    });
    expect(signedZero.status()).toBe(200);
    expect(Object.is((await signedZero.json()).result, -0)).toBe(true);
  });

  test("rejects non-real, zero-denominator, and overflow cases", async ({
    request,
  }) => {
    const cases = [
      '{"operation":"divide","a":0,"b":0}',
      '{"operation":"divide","a":1,"b":-0}',
      '{"operation":"power","a":0,"b":-1}',
      '{"operation":"power","a":0,"b":-2}',
      '{"operation":"power","a":-8,"b":0.3333333333333333}',
      '{"operation":"sqrt","a":-5e-324}',
      '{"operation":"add","a":1.7976931348623157e308,"b":1.7976931348623157e308}',
      '{"operation":"subtract","a":1.7976931348623157e308,"b":-1.7976931348623157e308}',
      '{"operation":"multiply","a":1.7976931348623157e308,"b":2}',
      '{"operation":"divide","a":1.7976931348623157e308,"b":5e-324}',
      '{"operation":"power","a":2,"b":1024}',
      '{"operation":"percentage","a":200,"b":1.7976931348623157e308}',
    ];

    for (const data of cases) {
      const response = await request.post("/api/calculate", { data, headers });
      expect(response.status()).toBe(422);
      expect((await response.json()).error.code).toBe("domain_error");
    }

    for (const data of [
      '{"operation":"sqrt","a":1e400}',
      '{"operation":"sqrt","a":-1e400}',
    ]) {
      const lexicalOverflow = await request.post("/api/calculate", {
        data,
        headers,
      });
      expect(lexicalOverflow.status()).toBe(422);
      expect((await lexicalOverflow.json()).error.code).toBe(
        "validation_error",
      );
    }
  });
});
