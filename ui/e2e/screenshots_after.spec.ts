import { test, expect, type Page } from '@playwright/test';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const SCREENSHOT_DIR = path.resolve(__dirname, '../../docs/screenshots/after');

test('capture new UI screenshots', async ({ browser }) => {
    const context = await browser.newContext();
    const page = await context.newPage();

    // 1. Welcome Screen
    await page.goto('/');
    await expect(page.locator('.hero-title')).toBeVisible();
    await page.waitForTimeout(500);
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '01_welcome.png'), fullPage: true });

    // 2. New Game Screen
    await page.getByRole('button', { name: /Create New Game/i }).click();
    await expect(page.getByRole('heading', { name: 'Create a Game' })).toBeVisible();
    await page.waitForTimeout(500);
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '02_new_game.png'), fullPage: true });

    // 3. Create Game -> Lobby (Empty)
    await page.locator('#qs-select').selectOption({ index: 1 });
    await page.locator('.submit-btn').click();
    
    await expect(page.locator('.lobby-header')).toBeVisible();
    await page.waitForTimeout(500);
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '03_lobby_empty.png'), fullPage: true });

    // 4. Join Lobby as Host
    const gameId = (await page.locator('.game-code-value').textContent()) ?? '';
    await page.locator('input[placeholder="Enter your name…"]').fill('HostPlayer');
    await page.getByRole('button', { name: /Join Game/i }).click();
    await expect(page.getByText(/You're in as/)).toBeVisible();

    // 5. Add second player for "Ready" state
    const context2 = await browser.newContext();
    const page2 = await context2.newPage();
    await page2.goto('/');
    await page2.getByRole('button', { name: /Join Existing Game/i }).click();
    // Use first card, wait for it
    await expect(page2.locator('.game-card').first()).toBeVisible();
    await page2.locator('.game-card').filter({ hasText: gameId }).first().click();
    await page2.locator('input[placeholder="Enter your name…"]').fill('GuestPlayer');
    await page2.getByRole('button', { name: /Join Game/i }).click();

    // Back to Host - Shuffle
    await page.bringToFront();
    await page.getByRole('button', { name: /Shuffle Teams/i }).click();
    await page.waitForTimeout(1000); 
    
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '04_lobby_ready.png'), fullPage: true });

    // 6. Game View (Team A Perspective)
    await page.getByRole('button', { name: /Start Game/i }).click();
    await expect(page.locator('.game-board')).toBeVisible();
    await page.waitForTimeout(2000); // Wait for animations
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '05_game_view_team_a.png'), fullPage: true });

    // 7. Game View (Team B Perspective)
    await page2.bringToFront();
    await expect(page2.locator('.game-board')).toBeVisible();
    await page2.waitForTimeout(2000);
    await page2.screenshot({ path: path.join(SCREENSHOT_DIR, '06_game_view_team_b.png'), fullPage: true });

    await context.close();
    await context2.close();
});