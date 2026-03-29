<script lang="ts">
  import { send } from '$lib/transport/socket';
  import { mode, currentQuestion, gameState, myTeam, me, winner } from '$lib/store';
  import type { Team } from '$lib/models/Team';

  let currentAnswer = '';

  const getClearedCupCount = (team: Team, winnerName: string | null) => {
    if (winnerName && team.name === winnerName) {
      return team.players.length;
    }

    return Math.min(team.turn, team.players.length);
  };

  const submitAnswer = () => {
    if (!currentAnswer.trim()) return;
    send({ type: 'check_answer', payload: { answer: currentAnswer } });
    currentAnswer = '';
  };

  const resetGame = () => send({ type: 'restart_game' });

  const quitGame = () => {
    if (confirm('Are you sure you want to quit the game? This will clear your session.')) {
      if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem('flipcup_player_id');
        sessionStorage.removeItem('flipcup_game_id');
      }
      window.location.reload();
    }
  };

  const handleKey = (e: KeyboardEvent) => {
    if (e.key === 'Enter') submitAnswer();
  };
</script>

{#if $mode === 'game'}
  <div class="game-wrap">
    <button class="quit-btn-floating" on:click={quitGame}>Leave Table</button>

    {#if $me?.isMyTurn && $currentQuestion && !$winner}
      <div class="question-card">
        <div class="question-meta">
          <span class="your-turn-badge">Your Cup</span>
        </div>
        <p class="question-text">{$currentQuestion}</p>
        <div class="answer-row">
          <input
            class="answer-input"
            type="text"
            bind:value={currentAnswer}
            placeholder="Type your answer…"
            on:keydown={handleKey}
          />
          <button class="submit-btn" on:click={submitAnswer} disabled={!currentAnswer.trim()}>
            Flip It
          </button>
        </div>
      </div>
    {:else if !$winner}
      <div class="waiting-card">
        <div class="waiting-dot"></div>
        <span>Waiting for the next player to step up…</span>
      </div>
    {/if}

    <div
      class="game-board"
      class:team-a-view={$myTeam?.name === $gameState.teamA.name}
      class:team-b-view={$myTeam?.name === $gameState.teamB.name}
    >
      <div class="board-label">
        {#if $myTeam?.name === $gameState.teamA.name}
          Table View (Team A Side)
        {:else if $myTeam?.name === $gameState.teamB.name}
          Table View (Team B Side)
        {:else}
          Table View (Spectator Side)
        {/if}
      </div>

      <div class="board-teams">
        <div class="team-col">
          <div class="team-name-row">
            <span class="team-badge a">A</span>
            <span class="team-name-txt">{$gameState.teamA.name}</span>
          </div>
          <div class="cups-list">
            {#each $gameState.teamA.players as player, i}
              {@const flipped = $gameState.teamA.turn > i}
              {@const isCurrent = $gameState.teamA.turn === i}
              <div class="cup-row" class:current={isCurrent}>
                <div class="cup-wrapper">
                  <div class="cup" class:flipped class:current={isCurrent} title={player.name}></div>
                  {#if isCurrent}
                    <div class="cup-glow"></div>
                  {/if}
                </div>
                <span class="cup-player-name" class:active={isCurrent} class:done={flipped}>
                  {player.name}
                  {#if flipped}<span class="done-icon">✓</span>{/if}
                </span>
              </div>
            {/each}
          </div>
        </div>

        <div class="vs-col">
          <div class="vs-bar"></div>
          <span class="vs-text">table</span>
          <div class="vs-bar"></div>
        </div>

        <div class="team-col">
          <div class="team-name-row right">
            <span class="team-name-txt">{$gameState.teamB.name}</span>
            <span class="team-badge b">B</span>
          </div>
          <div class="cups-list">
            {#each $gameState.teamB.players as player, i}
              {@const flipped = $gameState.teamB.turn > i}
              {@const isCurrent = $gameState.teamB.turn === i}
              <div class="cup-row right" class:current={isCurrent}>
                <span class="cup-player-name right" class:active={isCurrent} class:done={flipped}>
                  {#if flipped}<span class="done-icon">✓</span>{/if}
                  {player.name}
                </span>
                <div class="cup-wrapper">
                  <div class="cup" class:flipped class:current={isCurrent} title={player.name}></div>
                  {#if isCurrent}
                    <div class="cup-glow"></div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}

{#if $winner}
  <div class="game-over">
    <div class="game-over-card">
      {#if $myTeam && $winner === $myTeam.name}
        <p class="game-over-label">Table cleared</p>
        <h2 class="game-over-winner">You ran the table!</h2>
      {:else if $myTeam}
        <p class="game-over-label">Next round</p>
        <h2 class="game-over-winner">{$winner} ran the table</h2>
      {:else}
        <p class="game-over-label">Winner</p>
        <h2 class="game-over-winner">{$winner} cleared the table</h2>
      {/if}

      {#if $gameState}
        <div class="game-over-table">
          {#each [$gameState.teamA, $gameState.teamB] as team, teamIndex}
            {@const clearedCupCount = getClearedCupCount(team, $winner)}
            <section
              class="game-over-team"
              class:winner-team={team.name === $winner}
              class:my-team={$myTeam?.name === team.name}
            >
              <div class="game-over-team-header">
                <span class="game-over-team-name">{team.name}</span>
                <span class="game-over-team-status">
                  {team.name === $winner ? 'Table cleared' : `${team.turn} of ${team.players.length} cups cleared`}
                </span>
              </div>

              <div class="game-over-cups" class:right={teamIndex === 1}>
                {#each team.players as player, index}
                  {@const flipped = index < clearedCupCount}
                  <div class="game-over-cup-slot">
                    <div class="cup game-over-cup" class:flipped class:filled={!flipped} title={player.name}></div>
                    <span class="game-over-cup-name">{player.name}</span>
                  </div>
                {/each}
              </div>
            </section>
          {/each}
        </div>
      {/if}

      <button class="restart-btn" on:click={resetGame}>
        Play Again
      </button>
    </div>
  </div>
{/if}

<style>
  .game-wrap {
    max-width: 660px;
    margin: 0 auto;
    padding: 1.5rem 1.25rem 3rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    position: relative;
  }

  .quit-btn-floating {
    position: absolute;
    top: 0;
    right: 0;
    margin-top: -30px;
    padding: 0.4rem 0.8rem;
    background: transparent;
    color: var(--text-muted);
    border: 1px solid var(--border);
    border-radius: var(--r-md);
    font-size: 0.7rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .quit-btn-floating:hover {
    color: #fecaca;
    border-color: var(--danger-border);
    background: rgba(239, 68, 68, 0.1);
  }

  .question-card {
    background: var(--bg-card);
    border: 1px solid var(--accent-border);
    border-radius: var(--r-xl);
    padding: 1.5rem 1.75rem;
    box-shadow: 0 0 0 1px var(--accent-border), 0 4px 20px rgba(220, 38, 38, 0.14);
    animation: slideDown 0.3s var(--ease);
  }

  @keyframes slideDown {
    from { opacity: 0; transform: translateY(-8px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .question-meta {
    margin-bottom: 0.75rem;
  }

  .your-turn-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.72rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--accent);
    background: var(--accent-dim);
    padding: 0.25rem 0.625rem;
    border-radius: var(--r-full);
  }

  .question-text {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.55;
    margin-bottom: 1.25rem;
  }

  .answer-row {
    display: flex;
    gap: 0.625rem;
  }

  .answer-input {
    flex: 1;
    padding: 0.65rem 0.875rem;
    background: var(--bg-surface);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-md);
    color: var(--text-primary);
    font-size: 0.9375rem;
    font-weight: 500;
    outline: none;
    transition: border-color 0.15s, box-shadow 0.15s;
  }

  .answer-input:focus {
    border-color: var(--accent-border);
    box-shadow: 0 0 0 3px var(--accent-dim);
  }

  .submit-btn {
    padding: 0.65rem 1.25rem;
    font-size: 0.875rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
    color: #fff;
    border-radius: var(--r-md);
    border: none;
    white-space: nowrap;
    box-shadow: 0 3px 10px rgba(220, 38, 38, 0.28);
    transition: all 0.2s var(--ease);
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 5px 15px rgba(220, 38, 38, 0.38);
  }

  .submit-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .waiting-card {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    padding: 0.875rem 1.25rem;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--r-lg);
    font-size: 0.875rem;
    color: var(--text-muted);
  }

  .waiting-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--warning);
    animation: pulse 1.4s ease-in-out infinite;
    flex-shrink: 0;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; transform: scale(1); }
    50% { opacity: 0.5; transform: scale(0.8); }
  }

  .game-board {
    background:
      linear-gradient(180deg, rgba(84, 45, 32, 0.96), rgba(58, 28, 20, 0.98)),
      #4a281d;
    border: 8px solid #6b3428;
    border-radius: var(--r-lg);
    padding: 1.5rem 1.25rem 1.75rem;
    box-shadow:
      inset 0 0 40px rgba(0,0,0,0.16),
      0 10px 30px rgba(0,0,0,0.3);
    position: relative;
    transition: background-color 0.5s ease;
  }

  .game-board.team-a-view {
    background:
      linear-gradient(180deg, rgba(112, 39, 30, 0.96), rgba(73, 24, 19, 0.98)),
      #5f251d;
    border-color: #8b3326;
  }

  .game-board.team-b-view {
    background:
      linear-gradient(180deg, rgba(93, 58, 18, 0.96), rgba(64, 39, 13, 0.98)),
      #5f3b12;
    border-color: #8c5a18;
  }

  .game-board::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(90deg, rgba(255,255,255,0.03), transparent 18%, rgba(0,0,0,0.08) 50%, transparent 78%),
      repeating-linear-gradient(0deg, transparent, transparent 17px, rgba(255,255,255,0.02) 17px, rgba(255,255,255,0.02) 18px);
    pointer-events: none;
    opacity: 0.5;
    border-radius: var(--r-lg);
  }

  .board-label {
    font-size: 0.8rem;
    font-weight: 800;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: rgba(255, 237, 213, 0.78);
    text-align: center;
    margin-bottom: 1.25rem;
    text-shadow: 0 1px 0 rgba(0,0,0,0.5);
  }

  .board-teams {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .vs-col {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0 0.25rem;
    padding-top: 40px;
    gap: 0.375rem;
    flex-shrink: 0;
  }

  .vs-text {
    font-size: 0.65rem;
    font-weight: 900;
    letter-spacing: 0.12em;
    color: rgba(255, 237, 213, 0.56);
    text-transform: uppercase;
  }

  .vs-bar {
    width: 1px;
    height: 30px;
    background: rgba(255, 237, 213, 0.14);
  }

  .team-col {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    min-width: 0;
    padding: 1rem 0.5rem;
    border-radius: var(--r-md);
    background: rgba(24, 12, 10, 0.16);
  }

  .team-name-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.375rem;
  }

  .team-name-row.right {
    justify-content: flex-end;
  }

  .team-badge {
    width: 22px;
    height: 22px;
    border-radius: var(--r-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.65rem;
    font-weight: 800;
    flex-shrink: 0;
  }

  .team-badge.a { background: rgba(220, 38, 38, 0.24); color: #fecaca; }
  .team-badge.b { background: rgba(245, 158, 11, 0.22); color: #fde68a; }

  .team-name-txt {
    font-size: 0.8rem;
    font-weight: 700;
    color: rgba(255, 237, 213, 0.9);
    text-transform: uppercase;
    letter-spacing: 0.04em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .cups-list {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .cup-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.375rem 0.5rem;
    border-radius: var(--r-md);
    transition: background 0.2s;
  }

  .cup-row.right { flex-direction: row-reverse; }
  .cup-row.current { background: rgba(248, 113, 113, 0.14); }

  .cup-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .cup {
    width: 44px;
    height: 44px;
    position: relative;
    overflow: hidden;
    background-image: url('/solo-cup.png');
    background-size: cover;
    background-repeat: no-repeat;
    background-position: center;
    border-radius: 4px;
    transition: transform 0.4s var(--ease), opacity 0.4s var(--ease), filter 0.4s var(--ease);
  }

  .cup.flipped {
    transform: rotate(180deg);
    opacity: 0.3;
    filter: grayscale(0.6);
  }

  .cup.current {
    filter: drop-shadow(0 0 6px rgba(248, 113, 113, 0.82));
  }

  .cup-glow {
    position: absolute;
    inset: -4px;
    border-radius: 8px;
    background: radial-gradient(circle, rgba(248, 113, 113, 0.3) 0%, transparent 70%);
    pointer-events: none;
    animation: glow-pulse 1.6s ease-in-out infinite;
  }

  @keyframes glow-pulse {
    0%, 100% { opacity: 0.7; }
    50% { opacity: 0.2; }
  }

  .cup-player-name {
    font-size: 0.78rem;
    font-weight: 500;
    color: rgba(255, 237, 213, 0.66);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
    min-width: 0;
    transition: color 0.2s;
  }

  .cup-player-name.right { text-align: right; }
  .cup-player-name.active { color: #fff7ed; font-weight: 700; }
  .cup-player-name.done { color: #bbf7d0; opacity: 0.82; }

  .done-icon {
    font-size: 0.7rem;
    margin-right: 0.2rem;
    color: var(--success);
  }

  .cup-player-name.right .done-icon {
    margin-right: 0;
    margin-left: 0.2rem;
  }

  .game-over {
    position: fixed;
    inset: 0;
    background: rgba(11, 11, 24, 0.9);
    backdrop-filter: blur(14px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 300;
    animation: fadeIn 0.35s var(--ease);
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  .game-over-card {
    text-align: center;
    background: var(--bg-card);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-xl);
    padding: 3rem 4rem;
    box-shadow: var(--shadow-lg), 0 0 80px rgba(220, 38, 38, 0.18);
    animation: popIn 0.4s var(--ease);
  }

  @keyframes popIn {
    from { opacity: 0; transform: scale(0.88); }
    to { opacity: 1; transform: scale(1); }
  }

  .game-over-label {
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--text-muted);
    margin-bottom: 0.5rem;
  }

  .game-over-winner {
    font-size: 2.5rem;
    font-weight: 900;
    letter-spacing: -0.04em;
    background: linear-gradient(135deg, #f87171, #fb923c, #fbbf24);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: 1.5rem;
  }

  .game-over-table {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
    text-align: left;
  }

  .game-over-team {
    padding: 1rem 1.125rem;
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: var(--r-lg);
    background: rgba(19, 24, 39, 0.58);
  }

  .game-over-team.winner-team {
    border-color: rgba(251, 191, 36, 0.32);
    box-shadow: inset 0 0 0 1px rgba(251, 191, 36, 0.08);
  }

  .game-over-team.my-team {
    background: rgba(32, 19, 17, 0.72);
  }

  .game-over-team-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.75rem;
    margin-bottom: 0.9rem;
  }

  .game-over-team-name {
    font-size: 0.9rem;
    font-weight: 800;
    color: #fff7ed;
    letter-spacing: 0.03em;
    text-transform: uppercase;
  }

  .game-over-team-status {
    font-size: 0.78rem;
    color: rgba(255, 237, 213, 0.72);
  }

  .game-over-cups {
    display: flex;
    flex-wrap: wrap;
    gap: 0.9rem;
    justify-content: flex-start;
  }

  .game-over-cups.right {
    justify-content: flex-end;
  }

  .game-over-cup-slot {
    width: 68px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.45rem;
  }

  .game-over-cup {
    width: 52px;
    height: 52px;
    filter: drop-shadow(0 8px 14px rgba(0, 0, 0, 0.28));
  }

  .game-over-cup.filled::after {
    content: '';
    position: absolute;
    left: 10px;
    right: 10px;
    top: 8px;
    height: 12px;
    border-radius: 999px 999px 10px 10px;
    background: linear-gradient(180deg, rgba(250, 204, 21, 0.92), rgba(217, 119, 6, 0.88));
    box-shadow:
      0 1px 0 rgba(255, 255, 255, 0.18) inset,
      0 0 0 1px rgba(120, 53, 15, 0.24);
  }

  .game-over-cup-name {
    font-size: 0.76rem;
    font-weight: 600;
    color: rgba(255, 237, 213, 0.86);
    text-align: center;
    line-height: 1.25;
    word-break: break-word;
  }

  .restart-btn {
    padding: 0.8rem 2rem;
    font-size: 0.9375rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
    color: #fff;
    border: none;
    border-radius: var(--r-md);
    box-shadow: 0 4px 16px rgba(220, 38, 38, 0.3);
    transition: all 0.2s var(--ease);
  }

  .restart-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 22px rgba(220, 38, 38, 0.42);
  }

  @media (max-width: 640px) {
    .game-over-card {
      width: min(100% - 1.5rem, 560px);
      padding: 2rem 1.25rem;
    }

    .game-over-winner {
      font-size: 2rem;
    }

    .game-over-team-header {
      flex-direction: column;
      align-items: flex-start;
      margin-bottom: 0.75rem;
    }

    .game-over-cups.right {
      justify-content: flex-start;
    }
  }
</style>
