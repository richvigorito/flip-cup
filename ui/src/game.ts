// game.ts
type Player = { name: string; team: number | null };
type GameState = {
  players: Player[];
  teams: [Player[], Player[]];
  cups: { [team: number]: boolean[] };
  currentTurn: number;
  winner: string | null;
  question: string | null;
};

//const socket = new WebSocket('ws://localhost:3000');
const socket = new WebSocket('ws://localhost:8080/ws');
const state: GameState = {
  players: [],
  teams: [[], []],
  cups: { 1: [], 2: [] },
  currentTurn: 1,
  winner: null,
  question: null,
};

const listeners: Function[] = [];

socket.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  switch (msg.type) {
    case 'players':
      state.players = msg.players;
      break;
    case 'teams_assigned':
      state.teams = msg.teams;
      state.cups = {
        1: new Array(msg.teams[0].length).fill(false),
        2: new Array(msg.teams[1].length).fill(false),
      };
      break;
    case 'start':
      state.winner = null;
      break;
    case 'question':
      state.question = msg.name;
      break;
    case 'flip':
      state.cups[msg.team][msg.index] = true;
      state.currentTurn = msg.nextTurn;
      break;
    case 'winner':
      state.winner = msg.name;
      break;
  }
  listeners.forEach((fn) => fn());
};

export function subscribe(fn: () => void) {
  listeners.push(fn);
  return () => listeners.splice(listeners.indexOf(fn), 1);
}

export function getState() {
  return state;
}

export function send(msg: any) {
  socket.send(JSON.stringify(msg));
}

