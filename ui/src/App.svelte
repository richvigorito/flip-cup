<script lang="ts">
    import '$styles/App.css';
    import { socket, connectSocket } from '$lib/transport/socket';
    import { onMount } from 'svelte';
    import { mode } from '$lib/store';

    import Welcome     from './components/Welcome.svelte';
    import NewGame     from './components/NewGame.svelte';
    import JoinGame    from './components/JoinGame.svelte';
    import Lobby       from './components/Lobby.svelte';
    import GameView    from './components/GameView.svelte';
    import EventLog    from './components/EventLog.svelte';
    import Instructions from './components/Instructions.svelte';

    import { fetchQuizzes } from '$lib/transport/http/Quizzes';
    import { questionSets, gameId } from '$lib/store';

    onMount(async () => {
        const tmp = await fetchQuizzes();
        questionSets.set(tmp);

        // Attempt reconnect if session exists
        if (typeof sessionStorage !== 'undefined') {
            const pid = sessionStorage.getItem('flipcup_player_id');
            const gid = sessionStorage.getItem('flipcup_game_id');
            if (pid && gid) {
                connectSocket({
                    type: 'join_existing_game',
                    payload: {
                        game_id: gid,
                        player_id: pid
                    }
                });

                mode.set('lobby'); 
                gameId.set(gid); // Ensure gameId is set for Lobby display
            }
        }
    });
</script>

<svelte:head>
  <title>FlipCup — Answer. Flip. Win.</title>
</svelte:head>

<!-- ── Fixed Header ── -->
<header class="site-header">
  <div class="header-inner">
    <span class="logo-link" role="button" tabindex="0" on:click={() => mode.set('welcome')} on:keydown={(e) => e.key === 'Enter' && mode.set('welcome')}>
      <img src="/solo-cup.png" alt="Red solo cup" class="logo-icon-img" />
      <span class="logo-text">Flip<span class="logo-accent">Cup</span></span>
    </span>
    <div class="header-right">
      <Instructions />
    </div>
  </div>
</header>

<!-- ── Main Content ── -->
<main class="main-content" class:has-sidebar={['lobby', 'game'].includes($mode)}>
  {#if $mode === 'welcome'}
    <Welcome />
  {:else if $mode === 'new'}
    <NewGame />
  {:else if $mode === 'join'}
    <JoinGame />
  {:else if $mode === 'lobby'}
    <Lobby />
  {:else if $mode === 'game'}
    <GameView />
  {:else}
    <p style="color: var(--text-muted); text-align: center;">Unknown mode: {$mode}</p>
  {/if}
</main>

<!-- ── Sidebar (Lobby & Game only) ── -->
{#if ['lobby', 'game'].includes($mode)}
  <EventLog />
{/if}

<style>
  .main-content {
    padding-top: 64px;
    min-height: 100vh;
  }
  .main-content.has-sidebar {
    padding-right: 264px;
  }
  .header-right {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
</style>
