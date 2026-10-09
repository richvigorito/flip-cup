import { test, expect, type BrowserContext, type Page } from '@playwright/test';
import { answerIfMyTurn, createGame, getLobbyGameId, joinExistingGame, joinLobby } from './helpers/game';

// Two players end up one per team, so each team's turn order is just that player.
async function startTwoPlayerGame(
  host: Page,
  guest: Page,
  questionsPerPlayer?: number,
): Promise<void> {
  await createGame(host);
  const gameId = await getLobbyGameId(host);
  await joinExistingGame(guest, gameId);
  await Promise.all([joinLobby(host, 'Alice'), joinLobby(guest, 'Bob')]);

  await host.getByRole('button', { name: /Mix Teams/i }).click();
  await host.waitForTimeout(1_000);

  if (questionsPerPlayer) {
    await host.locator('#qpp-select').selectOption(String(questionsPerPlayer));
    await expect(guest.locator('#qpp-select')).toHaveValue(String(questionsPerPlayer));
  }

  const startButton = host.getByRole('button', { name: /Rack Cups & Start/i });
  await expect(startButton).toBeEnabled({ timeout: 5_000 });
  await startButton.click();

  await Promise.all([host, guest].map((page) => expect(page.locator('.question-card')).toBeVisible({ timeout: 10_000 })));
}

test.describe('Skip and questions-per-player', () => {
  let ctx1: BrowserContext;
  let ctx2: BrowserContext;
  let host: Page;
  let guest: Page;

  test.beforeEach(async ({ browser }) => {
    ctx1 = await browser.newContext();
    ctx2 = await browser.newContext();
    host = await ctx1.newPage();
    guest = await ctx2.newPage();
  });

  test.afterEach(async () => {
    await ctx1.close();
    await ctx2.close();
  });

  test('a player must answer N questions before their team wins', async () => {
    await startTwoPlayerGame(host, guest, 2);

    const firstQuestion = ((await host.locator('.question-text').textContent()) ?? '').trim();
    await answerIfMyTurn(host);

    // Not done yet: the same player gets a second, different question and nobody has won.
    await expect(host.locator('.question-text')).not.toHaveText(firstQuestion);
    await expect(host.locator('.question-card')).toBeVisible();
    await expect(host.locator('.game-over')).toHaveCount(0);

    await answerIfMyTurn(host);
    await expect(host.locator('.game-over')).toBeVisible({ timeout: 10_000 });
    await expect(guest.locator('.game-over')).toBeVisible({ timeout: 10_000 });
  });

  test('skipping notifies every player and gives the skipper a new question', async () => {
    await startTwoPlayerGame(host, guest);

    const before = ((await host.locator('.question-text').textContent()) ?? '').trim();
    await host.getByRole('button', { name: 'Skip' }).click();

    await expect(host.getByTestId('skip-banner')).toContainText('Alice is stuck and skipped!');
    await expect(guest.getByTestId('skip-banner')).toContainText('Alice is stuck and skipped!');
    await expect(host.locator('.question-text')).not.toHaveText(before);
    await expect(host.locator('.question-card')).toBeVisible();
  });

  test('a skip delays the next question after the correct answer', async () => {
    await startTwoPlayerGame(host, guest, 2);

    const before = ((await host.locator('.question-text').textContent()) ?? '').trim();
    await host.getByRole('button', { name: 'Skip' }).click();
    await expect(host.locator('.question-text')).not.toHaveText(before);
    await answerIfMyTurn(host);

    // The next question is held back for the penalty, then arrives.
    await expect(host.getByTestId('penalty-wait')).toBeVisible();
    await expect(host.locator('.question-card')).toHaveCount(0);
    await expect(host.locator('.question-card')).toBeVisible({ timeout: 8_000 });
  });

  test('a skip delays the win when it happens on the final question', async () => {
    await startTwoPlayerGame(host, guest);

    const before = ((await host.locator('.question-text').textContent()) ?? '').trim();
    await host.getByRole('button', { name: 'Skip' }).click();
    await expect(host.locator('.question-text')).not.toHaveText(before);
    await answerIfMyTurn(host); // last answer for a one-player team

    await expect(host.getByTestId('penalty-wait')).toBeVisible();
    await expect(host.locator('.game-over')).toHaveCount(0);
    await expect(host.locator('.game-over')).toBeVisible({ timeout: 8_000 });
    await expect(guest.locator('.game-over')).toBeVisible({ timeout: 8_000 });
  });
});
