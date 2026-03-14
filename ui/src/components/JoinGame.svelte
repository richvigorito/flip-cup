<script lang="ts">
  import { onMount } from 'svelte';
  import { connectSocket } from '$lib/transport/socket';
  import { fetchGames } from '$lib/transport/http/Games';
  import type { GameState } from '$lib/models/GameState';
  import { mode, gameId, currentPlayerName } from '$lib/store';

  let availableGames: GameState[] = [];
  let loading = true;

  function joinGame(id: string) {
    gameId.set(id);
    currentPlayerName.set('');
    connectSocket({
      type: 'join_existing_game',
      payload: { game_id: id }
    });
    mode.set('lobby');
  }

  const goBack = () => mode.set('welcome');

  onMount(async () => {
    availableGames = await fetchGames();
    loading = false;
  });
</script>

<div class="page-wrap">

  <button class="back-btn" on:click={goBack}>← Back</button>

  <div class="join-container">

    <div class="page-header">
      <div class="page-icon">🎮</div>
      <h2 class="page-title">Join a Game</h2>
      <p class="page-subtitle">Select an open game below to jump in</p>
    </div>

    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Scanning for games…</p>
      </div>

    {:else if availableGames.length === 0}
      <div class="empty-state">
        <div class="empty-icon">🕵️</div>
        <p class="empty-title">No open games found</p>
        <p class="empty-sub">Be the first — create one!</p>
        <button class="create-btn" on:click={() => mode.set('new')}>
          Create New Game →
        </button>
      </div>

    {:else}
      <div class="game-grid">
        {#each availableGames as game}
          <button class="game-card" on:click={() => joinGame(game.id)}>
            <div class="game-card-id">#{game.id}</div>
            <div class="teams-preview">
              <div class="team-preview">
                <span class="team-badge team-a">A</span>
                <div class="team-info">
                  <span class="team-label">{game.teamA.name}</span>
                  <span class="team-players">
                    {#if game.teamA.players.length > 0}
                      {game.teamA.players.map(p => p.name).join(', ')}
                    {:else}
                      No players yet
                    {/if}
                  </span>
                </div>
              </div>

              <div class="team-vs">vs</div>

              <div class="team-preview">
                <span class="team-badge team-b">B</span>
                <div class="team-info">
                  <span class="team-label">{game.teamB.name}</span>
                  <span class="team-players">
                    {#if game.teamB.players.length > 0}
                      {game.teamB.players.map(p => p.name).join(', ')}
                    {:else}
                      No players yet
                    {/if}
                  </span>
                </div>
              </div>
            </div>
            <div class="join-hint">Tap to join →</div>
          </button>
        {/each}
      </div>
    {/if}

  </div>
</div>

<style>
  .page-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    min-height: calc(100vh - 64px);
    padding: 2rem 1.5rem;
    gap: 1.25rem;
    width: 100%;
    max-width: 900px;
    margin: 0 auto;
  }

  .back-btn {
    align-self: flex-start;
    color: var(--text-muted);
    font-size: 0.875rem;
    font-weight: 600;
    padding: 0.25rem 0;
    background: none;
    border: none;
    transition: color 0.15s;
    cursor: pointer;
  }
  .back-btn:hover { color: var(--text-primary); }

  .join-container {
    width: 100%;
  }

  /* Header */
  .page-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  .page-icon { font-size: 2.5rem; margin-bottom: 0.75rem; }
  .page-title {
    font-size: 1.75rem;
    font-weight: 800;
    letter-spacing: -0.035em;
    margin-bottom: 0.375rem;
  }
  .page-subtitle {
    font-size: 0.9rem;
    color: var(--text-secondary);
  }

  /* Loading */
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 3rem;
    color: var(--text-muted);
    font-size: 0.9rem;
  }
  .spinner {
    width: 28px;
    height: 28px;
    border: 2px solid var(--border);
    border-top-color: var(--accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  /* Empty state */
  .empty-state {
    text-align: center;
    padding: 3rem 2rem;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--r-xl);
  }
  .empty-icon { font-size: 3rem; margin-bottom: 1rem; }
  .empty-title {
    font-size: 1.1rem;
    font-weight: 700;
    margin-bottom: 0.375rem;
  }
  .empty-sub {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1.5rem;
  }
  .create-btn {
    display: inline-flex;
    padding: 0.7rem 1.5rem;
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    font-weight: 700;
    border-radius: var(--r-md);
    box-shadow: 0 4px 14px rgba(124,58,237,0.35);
    transition: all 0.2s;
    border: none;
    cursor: pointer;
  }
  .create-btn:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(124,58,237,0.5); }

  /* Game grid */
  .game-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 0.875rem;
  }

  /* Game card */
  .game-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--r-lg);
    padding: 1.25rem 1.25rem 1rem;
    text-align: left;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    transition: all 0.2s var(--ease);
    cursor: pointer;
    width: 100%;
  }
  .game-card:hover {
    border-color: var(--accent-border);
    background: var(--bg-surface);
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0,0,0,0.3), 0 0 0 1px var(--accent-border);
  }

  .game-card-id {
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: var(--accent);
    font-variant-numeric: tabular-nums;
  }

  .teams-preview {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .team-preview {
    display: flex;
    align-items: center;
    gap: 0.625rem;
  }

  .team-badge {
    width: 22px;
    height: 22px;
    border-radius: var(--r-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.7rem;
    font-weight: 800;
    flex-shrink: 0;
  }
  .team-a { background: rgba(124, 58, 237, 0.25); color: #a78bfa; }
  .team-b { background: rgba(56, 189, 248, 0.2);  color: #38bdf8; }

  .team-info {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
    min-width: 0;
  }
  .team-label {
    font-size: 0.8rem;
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .team-players {
    font-size: 0.75rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .team-vs {
    font-size: 0.7rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--text-muted);
    padding-left: 0.25rem;
  }

  .join-hint {
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--text-muted);
    border-top: 1px solid var(--border);
    padding-top: 0.625rem;
    transition: color 0.15s;
  }
  .game-card:hover .join-hint { color: var(--accent); }
</style>