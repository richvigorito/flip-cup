<script lang="ts">
    import { get } from 'svelte/store';
    import { send } from '$lib/transport/socket';
    import { mode, gameId, currentQuestion, currentPlayerName, gameState, myTeam, me, winner} from '$lib/store';

    let currentAnswer: string = '';

    $: tableColor = $myTeam.color;

    const submitAnswer = () => {
        send({ 
            type: 'check_answer',
            payload: {answer: currentAnswer }
        });
    };

    const resetGame = () => {
        send({ type: 'restart_game' });
    };

</script>
<div class="game-layout">
{#if $mode == 'game' }
  <div class="game">

    {#if $me.isMyTurn && $currentQuestion && !$winner}
      <div class="question">
        <p>{$currentQuestion}</p>
        <input type="text" bind:value={currentAnswer} placeholder="Your answer" />
        <button on:click={submitAnswer}>Submit Answer</button>
      </div>
    {/if}

    <!-- Team Display -->
    <div class="table" style="--table-color: {tableColor}">
      <!-- Team A -->
      <div class="team-column">
        <div class>{$gameState.teamA.name}</div>
        {#each $gameState.teamA.players as player, index}
          <div class="team-cups">
            <div class="player-name">{player.name}</div>
            <div class="cup {$gameState.teamA.turn > index ? 'flipped' : ''} {$gameState.teamA.turn == (index) ? "current" : ''}">
            </div>
          </div>
          <hr>
        {/each}
      </div>

      <!-- Team B -->
      <div class="team-column">
        <div class>{$gameState.teamB.name}</div>
        {#each $gameState.teamB.players as player, index}
          <div class="team-cups">
            <div class="player-name">{player.name}</div>
            <div class="cup {$gameState.teamB.turn > index ? 'flipped' : ''} {$gameState.teamB.turn == (index) ? "current" : ''}"></div>
          </div>
          <hr>
        {/each}
      </div> 

    </div><!-- end Team Display -->

  </div>
{/if}
</div>

<!-- Game Over Section -->
{#if $winner }
  <div class="game-over">
    <h2>Game Over!</h2>
    <p>{$winner} wins!!</p>
    <button on:click={resetGame}>Restart Game</button>
  </div>
{/if}
