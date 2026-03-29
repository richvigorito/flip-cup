# SKILLS.md

Practical task guide for AI assistants working on `FlipCup`.

## Skill: understand the app quickly

Read in this order:

1. `README.md`
2. `gameflow.md`
3. `docs/ai-approach.md`
4. The specific files involved in the task

What you should learn fast:

- how players create and join a game
- when teams are assigned
- how turns and questions progress
- where websocket messages are handled
- how the real Flip Cup mechanic works when it affects UI, copy, or animation: drink from an upright cup, set it upside down on the table edge, then pop-flip it upright

## Skill: run the project locally

### Full stack via Docker

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Backend only

```bash
cd game-server
go run cmd/flipcup/main.go
```

### Frontend only

```bash
cd ui
npm install
npm run dev
```

The frontend uses `VITE_WS_URL`, and local development usually points it to `localhost:8080`.

## Skill: make a backend change safely

Checklist:

- inspect the relevant package under `game-server/internal/`
- understand whether the change affects game state, transport, or quiz loading
- check for websocket payload implications
- run:

```bash
cd game-server
go test ./...
```

Extra caution areas:

- nil dereferences
- turn index bounds
- player joins/leaves
- game cleanup and restart flow

## Skill: make a frontend change safely

Checklist:

- inspect the target Svelte component and related utilities
- confirm whether the change is visual, state-related, or websocket-related
- keep existing component structure and styling conventions where possible
- for flip-cup-themed visuals or copy, preserve the real drinking-game sequence instead of generic spinning cup motion
- validate with:

```bash
cd ui
npm run build
```

## Skill: work on end-to-end tests

Tests live in `ui/e2e/`.

Useful commands:

```bash
cd ui
npm run test:e2e
npm run test:e2e:ui
npm run test:e2e:report
```

Notes:

- some tests require the backend and frontend to be running
- multiplayer tests may use separate browser contexts to simulate different players

## Skill: trace the real-time flow

When debugging a gameplay issue, follow this path:

1. UI action in `ui/src`
2. socket or transport helper in the frontend
3. websocket/API handler in `game-server/internal/transport/`
4. game state logic in `game-server/internal/game/`
5. outbound state/question/winner events back to the UI

## Skill: write good changes here

Prefer:

- small diffs
- explicit error handling
- keeping contracts aligned across UI and server
- updating docs/tests when behavior changes

Avoid:

- speculative abstractions
- changing both stacks when one stack is enough
- editing unrelated files during a focused task
