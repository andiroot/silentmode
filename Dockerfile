FROM golang:1.25-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ===== SERVER =====
FROM base AS server
RUN go build -o server ./server

# ===== CLIENT =====
FROM base AS client
RUN go build -o client ./client
