# Game Flow & WebSocket Protocol

This document outlines the WebSocket message flow between the client (`ui/`) and server (`game-server/`).

## 1. Connection & Setup
**Client** connects to `ws://host/ws`.
**Server** upgrades connection and creates a `Player` instance.

## 2. Join Lobby
**Client** sends:
```json
{
  "type": "add_player",
  "payload": { "name": "Alice" }
}
```

**Server** responds:
1. To **All**: `{"type": "player_joined", "name": "Alice"}`
2. To **Alice**: `{"type": "joined_success", "name": "Alice"}`

## 3. Team Assignment
**Client** (Host) sends:
```json
{ "type": "assign_teams" }
```
*(Only allowed when game is inactive)*

**Server** responds:
1. To **All** (x2): `{"type": "my_current_team", "payload": { ... }}` (updates team view for A and B)
2. To **All**: `{"type": "teams_assigned", "payload": { "game_snapshot": ... }}`

## 4. Start Game
**Client** (Host) sends:
```json
{ "type": "start" }
```

**Server** responds:
1. To **All**: `{"type": "game_started", "payload": { ... }}`
2. To **Active Player**: `{"type": "question", "name": "What is 2+2?"}`

## 5. Gameplay Loop
**Client** (Active Player) sends:
```json
{
  "type": "check_answer",
  "payload": { "answer": "4" }
}
```

**Server** logic:
- **If Correct**:
  - Broadcasts `{"type": "answered_correctly", "payload": ...}`
  - Sends next question to next player: `{"type": "question", "name": "Next Q..."}`
- **If Incorrect**:
  - Broadcasts `{"type": "incorrect_answer", "name": "Alice"}` (Client shows error feedback)

## 6. Game Over
**Server** detects end of quiz for a team.
**Server** broadcasts:
```json
{ "type": "winner", "name": "A-Team" }
```

## 7. Restart
**Client** sends:
```json
{ "type": "restart_game" }
```
**Server** resets state and broadcasts `{"type": "game_restarted", ...}`.

---

### Message Types Summary

| Type | Direction | Payload | Description |
|------|-----------|---------|-------------|
| `add_player` | Inbound | `{name}` | Set player name |
| `assign_teams` | Inbound | - | Shuffle teams |
| `start` | Inbound | - | Start game |
| `check_answer` | Inbound | `{answer}` | Submit answer |
| `join_existing_game` | Inbound | `{game_id, player_id}` | Reconnect/Join specific game |
| `question` | Outbound | `{name}` | Question prompt |
| `winner` | Outbound | `{name}` | Winning team name |
| `my_current_team` | Outbound | `{players, ...}` | Team state update |
