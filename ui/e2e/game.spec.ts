import { test, expect, type Page, type BrowserContext } from '@playwright/test';
import { answerIfMyTurn, createGame, getLobbyGameId, joinExistingGame, joinLobby } from './helpers/game';

/**
 * Play through the entire game, alternating turns between page1 and page2.
 * Returns when a winner is declared on either page.
 */
async function playGame(page1: Page, page2: Page): Promise<string> {
  const maxRounds = 60;

  for (let round = 0; round < maxRounds; round++) {
    // Check if game is over on either page
    for (const page of [page1, page2]) {
      const gameOver = page.locator('.game-over');
      if (await gameOver.isVisible()) {
        const winnerText = (await page.locator('.game-over-winner').textContent()) ?? '';
        return winnerText.trim();
      }
    }

    // Try to answer on either page (whoever has "Your Turn")
    const answered = await answerIfMyTurn(page1) || await answerIfMyTurn(page2);

    if (!answered) {
      // Neither player has a turn yet — wait a bit for server to send question
      await page1.waitForTimeout(500);
    }
  }

  throw new Error('Game did not finish within the expected number of rounds');
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

test.describe('Full multiplayer game', () => {
  let ctx1: BrowserContext;
  let ctx2: BrowserContext;
  let player1: Page;
  let player2: Page;

  test.beforeEach(async ({ browser }) => {
    ctx1 = await browser.newContext();
    ctx2 = await browser.newContext();
    player1 = await ctx1.newPage();
    player2 = await ctx2.newPage();
  });

  test.afterEach(async () => {
    await ctx1.close();
    await ctx2.close();
  });

  test('two players can create, join, start, and finish a game', async () => {
    // --- Step 1: Player 1 creates a game ---
    await createGame(player1);
    const gameId = await getLobbyGameId(player1);
    expect(gameId).toBeTruthy();

    // --- Step 2: Player 2 joins that game ---
    await joinExistingGame(player2, gameId);

    // --- Step 3: Both players set their names ---
    // Do in parallel — both enter names concurrently
    await Promise.all([
      joinLobby(player1, 'Alice'),
      joinLobby(player2, 'Bob'),
    ]);

    // --- Step 4: Player 1 mixes teams then starts the game ---
    await player1.getByRole('button', { name: /Mix Teams/i }).click();

    // Wait for team assignments to reflect
    await player1.waitForTimeout(1_000);

    const startBtn = player1.getByRole('button', { name: /Rack Cups & Start/i });
    await expect(startBtn).toBeEnabled({ timeout: 5_000 });
    await startBtn.click();

    // Both pages should transition into the game view
    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });
    await expect(player2.locator('.game-board')).toBeVisible({ timeout: 10_000 });

    // --- Step 5: Play the game to completion ---
    const winner = await playGame(player1, player2);
    expect(winner).toBeTruthy();

    // Winner modal visible on at least one page
    const p1GameOver = player1.locator('.game-over');
    const p2GameOver = player2.locator('.game-over');
    const anyGameOver =
      (await p1GameOver.isVisible()) || (await p2GameOver.isVisible());
    expect(anyGameOver).toBe(true);

    const gameOverPage = (await p1GameOver.isVisible()) ? player1 : player2;
    await expect(gameOverPage.locator('.game-over-team')).toHaveCount(2);
    await expect(gameOverPage.locator('.game-over-cup-name')).toHaveCount(2);
    await expect(gameOverPage.locator('.game-over-cup.flipped')).toHaveCount(1);
    await expect(gameOverPage.locator('.game-over-cup.filled')).toHaveCount(1);
  });

  test('Play Again resets the game and returns to lobby', async () => {
    // Fast path — create game, single player, can't start without 2 but
    // we verify the lobby reset UI works after a game-over state.
    await createGame(player1);
    const gameId = await getLobbyGameId(player1);

    await joinExistingGame(player2, gameId);

    await Promise.all([
      joinLobby(player1, 'Alice'),
      joinLobby(player2, 'Bob'),
    ]);

    await player1.getByRole('button', { name: /Mix Teams/i }).click();
    await player1.waitForTimeout(1_000);
    await player1.getByRole('button', { name: /Rack Cups & Start/i }).click();

    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });

    await playGame(player1, player2);

    // Click "Play Again" on whichever page shows it
    const restartBtns = [
      player1.getByRole('button', { name: /Play Again/i }),
      player2.getByRole('button', { name: /Play Again/i }),
    ];

    for (const btn of restartBtns) {
      if (await btn.isVisible()) {
        await btn.click();
        break;
      }
    }

    // Should be back in lobby mode
    await expect(player1.locator('.lobby-icon')).toBeVisible({ timeout: 8_000 });
  });
});

test.describe('Lobby behaviour', () => {
  test('Rack Cups & Start warns when both teams do not have players', async ({ browser }) => {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    await createGame(page);
    await joinLobby(page, 'Solo');

    // Handle dialog explicitly
    page.on('dialog', async dialog => {
        expect(dialog.message()).toContain('Teams must have at least 1 player each before starting.');
        await dialog.accept();
    });

    const startBtn = page.getByRole('button', { name: /Rack Cups & Start/i });
    await startBtn.click();
    
    // Ensure we are still in lobby (game didn't start)
    await expect(page.locator('.lobby-icon')).toBeVisible();

    await ctx.close();
  });

  test('Game ID is displayed prominently in lobby', async ({ browser }) => {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    await createGame(page);
    const id = await getLobbyGameId(page);
    expect(id.length).toBeGreaterThan(0);

    await ctx.close();
  });

  test('Join button is disabled until name is entered', async ({ browser }) => {
    const ctx1 = await browser.newContext();
    const ctx2 = await browser.newContext();
    const host = await ctx1.newPage();
    const guest = await ctx2.newPage();

    await createGame(host);
    const gameId = await getLobbyGameId(host);
    await joinExistingGame(guest, gameId);

    const joinBtn = guest.getByRole('button', { name: /Join Game/i });
    await expect(joinBtn).toBeDisabled();

    await guest.locator('input[placeholder="Enter your name…"]').fill('Bob');
    await expect(joinBtn).toBeEnabled();

    await ctx1.close();
    await ctx2.close();
  });

  test('Creating a new game after joining still requires entering a name', async ({ browser }) => {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    await createGame(page);
    await joinLobby(page, 'Alice');

    await page.locator('.logo-link').click();
    await page.getByRole('button', { name: /Create New Game/i }).click();
    await page.locator('#qs-select').selectOption({ index: 1 });
    await page.locator('.submit-btn').click();

    await expect(page.locator('input[placeholder="Enter your name…"]')).toBeVisible();
    await expect(page.getByRole('button', { name: /Join Game/i })).toBeDisabled();

    await ctx.close();
  });
});
