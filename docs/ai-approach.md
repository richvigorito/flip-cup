# AI-Assisted Development Approach

This document describes how AI tooling (GitHub Copilot) was used to improve this project — what was delegated, how decisions were made, and what the AI actually changed.

---

## Overview

Two separate workstreams were handed to AI agents:

| Branch | What changed |
|--------|-------------|
| `feat/ui-redesign` | Full visual redesign of the Svelte frontend |
| `feat/playwright-e2e-tests` | End-to-end test suite using Playwright |

Both were produced in a single session using [GitHub Copilot CLI](https://githubnext.com) with a `dev-problem-coordinator` agent that broke work into tasks and delegated to specialised sub-agents.

---

## Step 1 — Codebase Audit

Before any changes, the AI was asked to:

1. Explain how to get the project running
2. Identify what could be improved

The agent explored all source files, read `gameflow.md`, inspected Docker config, WebSocket message types, and Go game logic. It produced a ranked list of issues across five areas:

| Severity | Category | Examples found |
|----------|----------|----------------|
| 🔴 Critical | Server crashes | `must.go` panics on bad JSON; nil deref in `handler.go` |
| 🔴 Critical | Race conditions | Turn state; concurrent shuffle |
| 🟠 High | Reliability | No reconnection logic; dead games never cleaned up |
| 🟡 Medium | UX gaps | No mobile CSS; no answer feedback; awkward quiz picker |
| 🟢 Low | Features | No leaderboard; no spectator mode |

This audit became the backlog for subsequent work.

---

## Step 2 — UI Redesign (`feat/ui-redesign`)

### Why

The original UI had no design system, inconsistent spacing, a large banner image header, and zero mobile responsiveness. The request was blunt: *"the ui on it is trash, would like to look waaay more cleaner"*.

### What the AI did

The `frontend-dev-expert` sub-agent was given full autonomy over the `ui/src` directory. It made no changes to the Go backend.

**Design system (`app.css`)**
- Introduced CSS custom properties for the entire token set: colours, radii, shadows, motion easing, font scale
- Colour palette: deep navy base (`#0b0b18`), vivid purple/indigo brand accent, semantic green/red/amber for game states
- Typography: Inter (Google Fonts), weight 400–900, tight letter-spacing on headings

**Layout changes**
- Replaced the old banner image header with a frosted-glass fixed header containing a gradient logo
- `EventLog` sidebar properly constrained below the header via `top: 64px`

**Component-by-component changes**

| Component | Before | After |
|-----------|--------|-------|
| `Welcome.svelte` | Plain buttons on white | Hero layout, floating animated cup emoji, gradient wordmark, pill tags |
| `NewGame.svelte` | Unstyled form | Card with custom styled `<select>`, disabled submit until selection |
| `JoinGame.svelte` | Flat list | Responsive card grid, team badges, hover glow, empty state |
| `Lobby.svelte` | Raw inputs | Two-step flow (name → team preview), A/B columns, colour-coded badges |
| `GameView.svelte` | Basic layout | Two-column VS board, glowing active cups, smooth flip animation |
| `EventLog.svelte` | Unstyled list | Slide-in entries, colour-coded by type (info/success/error) |
| `Instructions.svelte` | Inline text | Native `<dialog>` with blur backdrop, structured sections |

### What the AI did NOT change

- No significant game logic (only perspective fixes and store updates)
- No dependencies were added (purely CSS/markup changes within existing Svelte components)
- `Playwright` was added for testing, but no production dependencies.

### Decisions made autonomously

The AI chose the colour palette, component layout patterns, and animation approach without being asked. The only human constraint was "cleaner."

---

## Step 3 — E2E Tests (`feat/playwright-e2e-tests`)

### Why

There were zero tests in the project. Any change to game logic or UI had no safety net.

### What the AI did

Installed Playwright (`@playwright/test`) into the `ui/` package and wrote **12 tests** across two files.

**`welcome.spec.ts`** — 7 tests, no server required
- Logo, tagline, and CTAs render correctly
- Navigation to Create / Join screens
- Back button and logo return to Welcome
- Submit button disabled until a quiz category is selected
- How to Play dialog opens and closes

**`game.spec.ts`** — 5 tests, requires game server running
- **Full 2-player game**: two browser contexts (Player 1, Player 2) create → join → enter names → shuffle teams → start game → answer questions turn-by-turn using known Q&A pairs → verify winner declared
- **Play Again**: winner screen → lobby reset
- Start Game disabled with only 1 player
- Game ID displayed prominently
- Join button disabled until name is typed

**`answers.ts`** — known question/answer pairs from `_.default.yaml` used to automate actual gameplay.

### How the multiplayer test works

Playwright's `browser.newContext()` creates isolated browser sessions. Two contexts simulate two independent players connecting over WebSocket:

```
Context 1 (Alice) ──── creates game ──────────────── gets gameId
Context 2 (Bob)   ──── joins game by gameId ─────────────────────
                                          ↓
                                 both enter names
                                 Alice shuffles teams
                                 Alice starts game
                                          ↓
                        loop: check which page has "Your Turn"
                              answer question on that page
                              until .game-over appears
```

### Test runner

```bash
cd ui

# Run all tests (requires: docker-compose up -d first)
npm run test:e2e

# Interactive mode with browser UI
npm run test:e2e:ui

# View last run report
npm run test:e2e:report
```

---

## Reflections

### What worked well

- **Codebase audit first** — having the AI map the whole project before touching anything meant changes were targeted and didn't introduce regressions into areas it didn't understand
- **Separation of concerns** — UI work and test work were kept on separate branches, making review and rollback clean
- **Specificity of prompts** — vague asks ("make it cleaner") worked fine for aesthetic work; functional asks ("write tests that actually play the game") needed more structure

### What to watch

- The AI redesigned components based on class names and markup it could see — if the underlying component structure changes significantly, styles may need revisiting
- The game tests depend on the `_.default.yaml` question file staying intact; if that file is removed or renamed the full-game tests will need updating
- Several critical server bugs (panics, nil dereferences, race conditions) were identified in the audit but not fixed — those are tracked separately as the next priority

---

## Next Steps (from the audit)

These were identified but not yet addressed:

1. Replace `panic()` calls in `must.go` with graceful error returns
2. Fix nil dereference in `handler.go` (check `g != nil` *before* calling methods on it)
3. Fix uninitialized player name (`name` var never populated before `AddPlayer`)
4. Add bounds checking before `g.QuestionFile.Questions[t.Turn]`
5. Add WebSocket reconnection logic in `socket.ts`
6. Call `DeleteGameByID()` on game completion to prevent memory leak
7. Add `@media` queries for mobile responsiveness
