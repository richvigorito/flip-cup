#!/bin/bash

echo '{"type":"join", "name":"Player 1"}' | websocat ws://localhost:8080/ws

