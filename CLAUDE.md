# CLAUDE.md

This is the canonical AI-assistant instruction file for this repository.

Use `SKILLS.md` as the supplemental task and command reference.

## Start here

1. Read `README.md`
2. For gameplay or socket work, read `gameflow.md`
3. Read `SKILLS.md` if you need task-specific commands or checklists
4. Inspect only the files needed for the task before editing

## What this repo is

`FlipCup` is a real-time multiplayer quiz game:

- Go backend in `game-server/`
- Svelte frontend in `ui/`
- WebSocket-driven gameplay
- Playwright E2E coverage in `ui/e2e/`

## Claude-specific expectations

- Be concise in explanations, but complete in implementation.
- Prefer exact edits over speculative refactors.
- Reuse existing patterns before adding helpers or abstractions.
- If a task touches gameplay, reason through multiplayer side effects.
- If a task touches UI-only styling, avoid backend churn.
- If a task references the real-world Flip Cup game, keep the mechanic accurate: players drink from an upright cup, place it upside down on the table edge, then fingertip-pop it so it lands upright.

## Commands you will likely need

### Full stack

```bash
docker-compose up -d
```

### Backend

```bash
cd game-server
go run cmd/flipcup/main.go
go test ./...
```

### Frontend

```bash
cd ui
npm install
npm run dev
npm run build
```

### E2E

```bash
cd ui
npm run test:e2e
```

## Good entry points by task

- Bug in lobby/game lifecycle: `game-server/internal/game/`, `game-server/internal/transport/ws/`
- API or websocket handler bug: `game-server/internal/transport/`
- UI state or screens: `ui/src/components/`, `ui/src/lib/`
- Test fixes: `ui/e2e/`

## Avoid

- Large-scale renames unless required
- Silent behavior changes
- Duplicating logic or instructions unnecessarily across files

## Definition of done

- The requested change is implemented
- Relevant validation commands were run when code changed
- Any directly affected docs or tests were updated
