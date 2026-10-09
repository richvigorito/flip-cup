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

### Questions per player
In the lobby the host sends `{"type": "set_questions_per_player", "payload": {"count": 2}}` (clamped to 1–5, default 1, ignored once the game has started). The server broadcasts `questions_per_player_updated` with the same payload, and `game_snapshot.questionsPerPlayer` carries the value. Players take turns in rotation; a team wins after `players × count` correct answers.

### Skipping
The active player may send `{"type": "skip_question"}` (ignored if it is not their turn, the game is inactive, or their team is serving a penalty).

**Server** logic:
- Broadcasts `{"type": "question_skipped", "name": "Alice", "payload": {"penalty_seconds": 7, "total_penalty_seconds": 14}}` to everyone.
- Sends the same player a different question. The turn does not advance.
- After that player answers correctly, the team's next question is delayed by 7s × number of skips that turn. The server broadcasts `{"type": "penalty_wait", "payload": {"team": "A-Team", "seconds": 14}}` and sends the next `question` when the wait ends. If the answer was the team's last, the `winner` message is delayed by the same penalty (unless the other team wins first). Answers and skips from that team are ignored meanwhile.
- The penalty is `FLIPCUP_SKIP_PENALTY` (default `7s`); it is a server setting, not configurable per game.

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
| `skip_question` | Inbound | - | Skip current question (penalty applies) |
| `set_questions_per_player` | Inbound | `{count}` | Lobby only; 1–5 |
| `join_existing_game` | Inbound | `{game_id, player_id}` | Reconnect/Join specific game |
| `question` | Outbound | `{name}` | Question prompt |
| `question_skipped` | Outbound | `{penalty_seconds, total_penalty_seconds}` + `name` | Someone skipped |
| `penalty_wait` | Outbound | `{team, seconds}` | Team waits before next question |
| `questions_per_player_updated` | Outbound | `{count}` | Lobby setting changed |
| `winner` | Outbound | `{name}` | Winning team name |
| `my_current_team` | Outbound | `{players, ...}` | Team state update |
