import { test, expect, type BrowserContext, type Page } from '@playwright/test';
import { answerIfMyTurn, createGame, getLobbyGameId, joinExistingGame, joinLobby } from './helpers/game';

// Generous: the suite shares one server and runs after 10-player scenarios.
const POLL_TIMEOUT = 20_000;
const NAMES = ['Alice', 'Bob', 'Carol', 'Dave'];

async function questionText(page: Page): Promise<string> {
  const card = page.locator('.question-text');
  return (await card.isVisible()) ? ((await card.textContent()) ?? '').trim() : '';
}

async function isTeamA(page: Page): Promise<boolean> {
  return ((await page.locator('.game-board').getAttribute('class')) ?? '').includes('team-a-view');
}

// Returns the page (among `pages`) whose player currently holds a question.
async function waitForActive(pages: Page[]): Promise<Page> {
  let active: Page | undefined;
  await expect
    .poll(async () => {
      for (const page of pages) {
        if (await page.locator('.question-card').isVisible()) {
          active = page;
          return true;
        }
      }
      return false;
    }, { timeout: POLL_TIMEOUT })
    .toBe(true);
  return active!;
}

// Answers on `page` and waits until its question changes or goes away.
async function answerAndWait(page: Page): Promise<void> {
  const before = await questionText(page);
  await answerIfMyTurn(page);
  await expect.poll(() => questionText(page), { timeout: POLL_TIMEOUT }).not.toBe(before);
}

async function teamPages(pages: Page[], withPage: Page): Promise<Page[]> {
  const wantA = await isTeamA(withPage);
  const result: Page[] = [];
  for (const page of pages) {
    if ((await isTeamA(page)) === wantA) result.push(page);
  }
  return result;
}

async function startFourPlayerGame(pages: Page[], questionsPerPlayer?: number): Promise<void> {
  const [host, ...guests] = pages;
  await createGame(host);
  const gameId = await getLobbyGameId(host);
  await Promise.all(guests.map((page) => joinExistingGame(page, gameId)));
  await Promise.all(pages.map((page, i) => joinLobby(page, NAMES[i])));

  await host.getByRole('button', { name: /Mix Teams/i }).click();
  await host.waitForTimeout(1_000);

  if (questionsPerPlayer) {
    await host.locator('#qpp-select').selectOption(String(questionsPerPlayer));
    await Promise.all(pages.map((page) => expect(page.locator('#qpp-select')).toHaveValue(String(questionsPerPlayer))));
  }

  const startButton = host.getByRole('button', { name: /Rack Cups & Start/i });
  await expect(startButton).toBeEnabled({ timeout: 5_000 });
  await startButton.click();

  await Promise.all(pages.map((page) => expect(page.locator('.game-board')).toBeVisible({ timeout: 10_000 })));
}

test.describe('Skip and questions-per-player (2v2)', () => {
  let contexts: BrowserContext[];
  let pages: Page[];

  test.beforeEach(async ({ browser }) => {
    contexts = await Promise.all(NAMES.map(() => browser.newContext()));
    pages = await Promise.all(contexts.map((context) => context.newPage()));
  });

  test.afterEach(async () => {
    await Promise.all(contexts.map((context) => context.close()));
  });

  test('each player answers N questions before their team wins', async () => {
    await startFourPlayerGame(pages, 2);
    const team = await teamPages(pages, pages[0]);
    expect(team).toHaveLength(2);

    // 2 players x 2 questions = 4 correct answers; no win before the last.
    for (let answered = 1; answered <= 4; answered++) {
      const active = await waitForActive(team);
      if (answered < 4) {
        await answerAndWait(active);
        await expect(active.locator('.game-over')).toHaveCount(0);
      } else {
        await answerIfMyTurn(active);
      }
    }

    await Promise.all(pages.map((page) => expect(page.locator('.game-over')).toBeVisible({ timeout: 10_000 })));
  });

  test('skipping notifies all four players and gives the skipper a new question', async () => {
    await startFourPlayerGame(pages);

    const skipper = await waitForActive(pages);
    const before = await questionText(skipper);
    await skipper.getByRole('button', { name: 'Skip' }).click();

    await Promise.all(
      pages.map((page) => expect(page.getByTestId('skip-banner')).toContainText('is stuck and skipped!')),
    );
    await expect.poll(() => questionText(skipper), { timeout: POLL_TIMEOUT }).not.toBe(before);
    await expect(skipper.locator('.question-card')).toBeVisible();
  });

  test('a skip delays the teammate\'s next question and shows the penalty to the team', async () => {
    await startFourPlayerGame(pages, 2);
    const team = await teamPages(pages, pages[0]);

    const skipper = await waitForActive(team);
    const teammate = team.find((page) => page !== skipper)!;
    const before = await questionText(skipper);
    await skipper.getByRole('button', { name: 'Skip' }).click();
    await expect.poll(() => questionText(skipper), { timeout: POLL_TIMEOUT }).not.toBe(before);
    await answerIfMyTurn(skipper);

    await expect(teammate.getByTestId('penalty-wait')).toBeVisible();
    await expect(teammate.locator('.question-card')).toHaveCount(0);

    // After the penalty the teammate is asked, and the penalty message clears for them.
    await expect(teammate.locator('.question-card')).toBeVisible({ timeout: 8_000 });
    await expect(skipper.getByTestId('penalty-wait')).toHaveCount(0);
  });

  test('a skip on the team\'s final question delays the win', async () => {
    await startFourPlayerGame(pages);
    const team = await teamPages(pages, pages[0]);

    // Last player on the team skips and then answers: the win, not a question, follows the penalty.
    const first = await waitForActive(team);
    await answerAndWait(first);
    const last = await waitForActive(team);
    const other = team.find((page) => page !== last)!;

    const before = await questionText(last);
    await last.getByRole('button', { name: 'Skip' }).click();
    await expect.poll(() => questionText(last), { timeout: POLL_TIMEOUT }).not.toBe(before);
    await answerIfMyTurn(last);

    await expect(other.getByTestId('penalty-wait')).toBeVisible();
    await expect(last.locator('.game-over')).toHaveCount(0);
    await Promise.all(pages.map((page) => expect(page.locator('.game-over')).toBeVisible({ timeout: 8_000 })));
  });
});
