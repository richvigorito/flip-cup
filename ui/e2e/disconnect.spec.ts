import { test, expect, type Page, type BrowserContext } from '@playwright/test';

// ---------------------------------------------------------------------------
// Helpers (Copied/Adapted from game.spec.ts)
// ---------------------------------------------------------------------------

async function createGame(page: Page): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Create New Game/i }).click();
  await expect(page.locator('#qs-select')).toBeVisible();
  await page.locator('#qs-select').selectOption({ index: 1 });
  await page.locator('.submit-btn').click();
  await expect(page.locator('.lobby-icon')).toBeVisible();
}

async function getLobbyGameId(page: Page): Promise<string> {
  const el = page.locator('.game-code-value');
  await expect(el).toBeVisible();
  return (await el.textContent()) ?? '';
}

async function joinLobby(page: Page, name: string): Promise<void> {
  const input = page.locator('input[placeholder="Enter your name…"]');
  await expect(input).toBeVisible();
  await input.fill(name);
  await page.getByRole('button', { name: /Join Game/i }).click();
  await expect(page.getByRole('heading', { name: /^Lobby$/i })).toBeVisible({ timeout: 8_000 });
  await expect(page.getByText(new RegExp(`You're in as\\s+${name}`, 'i'))).toBeVisible({ timeout: 8_000 });
}

async function joinExistingGame(page: Page, gameId: string): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Join Existing Game/i }).click();
  const card = page.locator('.game-card', { hasText: gameId });
  await expect(card).toBeVisible({ timeout: 10_000 });
  await card.click();
  await expect(page.locator('.lobby-icon')).toBeVisible();
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

test.describe('WebSocket Disconnection Handling', () => {
  let ctx1: BrowserContext;
  let ctx2: BrowserContext;
  let player1: Page;
  let player2: Page;

  test.beforeEach(async ({ browser }) => {
    ctx1 = await browser.newContext();
    ctx2 = await browser.newContext();
    player1 = await ctx1.newPage();
    player2 = await ctx2.newPage();
    
    // Log console messages
    player1.on('console', msg => console.log(`P1 Console: ${msg.text()}`));
    player2.on('console', msg => console.log(`P2 Console: ${msg.text()}`));
  });

  test.afterEach(async () => {
    await ctx1.close();
    await ctx2.close();
  });

  test('Player can reconnect after refreshing the page', async () => {
    // --- Setup: Start a game with 2 players ---
    await createGame(player1);
    const gameId = await getLobbyGameId(player1);
    
    // Player 2 joins
    await joinExistingGame(player2, gameId);

    // Both join lobby
    await Promise.all([
      joinLobby(player1, 'Alice'),
      joinLobby(player2, 'Bob'),
    ]);

    // Mix teams and start
    await player1.getByRole('button', { name: /Mix Teams/i }).click();
    await player1.waitForTimeout(1000); 
    await player1.getByRole('button', { name: /Rack Cups & Start/i }).click();

    // Verify game started
    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });
    await expect(player2.locator('.game-board')).toBeVisible({ timeout: 10_000 });

    // --- Action: Player 1 refreshes the page ---
    console.log('Reloading Player 1 page...');
    await player1.reload();

    // --- Expectation: Seamless Reconnect ---
    // The game should restore state from sessionStorage and reconnect socket
    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });
      
    // Check perspective is restored (not generic view)
    // Team A view has specific class
    const board = player1.locator('.game-board');
    await expect(board).toHaveClass(/team-(a|b)-view/);

    // Check if it's my turn (should see input)
    // Note: It might NOT be player 1's turn if they aren't first in the shuffled order.
    // But we at least expect to be "in the game" (not welcome screen)
  });
});
