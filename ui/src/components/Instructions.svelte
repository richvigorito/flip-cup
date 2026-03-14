<script lang="ts">
  let dialog: HTMLDialogElement;
  const open  = () => dialog.showModal();
  const close = () => dialog.close();
</script>

<button class="instructions-trigger" on:click={open}>
  📘 How to Play
</button>

<dialog bind:this={dialog}>
  <div class="dialog-inner">

    <div class="dialog-header">
      <h2 class="dialog-title">How to Play FlipQuiz</h2>
      <button class="dialog-close" on:click={close} aria-label="Close">✕</button>
    </div>

    <div class="dialog-body">
      <p><strong>Goal:</strong> Answer trivia questions correctly — first team to flip all their cups wins!</p>

      <h3>🎯 Starting a Game</h3>
      <ol>
        <li>Click <strong>Create New Game</strong> and pick a question set</li>
        <li>Share the Game ID with friends</li>
        <li>Everyone enters their name and clicks <strong>Join Game</strong></li>
      </ol>

      <h3>🎮 Joining a Friend's Game</h3>
      <ol>
        <li>Click <strong>Join Existing Game</strong></li>
        <li>Find the game card and click it</li>
        <li>Enter your name and join</li>
      </ol>

      <h3>⚡ Before You Start</h3>
      <ul>
        <li>Click <strong>Shuffle Teams</strong> to auto-assign players</li>
        <li>The host clicks <strong>Start Game</strong> when ready</li>
        <li>Both teams need at least 1 player</li>
      </ul>

      <h3>🏆 During the Game</h3>
      <ul>
        <li>Each player takes a turn answering a question</li>
        <li>Get it right → your cup flips!</li>
        <li>First team with all cups flipped wins 🎉</li>
      </ul>
    </div>

    <div class="dialog-footer">
      <button class="close-btn" on:click={close}>Got it!</button>
    </div>

  </div>
</dialog>

<style>
  .instructions-trigger {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.375rem 0.875rem;
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--text-secondary);
    background: var(--bg-surface);
    border: 1px solid var(--border);
    border-radius: var(--r-full);
    cursor: pointer;
    transition: all 0.15s;
  }
  .instructions-trigger:hover {
    color: var(--text-primary);
    border-color: var(--border-strong);
    background: var(--bg-elevated);
  }

  dialog {
    background: var(--bg-card);
    color: var(--text-primary);
    border: 1px solid var(--border-strong);
    border-radius: var(--r-xl);
    padding: 0;
    max-width: 480px;
    width: 90vw;
    box-shadow: 0 20px 60px rgba(0,0,0,0.6);
  }

  dialog::backdrop {
    background: rgba(11, 11, 24, 0.75);
    backdrop-filter: blur(6px);
  }

  .dialog-inner {
    padding: 1.75rem 2rem;
  }

  .dialog-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.25rem;
    padding-bottom: 1rem;
    border-bottom: 1px solid var(--border);
  }

  .dialog-title {
    font-size: 1.1rem;
    font-weight: 800;
    letter-spacing: -0.025em;
  }

  .dialog-close {
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--r-full);
    color: var(--text-muted);
    background: var(--bg-surface);
    border: 1px solid var(--border);
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.15s;
  }
  .dialog-close:hover {
    color: var(--text-primary);
    background: var(--bg-elevated);
  }

  .dialog-body {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .dialog-body p {
    font-size: 0.875rem;
    line-height: 1.65;
    color: var(--text-secondary);
    margin-bottom: 0.75rem;
  }
  .dialog-body p strong { color: var(--text-primary); }

  .dialog-body h3 {
    font-size: 0.75rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--accent);
    margin-top: 1.25rem;
    margin-bottom: 0.5rem;
  }

  .dialog-body ol,
  .dialog-body ul {
    padding-left: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .dialog-body li {
    font-size: 0.875rem;
    color: var(--text-secondary);
    line-height: 1.5;
  }
  .dialog-body li strong { color: var(--text-primary); }

  .dialog-footer {
    margin-top: 1.5rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border);
    display: flex;
    justify-content: flex-end;
  }

  .close-btn {
    padding: 0.6rem 1.5rem;
    font-size: 0.875rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    border: none;
    border-radius: var(--r-md);
    cursor: pointer;
    box-shadow: 0 3px 10px rgba(124,58,237,0.35);
    transition: all 0.2s;
  }
  .close-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 5px 16px rgba(124,58,237,0.5);
  }
</style>

