# FlipQuiz 🎉

FlipQuiz is a real-time team quiz game inspired by the college party classic **Flip Cup** — but instead of just flipping cups, you have to **answer questions correctly** before flipping!

## What is this?

At its core, FlipQuiz is a mix of:
- Team-based trivia
- Fast-paced racing (via "flip" mechanics)

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
  - Lobby screen functional
  - Team assignment functional
  - Game table (flip cup style) layout work-in-progress
  - Cup flipping / game flow is the current focus!

---
## What's next?

- Full game play loop working end-to-end on the frontend
- Polish the UI / UX for mobile and desktop
- Add sound effects, timers, and animations 🎵⏳✨
- Deploy a dockerized, hosted version for people to play with friends

---
> PRs, ideas, and feedback always welcome — this is a learning project, so anything that makes it cooler is a win. Cheers! 🍻

> Lastly, for the record, no i was not in a frat. 
