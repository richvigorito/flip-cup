<script lang="ts">
    import { get } from 'svelte/store';
    import { mode, gameId, joined, currentPlayerName, gameState, gamesCompleted } from '$lib/store';
    import { send } from '$lib/transport/socket';

    let localPlayerName = '';
    
    import QuestionSetDropdown from './QuestionSetDropdown.svelte';

    const assignTeams = () => {
        send({ type: 'assign_teams' });
    };

    console.log('gsx', $gameState);

    const startGame = () => {
        if ($gameState.canStart()){
            console.log("Sending start message");
            send({ type: 'start' });
            mode.set('game');
        } else {
            console.log("Please assign teams before starting the game.");
            alert('Please assign teams before starting the game. Teams must have an equal number.');
        }
    };

    const joinGame = (playerName: string) => {
        const newPlayer = { name: playerName, team: null };

        joined.set(true);
        console.log(JSON.stringify({ type: "add_player", name: playerName }));
        send({ 
            type: "add_player", 
            payload: {
                name: playerName 
            }
        });
    };
    console.log('completed', $gamesCompleted, $gamesCompleted > 0 );

    function handleNewQuiz(selectedFile: string) {
        console.log("New quiz selected:", selectedFile);
        send({
            type: "update_quiz",
            payload: {
                quizfile: selectedFile
            }
        });

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

    {#if $gamesCompleted > 0 }
        <br/>
        <label for="qs-select" class="quiz-label">
            Current Quiz (you can change it):
        </label>
        <QuestionSetDropdown on:select={(e) => handleNewQuiz(e.detail)} />
    {/if}

</div>
