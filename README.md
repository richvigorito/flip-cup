# FlipCup 🎉

This FlipCup game is a real-time team quiz game inspired by the college party classic **Flip Cup** — but instead of just flipping cups, you have to **answer questions correctly** before flipping!

## What is this?
Team-based flip-cup

Players join a lobby, get assigned to teams, and race through a set of quiz questions. Answer right → your cup flip's → pass the turn. First team to finish wins!

The app has two main parts:
- A **game server** (written in Go) handling players, lobbies, teams, and real-time communication
- A **web UI** (built with Svelte) managing lobby creation, team assignments, and (soon) full game play

---
## Why am I building this?
- I'm currently unemployed and wanted to dive into a **fun, challenging project** 🎯
- I wanted an excuse to **learn Go** 
- I wanted to **level up** my Node.js + websockets knowledge

---
## Current Progress 🚀
- ✅ **Game server:**  
  - Built in Go
  - Ingests and stores quiz questions
  - Accepts players into a lobby
  - Randomly assigns teams
  - Transmits questions to the correct player
  - Accepts player answers
  - Tracks team progress
  - Determines and announces the winning team
- 🛠️ **Frontend (Svelte):**
  - Game Log
  - Lobby screen functional
  - Team assignment functional
  - Game table (flip cup style) layout work-in-progress
  - Cups flips and next question/player is broadcasted

---
## To Play Local
### edit env 
add ``VITE_WS_URL=<your_ip>:8080`` to ``ui/.env``
### run docker
rootdir> docker-compose build --no-cache
rootdir> docker-compose up -d

### run servers
```bash
cd game-server
go run main.go
cd ../ui
npm run dev
```
### game play
find an equal number of friends, go to your ui
1) each person joins game
1) after everyone joins 'assign' teams
1) start game
1) the ui defintely done by a BE engineer but should be obvious whose turn it is
1) (as of right now, restart functionality doesnt work, need to restart go server)


---
## What's next?

- Better error handling with websocket disconnects
- Cleanup up game/team assignment
- Full game play loop working end-to-end on the frontend
-- restart needs to clear room
- Add abilty to create multiple games with different questions
- Polish the UI / UX 
-- Add sound effects, timers, and animations 🎵⏳✨
- Dockerize
- Possibley deploy, hosted version for people to play with friends

---
> PRs, ideas, and feedback always welcome — this is a learning project, so anything that makes it cooler is a win. Cheers! 🍻
> Lastly, for the record, no i was not in a frat. 
