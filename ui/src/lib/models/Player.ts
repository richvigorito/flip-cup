export interface RawPlayer {
  id: string;
  name: string;
  isMyTurn: boolean;
}

export class Player {
  id: string;
  name: string;
  isMyTurn: boolean;
 
  constructor(data: RawPlayer | {player_id: string; name: string }) {
    this.id = 'player_id' in data ? data.player_id : data.id;
    this.name = data.name;
    this.isMyTurn = false;
  }
}

