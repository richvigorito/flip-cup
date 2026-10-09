# Spec: N questions per player + skip with penalty

Status: approved

## Goals

1. A game is more than one question per player when few people play.
2. A stuck player can skip, at a cost to their team's tempo, and everyone sees it.

## Non-goals (do not build)

- Configurable penalty length, skip limits, or per-player skip tracking across turns
- Persisting settings between games, new quiz formats, scoring changes
- Handling reconnect mid-penalty beyond not crashing (the timer keeps running server-side)

## Behavior

### A. Questions per player (N)

- The host sets N in the lobby. Range 1-5, default 1. Locked once the game starts.
- Each player answers N questions in rotation: P1, P2, ..., Pk, P1, ... until every player has answered N.
- A team wins when it has correctly answered `players x N` questions. First team to finish wins (existing race rule).
- If `players x N` plus skips exceeds the quiz size, questions wrap around to the start of the shuffled list.
- Server clamps N into 1-5; invalid input is ignored.

### B. Skip

- Only the current player (the one being asked) can skip, only during an active game, and not while a penalty wait is running.
- Skipping gives that player a different question (the team's next unused question). The skipped question is not scored.
- Skips within one turn stack: each adds `SkipPenalty` (7s, server constant) to a pending penalty.
- After the player answers correctly, the pending penalty is served: the team waits `7s x skips` before the next player is asked. Then the pending penalty resets to 0.
- If the player skipped zero times, there is no wait (unchanged from today).
- The other team is never delayed.
- Everyone (both teams) is notified at skip time.

## Protocol changes (docs/gameflow.md must be updated)

Inbound:

| Type | Payload | Notes |
| --- | --- | --- |
| `set_questions_per_player` | `{count}` | lobby only, game inactive |
| `skip_question` | - | current player, active game |

Outbound:

| Type | Payload | Sent to |
| --- | --- | --- |
| `question_skipped` | `{name, penalty_seconds, total_penalty_seconds}` | all players |
| `penalty_wait` | `{seconds}` | all players, sent when the wait begins |
| `questions_per_player_updated` | `{count}` | all players |
| `question` | existing | the skipper's replacement question, and the next player after the wait |

`GameSnapshot` gains `questionsPerPlayer`. `TeamSnapshot.turn` keeps meaning "answers completed by this team" (the UI uses it to compare team progress); the current player is `turn % len(players)`.

Skip announcement text (shown in the event log and as a banner): `"<Name> is stuck and skipped! +7s penalty"`.

## Acceptance criteria

Unit (Go):

1. N=1 behaves exactly as today (regression).
2. N=2 with 2 players per team: players alternate P1,P2,P1,P2 and the team wins after 4 correct answers.
3. N is clamped to 1-5 and cannot change after the game starts.
4. Skip sends a different question to the skipper and broadcasts `question_skipped` to all players.
5. Two skips in one turn then a correct answer waits 14s before the next question; one skip waits 7s; none waits 0.
6. The wait uses an injectable clock/delay so tests do not sleep 7 real seconds.
7. Skip is rejected for a non-current player, an inactive game, and during a penalty wait.
8. Restarting the game during a penalty wait cancels the pending question.
9. Question index wraps instead of panicking when skips exhaust the quiz.

E2E (Playwright):

1. Host sets N=2, 2v2 game: a team needs 4 correct answers to win.
2. A player skips: all four clients see the skip notification; the skipper gets a new prompt.
3. After skip + correct answer, the next player does not receive a question until the penalty passes (test uses a shortened penalty through the same server constant override used for tests, e.g. env `FLIPCUP_SKIP_PENALTY`).

## Risks

- `Team.Turn` is shared between player index and question index today; separating them touches `NextQuestion`, `handleCheckAnswer`, and the UI turn comparison. Cover with the regression test in Unit 1.
- Skip penalty timer must not leak across restarts or fire into an ended game (Unit 8).
