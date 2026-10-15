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
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    {
      name: "mobile-chromium",
      use: {
        ...devices["iPhone 13"],
        browserName: "chromium",
        viewport: { width: 320, height: 700 },
      },
    },
    { name: "firefox", use: { ...devices["Desktop Firefox"] } },
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
