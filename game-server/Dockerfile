# ---------- DEVELOPMENT STAGE ----------
FROM golang:1.24 AS dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "cmd/flipcup/main.go"]

# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o game-server ./cmd/flipcup

# ---------- PRODUCTION STAGE ----------
FROM alpine:latest AS prod
WORKDIR /app
COPY --from=builder /app/game-server .
EXPOSE 8080
CMD ["./game-server"]

