import { test, expect, type Browser, type BrowserContext, type Page } from '@playwright/test';
import { answerIfMyTurn, createGame, getLobbyGameId, joinExistingGame, joinLobby } from './helpers/game';

type TeamKey = 'A' | 'B';

type Participant = {
  name: string;
  context: BrowserContext;
  page: Page;
  assignedTeam?: TeamKey;
};

const TEAM_A_NAME = 'A-Team';
const TEAM_B_NAME = 'B-squad';
const TEAM_SIZE = 5;
const LOSING_TEAM_SCENARIOS = [1, 2, 3, 4] as const;

function buildParticipantPlan(): { name: string }[] {
  return Array.from({ length: TEAM_SIZE }, (_, index) => [
    { name: `Alpha ${index + 1}` },
    { name: `Bravo ${index + 1}` },
  ]).flat();
}

async function setupFiveVsFiveGame(browser: Browser): Promise<Participant[]> {
  const plan = buildParticipantPlan();
  const contexts = await Promise.all(plan.map(() => browser.newContext()));
  const pages = await Promise.all(contexts.map((context) => context.newPage()));

  const participants = plan.map((player, index) => ({
    ...player,
    context: contexts[index],
    page: pages[index],
  }));

  await createGame(participants[0].page);
  const gameId = await getLobbyGameId(participants[0].page);
  expect(gameId).toBeTruthy();

  await Promise.all(participants.slice(1).map(({ page }) => joinExistingGame(page, gameId)));
  await Promise.all(participants.map(({ page, name }) => joinLobby(page, name)));

  await participants[0].page.getByRole('button', { name: /Mix Teams/i }).click();
  await participants[0].page.waitForTimeout(1_000);

  const startButton = participants[0].page.getByRole('button', { name: /Rack Cups & Start/i });
  await expect(startButton).toBeEnabled({ timeout: 5_000 });
  await startButton.click();

  await Promise.all(
    participants.map(({ page }) => expect(page.locator('.game-board')).toBeVisible({ timeout: 10_000 })),
  );

  await assignRuntimeTeams(participants);

  return participants;
}

async function closeParticipants(participants: Participant[]): Promise<void> {
  await Promise.all(participants.map(({ context }) => context.close()));
}

async function assignRuntimeTeams(participants: Participant[]): Promise<void> {
  for (const participant of participants) {
    const gameBoard = participant.page.locator('.game-board');
    const className = await gameBoard.getAttribute('class');

    if (className?.includes('team-a-view')) {
      participant.assignedTeam = 'A';
      continue;
    }

    if (className?.includes('team-b-view')) {
      participant.assignedTeam = 'B';
      continue;
    }

    throw new Error(`Could not resolve team assignment for ${participant.name}`);
  }
}

function getTeamParticipants(participants: Participant[], team: TeamKey): Participant[] {
  return participants.filter((participant) => participant.assignedTeam === team);
}

async function playScenario(
  winningTeamPlayers: Participant[],
  losingTeamPlayers: Participant[],
  losingTeamCorrectAnswers: number,
): Promise<void> {
  const allParticipants = [...winningTeamPlayers, ...losingTeamPlayers];
  let losingTeamCorrectCount = 0;

  for (let tick = 0; tick < 160; tick++) {
    if (await isGameOver(allParticipants)) {
      return;
    }

    let progressed = false;

    for (const participant of winningTeamPlayers) {
      if (await answerIfMyTurn(participant.page)) {
        progressed = true;
        break;
      }
    }

    if (losingTeamCorrectCount < losingTeamCorrectAnswers) {
      for (const participant of losingTeamPlayers) {
        if (await answerIfMyTurn(participant.page)) {
          losingTeamCorrectCount += 1;
          progressed = true;
          break;
        }
      }
    }

    if (!progressed) {
      await allParticipants[0].page.waitForTimeout(200);
    }
  }

  throw new Error(`Game did not finish while validating losing team with ${losingTeamCorrectAnswers} correct answers`);
}

async function isGameOver(participants: Participant[]): Promise<boolean> {
  for (const { page } of participants) {
    if (await page.locator('.game-over').isVisible()) {
      return true;
    }
  }

  return false;
}

function teamSection(page: Page, teamName: string) {
  return page.locator('.game-over-team').filter({
    has: page.locator('.game-over-team-name', { hasText: teamName }),
  });
}

async function expectTeamCupState(page: Page, teamName: string, flipped: number, filled: number, names: string[]): Promise<void> {
  const section = teamSection(page, teamName);
  await expect(section).toHaveCount(1);
  await expect(section.locator('.game-over-cup.flipped')).toHaveCount(flipped);
  await expect(section.locator('.game-over-cup.filled')).toHaveCount(filled);

  for (const name of names) {
    await expect(section.getByText(name, { exact: true })).toBeVisible();
  }
}

test.describe('Game over cup scenarios', () => {
  test.setTimeout(120_000);

  for (const losingTeamCorrectAnswers of LOSING_TEAM_SCENARIOS) {
    test(`5v5 game shows ${losingTeamCorrectAnswers} losing cups flipped when Team A wins`, async ({ browser }) => {
      const participants = await setupFiveVsFiveGame(browser);

      try {
        const teamAPlayers = getTeamParticipants(participants, 'A');
        const teamBPlayers = getTeamParticipants(participants, 'B');
        expect(teamAPlayers).toHaveLength(TEAM_SIZE);
        expect(teamBPlayers).toHaveLength(TEAM_SIZE);

        await playScenario(teamAPlayers, teamBPlayers, losingTeamCorrectAnswers);

        const teamAPage = teamAPlayers[0].page;
        const teamBPage = teamBPlayers[0].page;
        const teamAPlayerNames = teamAPlayers.map(({ name }) => name);
        const teamBPlayerNames = teamBPlayers.map(({ name }) => name);

        await expect(teamAPage.locator('.game-over')).toBeVisible();
        await expect(teamBPage.locator('.game-over')).toBeVisible();
        await expect(teamAPage.locator('.game-over-winner')).toHaveText('You ran the table!');
        await expect(teamBPage.locator('.game-over-winner')).toHaveText(`${TEAM_A_NAME} ran the table`);

        for (const page of [teamAPage, teamBPage]) {
          await expect(page.locator('.game-over-cup-name')).toHaveCount(TEAM_SIZE * 2);
          await expectTeamCupState(page, TEAM_A_NAME, TEAM_SIZE, 0, teamAPlayerNames);
          await expectTeamCupState(
            page,
            TEAM_B_NAME,
            losingTeamCorrectAnswers,
            TEAM_SIZE - losingTeamCorrectAnswers,
            teamBPlayerNames,
          );
        }
      } finally {
        await closeParticipants(participants);
      }
    });
  }
});
