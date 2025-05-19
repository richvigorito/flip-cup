// src/lib/store.ts
import { writable } from 'svelte/store';
import { GameState } from '$lib/models/GameState';
import { Team } from '$lib/models/Team';
import { Player } from '$lib/models/Team';

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
export const gamesCompleted = writable<int | null >(0);

export const eventLog = writable<{ message: string; type: 'info' | 'success' | 'error' }[]>([]);

export const gameState = writable<GameState | null >(null);
export const myTeam = writable<Team | null >(null);
export const me = writable<Player | null >(null);
