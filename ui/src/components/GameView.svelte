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


             console.log(Array($gameState.cups));
             console.log($gameState.cups);
             console.log('gs', $gameState);



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
            {#each Array($gameState.cups) as _, cupRound}
        {#each $gameState.teamA.players as player, index}
          <div class="team-cups">
            <div class="player-name">{player.name}</div>
            <div class="cup 
                {(cupRound * $gameState.teamA.turn) > (cupRound * index) ? 'flipped' : ''} 
                {(cupRound * $gameState.teamA.turn) == (cupRound * index) ? "current" : ''}"
            >
            </div>
          </div>
          <hr>
        {/each}
        {/each}
      </div>

      <!-- Team B -->
      <div class="team-column">
        <div class>{$gameState.teamB.name}</div>
            {#each Array($gameState.cups) as _, cupRound}
        {#each $gameState.teamB.players as player, index}
          <div class="team-cups">
            <div class="player-name">{player.name}</div>
            <div class="cup 
                {(cupRound * $gameState.teamB.turn) > (cupRound * index) ? 'flipped' : ''} 
                {(cupRound * $gameState.teamB.turn) == (cupRound * index) ? "current" : ''}">
            </div>
          </div>
          <hr>
        {/each}
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
