FROM golang:1.25-alpine AS builder
WORKDIR /app

# 1. Copy mod files from root context
COPY server/go.mod server/go.sum ./server/
COPY proto/go.mod proto/go.sum ./proto/

# 2. Download dependencies 
# We cd into server because that's where the main go.mod is
WORKDIR /app/server
RUN go mod download

# 3. Copy everything from root
WORKDIR /app
COPY . .

# 4. Build the server
WORKDIR /app/server
RUN go build -o main .

# Stage 2: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server/main .
RUN mkdir -p server-data
CMD ["./main"]