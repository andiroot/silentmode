FROM golang:1.25-alpine AS builder
WORKDIR /app

# 1. Copy mod files from root context
COPY client/go.mod client/go.sum ./client/
COPY proto/go.mod proto/go.sum ./proto/

# 2. Download dependencies
WORKDIR /app/client
RUN go mod download

# 3. Copy everything
WORKDIR /app
COPY . .

# 4. Build the client
WORKDIR /app/client
RUN go build -o main .

# Stage 2: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/client/main .
CMD ["./main"]