import { expect, test } from "@playwright/test";

const applicationOrigin = "http://127.0.0.1:18080";

async function blockExternalRequests(page: import("@playwright/test").Page) {
  const externalRequests: string[] = [];
  await page.route("**/*", (route) => {
    if (new URL(route.request().url()).origin === applicationOrigin) {
      return route.continue();
    }
    externalRequests.push(route.request().url());
    return route.abort();
  });
  return externalRequests;
}

async function calculationOperation(page: import("@playwright/test").Page) {
  const operation = page
    .locator(".opblock", { hasText: "Calculate a finite real result" })
    .first();
  await operation.locator(".opblock-summary").click();
  await expect(operation.locator(".try-out__btn")).toBeVisible();
  return operation;
}

async function executeCalculation(
  page: import("@playwright/test").Page,
  operation: import("@playwright/test").Locator,
  body: string,
) {
  const tryItOut = operation.getByRole("button", {
    name: "Try it out",
    exact: true,
  });
  if (await tryItOut.count()) await tryItOut.click();
  await operation.locator(".body-param textarea").fill(body);
  const response = page.waitForResponse(
    (candidate) =>
      candidate.url() === `${applicationOrigin}/api/calculate` &&
      candidate.request().method() === "POST",
  );
  await operation.getByRole("button", { name: "Execute", exact: true }).click();
  return response;
}

test.describe("bundled Swagger UI", () => {
  test("loads the real same-origin contract without external requests", async ({
    page,
  }) => {
    const externalRequests = await blockExternalRequests(page);
    await page.goto("/docs/");
    await expect(page.locator(".swagger-ui")).toBeVisible();
    await expect(page.getByText("Sezzle Calculator API")).toBeVisible();
    await expect(
      page.locator(".opblock-summary-description", {
        hasText: "Calculate a finite real result",
      }),
    ).toBeVisible();
    await expect(page.locator(".swagger-ui")).toContainText("/api/calculate");
    expect(externalRequests).toEqual([]);
  });

  test("redirects the API documentation route to the bundled UI", async ({
    page,
  }) => {
    await page.goto("/api/docs");
    await expect(page).toHaveURL(/\/docs\/$/);
    await expect(page.locator(".swagger-ui")).toBeVisible();
  });

  test("executes real success and domain-error calculations", async ({
    page,
  }) => {
    const externalRequests = await blockExternalRequests(page);
    await page.goto("/docs/");
    const operation = await calculationOperation(page);

    const success = await executeCalculation(
      page,
      operation,
      '{"operation":"add","a":2,"b":3}',
    );
    expect(success.status()).toBe(200);
    expect(await success.json()).toEqual({ result: 5 });
    await expect(operation.locator(".live-responses-table")).toContainText(
      '{ "result": 5 }',
    );

    const domainError = await executeCalculation(
      page,
      operation,
      '{"operation":"divide","a":1,"b":0}',
    );
    expect(domainError.status()).toBe(422);
    expect((await domainError.json()).error.code).toBe("domain_error");
    await expect(operation.locator(".live-responses-table")).toContainText(
      "domain_error",
    );
    expect(externalRequests).toEqual([]);
  });
});
