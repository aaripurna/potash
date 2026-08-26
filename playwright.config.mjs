import { defineConfig, devices } from "@playwright/test";

const GO_PORT = 3001;
const VITE_PORT = 5173;

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? "list" : [["list"], ["html", { open: "never" }]],

  use: {
    baseURL: `http://localhost:${GO_PORT}`,
    trace: "on-first-retry",
  },

  // Uses the Chrome already installed on the machine instead of downloading
  // Playwright's bundled chromium. Drop `channel` and run
  // `bunx playwright install chromium` to use the bundled one instead.
  projects: [{ name: "chrome", use: { ...devices["Desktop Chrome"], channel: "chrome" } }],

  // Dev-mode pair: Go renders the page, vite serves the modules it points at.
  // Both are reused if already running locally.
  webServer: [
    {
      command: "bunx vite",
      port: VITE_PORT,
      reuseExistingServer: !process.env.CI,
      stdout: "ignore",
    },
    {
      command: `go run . serve -p ${GO_PORT}`,
      port: GO_PORT,
      reuseExistingServer: !process.env.CI,
      stdout: "ignore",
    },
  ],
});
