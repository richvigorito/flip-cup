<script lang="ts">
  import { mode, resetClientGameState } from '$lib/store';
  import { connectSocket } from '$lib/transport/socket';
  import { questionSets } from '$lib/store';

  let selectedQuestionSet = '';

  const createGame = () => {
    if (!selectedQuestionSet) return;
    if (typeof sessionStorage !== 'undefined') {
      sessionStorage.removeItem('flipcup_player_id');
      sessionStorage.removeItem('flipcup_game_id');
    }
    resetClientGameState();
    connectSocket({
      type: 'create_game',
      payload: { file: selectedQuestionSet }
    });
    mode.set('lobby');
  };

  const goBack = () => mode.set('welcome');
</script>

<div class="page-wrap">
  <button class="back-btn" on:click={goBack}>
    ← Back
  </button>

  <div class="form-card">
    <div class="card-header">
      <img src="/solo-cup.png" alt="" class="card-icon" />
      <h2 class="card-title">Create a Game</h2>
      <p class="card-subtitle">Pick a category and set the table.</p>
    </div>

    <div class="field">
      <label class="field-label" for="qs-select">Question Set</label>
      <div class="qs-select-wrapper">
        <select id="qs-select" bind:value={selectedQuestionSet}>
          <option value="" disabled>— Choose a category —</option>
          {#each Object.entries(
            $questionSets.reduce((groups, qs) => {
              const cat = qs.category || 'Other';
              (groups[cat] ||= []).push(qs);
              return groups;
            }, {})
          ) as [category, sets]}
            <optgroup label={category}>
              {#each sets as { file, label }}
                <option value={file}>{label}</option>
              {/each}
            </optgroup>
          {/each}
        </select>
      </div>
    </div>

    <button class="submit-btn" on:click={createGame} disabled={!selectedQuestionSet}>
      Create Game →
    </button>
  </div>
</div>

<style>
  .page-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 64px);
    padding: 2rem 1.5rem;
    gap: 1.25rem;
  }

  .back-btn {
    align-self: flex-start;
    max-width: 440px;
    width: 100%;
    color: var(--text-muted);
    font-size: 0.875rem;
    font-weight: 600;
    padding: 0.25rem 0;
    background: none;
    border: none;
    text-align: left;
    transition: color 0.15s;
    cursor: pointer;
  }

  .back-btn:hover { color: var(--text-primary); }

  .form-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--r-xl);
    padding: 2.5rem 2rem;
    width: 100%;
    max-width: 440px;
    box-shadow: var(--shadow-card);
  }

  .card-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .card-icon {
    width: 72px;
    height: auto;
    margin-bottom: 0.75rem;
    filter: drop-shadow(0 10px 18px rgba(220, 38, 38, 0.18));
  }

  .card-title {
    font-size: 1.5rem;
    font-weight: 800;
    letter-spacing: -0.03em;
    margin-bottom: 0.375rem;
  }

  .card-subtitle {
    font-size: 0.875rem;
    color: var(--text-secondary);
  }

  .field {
    margin-bottom: 1.5rem;
  }

  .field-label {
    display: block;
    font-size: 0.8rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-muted);
    margin-bottom: 0.5rem;
  }

  .qs-select-wrapper {
    position: relative;
  }

  .qs-select-wrapper::after {
    content: '▾';
    position: absolute;
    right: 0.875rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-muted);
    pointer-events: none;
    font-size: 0.85rem;
  }

  .qs-select-wrapper select {
    width: 100%;
    padding-right: 2.5rem;
    cursor: pointer;
    background: var(--bg-surface);
    color: var(--text-primary);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-md);
    font-size: 0.9375rem;
    font-weight: 500;
    padding-top: 0.625rem;
    padding-bottom: 0.625rem;
    padding-left: 0.875rem;
    appearance: none;
    -webkit-appearance: none;
    outline: none;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .qs-select-wrapper select:focus {
    border-color: var(--accent-border);
    box-shadow: 0 0 0 3px var(--accent-dim);
  }

  .submit-btn {
    width: 100%;
    padding: 0.875rem 1.5rem;
    font-size: 1rem;
    font-weight: 700;
    border-radius: var(--r-md);
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
    color: #fff;
    border: none;
    box-shadow: 0 4px 16px rgba(220, 38, 38, 0.28);
    transition: all 0.2s var(--ease);
    cursor: pointer;
  }

  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 22px rgba(220, 38, 38, 0.38);
  }

  .submit-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
