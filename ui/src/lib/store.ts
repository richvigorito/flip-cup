// src/lib/store.ts
import { writable, derived } from 'svelte/store';
import { GameState } from '$lib/models/GameState';
import { Team } from '$lib/models/Team';
import { Player } from '$lib/models/Player';

import type { QuestionSet } from '$lib/types/QuestionSet'; 
export const questionSets = writable<QuestionSet[]>([]);


export const mode = writable<'welcome' | 'new' | 'join' | 'lobby' | 'game'>('welcome');
export const gameId = writable<string | null>(null);

export const teamId = writable<string | null>(null);
export const playerId = writable<string | null>(null);

export const playerName = writable<string | null>(null);
export const joined = writable<boolean>(false);
export const availableGames = writable<any[]>([]);
export const loadingGames = writable<boolean>(false);
export const currentPlayerName = writable<string | null >(null);

export const currentQuestion = writable<string | null >(null);
export const winner = writable<string | null >(null);
export const gamesCompleted = writable<number | null >(0);

export const eventLog = writable<{ message: string; type: 'info' | 'success' | 'error' }[]>([]);

export const gameState = writable<GameState | null >(null);
export const me = writable<Player | null >(null);

export function resetClientGameState() {
  gameId.set(null);
  joined.set(false);
  currentPlayerName.set(null);
  currentQuestion.set(null);
  winner.set(null);
  gamesCompleted.set(0);
  eventLog.set([]);
  gameState.set(null);
  me.set(null);
}

// Derived store to automatically determine myTeam based on gameState and me
export const myTeam = derived([gameState, me], ([$gameState, $me]) => {
  if (!$gameState || !$me) return null;
  
  const inTeamA = $gameState.teamA.players.some(p => p.id === $me.id);
  if (inTeamA) return $gameState.teamA;
  
  const inTeamB = $gameState.teamB.players.some(p => p.id === $me.id);
  if (inTeamB) return $gameState.teamB;
  
  return null;
});
