# FlipCup

A multiplayer quiz game where teams compete by flipping cups—by answering questions correctly. Inspired by Flip Cup, minus the beers. 🥤😎

## 🚀 Live Demo

Try it out here: https://flipcup.fly.dev

## 🎯 Purpose

This was a fun side project to explore Go and Svelte — my first functionaly project in either. It's still a work in progress, so be gentle with the feedback 😄. That said, all contributions and ideas are welcome!

## 🛠 Project Overview
```text
├── game-server                      // houses backend Go game server
│   ├── cmd
│   │   └── flipcup                     -- entrypoint
│   ├── internal
│   │   ├── game                        -- game-domain model/game-play files (game, team, player, etc)
│   │   ├── quiz                        -- quiz-domain models 
│   │   ├── transport                   -- http routing for rest endpoints and websocket handler for game play
│   │   │   ├── api
│   │   │   ├── types
│   │   │   └── ws
│   │   └── utils
│   ├── public                          -- for single container depoyments stores ui build
│   └── questions                       -- stores all question yaml files
└── ui                              // houses Svelte-kit frontend app
    ├── public
    └── src
        ├── assets
        │   └── fonts
        ├── components
        ├── lib
        │   ├── models
        │   ├── transport
        │   │   └── http
        │   ├── types
        │   └── utils
        └── styles

```


## 🛠 Local Development

To run the game locally using Docker Compose:

1. Clone the repo:

   git clone https://github.com/yourname/flip-cup.git  
   cd flip-cup

2. Update your `.env` file with your machine’s local IP address (needed for WebSocket connection):

```bash
cd ui
cat "VITE_WS_URL=<your-local-ip>:8080" > .env
```

###   Example:
``VITE_WS_URL=192.168.1.12:8080``

3. Start the app locally:

   docker-compose down && docker-compose build --no-cache && docker-compose up -d

4. Open your browser and go to:

   http://<your-local-ip>:5173 or 
   http://localhost:5173

## 🤖 AI-Assisted Development

Recent improvements to this project (UI redesign, E2E test suite) were built with the help of GitHub Copilot. A full breakdown of what was delegated, how decisions were made, and what changed lives here:

👉 **[docs/ai-approach.md](docs/ai-approach.md)**

## 🤝 Contributing
Got feedback, ideas, or issues? Open an issue or a pull request — would love to hear what you think! Lastly, for the record, no i was not in a frat. 
