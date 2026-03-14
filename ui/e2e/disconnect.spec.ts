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
  await expect(page.locator('.teams-preview')).toBeVisible({ timeout: 8_000 });
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
    
    await joinExistingGame(player2, gameId);

    await Promise.all([
      joinLobby(player1, 'Alice'),
      joinLobby(player2, 'Bob'),
    ]);

    // Start the game
    await player1.getByRole('button', { name: /Shuffle Teams/i }).click();
    await player1.waitForTimeout(1000); // Wait for shuffle
    await player1.getByRole('button', { name: /Start Game/i }).click();

    // Verify game started
    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });
    await expect(player2.locator('.game-board')).toBeVisible({ timeout: 10_000 });

    // --- Action: Player 1 refreshes the page ---
    console.log('Reloading Player 1 page...');
    await player1.reload();

    // --- Expectation: Player 1 should be back in the game or able to rejoin easily ---
    // If the game handles reconnects seamlessly, we expect to see the game board again.
    // If it requires manual rejoin, we might land on welcome/join screen.
    
    // For now, let's assert what currently happens or what we WANT to happen.
    // Ideally: seamless reconnect to game board.
    try {
      await expect(player1.locator('.game-board')).toBeVisible({ timeout: 5_000 });
      console.log('Success: Player 1 seamlessly reconnected to game board.');
      
      // Additional checks for bug regression (background black, question missing)
      // Check input is visible (means isMyTurn is true)
      await expect(player1.locator('.answer-input')).toBeVisible();
      
      // Check table color is NOT black (means myTeam is set)
      const table = player1.locator('.table');
      const style = await table.getAttribute('style');
      expect(style).not.toContain('--table-color: #000');
      expect(style).not.toContain('--table-color: rgb(0, 0, 0)');

    } catch (e) {
      console.log('Failed: Player 1 did not see game board immediately after reload.');
      
      // Check where we are
      if (await player1.locator('.welcome-screen').isVisible()) {
        console.log('Current state: Player 1 is on Welcome Screen.');
      } else if (await player1.locator('.lobby-icon').isVisible()) {
         console.log('Current state: Player 1 is in Lobby.');
      } else {
         console.log('Current state: Unknown.');
      }
      
      // Fail the test to signal we need to implement reconnect logic
      throw new Error('Player 1 lost game state after reload');
    }
  });
});
