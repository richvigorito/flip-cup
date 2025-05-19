# FlipCup

A multiplayer quiz game where teams compete by flipping cups—by answering questions correctly. Inspired by Flip Cup, minus the beers. 🧠🥤

## 🚀 Live Demo

Try it out here: https://flipcup.fly.dev

## 🎯 Purpose

This was a fun side project to explore Go and Svelte — my first functionaly project in either. It's still a work in progress, so be gentle with the feedback 😄. That said, all contributions and ideas are welcome!

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

## 🤝 Contributing
Got feedback, ideas, or issues? Open an issue or a pull request — would love to hear what you think! Lastly, for the record, no i was not in a frat. 
