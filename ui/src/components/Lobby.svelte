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
        if ($gameState && $gameState.canStart()){
            console.log("Sending start message");
            send({ type: 'start' });
            mode.set('game');
        } else {
            console.log("Please assign teams before starting the game.");
            alert('Please assign teams before starting the game. Teams must have at least 1 player.');
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

    const leaveGame = () => {
        if (confirm('Are you sure you want to leave the game? This will clear your session.')) {
            sessionStorage.removeItem('flipcup_player_id');
            sessionStorage.removeItem('flipcup_game_id');
            window.location.reload();
        }
    };

</script>

<div class="lobby">
    <div class="lobby-header">
        <span class="lobby-icon" on:click={leaveGame} style="cursor: pointer;" title="Leave Game">🏠</span>
        <h2>Lobby Code: <span class="game-code-value">{$gameId}</span></h2>
    </div>

    {#if !$joined && $mode == 'lobby' }
        <input type="text" placeholder="Enter your name…" bind:value={localPlayerName} />

        <button disabled={!localPlayerName} on:click={() => {
            currentPlayerName.set(localPlayerName); 
            joinGame(localPlayerName);
        }}>Join Game</button>

    {:else}
        <div class="teams-preview">
            <div class="team">
                <h3>{$gameState?.teamA?.name || 'Team A'}</h3>
                <ul>
                    {#each $gameState?.teamA?.players || [] as player}
                        <li>{player.name}</li>
                    {/each}
                </ul>
            </div>
            <div class="team">
                <h3>{$gameState?.teamB?.name || 'Team B'}</h3>
                <ul>
                    {#each $gameState?.teamB?.players || [] as player}
                        <li>{player.name}</li>
                    {/each}
                </ul>
            </div>
        </div>
        <button on:click={assignTeams}>Shuffle Teams</button>
        <button disabled={!$gameState || !$gameState.canStart()} on:click={startGame}>Start Game</button>
        <div style="margin-top: 1rem;">
             <button class="secondary" on:click={leaveGame}>Leave Game</button>
        </div>
    {/if}

    {#if $gamesCompleted > 0 }
        <br/>
        <label for="qs-select" class="quiz-label">
            Current Quiz (you can change it):
        </label>
        <QuestionSetDropdown on:select={(e) => handleNewQuiz(e.detail)} />
    {/if}

</div>
