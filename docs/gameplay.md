# Gameplay guide

This document explains what the player flow is supposed to feel like and where that behavior lives in the codebase.

If you need the lower-level message contract, use [`../gameflow.md`](../gameflow.md). This page is the friendlier overview.

## Player journey

The normal flow is:

1. a player creates or joins a game
2. players gather in the lobby
3. the host assigns teams
4. the host starts the game
5. the active player answers the current question
6. the server advances the game until one team wins
7. players can restart into a new round

## Lobby phase

The lobby is where most session setup happens:

- players join and pick names
- the host waits for enough players
- teams are assigned before the game starts

If something feels off here, look at:

- backend: `game-server/internal/game/`
- transport: `game-server/internal/transport/ws/`
- frontend: `ui/src/components/Lobby.svelte`

## Active game phase

Once the match starts:

- the server tracks whose turn it is
- the current player gets the question prompt
- correct answers advance the team
- incorrect answers broadcast feedback and preserve turn/game rules

This is mostly coordinated through:

- `game-server/internal/game/`
- `game-server/internal/transport/ws/`
- `ui/src/components/GameView*.svelte`

## Reconnect and recovery

Because this is a multiplayer WebSocket app, reconnect behavior matters.

The current design tries to preserve enough client identity that a page refresh can reconnect the player to the right game instead of dumping them into a blank start screen.

Useful places:

- `ui/src/lib/`
- `ui/e2e/disconnect.spec.ts`
- backend game/session handling under `game-server/internal/`

## Endgame

When a team wins:

- the backend emits a winner event
- the UI renders the game-over state
- players can start another round with a restart action

Useful places:

- `game-server/internal/game/`
- `ui/src/components/`
- `ui/e2e/game-over-scenarios.spec.ts`

## Stale cleanup

FlipCup also has cleanup logic for inactive in-memory games.

That matters operationally because:

- the server keeps state in memory
- dead games should not accumulate forever
- staging/ops workflows may want to trigger or rely on cleanup behavior

Useful places:

- `game-server/internal/game/cleanup.go`
- related tests under `game-server/internal/game/`

## If you need the protocol details

Use [`../gameflow.md`](../gameflow.md) when you need:

- exact message names
- inbound vs outbound message direction
- payload structure examples
- the WebSocket event sequence
