# ---------- UI BUILD ----------
FROM node:22-alpine AS ui-builder
WORKDIR /ui

COPY ui/package*.json ./
RUN npm ci

COPY ui ./
RUN npm run build

# ---------- GO BUILD ----------
FROM golang:1.24-alpine AS go-builder
WORKDIR /app

COPY game-server/go.mod game-server/go.sum ./
RUN go mod download

COPY game-server ./
RUN CGO_ENABLED=0 GOOS=linux go build -o game-server ./cmd/flipcup

# ---------- FINAL IMAGE ----------
FROM alpine:3.21
WORKDIR /app

COPY --from=go-builder /app/game-server ./game-server
COPY --from=go-builder /app/questions ./questions
COPY --from=ui-builder /ui/dist ./public

ENV PORT=8080
EXPOSE 8080
CMD ["./game-server"]
