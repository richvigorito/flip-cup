<script lang="ts">
    import { get } from 'svelte/store';
    import { mode, gameId, joined, currentPlayerName, gameState } from '$lib/store';
    import { send } from '$lib/transport/socket';

    let localPlayerName = '';

    const assignTeams = () => {
        send({ type: 'assign_teams' });
    };

    const startGame = () => {
        if ($gameState.canStart()){
            console.log("Sending start message");
            send({ type: 'start' });
            mode.set('game');
        } else {
            console.log("Please assign teams before starting the game.");
            alert('Please assign teams before starting the game. Teams must have an equal number.');
        }
    };

    const joinGame = (playerName: string) => {
        const newPlayer = { name: playerName, team: null };

        joined.set(true);
        console.log(JSON.stringify({ type: "add_player", name: playerName }));
        send({ 
            type: "add_player", 
            payload: {
                name: playerName 
            }
        });
    };

</script>

<div class="lobby">
    {#if !$joined && $mode == 'lobby' }
        <input type="text" placeholder="Enter name" bind:value={localPlayerName} />

        <button on:click={() => {
            currentPlayerName.set(localPlayerName); 
            joinGame(localPlayerName);
        }}>Join Game</button>

    {:else}
      <button on:click={assignTeams}>Assign Teams</button>
            <button on:click={startGame}>Start Game</button>
    {/if}
  </div>
