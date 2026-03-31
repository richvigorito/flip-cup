# FlipCup

FlipCup is a real-time multiplayer trivia game built with a Go backend and a Svelte frontend. Players join a shared lobby, get split into teams, and race through quiz questions over WebSockets until one side wins.

The repo now has enough moving parts that the README is best treated as a guided index: it should help you understand what is here, where to look next, and how the pieces fit together without forcing every detail into one page.

## What is in this repo?

At a high level, FlipCup is made of five main parts:

- `game-server/` — the Go backend that owns HTTP routes, WebSocket sessions, game state, question loading, and stale-game cleanup
- `ui/` — the Svelte frontend that renders the lobby/game screens and talks to the backend over HTTP + WebSockets
- `ui/e2e/` — Playwright coverage for the multiplayer flows
- `deploy/` — staging deployment assets, currently centered on a Nomad job definition for the homelab Pi cluster
- `docs/` — project guides for architecture, gameplay, development, testing, deployment, and AI-assisted work

Quick repo map:

```text
.
├── game-server/                # Go service, websocket/game logic, quiz loading
├── ui/                         # Svelte client
├── deploy/nomad/              # Staging Nomad job definition
├── docs/                      # Project guides and supporting docs
├── docker-compose.yml         # Local full-stack dev entry point
├── Dockerfile                 # Shared deployment image for staging / Fly
├── fly.toml                   # Fly.io production config
└── gameflow.md                # Message-level websocket/game flow reference
```

## Documentation map

If you did not write most of this code, start here:

- [`docs/architecture.md`](docs/architecture.md) — codebase tour and component responsibilities
- [`docs/gameplay.md`](docs/gameplay.md) — the game lifecycle from lobby to winner, with links to protocol details
- [`docs/development.md`](docs/development.md) — local setup, common workflows, and where to start debugging
- [`docs/testing.md`](docs/testing.md) — backend, frontend, Playwright, and CI coverage
- [`docs/deployment.md`](docs/deployment.md) — local Docker Compose, staging on the Pi/Nomad stack, and Fly production
- [`docs/ai-approach.md`](docs/ai-approach.md) — how Copilot was used across implementation, testing, docs, and repo operations
- [`gameflow.md`](gameflow.md) — lower-level WebSocket message flow reference

## Components at a glance

### Game server

The backend is a Go service under `game-server/` using Gorilla Mux and Gorilla WebSocket. It serves the API, upgrades `/ws` connections, manages game state in memory, loads quiz YAML, and handles cleanup of stale games.

Start with:

- `game-server/cmd/flipcup/`
- `game-server/internal/game/`
- `game-server/internal/transport/`

### UI

The frontend lives in `ui/` and is built with Svelte + Vite. It renders the welcome flow, lobby, active game, reconnect state, and endgame screens while deriving the runtime HTTP/WS host from either `VITE_WS_URL` or the current browser host.

Start with:

- `ui/src/components/`
- `ui/src/lib/`
- `ui/src/lib/utils/config.ts`

### Tests

There are three practical layers of validation:

- Go unit/integration tests under `game-server/`
- frontend production build validation in `ui/`
- Playwright end-to-end tests in `ui/e2e/`

GitHub Actions runs the same test pipeline on pushes and pull requests. See [`docs/testing.md`](docs/testing.md) for the suite breakdown.

### Deployments and environments

FlipCup currently has three important environments:

- local development via `docker-compose.yml`
- staging on a Raspberry Pi homelab using a self-hosted GitHub runner, Nomad, Consul, Traefik, and Vault-backed runtime config
- production on Fly.io via `fly.toml`

The repo-root `Dockerfile` is the shared deployment image for staging and Fly. Local Compose does **not** use that file; it builds the backend from `game-server/` and the frontend from `ui/`, both in dev mode.

See [`docs/deployment.md`](docs/deployment.md) for the environment-by-environment story.

## Quick start

### Local full stack with Docker Compose

```bash
docker-compose up -d
```

That starts:

- the Go server on `http://localhost:8080`
- the Svelte dev server on `http://localhost:5173`

### Run the apps directly

Backend:

```bash
cd game-server
go run cmd/flipcup/main.go
```

Frontend:

```bash
cd ui
npm install
npm run dev
```

## Testing commands

Backend tests:

```bash
cd game-server
go test ./...
```

Frontend build:

```bash
cd ui
npm ci
npm run build
```

Playwright:

```bash
cd ui
npx playwright install --with-deps chromium
npm run test:e2e
```

More detail on what these actually cover lives in [`docs/testing.md`](docs/testing.md).

## AI-assisted development

This project includes a meaningful amount of GitHub Copilot-assisted work: UI redesign, reconnect handling, stale-game cleanup, Playwright coverage, GitHub Actions wiring, staging deployment setup, and documentation updates.

The detailed write-up is in [`docs/ai-approach.md`](docs/ai-approach.md).

## Contributing / navigating

If you are exploring the codebase for the first time:

1. read [`docs/architecture.md`](docs/architecture.md)
2. skim [`docs/gameplay.md`](docs/gameplay.md)
3. use [`docs/development.md`](docs/development.md) for local setup
4. use [`docs/testing.md`](docs/testing.md) before changing behavior
5. use [`docs/deployment.md`](docs/deployment.md) before touching staging or Fly

And yes, for the record: no, this project did not require being in a frat.
