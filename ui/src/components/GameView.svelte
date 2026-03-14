<script lang="ts">
    import { send } from '$lib/transport/socket';
    import { mode, gameId, currentQuestion, currentPlayerName, gameState, myTeam, me, winner} from '$lib/store';

    let currentAnswer: string = '';

    $: tableColor = $myTeam ? $myTeam.color : '#000';

    const submitAnswer = () => {
        send({ 
            type: 'check_answer',
            payload: {answer: currentAnswer }
        });
    };

    const resetGame = () => {
        send({ type: 'restart_game' });
    };

    const quitGame = () => {
        if (confirm('Are you sure you want to quit the game? This will clear your session.')) {
            sessionStorage.removeItem('flipcup_player_id');
            sessionStorage.removeItem('flipcup_game_id');
            window.location.reload();
        }
    };

</script>
<div class="game-layout">
<button class="quit-btn" on:click={quitGame} style="position: absolute; top: 10px; right: 10px; z-index: 1000;">Quit Game</button>
{#if $mode == 'game' }
  <div class="game game-board">

    {#if $me && $me.isMyTurn && $currentQuestion && !$winner}
      <div class="question question-card">
        <p class="question-text">{$currentQuestion}</p>
        <input class="answer-input" type="text" bind:value={currentAnswer} placeholder="Your answer" />
        <button class="submit-btn" on:click={submitAnswer}>Submit Answer</button>
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
    <p class="winner-name">{$winner} wins!!</p>
    <button on:click={resetGame}>Restart Game</button>
  </div>
{/if}
