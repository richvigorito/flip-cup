# Development guide

This is the practical guide for running FlipCup locally and figuring out where to start when you want to change something.

## Fastest local setup

The quickest full-stack loop is Docker Compose:

```bash
docker-compose up -d
```

That uses:

- `game-server/Dockerfile` with the `dev` target
- `ui/Dockerfile` with the `dev` target

Important: local Compose does **not** use the repo-root `Dockerfile`. That root image exists for deployment packaging, not for your normal inner loop.

## Running the apps directly

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

Useful local URLs:

- backend: `http://localhost:8080`
- frontend: `http://localhost:5173`

## Runtime host / WebSocket behavior

The frontend runtime URL logic lives in:

- `ui/src/lib/utils/config.ts`

Behavior summary:

- local dev usually uses `VITE_WS_URL=localhost:8080`
- deployed environments can leave `VITE_WS_URL` unset
- when unset, the UI falls back to the current browser host and builds HTTP/WS URLs from there

That lets staging/Fly avoid hardcoding a specific deploy hostname into the frontend build.

## Recommended reading order for new contributors

If you are new to the repo:

1. read [`architecture.md`](architecture.md)
2. read [`gameplay.md`](gameplay.md)
3. skim [`../gameflow.md`](../gameflow.md) if you are changing WebSocket behavior
4. use [`testing.md`](testing.md) before changing anything risky

## Common change entry points

### Backend logic

Start in:

- `game-server/internal/game/`
- `game-server/internal/quiz/`

### API or WebSocket transport

Start in:

- `game-server/internal/transport/api/`
- `game-server/internal/transport/ws/`

### Frontend UI or state

Start in:

- `ui/src/components/`
- `ui/src/lib/`

### Deployment work

Start in:

- `deploy/nomad/flipcup.nomad.hcl`
- `.github/workflows/deploy-staging.yml`
- `Dockerfile`
- `fly.toml`

## A few practical cautions

- game state is in memory, so be careful about assumptions around scaling or multiple replicas
- if you touch reconnect logic, validate both frontend state and backend session handling
- if you touch deployment packaging, make sure you are not accidentally changing local Compose behavior
- if you touch gameplay, think through multiplayer side effects instead of only the happy path

## What "local development" really means here

There are two different workflows in this repo:

- **developer workflow** — Docker Compose or direct app commands
- **deployment workflow** — build the shared root `Dockerfile` for staging/Fly

Keeping those separate is intentional. It preserves a fast dev loop while still letting staging and production share one deployable image.
