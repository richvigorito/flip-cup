# FlipCup

A multiplayer quiz game where teams compete by flipping cups—by answering questions correctly. Inspired by Flip Cup, minus the beers. 🥤😎

## 🚀 Live Demo

Try it out here: https://flipcup.fly.dev

## 🎯 Purpose

This was a fun side project to explore Go and Svelte — my first functional project in either. It's still a work in progress, so be gentle with the feedback 😄. That said, all contributions and ideas are welcome!

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
│   ├── public                          -- for single container deployments stores ui build
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


## 🤖 AI-Assisted Development

Recent improvements to this project (UI redesign, E2E test suite) were built with the help of GitHub Copilot. A full breakdown of what was delegated, how decisions were made, and what changed lives here:

👉 **[docs/ai-approach.md](docs/ai-approach.md)**

## 🧪 Testing

We use [Playwright](https://playwright.dev/) for End-to-End (E2E) testing. The tests simulate real user interactions (creating games, joining, answering questions) to ensure the full stack works together.

GitHub Actions now runs the test pipeline on every push and pull request. The first-pass required check for protecting `main` is `tests`, and `main` should be configured to only accept changes through pull requests after that check passes.

### Prerequisites

Playwright now starts the backend and frontend automatically when you run the E2E suite, so the manual startup flow below is mainly useful for interactive local development.

If you want to run the app by hand, ensure both the backend and frontend are running locally:

1.  **Start the Backend**:
    ```bash
    cd game-server
    go run cmd/flipcup/main.go
    ```
2.  **Start the Frontend**:
    ```bash
    cd ui
    npm install
    npm run dev
    ```

### Running Tests

Run all tests:
```bash
cd ui
npx playwright test
```

Run a specific test file (e.g., disconnection logic):
```bash
npx playwright test e2e/disconnect.spec.ts
```

Run tests with a visual UI (great for debugging):
```bash
npx playwright test --ui
```

### Writing Tests

Tests are located in `ui/e2e/`. To add a new test:
1.  Create a file like `ui/e2e/my-feature.spec.ts`.
2.  Import `test` and `expect` from `@playwright/test`.
3.  Use `page` fixtures to interact with the game.

**Note**: We use `sessionStorage` for player state, allowing you to simulate multiple players in a single browser instance by using different `BrowserContext`s. See `e2e/disconnect.spec.ts` for an example of multi-player testing.

## 🤝 Contributing
Got feedback, ideas, or issues? Open an issue or a pull request — would love to hear what you think! Lastly, for the record, no i was not in a frat. 
