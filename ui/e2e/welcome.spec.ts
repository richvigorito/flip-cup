import { test, expect } from '@playwright/test';

test.describe('Welcome screen', () => {
  test('renders logo, tagline, and both CTAs', async ({ page }) => {
    await page.goto('/');

    await expect(page.locator('.logo-text')).toBeVisible();
    await expect(page.locator('.hero-tagline')).toContainText('Answer trivia');

    await expect(page.getByRole('button', { name: /Create New Game/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /Join Existing Game/i })).toBeVisible();
  });

  test('navigates to Create Game screen', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Create New Game/i }).click();

    await expect(page.locator('.card-title')).toHaveText('Create a Game');
    await expect(page.locator('#qs-select')).toBeVisible();
  });

  test('navigates to Join Game screen', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Join Existing Game/i }).click();

    await expect(page.locator('.card-title')).toHaveText('Join a Game');
  });

  test('back button returns to welcome', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Create New Game/i }).click();
    await page.getByRole('button', { name: /← Back/i }).click();

    await expect(page.getByRole('button', { name: /Create New Game/i })).toBeVisible();
  });

  test('logo click always returns to welcome', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Create New Game/i }).click();
    await page.locator('.logo-link').click();

    await expect(page.getByRole('button', { name: /Create New Game/i })).toBeVisible();
  });

  test('"Create Game" button is disabled until a quiz is selected', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Create New Game/i }).click();

    const submit = page.locator('.submit-btn');
    await expect(submit).toBeDisabled();

    await page.locator('#qs-select').selectOption({ index: 1 });
    await expect(submit).toBeEnabled();
  });

  test('How to Play dialog opens and closes', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /How to Play/i }).click();

    await expect(page.getByText('How to Play FlipQuiz')).toBeVisible();
    await page.getByRole('button', { name: /Got it/i }).click();
    await expect(page.getByText('How to Play FlipQuiz')).not.toBeVisible();
  });
});
