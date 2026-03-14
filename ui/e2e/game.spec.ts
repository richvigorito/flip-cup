import { test, expect, type Page, type BrowserContext } from '@playwright/test';
import { DEFAULT_QUIZ_ANSWERS } from './answers';

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Create a game as player 1 and return to the lobby page. */
async function createGame(page: Page): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Create New Game/i }).click();
  await expect(page.locator('#qs-select')).toBeVisible();

  // Select "General Quiz Set 1" (_.default) — first real option
  await page.locator('#qs-select').selectOption({ index: 1 });
  await page.locator('.submit-btn').click();

  // Should land in lobby
  await expect(page.locator('.lobby-icon')).toBeVisible();
}

/** Read the Game ID from the lobby. */
async function getLobbyGameId(page: Page): Promise<string> {
  const el = page.locator('.game-code-value');
  await expect(el).toBeVisible();
  return (await el.textContent()) ?? '';
}

/** Enter a player name and join the lobby. */
async function joinLobby(page: Page, name: string): Promise<void> {
  const input = page.locator('input[placeholder="Enter your name…"]');
  await expect(input).toBeVisible();
  await input.fill(name);
  await page.getByRole('button', { name: /Join Game/i }).click();
  // Confirm team preview appears (the player is now in the lobby)
  await expect(page.locator('.teams-preview')).toBeVisible({ timeout: 8_000 });
}

/** Join an existing game from the Join screen. */
async function joinExistingGame(page: Page, gameId: string): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Join Existing Game/i }).click();

  // Wait for game cards to load
  await expect(page.locator('.game-card').first()).toBeVisible({ timeout: 10_000 });

  // Click the card matching the game ID
  const card = page.locator('.game-card', { hasText: gameId });
  await expect(card).toBeVisible();
  await card.click();

  await expect(page.locator('.lobby-icon')).toBeVisible();
}

/**
 * Answer the question shown on `page` if it's that player's turn.
 * Returns true if an answer was submitted, false if it wasn't this player's turn.
 */
async function answerIfMyTurn(page: Page): Promise<boolean> {
  const questionCard = page.locator('.question-card');
  const isVisible = await questionCard.isVisible();
  if (!isVisible) return false;

  // Read the question text
  const questionText = (await page.locator('.question-text').textContent()) ?? '';

  // Find the answer (exact match first, then try all known answers)
  const answer = DEFAULT_QUIZ_ANSWERS[questionText.trim()] ?? Object.values(DEFAULT_QUIZ_ANSWERS)[0];

  const input = page.locator('.answer-input');
  await input.fill(answer);
  await page.locator('.submit-btn').click();

  return true;
}

/**
 * Play through the entire game, alternating turns between page1 and page2.
 * Returns when a winner is declared on either page.
 */
async function playGame(page1: Page, page2: Page): Promise<string> {
  const maxRounds = 30;

  for (let round = 0; round < maxRounds; round++) {
    // Check if game is over on either page
    for (const page of [page1, page2]) {
      const gameOver = page.locator('.game-over');
      if (await gameOver.isVisible()) {
        const winnerText = (await page.locator('.winner-name').textContent()) ?? '';
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

    // --- Step 4: Player 1 shuffles teams then starts the game ---
    await player1.getByRole('button', { name: /Shuffle Teams/i }).click();

    // Wait for team assignments to reflect
    await player1.waitForTimeout(1_000);

    const startBtn = player1.getByRole('button', { name: /Start Game/i });
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

    await player1.getByRole('button', { name: /Shuffle Teams/i }).click();
    await player1.waitForTimeout(1_000);
    await player1.getByRole('button', { name: /Start Game/i }).click();

    await expect(player1.locator('.game-board')).toBeVisible({ timeout: 10_000 });

    await playGame(player1, player2);

    // Click "Play Again" on whichever page shows it
    const restartBtns = [
      player1.getByRole('button', { name: /Restart Game/i }),
      player2.getByRole('button', { name: /Restart Game/i }),
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
  test('Start Game is disabled with only one player', async ({ browser }) => {
    const ctx = await browser.newContext();
    const page = await ctx.newPage();

    await createGame(page);
    await joinLobby(page, 'Solo');

    const startBtn = page.getByRole('button', { name: /Start Game/i });
    // With only one player and one team empty, start should be disabled
    await expect(startBtn).toBeDisabled();

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
});
