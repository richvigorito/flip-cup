<script lang="ts">
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';
    import { connectSocket } from '$lib/transport/socket';
    import { fetchGames } from '$lib/transport/http/Games';
    import type { GameState } from '$lib/models/GameState';

    import { mode, gameId, loadingGames, currentPlayerName } from '$lib/store';

    let availableGames: GameState[] = [];


    function joinExistingGame() {
        if (!$gameId) return alert('Please enter a game code');
        connectSocket({
            type: 'join_existing_game',
            payload: {
                game_id: $gameId
            }
        });
        mode.set('lobby');
    }

    onMount(async() => {availableGames = await fetchGames(); });

</script>

<div class="join-game">
  <h2>Available Games:</h2>
  {#if $loadingGames}
    <p>Loading games...</p>
  {:else if availableGames.length === 0}
    <p>No games found.</p>
  {:else}
    <div class="game-grid">
      {#each availableGames as game}
        <button on:click={() => {
            gameId.set(game.id);
            currentPlayerName.set('');
            joinExistingGame();
            }}>
                <strong>Join {game.id}</strong><br />
                {game.teamA.name}: 
                {#if game.teamA.players.length > 0}
                    {game.teamA.players.map(p => p.name).join(', ')}
                {:else}
                    — No players —
                {/if}
                <br />
                {game.teamB.name}: 
                {#if game.teamB.players.length > 0}
                    {game.teamB.players.map(p => p.name).join(', ')}
                {:else}
                    — No players —
                {/if}
        </button>
      {/each}
    </div>
  {/if} 
  </div>
