import { defineConfig, devices } from "@playwright/test";

const port = 18080;
const baseURL = `http://127.0.0.1:${port}`;

export default defineConfig({
  testDir: "./e2e",
  outputDir: "test-results",
  timeout: 20_000,
  expect: { timeout: 5_000 },
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? [["list"], ["html", { open: "never" }]] : "list",
  use: {
    baseURL,
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
      testIgnore: ["**/numeric-boundaries.spec.ts", "**/responsive.spec.ts"],
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
      testIgnore: ["**/numeric-boundaries.spec.ts", "**/responsive.spec.ts"],
    },
    {
      name: "webkit-iphone-emulation",
      use: {
        ...devices["iPhone 13"],
      },
      testIgnore: ["**/numeric-boundaries.spec.ts", "**/responsive.spec.ts"],
    },
    {
      name: "numeric-chromium",
      use: { ...devices["Desktop Chrome"] },
      testMatch: "**/numeric-boundaries.spec.ts",
    },
    {
      name: "responsive-webkit-390",
      use: { ...devices["iPhone 13"] },
      testMatch: "**/responsive.spec.ts",
    },
    ...[
      ["responsive-320", 320, 700],
      ["responsive-375", 375, 760],
      ["responsive-390", 390, 844],
      ["responsive-768", 768, 1024],
      ["responsive-1024", 1024, 768],
      ["responsive-landscape", 844, 390],
    ].map(([name, width, height]) => ({
      name: name as string,
      use: {
        ...devices["Desktop Chrome"],
        viewport: { width: width as number, height: height as number },
      },
      testMatch: "**/responsive.spec.ts",
    })),
  ],
  webServer: {
    command:
      "pnpm build && cd ../backend && PORT=18080 STATIC_DIR=../frontend/dist go run ./cmd/server",
    cwd: ".",
    url: `${baseURL}/api/health`,
    timeout: 120_000,
    reuseExistingServer: false,
  },
});
