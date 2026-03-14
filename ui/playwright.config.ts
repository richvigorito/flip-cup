import { defineConfig, devices } from '@playwright/test';

/**
 * Prerequisites before running tests:
 *   1. Game server:  cd game-server && go run cmd/flipcup/main.go
 *   2. UI dev server: cd ui && npm run dev
 * Or use docker-compose up -d
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

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
