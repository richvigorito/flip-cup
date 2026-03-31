# Testing guide

FlipCup uses a layered test strategy:

- Go tests for backend behavior
- frontend production build as a sanity check
- Playwright for end-to-end multiplayer flows
- GitHub Actions to run the same validation automatically

## Backend tests

Run:

```bash
cd game-server
go test ./...
```

This is the main safety net for:

- game state logic
- cleanup behavior
- API behavior
- backend regressions that do not need a browser

## Frontend build check

Run:

```bash
cd ui
npm ci
npm run build
```

This is the fastest way to catch:

- broken imports
- type/build issues
- Svelte/Vite bundling problems

## End-to-end tests

Run:

```bash
cd ui
npx playwright install --with-deps chromium
npm run test:e2e
```

Additional helpers:

```bash
cd ui
npm run test:e2e:ui
npm run test:e2e:report
```

The Playwright suite is the best safety net for:

- multiplayer flows
- reconnect behavior
- screen-level regressions
- "does the app actually work from the browser?" validation

Key test locations:

- `ui/e2e/game.spec.ts`
- `ui/e2e/disconnect.spec.ts`
- `ui/e2e/game-over-scenarios.spec.ts`
- screenshot/spec helpers under `ui/e2e/`

## CI

The main CI workflow is:

- `.github/workflows/ci.yml`

The reusable validation workflow is:

- `.github/workflows/validate.yml`

Those workflows run:

1. Go tests
2. UI dependency install
3. UI build
4. Playwright browser install
5. Playwright end-to-end tests

That same validation pipeline is reused by staging deploys so deployment does not skip the normal checks.

## What to run before different kinds of changes

### Backend-only change

Usually enough:

```bash
cd game-server
go test ./...
```

### Frontend-only change

Usually enough:

```bash
cd ui
npm run build
```

### Reconnect, gameplay, or cross-stack change

Run both:

```bash
cd game-server && go test ./...
cd ui && npm run build
cd ui && npm run test:e2e
```

## Screenshots and visual artifacts

The repo also keeps screenshot artifacts under `docs/screenshots/`. Those are more documentation/review assets than formal tests, but they are useful when you want to show UI evolution over time.
