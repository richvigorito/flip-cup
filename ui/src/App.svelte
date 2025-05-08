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

  $: userTeam = players.find(p => p.name === currentPlayerName)?.team;
  $: tableColor = userTeam === 'team1' ? '#563517' : '#9c6f44';
  
  // WebSocket connection setup
  onMount(() => {
    socket = new WebSocket('ws://localhost:8080/ws'); // Replace with your real URL

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

  const handleMessage = (message: any) => {
    console.log("handle incoming message:", JSON.stringify(message));
    switch (message.type) {

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
        break;
      case 'player_joined':
        if (!players.find(p => p.name === message.name)) {
          players = [...players, { name: message.name, team: null }];
        }
        break;
      case 'assign_teams':
        assignedTeams.team1 = message.team1;
        assignedTeams.team2 = message.team2;
        break;
      case 'teams_update':
        // Update teams based on the message data
        assignedTeams.team1 = message.answer.teamA.players;
        assignedTeams.team2 = message.answer.teamB.players;
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
        break; 
      case 'game_started':
        isGameStarted = true;
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
.table {
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
    <!-- <div class="table"> -->
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

<!-- Game Over Section -->
{#if gameOver}
  <div class="game-over">
    <h2>Game Over!</h2>
    <p>{assignedTeams.team1.length === 0 ? team2Name : team1Name} wins!!</p>
    <button on:click={resetGame}>Restart Game</button>
  </div>
{/if}

