# Architecture guide

This is the high-level codebase tour for FlipCup. If you did not write most of the project, this should give you the "what lives where?" map before you dive into implementation details.

## The big picture

FlipCup is a real-time multiplayer trivia app with:

- a Go backend for HTTP, WebSockets, game state, and quiz loading
- a Svelte frontend for the lobby/game UI
- Playwright coverage for multiplayer browser flows
- a shared deployment image for staging and Fly

The app currently keeps game state in memory, which explains a lot of the deployment choices:

- staging runs as a single Nomad allocation
- Fly should stay conservative about scaling
- reconnect behavior matters because deploys and browser refreshes can interrupt live sessions

## Codebase layout

### `game-server/`

This is the backend application.

Important areas:

- `cmd/flipcup/` — the executable entrypoint
- `internal/game/` — core game models and lifecycle logic
- `internal/quiz/` — quiz loading and domain logic
- `internal/transport/api/` — HTTP endpoints
- `internal/transport/ws/` — WebSocket handling
- `questions/` — YAML quiz content
- `public/` — built frontend assets used by the single-container deploy image

The backend is responsible for:

- creating and joining games
- assigning teams
- tracking active player/question flow
- broadcasting events over WebSockets
- handling reconnects and stale cleanup

### `ui/`

This is the frontend application built with Svelte + Vite.

Important areas:

- `src/components/` — major screens and UI components
- `src/lib/` — stores, helpers, and shared utilities
- `src/lib/utils/config.ts` — runtime host/URL resolution for HTTP and WebSocket traffic
- `e2e/` — Playwright tests

The frontend is responsible for:

- creating/joining games
- rendering the lobby and team state
- rendering the active game view
- reconnecting players to the correct game context
- surfacing winner/endgame state

### `deploy/`

This currently exists for staging.

- `deploy/nomad/flipcup.nomad.hcl` defines the staging job for the homelab Pi cluster

That job:

- runs a single service allocation
- uses a task-level Vault integration
- renders env vars from Vault into the app task
- registers the service for Traefik/Consul routing

### `docs/`

This directory holds the human-oriented guides:

- architecture
- gameplay
- development
- testing
- deployment
- AI-assisted development notes

### Root-level infrastructure files

- `docker-compose.yml` — local full-stack dev
- `Dockerfile` — shared deployment image for staging/Fly
- `fly.toml` — Fly.io production configuration
- `gameflow.md` — protocol-level message flow reference

## How requests and events move through the app

At a high level:

1. the browser hits the UI
2. the UI opens a WebSocket to `/ws`
3. the backend creates/loads in-memory game state
4. gameplay events are broadcast back to connected players
5. the UI updates from that state stream

For the detailed event reference, read [`../gameflow.md`](../gameflow.md).

## Good starting points by task

If you are fixing or adding something:

- lobby/game lifecycle bug — `game-server/internal/game/`
- API/WebSocket bug — `game-server/internal/transport/`
- UI/state issue — `ui/src/components/`, `ui/src/lib/`
- reconnect behavior — `ui/src/lib/`, `ui/e2e/disconnect.spec.ts`, backend transport/game files
- stale cleanup or timers — `game-server/internal/game/`
- staging deploy issue — `deploy/nomad/flipcup.nomad.hcl`, `.github/workflows/deploy-staging.yml`, [`deployment.md`](deployment.md)

## Constraints worth keeping in mind

- game state is in memory, so multiple replicas are not a free win
- WebSocket stability matters more than aggressive autoscaling
- local development is intentionally simple and should stay simple
- staging is intentionally close to production packaging, but runs inside the homelab

If you want the practical "how do I run this and work on it?" guide next, go to [`development.md`](development.md).
