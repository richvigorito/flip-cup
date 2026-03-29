<script lang="ts">
  import { mode } from '$lib/store';

  const startNewGame = () => mode.set('new');
  const joinExistingGame = () => mode.set('join');
</script>

<div class="welcome">
  <div class="hero">
    <div class="hero-icon-wrap" aria-hidden="true">
      <img src="/solo-cup.png" alt="" class="hero-icon" />
    </div>

    <h1 class="hero-title">
      Flip<span class="gradient-text">Cup</span>
    </h1>

    <p class="hero-tagline">
      Trivia with the energy of the drinking game. Line up your red cups, win your turns, and be the first team to clear the table.
    </p>

    <div class="hero-actions">
      <button class="btn btn-primary" on:click={startNewGame}>
        <span>🏁</span>
        Create New Game
      </button>
      <button class="btn btn-secondary" on:click={joinExistingGame}>
        <span>👥</span>
        Join Existing Game
      </button>
    </div>

    <div class="hero-tags">
      <span class="tag">Red Cup Relay</span>
      <span class="tag">House Party Trivia</span>
      <span class="tag">Team vs Team</span>
    </div>
  </div>
</div>

<style>
  .welcome {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 64px);
    padding: 2rem 1.5rem;
  }

  .hero {
    text-align: center;
    max-width: 560px;
    width: 100%;
  }

  .hero-icon-wrap {
    display: flex;
    justify-content: center;
    margin-bottom: 1.25rem;
  }

  .hero-icon {
    width: min(132px, 28vw);
    height: auto;
    animation: float 3.5s ease-in-out infinite;
    filter: drop-shadow(0 12px 24px rgba(220, 38, 38, 0.28));
  }

  @keyframes float {
    0%, 100% { transform: translateY(0px) rotate(-2deg); }
    50% { transform: translateY(-14px) rotate(2deg); }
  }

  .hero-title {
    font-size: clamp(3.25rem, 9vw, 5.5rem);
    font-weight: 900;
    letter-spacing: -0.05em;
    line-height: 1;
    color: var(--text-primary);
    margin-bottom: 1rem;
  }

  .gradient-text {
    background: linear-gradient(135deg, #f87171 0%, #fb923c 55%, #fbbf24 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .hero-tagline {
    font-size: 1.05rem;
    color: var(--text-secondary);
    line-height: 1.65;
    margin-bottom: 2.5rem;
    max-width: 460px;
    margin-left: auto;
    margin-right: auto;
  }

  .hero-actions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    max-width: 320px;
    margin: 0 auto 2rem;
  }

  .btn {
    width: 100%;
    padding: 0.875rem 1.75rem;
    font-size: 1rem;
    font-weight: 700;
    border-radius: var(--r-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    transition: all 0.2s var(--ease);
  }

  .btn-primary {
    background: linear-gradient(135deg, var(--accent), var(--indigo));
    color: #fff;
    border: none;
    box-shadow: 0 4px 20px rgba(220, 38, 38, 0.32);
  }

  .btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 28px rgba(220, 38, 38, 0.44);
  }

  .btn-primary:active { transform: translateY(0); }

  .btn-secondary {
    background: var(--bg-card);
    color: var(--text-primary);
    border: 1px solid var(--border-strong);
  }

  .btn-secondary:hover {
    background: var(--bg-surface);
    border-color: var(--accent-border);
    transform: translateY(-2px);
  }

  .hero-tags {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .tag {
    padding: 0.3rem 0.75rem;
    font-size: 0.78rem;
    font-weight: 600;
    color: var(--text-secondary);
    background: rgba(255, 237, 213, 0.04);
    border: 1px solid rgba(248, 113, 113, 0.14);
    border-radius: var(--r-full);
    letter-spacing: 0.01em;
  }
</style>
