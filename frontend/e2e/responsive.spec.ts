import { expect, test } from "@playwright/test";

async function expectWithinViewport(
  page: import("@playwright/test").Page,
  needsSecondNumber = true,
) {
  const viewport = page.viewportSize();
  expect(viewport).not.toBeNull();
  const horizontalOverflow = await page.evaluate(
    () => document.documentElement.scrollWidth > window.innerWidth,
  );
  expect(horizontalOverflow).toBe(false);

  const controls = [
    page.getByLabel("Operation"),
    page.getByLabel("First number"),
    page.getByRole("button", { name: "Calculate" }),
    page.locator("output"),
  ];
  if (needsSecondNumber) {
    controls.splice(2, 0, page.getByLabel(/Second number|Number to find/));
  }
  const result = page.locator("output strong");
  if (await result.count()) controls.push(result);

  for (const locator of controls) {
    await expect(locator).toBeVisible();
    const box = await locator.boundingBox();
    expect(box).not.toBeNull();
    expect(box!.x).toBeGreaterThanOrEqual(0);
    expect(box!.x + box!.width).toBeLessThanOrEqual(viewport!.width);
  }
}

test("keeps calculator controls and long results usable at each viewport", async ({
  page,
}, testInfo) => {
  await page.goto("/");
  await expectWithinViewport(page);

  await page.getByLabel("Operation").selectOption("multiply");
  await page.getByLabel("First number").fill("1.2345678901234567e154");
  await page.getByLabel("Second number").fill("1e154");
  await page.getByRole("button", { name: "Calculate" }).click();
  await expect(page.locator("output strong")).toHaveText(
    "1.2345678901234567e+308",
  );
  expect(
    (await page.locator("output strong").innerText()).length,
  ).toBeGreaterThanOrEqual(20);
  await expectWithinViewport(page);
  await page.screenshot({
    path: testInfo.outputPath(`${testInfo.project.name}-long-result.png`),
    fullPage: true,
  });

  await page.getByLabel("First number").fill("1.2345678901234567e-154");
  await page.getByLabel("Second number").fill("1e-154");
  await page.getByRole("button", { name: "Calculate" }).click();
  await expect(page.locator("output strong")).toContainText("e-308");
  await expectWithinViewport(page);

  await page.getByLabel("Operation").selectOption("percentage");
  await expect(
    page.getByLabel("Number to find the percentage of"),
  ).toBeVisible();
  await page.getByLabel("First number").fill("25");
  await page.getByLabel("Number to find the percentage of").fill("80");
  await page.getByRole("button", { name: "Calculate" }).click();
  await expect(page.locator("output strong")).toHaveText("20");
  await expectWithinViewport(page);

  await page.getByLabel("Operation").selectOption("sqrt");
  await page.getByLabel("First number").fill("-1");
  await page.getByRole("button", { name: "Calculate" }).click();
  await expect(page.getByRole("alert")).toBeVisible();
  await expectWithinViewport(page, false);
  await page.screenshot({
    path: testInfo.outputPath(`${testInfo.project.name}-error-state.png`),
    fullPage: true,
  });
});
