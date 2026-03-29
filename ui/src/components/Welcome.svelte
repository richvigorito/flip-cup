<script lang="ts">
  import { mode } from '$lib/store';

  const startNewGame = () => mode.set('new');
  const joinExistingGame = () => mode.set('join');
</script>

<div class="welcome">
  <div class="hero">
    <div class="hero-icon-wrap" aria-hidden="true">
      <div class="hero-cup-stage">
        <div class="hero-table"></div>
        <div class="hero-cup-scene">
          <div class="hero-cup-shadow"></div>
          <div class="hero-beer"></div>
          <img src="/solo-cup-hero.png" alt="" class="hero-icon" width="256" height="256" />
        </div>
      </div>
    </div>

    <h1 class="hero-title">
      Flip<span class="gradient-text">Cup</span>
    </h1>

    <p class="hero-tagline">
      Trivia with the energy of the drinking game. Line up your red cups, win your turns, and be the first team to clear the table.
    </p>

    <div class="hero-actions">
      <button class="btn btn-primary" on:click={startNewGame}>
        <span>🍺</span>
        Create New Game
      </button>
      <button class="btn btn-secondary" on:click={joinExistingGame}>
        <span>🍻</span>
        Join Existing Game
      </button>
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
    margin-bottom: 2rem;
  }

  .hero-cup-stage {
    position: relative;
    width: min(180px, 36vw);
    height: min(150px, 30vw);
  }

  .hero-cup-scene {
    position: relative;
    width: min(120px, 24vw);
    display: flex;
    justify-content: center;
    align-items: flex-end;
    height: 100%;
    margin: 0 auto;
    transform: translateY(-6px);
    transform-origin: 50% 50%;
    animation: flip-cup-sequence 6.2s cubic-bezier(0.22, 1, 0.36, 1) infinite;
  }

  .hero-table {
    position: absolute;
    left: 50%;
    bottom: 6px;
    width: 180%;
    height: 8px;
    transform: translateX(-50%);
    border-radius: 999px;
    background: linear-gradient(90deg, rgba(120, 53, 15, 0) 0%, rgba(180, 83, 9, 0.4) 18%, rgba(217, 119, 6, 0.52) 50%, rgba(180, 83, 9, 0.4) 82%, rgba(120, 53, 15, 0) 100%);
    z-index: 0;
    opacity: 0.9;
  }

  .hero-icon {
    width: 100%;
    height: auto;
    filter: drop-shadow(0 12px 24px rgba(220, 38, 38, 0.28));
    position: relative;
    z-index: 2;
  }

  .hero-cup-shadow {
    position: absolute;
    left: 50%;
    bottom: 2px;
    width: 62%;
    height: 14px;
    background: rgba(0, 0, 0, 0.22);
    border-radius: 999px;
    transform: translateX(-50%);
    filter: blur(8px);
    animation: cup-shadow 6.2s ease-in-out infinite;
  }

  .hero-beer {
    position: absolute;
    top: 25%;
    left: 50%;
    width: 46%;
    height: 12%;
    transform: translateX(-50%);
    border-radius: 50%;
    background:
      radial-gradient(ellipse at 50% 28%, rgba(255, 251, 191, 0.98) 0%, rgba(254, 240, 138, 0.98) 26%, rgba(245, 158, 11, 0.98) 72%, rgba(180, 83, 9, 0.98) 100%);
    box-shadow:
      inset 0 -6px 10px rgba(146, 64, 14, 0.18),
      0 0 0 2px rgba(255, 251, 235, 0.26),
      0 8px 18px rgba(217, 119, 6, 0.24);
    z-index: 4;
    animation: beer-level 6.2s ease-in-out infinite;
  }

  .hero-beer::before {
    content: '';
    position: absolute;
    left: 50%;
    top: 0;
    width: 74%;
    height: 16%;
    transform: translateX(-50%);
    border-radius: 999px;
    background: rgba(255, 251, 235, 0.92);
    box-shadow: 0 0 10px rgba(255, 251, 235, 0.32);
  }

  .hero-beer::after {
    content: '';
    position: absolute;
    right: 14%;
    top: 18%;
    width: 16%;
    height: 34%;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.28);
    filter: blur(1px);
  }

  @keyframes flip-cup-sequence {
    0%,
    14% {
      transform: translateY(-6px) rotate(-4deg);
    }
    20% {
      transform: translateY(-12px) rotate(2deg);
    }
    28%,
    31% {
      transform: translateY(-8px) rotate(0deg);
    }
    31.1%,
    50% {
      transform: translateY(26px) rotate(180deg);
    }
    61% {
      transform: translateY(24px) rotate(180deg);
    }
    62.5% {
      transform: translate(5px, 24px) rotate(192deg);
    }
    64.5% {
      transform: translate(11px, 2px) rotate(228deg);
    }
    66.5% {
      transform: translate(12px, -12px) rotate(276deg);
    }
    68.5% {
      transform: translate(7px, -22px) rotate(326deg);
    }
    70% {
      transform: translate(1px, -16px) rotate(350deg);
    }
    71.5%,
    100% {
      transform: translateY(-6px) rotate(360deg);
    }
  }

  @keyframes beer-level {
    0%,
    14% {
      opacity: 1;
      transform: translateX(-50%) scaleY(1);
    }
    20% {
      opacity: 0.96;
      transform: translateX(-50%) translateY(2px) scaleX(0.96) scaleY(0.66);
    }
    26% {
      opacity: 0.18;
      transform: translateX(-50%) translateY(5px) scaleX(0.82) scaleY(0.16);
    }
    30%,
    100% {
      opacity: 0;
      transform: translateX(-50%) scaleY(0.08);
    }
  }

  @keyframes cup-shadow {
    0%,
    28%,
    100% {
      opacity: 0.3;
      transform: translateX(-50%) scaleX(1);
    }
    31.1%,
    50% {
      opacity: 0.34;
      transform: translateX(-50%) scaleX(1.28);
    }
    61% {
      opacity: 0.22;
      transform: translateX(-50%) scaleX(1.08);
    }
    62.5% {
      opacity: 0.2;
      transform: translateX(-46%) scaleX(1.02);
    }
    64.5% {
      opacity: 0.17;
      transform: translateX(-41%) scaleX(0.84);
    }
    66.5% {
      opacity: 0.16;
      transform: translateX(-42%) scaleX(0.72);
    }
    68.5% {
      opacity: 0.2;
      transform: translateX(-47%) scaleX(0.76);
    }
    70% {
      opacity: 0.26;
      transform: translateX(-49%) scaleX(0.86);
    }
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
    margin: 0 auto;
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
    background: linear-gradient(135deg, var(--accent), var(--accent-secondary));
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

</style>
