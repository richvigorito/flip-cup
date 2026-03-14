<!-- src/lib/components/QuestionSetDropdown.svelte -->
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { gameState, questionSets } from '$lib/store';
  import type { QuestionSet } from '$lib/types/QuestionSet';

  let selectedQuestionSet = $gameState?.quizfile ?? '';

  const dispatch = createEventDispatcher();

  function handleChange(event: Event) {
    const value = (event.target as HTMLSelectElement).value;
    selectedQuestionSet = value;
    dispatch('select', value);
  }
</script>

<div class="dropdown-wrap">
  <select
    id="qs-select"
    bind:value={selectedQuestionSet}
    on:change={handleChange}
  >
    <option value="" disabled>— Choose a category —</option>

    {#each Object.entries(
      $questionSets.reduce<Record<string, QuestionSet[]>>((groups, qs) => {
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
  <span class="chevron">▾</span>
</div>

<style>
  .dropdown-wrap {
    position: relative;
  }
  select {
    width: 100%;
    padding: 0.625rem 2.25rem 0.625rem 0.875rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    background: var(--bg-surface);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-md);
    appearance: none;
    -webkit-appearance: none;
    cursor: pointer;
    outline: none;
    transition: border-color 0.15s, box-shadow 0.15s;
  }
  select:focus {
    border-color: var(--accent-border);
    box-shadow: 0 0 0 3px var(--accent-dim);
  }
  .chevron {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-muted);
    font-size: 0.8rem;
    pointer-events: none;
  }
</style>
