<script lang="ts">
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';
    import { mode } from '$lib/store';
    import { connectSocket } from '$lib/transport/socket';

    import { fetchQuizzes } from '$lib/transport/http/Quizzes';

    import type { QuestionSet } from '$lib/types/QuestionSet'; 

    let questionSets: QuestionSet[] = [];
    let selectedQuestionSet = '';
    let localName = '';

    const createGame = () => {
        connectSocket({
            type: 'create_game',
            payload: {
                file: selectedQuestionSet
            }
        });
        mode.set('lobby');
    };

    onMount(async () => {
        questionSets = await fetchQuizzes();
    });
</script>

<div class="new-game">
    <h2>Create New Game</h2>

    <label for="qs-select">Choose Question Set:</label>
    <select
      id="qs-select"
      bind:value={selectedQuestionSet}
    >
      <option value="" disabled>-- Select a set --</option>

      {#each Object.entries(
        // Group by category
        questionSets.reduce<Record<string, QuestionSet[]>>((groups, qs) => {
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

    <button
      on:click={createGame}
      disabled={!selectedQuestionSet}
    >
      Create Game
    </button>
</div>
