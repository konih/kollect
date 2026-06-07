import { defineConfig, devices } from "@playwright/test";

const port = 5173;
const baseURL = `http://127.0.0.1:${port}`;
const ci = !!process.env.CI;

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL,
    trace: "on-first-retry",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer: {
    // CI: preview serves pre-built assets (workflow build step); dev is too slow on cold runners.
    command: ci
      ? "npm run preview -- --host 127.0.0.1 --port 5173 --strictPort"
      : "npm run dev",
    url: baseURL,
    reuseExistingServer: !ci,
    timeout: ci ? 120_000 : 60_000,
    env: {
      VITE_MOCK_API: "true",
    },
  },
});
