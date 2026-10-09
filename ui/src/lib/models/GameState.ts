// lib/models/GameState.ts
import type { Team } from './Team';
import type { Player } from './Player';

export interface RawGameState {
  id: string;
  teamA: Team;
  teamB: Team;
  active: boolean;
  quizfile: string;
  questionsPerPlayer?: number;
}

export class GameState {
  id:       string;
  teamA:    Team;
  teamB:    Team;
  active:   boolean;
  quizfile: string;
  questionsPerPlayer: number;

  constructor(data: RawGameState) {
    this.id         = data.id;
    this.teamA      = data.teamA;
    this.teamB      = data.teamB;
    this.active     = data.active;
    this.quizfile   = data.quizfile;
    this.questionsPerPlayer = data.questionsPerPlayer ?? 1;
  }

  get currentTurn(): 'teamA' | 'teamB' {
    if (this.teamA.turn === this.teamB.turn) {
      return 'teamA';
    }
    return this.teamA.turn < this.teamB.turn ? 'teamA' : 'teamB';
  }

  canStart(): boolean {
    return this.teamA.players.length > 0 && this.teamB.players.length > 0;
  }
  
  get allPlayers(): Player[] {
    return [...this.teamA.players, ...this.teamB.players];
  }

  get isGameActive(): boolean {
    return this.active;
  }
}
