<script lang="ts">
  import { onMount } from 'svelte';

  let currentPlayerName: string = ''; // Name for the joining player
  let cupsPerTeam = 1;  // One cup per player
  let flippedTeam1: number[] = [];
  let flippedTeam2: number[] = [];
  let currentTurn: 1 | 2 = 1;
  let gameOver = false;
  let isGameStarted = false;
  let players: { name: string, team: string | null }[] = [];
  let socket: WebSocket;
  let question: string | null = null;
  let currentAnswer: string = '';
  let assignedTeams: { team1: string[], team2: string[] } = { team1: [], team2: [] };

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
      case 'question':
        question = message.name;
        break;
      case 'answer':
        if (currentTurn === 1 && flippedTeam1.length === cupsPerTeam) {
          gameOver = true;
        }
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
      case 'game_started':
        isGameStarted = true;
        // Broadcast game state to all players
        socket.send(JSON.stringify({ type: 'game_state', teams: assignedTeams }));
        break;
      case 'restart':
        resetGame();
        break;
      case 'flip':
        if (message.team === 1) {
          flippedTeam1 = [...flippedTeam1, message.index];
        } else {
          flippedTeam2 = [...flippedTeam2, message.index];
        }
        break;
    }
  };

  const resetGame = () => {
    flippedTeam1 = [];
    flippedTeam2 = [];
    currentTurn = 1;
    gameOver = false;
    isGameStarted = false;
    players = [];
    assignedTeams = { team1: [], team2: [] };
    socket.send(JSON.stringify({ type: 'restart' }));
  };

  const joinGame = (playerName: string) => {
    const newPlayer = { name: playerName, team: null };
    players = [...players, newPlayer];

    // Send message to the server
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ type: "join", name: playerName }));
      console.log("Sent join message:", playerName);
    } else {
      console.error("WebSocket not connected!");
    }

    console.log("Current players:", players);
  };

  const assignTeams = () => {
    if (players.length < 2) return alert('Need at least 2 players to assign teams!');
    const shuffledPlayers = [...players];
    shuffledPlayers.sort(() => Math.random() - 0.5);
    assignedTeams.team1 = shuffledPlayers.slice(0, Math.ceil(shuffledPlayers.length / 2));
    assignedTeams.team2 = shuffledPlayers.slice(Math.ceil(shuffledPlayers.length / 2));
    console.log("Sending assign_teams message");
    socket.send(JSON.stringify({ type: 'assign_teams', team1: assignedTeams.team1, team2: assignedTeams.team2 }));
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

  const flipCup = (team: number, index: number) => {
    if (gameOver) return;
    if (
      (team === 1 && currentTurn === 1 && !flippedTeam1.includes(index)) ||
      (team === 2 && currentTurn === 2 && !flippedTeam2.includes(index))
    ) {
      if (team === 1) flippedTeam1 = [...flippedTeam1, index];
      if (team === 2) flippedTeam2 = [...flippedTeam2, index];
      checkGameOver();
      currentTurn = currentTurn === 1 ? 2 : 1;
      socket.send(JSON.stringify({ type: 'flip', team, index }));
    }
  };

  const checkGameOver = () => {
    if (flippedTeam1.length === cupsPerTeam || flippedTeam2.length === cupsPerTeam) {
      gameOver = true;
      socket.send(
        JSON.stringify({ type: 'winner', name: flippedTeam1.length === cupsPerTeam ? 'A-Team' : 'B-Team' })
      );
    }
  };

  const submitAnswer = () => {
    socket.send(JSON.stringify({ type: 'answer', answer: currentAnswer }));
    currentAnswer = ''; 
  };
</script>

<style>
  .lobby {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
  }
  .game {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
    font-family: sans-serif;
  }
  .team {
    display: flex;
    gap: 1rem;
  }
  .cup {
    width: 50px;
    height: 75px;
    background-image: url('/solo-cup.png');
    background-size: contain;
    background-repeat: no-repeat;
    cursor: pointer;
    transform: rotate(0deg);
    transition: transform 0.4s;
  }
  .cup.flipped {
    transform: rotate(180deg);
    filter: grayscale(80%) opacity(0.4);
  }
  .turn {
    font-weight: bold;
    font-size: 1.2rem;
  }
  button {
    padding: 0.5rem 1rem;
    font-size: 1rem;
    cursor: pointer;
  }

.table {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 2rem;
}

.row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.cell {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 40px;
  text-align: center;
}

.cell.name {
  font-weight: bold;
}

.divider {
  width: 20px;
}

</style>

<!-- Lobby Section -->
{#if !isGameStarted}
  <div class="lobby">
    <input type="text" placeholder="Enter name" bind:value={currentPlayerName} />
    <button on:click={() => joinGame(currentPlayerName)}>Join Game</button>
    {#if players.length >= 2}
      <button on:click={assignTeams}>Assign Teams</button>
    {/if}
    {#if assignedTeams.team1.length > 0 && assignedTeams.team2.length > 0}
      <button on:click={startGame}>Start Game</button>
    {/if}
  </div>
{/if}

<!-- Game Section -->
{#if isGameStarted}
<div class="table">
  {#each Array(Math.max(assignedTeams.team1.length, assignedTeams.team2.length)) as _, i}
    <div class="row">
      <!-- Team 1 -->
      <div class="cell name">{assignedTeams.team1[i] ?? ''}</div>
      <div
        class="cell cup {flippedTeam1.includes(i) ? 'flipped' : ''}"
        on:click={() => flipCup(1, i)}
        role="button"
        tabindex="0"
        aria-label="Flip cup for Team 1, player {assignedTeams.team1[i]}"
      ></div>

      <div class="divider"></div>

      <!-- Team 2 -->
      <div
        class="cell cup {flippedTeam2.includes(i) ? 'flipped' : ''}"
        on:click={() => flipCup(2, i)}
        role="button"
        tabindex="0"
        aria-label="Flip cup for Team 2, player {assignedTeams.team2[i]}"
      ></div>
      <div class="cell name">{assignedTeams.team2[i] ?? ''}</div>
    </div>
  {/each}
</div>
{/if}
