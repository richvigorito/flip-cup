import { expect, type Page } from '@playwright/test';
import { DEFAULT_QUIZ_ANSWERS } from '../answers';

export async function createGame(page: Page): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Create New Game/i }).click();
  await expect(page.locator('#qs-select')).toBeVisible();
  await page.locator('#qs-select').selectOption({ index: 1 });
  await page.locator('.submit-btn').click();
  await expect(page.locator('.lobby-icon')).toBeVisible();
}

export async function getLobbyGameId(page: Page): Promise<string> {
  const el = page.locator('.game-code-value');
  await expect(el).toBeVisible();
  return (await el.textContent()) ?? '';
}

export async function joinLobby(page: Page, name: string): Promise<void> {
  const input = page.locator('input[placeholder="Enter your name…"]');
  await expect(input).toBeVisible();
  await input.fill(name);
  await page.getByRole('button', { name: /Join Game/i }).click();
  await expect(page.getByRole('heading', { name: /^Lobby$/i })).toBeVisible({ timeout: 8_000 });
  await expect(page.getByText(new RegExp(`You're in as\\s+${name}`, 'i'))).toBeVisible({ timeout: 8_000 });
}

export async function joinExistingGame(page: Page, gameId: string): Promise<void> {
  await page.goto('/');
  await page.getByRole('button', { name: /Join Existing Game/i }).click();
  await expect(page.locator('.game-card').first()).toBeVisible({ timeout: 10_000 });

  const card = page.locator('.game-card', { hasText: gameId });
  await expect(card).toBeVisible();
  await card.click();

  await expect(page.locator('.lobby-icon')).toBeVisible();
}

export async function answerIfMyTurn(page: Page, answer?: string): Promise<boolean> {
  const questionCard = page.locator('.question-card');
  if (!(await questionCard.isVisible())) {
    return false;
  }

  const questionText = ((await page.locator('.question-text').textContent()) ?? '').trim();
  const resolvedAnswer = answer ?? DEFAULT_QUIZ_ANSWERS[questionText] ?? Object.values(DEFAULT_QUIZ_ANSWERS)[0];

  const input = page.locator('.answer-input');
  await input.fill(resolvedAnswer);
  await page.locator('.submit-btn').click();

  return true;
}
