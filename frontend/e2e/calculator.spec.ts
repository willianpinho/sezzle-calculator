import { expect, test } from "@playwright/test";

const calculations = [
  ["add", "2", "3", "5"],
  ["subtract", "2", "3", "-1"],
  ["multiply", "-2", "3", "-6"],
  ["divide", "7", "2", "3.5"],
  ["power", "2", "8", "256"],
  ["sqrt", "9", undefined, "3"],
  ["percentage", "25", "80", "20"],
] as const;

test.describe("calculator", () => {
  for (const [operation, first, second, result] of calculations) {
    test(`calculates ${operation} through the live API`, async ({ page }) => {
      await page.goto("/");
      await page.getByLabel("Operation").selectOption(operation);
      await page.getByLabel("First number").fill(first);
      if (second)
        await page.getByLabel(/Second number|Number to find/).fill(second);
      const response = page.waitForResponse(
        (candidate) =>
          candidate.url().endsWith("/api/calculate") &&
          candidate.request().method() === "POST",
      );
      await page.getByRole("button", { name: "Calculate" }).click();
      expect((await response).ok()).toBe(true);
      await expect(page.locator("output strong")).toHaveText(result);
    });
  }

  test("accepts keyboard submission and clears a stale result after editing", async ({
    page,
  }) => {
    await page.goto("/");
    await page.getByLabel("First number").fill("2");
    await page.getByLabel("Second number").fill("3");
    await page.getByLabel("Second number").press("Enter");
    await expect(page.locator("output strong")).toHaveText("5");
    await page.getByLabel("First number").fill("4");
    await expect(page.locator("output strong")).toHaveCount(0);
    await expect(page.locator("output")).toContainText(
      "Enter numbers, then calculate.",
    );
  });

  test("shows form and server validation errors", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.getByRole("alert")).toContainText("valid number");

    await page.getByLabel("First number").fill("5");
    await page.getByLabel("Second number").fill("0");
    await page.getByLabel("Operation").selectOption("divide");
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.getByRole("alert")).toContainText(
      "no finite real result",
    );

    await page.getByLabel("Operation").selectOption("sqrt");
    await page.getByLabel("First number").fill("-1");
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.getByRole("alert")).toContainText(
      "no finite real result",
    );
  });

  test("shows an overflow error", async ({ page }) => {
    await page.goto("/");
    await page.getByLabel("Operation").selectOption("multiply");
    await page.getByLabel("First number").fill("1e308");
    await page.getByLabel("Second number").fill("1e308");
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.getByRole("alert")).toContainText(
      "no finite real result",
    );
  });

  test("recovers from an explicit network failure", async ({ page }) => {
    await page.goto("/");
    await page.getByLabel("First number").fill("2");
    await page.getByLabel("Second number").fill("3");
    await page.route("**/api/calculate", (route) => route.abort("failed"));
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.getByRole("alert")).toContainText(/failed|network/i);
    await page.unrouteAll({ behavior: "wait" });
    await page.getByRole("button", { name: "Calculate" }).click();
    await expect(page.locator("output strong")).toHaveText("5");
  });
});
