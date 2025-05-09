<script lang="ts">
  import { onMount } from 'svelte';

  let currentPlayerName: string = ''; // Name for the joining player
  let isGameStarted = false;
  let gameOver = false;
  let joined = false;
  let players: { name: string, team: string | null }[] = [];
  let assignedTeams: { team1: string[], team2: string[] } = { team1: [], team2: [] };
  let socket: WebSocket;
  let currentQuestion: string | null = null;
  let currentAnswer: string = '';

  let team1Name: string = ''  ;
  let team2Name: string = ''  ;
  let teamTurn1 = 1;
  let teamTurn2 = 1;

  let eventLog: { message: string, type: 'info' | 'success' | 'error' }[] = [];
  let logContainer: HTMLDivElement;

  $: {
    // Auto-scroll to bottom when new log is added
    if (logContainer) {
      setTimeout(() => {
        logContainer.scrollTop = logContainer.scrollHeight;
      }, 0);
    }
  }

  $: userTeam = players.find(p => p.name === currentPlayerName)?.team;
  $: tableColor = userTeam === 'team1' ? '#563517' : '#9c6f44';
  
  // WebSocket connection setup
onMount(() => {

    const wsUrl = 'ws://'+import.meta.env.VITE_WS_URL+'/ws' || 'ws://localhost:8080/ws';
    console.log(wsUrl)

    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('Connected to WebSocket');
    };

    socket.onmessage = (event) => {
      console.log('Raw event from WebSocket:', event);
      const message = JSON.parse(event.data);
      handleMessage(message);
    };

    socket.onclose = () => {
      console.log('Disconnected from WebSocket');
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
});

function arraysEqual(a: string[], b: string[]) {
  if (a.length !== b.length) return false;
  const sortedA = [...a].sort();
  const sortedB = [...b].sort();
  return sortedA.every((val, index) => val === sortedB[index]);
}

  const handleMessage = (message: any) => {
    console.log("handle incoming message:", JSON.stringify(message));
    switch (message.type) {

      case 'incorrect_answer':
          eventLog = [...eventLog, { message: `${message.name} submitted wrong answer`, type: 'error' }];
          break;
      case 'correct':
        currentQuestion = null;
        break;
      case 'question':
        currentQuestion = message.name;
        // myTeamTurn = message.team;
        // playerTurnIndex = message.index;
        break;
      case 'winner':
        gameOver = true;
        eventLog = [...eventLog, { message: `🏆🏆  Winner: ${message.name} 🏆🏆`, type: 'error' }];
        break;
      case 'player_joined':
        if (!players.find(p => p.name === message.name)) {
          players = [...players, { name: message.name, team: null }];
        }
        eventLog = [...eventLog, { message: `${message.name} joined the game.`, type: 'success' }];
        break;
      case 'assign_teams':
        assignedTeams.team1 = message.team1;
        assignedTeams.team2 = message.team2;
        break;
      case 'teams_update':
        // Update teams based on the message data
        const prevTeam1 = [...assignedTeams.team1];
        const prevTeam2 = [...assignedTeams.team2]
        assignedTeams.team1 = message.answer.teamA.players;
        assignedTeams.team2 = message.answer.teamB.players;
        const teamsUpdated = 
          !arraysEqual(prevTeam1, assignedTeams.team1) || 
          !arraysEqual(prevTeam2, assignedTeams.team2);

        team1Name = message.answer.teamA.name
        team2Name = message.answer.teamB.name

        teamTurn1 = message.answer.teamA.turn;
        teamTurn2 = message.answer.teamB.turn;

        // Update players[] with their assigned team
        players = players.map(p => {
          if (assignedTeams.team1.includes(p.name)) {
            return { ...p, team: 'team1' };
          } else if (assignedTeams.team2.includes(p.name)) {
            return { ...p, team: 'team2' };
          } else {
            return { ...p, team: null };
          }
        });

        if (teamsUpdated){
          eventLog = [...eventLog, { message: `teams assigned`, type: 'success' }];
          eventLog = [...eventLog, { message: `team1: ${assignedTeams.team1.join(', ')}`, type: 'info' }];
          eventLog = [...eventLog, { message: `team2: ${assignedTeams.team2.join(', ')}`, type: 'info' }];
        }
        break; 
      case 'game_started':
        isGameStarted = true;
        eventLog = [...eventLog, { message: `Game Started`, type: 'success' }];
        break;
      case 'restart':
        resetGame();
        break;
    }
  }; 

  const resetGame = () => {
    gameOver = false;
    isGameStarted = false;
    players = [];
    assignedTeams = { team1: [], team2: [] };
    socket.send(JSON.stringify({ type: 'start' }));
  };

  const joinGame = (playerName: string) => {
    const newPlayer = { name: playerName, team: null };
    players = [...players, newPlayer];

    // Send message to the server to join
    if (socket && socket.readyState === WebSocket.OPEN) {
      joined = true;
      socket.send(JSON.stringify({ type: "join", name: playerName }));
      console.log("Sent join message:", playerName);
    } else {
      console.error("WebSocket not connected!");
    }
  };

  const assignTeams = () => {
    if (players.length < 2) return alert('Need at least 2 players to assign teams!');
    socket.send(JSON.stringify({ type: 'assign_teams' }));
  };

  const startGame = () => {
    if (assignedTeams.team1.length > 0 && assignedTeams.team2.length > 0) {
      console.log("Sending start message");
      socket.send(JSON.stringify({ type: 'start' }));
    } else {
      console.log("Please assign teams before starting the game.");
      alert('Please assign teams before starting the game.');
    }
  };

  const submitAnswer = () => {
    socket.send(JSON.stringify({ type: 'answer', answer: currentAnswer }));
    currentAnswer = ''; 
  };
</script>

<style>
.game-layout {
  margin-right: 320px; /* give space for the log sidebar */
  padding: 1rem;
}

.event-log {
  position: fixed;
  top: 0;
  right: 0;
  height: 100vh;
  width: 75px;
  background: rgba(20, 20, 20, 0.95);
  color: white;
  padding: 1rem;
  border-left: 2px solid #333;
  font-size: 0.9rem;
  display: flex;
  flex-direction: column;
  z-index: 1000;
}

.event-log h3 {
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  color: #ddd;
  border-bottom: 1px solid #333;
  padding-bottom: 0.25rem;
}

.log-entries {
  flex: 1;
  overflow-y: auto;
  scrollbar-width: thin;
}

.log-entry {
  padding: 0.4rem;
  border-bottom: 1px solid #333;
  line-height: 1.0;
  font-weight: 300;
  font-size: x-small;
}

.log-entry.info {
  color: #aaa;
}

.log-entry.success {
  color: #22c55e;
}

.log-entry.error {
  color: #ef4444;
}

.table {
  min-width: 200px;
  width: 80vw;
  max-width: 300px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  padding: 1rem;
  border-radius: 12px;
  background-color: var(--table-color);
}

.team-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  border:2px black;
  gap: 1rem;
  border-radius: 25px;
  border: 2px solid #08080845;
  padding: 11px;
  width: 64px;
}
  
.team-cups {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.cup {
  width: 60px;
  height: 60px;
  background-image: url('/solo-cup.png');
  background-size: cover;
  transition: transform 0.3s ease;
}

.cup.flipped {
  transform: rotate(180deg);
  opacity: 0.6;
}


.cup.current {
  border:4px solid black;
}

.player-name {
  font-size: 0.85rem;
  text-align: center;
  color: white;
}
  
</style>

<!-- Lobby Section -->
{#if !isGameStarted}
  <div class="lobby">

    {#if !joined}
      <input type="text" placeholder="Enter name" bind:value={currentPlayerName} />
      <button on:click={() => joinGame(currentPlayerName)}>Join Game</button>
    {/if}
    {#if players.length >= 2 && 
      !(assignedTeams.team1.length > 0 && assignedTeams.team2.length > 0)
    }
      <button on:click={assignTeams}>Assign Teams</button>
    {/if}
    {#if assignedTeams.team1.length > 0 && assignedTeams.team2.length > 0}
      <button on:click={startGame}>Start Game</button>
    {/if}
  </div>
{/if}

<!-- Game Section -->
<div class="game-layout">
{#if isGameStarted}
  <div class="game">

    {#if !gameOver && currentQuestion}
      <div class="question">
        <p>{currentQuestion}</p>
        <input type="text" bind:value={currentAnswer} placeholder="Your answer" />
        <button on:click={submitAnswer}>Submit Answer</button>
      </div>
    {/if}

    <!-- Team Display -->
    <div class="table" style="--table-color: {tableColor}">
      <!-- Team A -->
      <div class="team-column">
        <div class>{team1Name}</div>
        {console.log(team1Name)}_
        {#each assignedTeams.team1 as player, index}
          <div class="team-cups">
            {console.log(teamTurn1, index, player)}
            <div class="player-name">{player}</div>
            <div class="cup {teamTurn1 > index ? 'flipped' : ''} {teamTurn1 == (index) ? "current" : ''}">
            </div>
          </div>
          <hr>
        {/each}
      </div>

      <!-- Team B -->
      <div class="team-column">
        <div class>{team2Name}</div>
        {console.log(team2Name)}_
        {#each assignedTeams.team2 as player, index}
          <div class="team-cups">
            {console.log(teamTurn2, index)}
            <div class="player-name">{player}</div>
            <div class="cup {teamTurn2 > index ? 'flipped' : ''} {teamTurn2 == (index) ? "current" : ''}"></div>
          </div>
          <hr>
        {/each}
      </div> 

    </div><!-- end Team Display -->

  </div>
{/if}
</div>

<div class="event-log">
  <h3>Game Log</h3>
  <div class="log-entries" bind:this={logContainer}>
    {#each eventLog as log}
      <div class="log-entry {log.type}">
        {log.message}
      </div>
    {/each}
  </div>
</div>

<!-- Game Over Section -->
{#if gameOver}
  <div class="game-over">
    <h2>Game Over!</h2>
    <p>{assignedTeams.team1.length === 0 ? team2Name : team1Name} wins!!</p>
    <button on:click={resetGame}>Restart Game</button>
  </div>
{/if}

