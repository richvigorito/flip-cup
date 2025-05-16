<script lang="ts">
    import { get } from 'svelte/store';
    import { mode, gameId, joined, currentPlayerName, gameState } from '$lib/store';
    import { socket } from '$lib/transport/socket';

    let localPlayerName = '';

    const ws = get(socket);

    const assignTeams = () => {
        ws.send(JSON.stringify({ type: 'assign_teams' }));
    };

    const startGame = () => {
        if ($gameState.canStart()){
            console.log("Sending start message");
            ws.send(JSON.stringify({ type: 'start' }));
            mode.set('game');
        } else {
            console.log("Please assign teams before starting the game.");
            alert('Please assign teams before starting the game. Teams must have an equal number.');
        }
    };

    const joinGame = (playerName: string) => {
        const newPlayer = { name: playerName, team: null };

        // Send message to the server to join
        if (ws && ws.readyState === WebSocket.OPEN) {
         joined.set(true);
        console.log(JSON.stringify({ type: "add_player", name: playerName }));
        ws.send(JSON.stringify({ 
          type: "add_player", 
          payload: {
              name: playerName 
          }
        }));
        console.log("Sent join message:", playerName);
        } else {
        console.error("WebSocket not connected!");
        }
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
