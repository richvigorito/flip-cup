<script lang="ts">

    import { mode, availableGames, loadingGames } from '$lib/store';

    const startNewGame        = () => mode.set('new');
    const joinExistingGame    = () => mode.set('join');

</script>

<div class="welcome">
    <h2>Welcome to Flip Quiz</h2>
    <button on:click={startNewGame}>Create New Game</button>
    <button on:click={() => {joinExistingGame(); }}>Join Existing Game</button>
  </div>
