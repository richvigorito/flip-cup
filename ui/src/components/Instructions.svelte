<script lang="ts">
  let instructionsDialog: HTMLDialogElement;

  function openInstructions() {
    instructionsDialog.showModal();
  }

  function closeInstructions() {
    instructionsDialog.close();
  }
</script>

<button class="instructions-button" on:click={openInstructions}>
  📘 Instructions
</button>

<dialog bind:this={instructionsDialog} class="instructions-dialog">
  <h2>How to Play Flip-Quiz</h2>

  <p><strong>Goal:</strong> Take turns answering questions. First team done wins!</p>

  <h3>Start a Game</h3>
  <ol>
    <li>Click <strong>Start a Game</strong></li>
    <li>Pick a quiz and click <strong>Create Game</strong></li>
    <li>Enter your name and click <strong>Join Game</strong></li>
  </ol>

  <h3>Join a Friend's Game</h3>
  <ol>
    <li>Click <strong>Join Existing Game</strong></li>
    <li>Find their game, enter your name, click <strong>Join Game</strong></li>
  </ol>

  <h3>Before Starting</h3>
  <ul>
    <li>Teams are assigned automatically</li>
    <li>Click <strong>Assign Teams</strong> to rebalance (optional)</li>
    <li>Click <strong>Start Game</strong></li>
  </ul>

  <h3>During Game</h3>
  <ul>
    <li>Each player answers a question on their turn</li>
    <li>Get it right to pass the turn!</li>
    <li>First team to finish wins</li>
  </ul>

  <form method="dialog">
    <button on:click={closeInstructions}>Close</button>
  </form>
</dialog>

