<script lang="ts">
  import { mode } from '$lib/store';
  import { connectSocket } from '$lib/transport/socket';
  import type { QuestionSet } from '$lib/types/QuestionSet';
  import { questionSets } from '$lib/store';

  let selectedQuestionSet = '';

  const createGame = () => {
    if (!selectedQuestionSet) return;
    connectSocket({
      type: 'create_game',
      payload: { file: selectedQuestionSet }
    });
    mode.set('lobby');
  };

  const goBack = () => mode.set('welcome');
</script>

<div class="page-wrap">

  <!-- Back button -->
  <button class="back-btn" on:click={goBack}>
    ← Back
  </button>

  <div class="form-card">

    <!-- Header -->
    <div class="card-header">
      <div class="card-icon">🎯</div>
      <h2 class="card-title">Create a Game</h2>
      <p class="card-subtitle">Pick a quiz category to get started</p>
    </div>

    <!-- Quiz Picker -->
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

    <!-- Action -->
    <button
      class="submit-btn"
      on:click={createGame}
      disabled={!selectedQuestionSet}
    >
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

  /* Card header */
  .card-header {
    text-align: center;
    margin-bottom: 2rem;
  }
  .card-icon {
    font-size: 2.5rem;
    margin-bottom: 0.75rem;
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

  /* Field */
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

  /* Wrapper with chevron */
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

  /* Submit */
  .submit-btn {
    width: 100%;
    padding: 0.875rem 1.5rem;
    font-size: 1rem;
    font-weight: 700;
    border-radius: var(--r-md);
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    border: none;
    box-shadow: 0 4px 16px rgba(124, 58, 237, 0.35);
    transition: all 0.2s var(--ease);
  }
  .submit-btn:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 22px rgba(124, 58, 237, 0.5);
  }
  .submit-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
</style>
