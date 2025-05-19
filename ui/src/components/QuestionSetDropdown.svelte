<!-- src/lib/components/QuestionSetDropdown.svelte -->
<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { writable } from 'svelte/store';
    import { mode, gameState } from '$lib/store';
    import { connectSocket } from '$lib/transport/socket';

    import { fetchQuizzes } from '$lib/transport/http/Quizzes';

    import type { QuestionSet } from '$lib/types/QuestionSet';
    import { questionSets } from '$lib/store';

    console.log('gs', $gameState)
    let selectedQuestionSet = $gameState?.quizfile ? $gameState.quizfile : '';

     export let onSelect: ((file: string) => void) | null = null;

  const dispatch = createEventDispatcher();

  function handleChange(event: Event) {
    const value = (event.target as HTMLSelectElement).value;
    selectedQuestionSet = value;
    dispatch('select', value);
  };


</script>
<div>
    <select
      id="qs-select"
      bind:value={selectedQuestionSet}
      on:change={handleChange}
    >
      <option value="" disabled>-- Select a set --</option>

      {#each Object.entries(
        // Group by category
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
</div>
