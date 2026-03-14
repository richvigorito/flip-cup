import { test, expect, type Page } from '@playwright/test';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const SCREENSHOT_DIR = path.resolve(__dirname, '../../docs/screenshots/before');

test('capture old UI screenshots', async ({ browser }) => {
    const context = await browser.newContext();
    const page = await context.newPage();

    // 1. Welcome Screen
    await page.goto('/');
    // Old UI likely has a title or button
    await page.waitForTimeout(1000);
    await page.screenshot({ path: path.join(SCREENSHOT_DIR, '01_welcome.png'), fullPage: true });

    // 2. New Game Screen
    // Old UI: "Create New Game" button?
    // Let's assume text is similar or check source if fails.
    // If fail, I can check old source via git show main:ui/src/components/Welcome.svelte
    const createBtn = page.getByRole('button', { name: /Create New Game/i });
    if (await createBtn.isVisible()) {
        await createBtn.click();
        await page.waitForTimeout(500);
        await page.screenshot({ path: path.join(SCREENSHOT_DIR, '02_new_game.png'), fullPage: true });
    }

    // 3. Create Game -> Lobby (Empty)
    // Old UI used a select for quiz?
    const select = page.locator('select');
    if (await select.isVisible()) {
        await select.selectOption({ index: 1 });
        await page.getByRole('button', { name: /Create Game/i }).click();
        await page.waitForTimeout(1000);
        await page.screenshot({ path: path.join(SCREENSHOT_DIR, '03_lobby_empty.png'), fullPage: true });
    }

    // 4. Join Lobby as Host
    // Old UI input for name
    const nameInput = page.locator('input[type="text"]');
    if (await nameInput.isVisible()) {
        const gameIdText = await page.textContent('body'); // Hacky to get ID?
        // Actually, let's just screenshot.
        await nameInput.fill('HostPlayer');
        await page.getByRole('button', { name: /Join Game/i }).click();
        await page.waitForTimeout(500);
    }
    
    // 5. Add second player
    const context2 = await browser.newContext();
    const page2 = await context2.newPage();
    await page2.goto('/');
    await page2.getByRole('button', { name: /Join Existing Game/i }).click();
    await page2.waitForTimeout(500);
    // Old UI might just list games or ask for ID.
    // Assuming list:
    const gameCard = page2.locator('.game-card, .card').first();
    if (await gameCard.isVisible()) {
        await gameCard.click();
        await page2.locator('input[type="text"]').fill('GuestPlayer');
        await page2.getByRole('button', { name: /Join Game/i }).click();
    }

    // Back to Host - Shuffle?
    await page.bringToFront();
    const shuffleBtn = page.getByRole('button', { name: /Shuffle Teams/i });
    if (await shuffleBtn.isVisible()) {
        await shuffleBtn.click();
        await page.waitForTimeout(500);
        await page.screenshot({ path: path.join(SCREENSHOT_DIR, '04_lobby_ready.png'), fullPage: true });
    
        // 6. Start Game
        await page.getByRole('button', { name: /Start Game/i }).click();
        await page.waitForTimeout(1000);
        await page.screenshot({ path: path.join(SCREENSHOT_DIR, '05_game_view_team_a.png'), fullPage: true });
    }

    // 7. Team B View
    await page2.bringToFront();
    await page2.waitForTimeout(500);
    await page2.screenshot({ path: path.join(SCREENSHOT_DIR, '06_game_view_team_b.png'), fullPage: true });

    await context.close();
    await context2.close();
});
