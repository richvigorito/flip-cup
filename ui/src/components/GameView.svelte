<script lang="ts">
  import { send } from '$lib/transport/socket';
  import { mode, currentQuestion, gameState, myTeam, me, winner, gameId, currentPlayerName } from '$lib/store';

  let currentAnswer = '';

  $: tableColor = $myTeam?.color ?? '#333';

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
    <!-- Quit button adapted for new UI -->
    <button class="quit-btn-floating" on:click={quitGame}>Quit Game</button>

    <!-- Question card (only shown to the active player) -->
    {#if $me?.isMyTurn && $currentQuestion && !$winner}
      <div class="question-card">
        <div class="question-meta">
          <span class="your-turn-badge">✦ Your Turn</span>
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
          <button
            class="submit-btn"
            on:click={submitAnswer}
            disabled={!currentAnswer.trim()}
          >
            Submit
          </button>
        </div>
      </div>

    {:else if !$winner}
      <!-- Waiting indicator -->
      <div class="waiting-card">
        <div class="waiting-dot"></div>
        <span>Waiting for another player's turn…</span>
      </div>
    {/if}

    <!-- Game Board -->
    <div class="game-board">
      <div class="board-label">Game Board</div>
      <div class="board-teams">

        <!-- Team A -->
        <div class="team-col">
          <div class="team-name-row">
            <span class="team-badge a">A</span>
            <span class="team-name-txt">{$gameState.teamA.name}</span>
          </div>
          <div class="cups-list">
            {#each $gameState.teamA.players as player, i}
              {@const flipped  = $gameState.teamA.turn > i}
              {@const isCurrent = $gameState.teamA.turn === i}
              <div class="cup-row" class:current={isCurrent}>
                <div class="cup-wrapper">
                  <div
                    class="cup"
                    class:flipped
                    class:current={isCurrent}
                    title={player.name}
                  ></div>
                  {#if isCurrent}
                    <div class="cup-glow"></div>
                  {/if}
                </div>
                <span
                  class="cup-player-name"
                  class:active={isCurrent}
                  class:done={flipped}
                >
                  {player.name}
                  {#if flipped}<span class="done-icon">✓</span>{/if}
                </span>
              </div>
            {/each}
          </div>
        </div>

        <!-- Divider -->
        <div class="vs-col">
          <div class="vs-bar"></div>
          <span class="vs-text">VS</span>
          <div class="vs-bar"></div>
        </div>

        <!-- Team B -->
        <div class="team-col">
          <div class="team-name-row right">
            <span class="team-name-txt">{$gameState.teamB.name}</span>
            <span class="team-badge b">B</span>
          </div>
          <div class="cups-list">
            {#each $gameState.teamB.players as player, i}
              {@const flipped   = $gameState.teamB.turn > i}
              {@const isCurrent = $gameState.teamB.turn === i}
              <div class="cup-row right" class:current={isCurrent}>
                <span
                  class="cup-player-name right"
                  class:active={isCurrent}
                  class:done={flipped}
                >
                  {#if flipped}<span class="done-icon">✓</span>{/if}
                  {player.name}
                </span>
                <div class="cup-wrapper">
                  <div
                    class="cup"
                    class:flipped
                    class:current={isCurrent}
                    title={player.name}
                  ></div>
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

<!-- Game Over Overlay -->
{#if $winner}
  <div class="game-over">
    <div class="game-over-card">
      <div class="trophy">🏆</div>
      <p class="game-over-label">Winner!</p>
      <h2 class="game-over-winner">{$winner}</h2>
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
      margin-top: -30px; /* Pull it up a bit */
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
      color: var(--error);
      border-color: var(--error);
      background: rgba(239, 68, 68, 0.1);
  }

  /* ── Question card ── */
  .question-card {
    background: var(--bg-card);
    border: 1px solid var(--accent-border);
    border-radius: var(--r-xl);
    padding: 1.5rem 1.75rem;
    box-shadow: 0 0 0 1px var(--accent-border), 0 4px 20px rgba(124,58,237,0.15);
    animation: slideDown 0.3s var(--ease);
  }
  @keyframes slideDown {
    from { opacity:0; transform: translateY(-8px); }
    to   { opacity:1; transform: translateY(0); }
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
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    border-radius: var(--r-md);
    border: none;
    white-space: nowrap;
    box-shadow: 0 3px 10px rgba(124,58,237,0.35);
    transition: all 0.2s var(--ease);
  }
  .submit-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 5px 15px rgba(124,58,237,0.5);
  }
  .submit-btn:disabled { opacity: 0.4; cursor: not-allowed; }

  /* ── Waiting indicator ── */
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
    50%       { opacity: 0.5; transform: scale(0.8); }
  }

  /* ── Game Board ── */
  .game-board {
    background: #e5e5e5; /* Cheap plastic table color */
    border: 8px solid #d4d4d4; /* Table edge */
    border-radius: var(--r-lg);
    padding: 1.5rem 1.25rem 1.75rem;
    box-shadow: 
      inset 0 0 40px rgba(0,0,0,0.05), /* Subtle texture/stain */
      0 10px 30px rgba(0,0,0,0.3); /* Table shadow */
    position: relative;
  }
  
  /* Wood grain / Garage texture hint */
  .game-board::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image: repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(0,0,0,0.02) 10px, rgba(0,0,0,0.02) 20px);
    pointer-events: none;
    opacity: 0.4;
    border-radius: var(--r-lg);
  }

  .board-label {
    font-size: 0.8rem;
    font-weight: 800;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: #555; /* Darker label for contrast on light table */
    text-align: center;
    margin-bottom: 1.25rem;
    text-shadow: 0 1px 0 rgba(255,255,255,0.8);
  }
  .board-teams {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
  }

  /* ── VS column ── */
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
    color: var(--text-muted);
  }
  .vs-bar {
    width: 1px;
    height: 30px;
    background: var(--border);
  }

  /* ── Team columns ── */
  .team-col {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    min-width: 0;
    padding: 1rem 0.5rem;
    border-radius: var(--r-md);
  }

  /* Distinct backgrounds for teams */
  .team-col:first-child { /* Team A */
    background: rgba(239, 68, 68, 0.15); /* Red tint - increased opacity */
    border: 2px solid rgba(239, 68, 68, 0.3);
  }
  .team-col:last-child { /* Team B */
    background: rgba(59, 130, 246, 0.15); /* Blue tint - increased opacity */
    border: 2px solid rgba(59, 130, 246, 0.3);
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
  .team-badge.a { background: rgba(124,58,237,0.25); color: #a78bfa; }
  .team-badge.b { background: rgba(56,189,248,0.2);  color: #38bdf8; }

  .team-name-txt {
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.04em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  /* ── Cups ── */
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
  .cup-row.current { background: var(--accent-dim); }

  .cup-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .cup {
    width: 44px;
    height: 44px;
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
    filter: drop-shadow(0 0 6px rgba(124,58,237,0.8));
  }

  .cup-glow {
    position: absolute;
    inset: -4px;
    border-radius: 8px;
    background: radial-gradient(circle, rgba(124,58,237,0.3) 0%, transparent 70%);
    pointer-events: none;
    animation: glow-pulse 1.6s ease-in-out infinite;
  }
  @keyframes glow-pulse {
    0%, 100% { opacity: 0.7; }
    50%       { opacity: 0.2; }
  }

  .cup-player-name {
    font-size: 0.78rem;
    font-weight: 500;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
    min-width: 0;
    transition: color 0.2s;
  }
  .cup-player-name.right { text-align: right; }
  .cup-player-name.active { color: var(--accent); font-weight: 700; }
  .cup-player-name.done   { color: var(--success); opacity: 0.75; }

  .done-icon {
    font-size: 0.7rem;
    margin-right: 0.2rem;
    color: var(--success);
  }
  .cup-player-name.right .done-icon { margin-right: 0; margin-left: 0.2rem; }

  /* ── Game Over ── */
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
    to   { opacity: 1; }
  }

  .game-over-card {
    text-align: center;
    background: var(--bg-card);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-xl);
    padding: 3rem 4rem;
    box-shadow: var(--shadow-lg), 0 0 80px rgba(124,58,237,0.2);
    animation: popIn 0.4s var(--ease);
  }
  @keyframes popIn {
    from { opacity: 0; transform: scale(0.88); }
    to   { opacity: 1; transform: scale(1); }
  }

  .trophy { font-size: 4rem; margin-bottom: 0.75rem; }

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
    background: linear-gradient(135deg, #a78bfa, #818cf8, #38bdf8);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: 2rem;
  }

  .restart-btn {
    padding: 0.8rem 2rem;
    font-size: 0.9375rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    border: none;
    border-radius: var(--r-md);
    box-shadow: 0 4px 16px rgba(124,58,237,0.4);
    transition: all 0.2s var(--ease);
  }
  .restart-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 22px rgba(124,58,237,0.55);
  }
</style>