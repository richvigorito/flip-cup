// src/lib/transport/http/Games.ts

import { RawGameState, GameState } from '$lib/models/GameState';

import { baseHttpUrl } from '$lib/utils/config';

export async function fetchGames(): Promise<Games[]>  {
    try {
        // there isnt a neeed for in active games
        // but api defined as :
        //      games/{?$activeOrInactive}
        //      
        const httpUrl = `${baseHttpUrl}/games/active`;

        const res = await fetch(httpUrl);
        const rawGames: RawGameState[] = await res.json();
        const games = rawGames.map(g => new GameState(g));
        console.log('✅ Fetched games:', games);
        return  games;
    } catch (err) {
        console.error('Failed to fetch games:', err);
        return [];
    } 
    /***
    finally {
        omit for now
        loadingGames.set(false);
    }
    */
};

