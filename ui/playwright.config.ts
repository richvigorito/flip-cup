import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright boots the app stack automatically for local runs and CI.
 */
export default defineConfig({
  testDir: './e2e',
  timeout: 45_000,
  expect: { timeout: 10_000 },
  retries: 0,
  workers: 1, // serial — tests share game server state

  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    video: 'on-first-retry',
  },

  webServer: [
    {
      command: 'cd ../game-server && go run cmd/flipcup/main.go',
      url: 'http://127.0.0.1:8080/quizzes',
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
    {
      command: 'npm run dev -- --host 127.0.0.1',
      url: 'http://127.0.0.1:5173',
      env: {
        ...process.env,
        VITE_WS_URL: '127.0.0.1:8080',
      },
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
