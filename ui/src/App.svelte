<header class="site-header">
  <img src="/flip-banner.png" alt="Flip-Quiz Banner" class="banner-image" />
</header>
<script lang="ts">
    import '$styles/App.css';
    import { socket } from '$lib/transport/socket';
    import { onMount } from 'svelte';
    import { mode } from '$lib/store';

    import Welcome from './components/Welcome.svelte';
    import NewGame from './components/NewGame.svelte';
    import JoinGame from './components/JoinGame.svelte';
    import Lobby from './components/Lobby.svelte';
    import GameView from './components/GameView.svelte';
    import EventLog from './components/EventLog.svelte';
    import Instructions from './components/Instructions.svelte';

    import { fetchQuizzes } from '$lib/transport/http/Quizzes';
    import { questionSets } from '$lib/store';


    onMount(async () => {
        const tmp = await fetchQuizzes();
        questionSets.set(tmp);
    });


</script>


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
    <p>Unknown mode: {$mode}</p>
{/if}


<Instructions />

<EventLog />

