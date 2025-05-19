import { get, writable } from 'svelte/store';
import { GameState } from '$lib/models/GameState';
import { Team } from '$lib/models/Team';
import { Player } from '$lib/models/Player';
import { gameState, myTeam, me, mode, joined, eventLog, currentQuestion, winner, gamesCompleted } from '$lib/store';


import { baseWsUrl } from '$lib/utils/config';

export const socket = writable<WebSocket | null>(null);

let ws: WebSocket;

export function connectSocket(initmsg: object) {

    ws = new WebSocket(baseWsUrl);

    socket.set(ws);
    
    ws.onopen = () => {
      console.log('Connected to WebSocket');
      console.log(JSON.stringify(initmsg));
      ws.send(JSON.stringify(initmsg));
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      handleMessage(message);
    };

    ws.onclose = () => console.log('Disconnected');
    ws.onerror = (err) => console.error('WebSocket error', err);
}

function handleMessage(message: any) {

 let newState;
 //let newTeam;
 let currentPlayer;

  console.log("↙️ handle incomin message:", JSON.stringify(message));
  switch (message.type) {

     case 'game_player_initialized':
        currentPlayer = new Player(message.payload);
        me.set(currentPlayer);
        break;

    case 'player_joined':
        handlePlayerJoined(message);
        break;
/*
    case 'assign_teams':
      assignedTeams.set({ team1: message.team1, team2: message.team2 });
      break;
*/
    case 'my_current_team':
        handleMyTeamAssignment(message);
        break;

    case 'teams_assigned':
        handleTeamAssignments(message);
        break;

/*
    case 'game_snapshot':
        newState = new GameState(message.payload.game_snapshot);
    
        // You can call logic here:
    console.log('Current State:', gameState);
    console.log('New State:', newState);
    
    gameState.set(newState);

    eventLog.update((log) => [
        ...log,
        { message: `teams assigned`, type: 'success' },
        { message: `team1: ${message.payload.teamA.players.join(', ')}`, type: 'info' },
        { message: `team2: ${message.payload.teamB.players.join(', ')}`, type: 'info' },
      ]);

      eventLog.update((log) => [
        ...log,
        { message: `teams assigned`, type: 'success' },
        { message: `team1: ${message.answer.teamA.players.join(', ')}`, type: 'info' },
        { message: `team2: ${message.answer.teamB.players.join(', ')}`, type: 'info' },
      ]);

      break;
*/

    case 'game_started':
      handleGameStarted(message)
      break;

    case 'question':
        handleAdministerQuestion(message)
        break;

    case 'answered_correctly':
        handleQuestionAnsweredCorrectly(message)
        break;

    case 'incorrect_answer':
      //eventLog.update((log) => [...log, { message: `${message.name} submitted wrong answer`, type: 'error' }]);
      logEvent(`${message.name} submitted wrong answer`, 'error');
      break;

    case 'winner':
      winner.set(message.name);
      var tmp = 1 + get(gamesCompleted);
      gamesCompleted.set(tmp);
      //eventLog.update((log) => [...log, { message: `🏆 Winner: ${message.name}`, type: 'success' }]);
      logEvent(`🏆 Winner: ${message.name}`,'success');
      break;

    case 'game_restarted':
        handleGameRestarted(message)    
        break;
  }
}

export const send = (msg: object) => {
  console.log("↗️ send outbound message:", JSON.stringify(msg));
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(msg));
  } else {
    console.warn('⚠️ WebSocket not connected.');
  }
};

const logEvent = (message: string, type: 'success' | 'error' | 'info') => {
    eventLog.update((log) => [...log, { message, type }]);
};

const handleGameRestarted = (message: any) => {
    winner.set(null);
    const newState = new GameState(message.payload.game_snapshot);
    gameState.set(newState);
    mode.set('lobby');
    logEvent(`Game Restarted.`, 'info');
};

const handleGameStarted = (message: any) => {
    mode.set('game');
    const newState = new GameState(message.payload.game_snapshot);
    gameState.set(newState);
    logEvent(`Game Started`, 'success');
};

const handleMyTeamAssignment = (message: any) => {
    const newTeam = new Team(message.payload);
    myTeam.set(newTeam);
    logEvent(`joined team: ${message.payload.name}`, 'info');
};

const handlePlayerJoined = (message: any) => {
//    console.log('handle player_joined');
    const currentPlayer = get(me);
    if (currentPlayer) {
        currentPlayer.name = message.name;
        me.set(currentPlayer); 
    } else {
        console.warn('⚠️ No player instance to update');
    }
    logEvent(`${message.name} joined the game`, 'success');
        
};

const handleTeamAssignments = (message: any) => {
    const newState = new GameState(message.payload.game_snapshot);
    gameState.set(newState);
    logEvent(`teams assigned`, 'success');
    logEvent(`${newState.teamA.name}: ${newState.teamA.players.map(p => p.name).join(', ')}`, 'info');
    logEvent(`${newState.teamB.name}: ${newState.teamA.players.map(p => p.name).join(', ')}`, 'info');
};


const handleAdministerQuestion = (message: any) => {
    const currentPlayer = get(me);
    currentPlayer.isMyTurn = true;
    me.set(currentPlayer);
    currentQuestion.set(message.name);
};

const handleQuestionAnsweredCorrectly = (message: any) => {
    // if $me answered correctly, then $me needs to relinquish the turn
    if(message.payload.action_performed_by.id == get(me).id){
        const currentPlayer = get(me);
        currentPlayer.isMyTurn = false;
        me.set(currentPlayer);
    }

    const newState = new GameState(message.payload.game_snapshot);
    gameState.set(newState);
};

