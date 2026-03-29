<script lang="ts">
  import { mode, joined, currentPlayerName, gameState, gamesCompleted, gameId } from '$lib/store';
  import { send } from '$lib/transport/socket';
  import QuestionSetDropdown from './QuestionSetDropdown.svelte';

  let localPlayerName = '';

  const assignTeams = () => send({ type: 'assign_teams' });

  const startGame = () => {
    if ($gameState?.canStart()) {
      send({ type: 'start' });
      mode.set('game');
    } else {
      alert('Teams must have at least 1 player each before starting.');
    }
  };

  const joinGame = (playerName: string) => {
    joined.set(true);
    send({ type: 'add_player', payload: { name: playerName } });
  };

  function handleNewQuiz(selectedFile: string) {
    send({ type: 'update_quiz', payload: { quizfile: selectedFile } });
  }

  const leaveGame = () => {
    if (confirm('Are you sure you want to leave the game? This will clear your session.')) {
      if (typeof sessionStorage !== 'undefined') {
        sessionStorage.removeItem('flipcup_player_id');
        sessionStorage.removeItem('flipcup_game_id');
      }
      window.location.reload();
    }
  };
</script>

<div class="lobby-wrap">
  <div class="lobby-card">
    <div class="lobby-header">
      <img src="/solo-cup.png" alt="" class="lobby-icon" />
      {#if $gameState}
        <div class="game-code">
          <span class="game-code-label">Game ID</span>
          <span class="game-code-value">{$gameState.id}</span>
        </div>
      {:else if $gameId}
        <div class="game-code">
          <span class="game-code-label">Game ID</span>
          <span class="game-code-value">{$gameId}</span>
        </div>
      {/if}
      <button class="leave-btn-icon" on:click={leaveGame} title="Leave table">Leave</button>
    </div>

    {#if !$joined && $mode === 'lobby'}
      <div class="step">
        <h2 class="step-title">What's your name?</h2>
        <p class="step-sub">Claim your cup and get into the lineup.</p>

        <div class="name-form">
          <input
            type="text"
            placeholder="Enter your name…"
            bind:value={localPlayerName}
            on:keydown={(e) => {
              if (e.key === 'Enter' && localPlayerName.trim()) {
                currentPlayerName.set(localPlayerName.trim());
                joinGame(localPlayerName.trim());
              }
            }}
          />
          <button
            class="join-btn"
            disabled={!localPlayerName.trim()}
            on:click={() => {
              const name = localPlayerName.trim();
              currentPlayerName.set(name);
              joinGame(name);
            }}
          >
            Join Game
          </button>
        </div>
      </div>
    {:else}
      <div class="step">
        <h2 class="step-title">Lobby</h2>
        <p class="step-sub">
          {#if $currentPlayerName}
            You're in as <strong>{$currentPlayerName}</strong>. Get both sides set, then let the cups fly.
          {:else}
            Waiting for players to step up to the table…
          {/if}
        </p>

        {#if $gameState}
          <div class="team-preview">
            <div class="team-slot">
              <span class="team-badge a">A</span>
              <div class="team-slot-info">
                <span class="team-slot-name">{$gameState.teamA.name}</span>
                <span class="team-slot-players">
                  {$gameState.teamA.players.length > 0
                    ? $gameState.teamA.players.map((p) => p.name).join(', ')
                    : 'No players yet'}
                </span>
              </div>
              <span class="player-count">{$gameState.teamA.players.length}</span>
            </div>

            <div class="team-divider">table</div>

            <div class="team-slot">
              <span class="team-badge b">B</span>
              <div class="team-slot-info">
                <span class="team-slot-name">{$gameState.teamB.name}</span>
                <span class="team-slot-players">
                  {$gameState.teamB.players.length > 0
                    ? $gameState.teamB.players.map((p) => p.name).join(', ')
                    : 'No players yet'}
                </span>
              </div>
              <span class="player-count">{$gameState.teamB.players.length}</span>
            </div>
          </div>
        {/if}

        <div class="lobby-actions">
          <button class="action-btn secondary" on:click={assignTeams}>
            Mix Teams
          </button>
          <button class="action-btn primary" on:click={startGame}>
            Rack Cups & Start
          </button>
        </div>

        {#if $gamesCompleted > 0}
          <div class="change-quiz">
            <label class="quiz-label" for="qs-select">Change Question Set</label>
            <QuestionSetDropdown on:select={(e) => handleNewQuiz(e.detail)} />
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .leave-btn-icon {
    padding: 0.45rem 0.8rem;
    background: var(--bg-surface);
    border: 1px solid var(--border);
    color: var(--text-secondary);
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    cursor: pointer;
    margin-left: 1rem;
    transition: all 0.2s;
  }

  .leave-btn-icon:hover {
    transform: translateY(-1px);
    border-color: var(--danger-border);
    color: #fecaca;
    background: rgba(127, 29, 29, 0.25);
  }

  .lobby-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 64px);
    padding: 2rem 1.5rem;
  }

  .lobby-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--r-xl);
    padding: 2.5rem 2rem;
    width: 100%;
    max-width: 480px;
    box-shadow: var(--shadow-card);
  }

  .lobby-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 2rem;
    padding-bottom: 1.25rem;
    border-bottom: 1px solid var(--border);
  }

  .lobby-icon {
    width: 40px;
    height: auto;
    filter: drop-shadow(0 8px 16px rgba(220, 38, 38, 0.22));
  }

  .game-code {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 0.125rem;
  }

  .game-code-label {
    font-size: 0.7rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--text-muted);
  }

  .game-code-value {
    font-size: 1rem;
    font-weight: 800;
    font-variant-numeric: tabular-nums;
    letter-spacing: -0.01em;
    color: var(--accent);
  }

  .step-title {
    font-size: 1.5rem;
    font-weight: 800;
    letter-spacing: -0.03em;
    margin-bottom: 0.375rem;
  }

  .step-sub {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1.5rem;
  }

  .step-sub strong { color: var(--text-primary); }

  .name-form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .name-form input {
    width: 100%;
    padding: 0.75rem 1rem;
    font-size: 1rem;
    background: var(--bg-surface);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-md);
    color: var(--text-primary);
    outline: none;
    transition: border-color 0.15s, box-shadow 0.15s;
  }

  .name-form input:focus {
    border-color: var(--accent-border);
    box-shadow: 0 0 0 3px var(--accent-dim);
  }

  .join-btn {
    padding: 0.8rem 1.5rem;
    font-size: 0.9375rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
    color: #fff;
    border: none;
    border-radius: var(--r-md);
    box-shadow: 0 4px 14px rgba(220, 38, 38, 0.28);
    transition: all 0.2s var(--ease);
  }

  .join-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(220, 38, 38, 0.38);
  }

  .join-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .team-preview {
    background: var(--bg-subtle);
    border: 1px solid var(--border);
    border-radius: var(--r-lg);
    padding: 0.75rem 1rem;
    margin-bottom: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .team-slot {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 0;
  }

  .team-badge {
    width: 28px;
    height: 28px;
    border-radius: var(--r-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 800;
    flex-shrink: 0;
  }

  .team-badge.a { background: rgba(220, 38, 38, 0.24); color: #fca5a5; }
  .team-badge.b { background: rgba(245, 158, 11, 0.22); color: #fcd34d; }

  .team-slot-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    min-width: 0;
  }

  .team-slot-name {
    font-size: 0.825rem;
    font-weight: 700;
    color: var(--text-primary);
  }

  .team-slot-players {
    font-size: 0.75rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .player-count {
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--text-muted);
    background: var(--bg-elevated);
    border-radius: var(--r-full);
    min-width: 22px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .team-divider {
    font-size: 0.7rem;
    font-weight: 800;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-muted);
    text-align: center;
    padding: 0.25rem 0;
    border-top: 1px solid var(--border);
    border-bottom: 1px solid var(--border);
  }

  .lobby-actions {
    display: flex;
    gap: 0.75rem;
  }

  .action-btn {
    flex: 1;
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
    font-weight: 700;
    border-radius: var(--r-md);
    transition: all 0.2s var(--ease);
    cursor: pointer;
    border: none;
  }

  .action-btn.primary {
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
    color: #fff;
    box-shadow: 0 4px 12px rgba(220, 38, 38, 0.28);
  }

  .action-btn.primary:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 18px rgba(220, 38, 38, 0.38);
  }

  .action-btn.secondary {
    background: var(--bg-surface);
    color: var(--text-primary);
    border: 1px solid var(--border-strong);
  }

  .action-btn.secondary:hover {
    background: var(--bg-elevated);
    border-color: var(--border-strong);
  }

  .change-quiz {
    margin-top: 1.5rem;
    padding-top: 1.25rem;
    border-top: 1px solid var(--border);
  }

  .quiz-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-muted);
    margin-bottom: 0.5rem;
  }
</style>
