import type { Player } from './Player';

export interface RawTeam {
  name: string;
  turn: number;
  players: Player[];
  color: string;
}

export class Team {
  name: string;
  turn: number;
  players: Player[];
  color: string;

  constructor(data: RawTeam) {
    this.name = data.name;
    this.turn = data.turn;
    this.players = data.players;
    this.color = Team.setTableColor(this.name);
  }

  // 
  // maybe kinda weird but i didnt want have to track color state
  // so when a team is created we just dynamically  generate a 
  // brown-ish color based on teamname
  // 
  //  (this will represent their table color) 
  private static setTableColor(input: string): string {
    let hash = 5381;
    for (let i = 0; i < input.length; i++) {
        hash = ((hash << 5) + hash) + input.charCodeAt(i); // hash * 33 + c
    }

    // Use the hash to vary hue, sat, light
    const hue = 20 + (Math.abs(hash) % 20);         // 20–40° (brown)
    const saturation = 40 + (Math.abs(hash >> 3) % 30); // 40–70%
    const lightness = 30 + (Math.abs(hash >> 5) % 30); // 30–60%

    return Team.hslToHex(hue, saturation, lightness); 
  }

  private static hslToHex(h: number, s: number, l: number): string {
    s /= 100;
    l /= 100;
    const k = (n: number) => (n + h / 30) % 12;
    const a = s * Math.min(l, 1 - l);
    const f = (n: number) =>
      Math.round(255 * (l - a * Math.max(-1, Math.min(k(n) - 3, Math.min(9 - k(n), 1)))));
    return `#${[f(0), f(8), f(4)].map(x => x.toString(16).padStart(2, '0')).join('')}`;
  }
}



