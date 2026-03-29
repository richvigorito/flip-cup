import { test, expect } from '@playwright/test';

test.describe('Welcome screen', () => {
  test('renders logo, tagline, and both CTAs', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('heading', { name: /FlipCup/i })).toBeVisible();
    await expect(page.locator('.hero-tagline')).toContainText('red cups');
    await expect(page.getByRole('button', { name: /Create New Game/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /Join Existing Game/i })).toBeVisible();
  });

  test('navigates to Create Game screen', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Create New Game/i }).click();

    await expect(page.getByRole('heading', { name: /Create a Game/i })).toBeVisible();
    await expect(page.locator('#qs-select')).toBeVisible();
  });

  test('navigates to Join Game screen', async ({ page }) => {
    await page.goto('/');
    await page.getByRole('button', { name: /Join Existing Game/i }).click();

    await expect(page.getByRole('heading', { name: /Join a Game/i })).toBeVisible();
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

    await expect(page.getByText('How to Play FlipCup')).toBeVisible();
    await page.getByRole('button', { name: /Close/i }).click();
    await expect(page.getByText('How to Play FlipCup')).not.toBeVisible();
  });
});
