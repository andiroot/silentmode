#!/bin/bash

# Define your module prefix (change this to your actual project name)
MODULE_PREFIX="silentmode"

echo "generating clients file..."

mkdir -p client-data server-data
go run utils/generate.go

echo "Starting Go Workspace Setup..."

# 1. Initialize Proto if needed
echo "Tidying [proto]..."
cd proto
if [ ! -f go.mod ]; then
    go mod init ${MODULE_PREFIX}/proto
fi
go mod tidy
cd ..

# 2. Setup Server
echo "Tidying [server]..."
cd server
if [ ! -f go.mod ]; then
    go mod init ${MODULE_PREFIX}/server
fi
go mod edit -replace ${MODULE_PREFIX}/proto=../proto
go mod tidy
cd ..

# 3. Setup Client
echo "Tidying [client]..."
cd client
if [ ! -f go.mod ]; then
    go mod init ${MODULE_PREFIX}/client
fi
go mod edit -replace ${MODULE_PREFIX}/proto=../proto
go mod tidy
cd ..

# 4. Optional: Create a Go Workspace for VSCode/IDE support
echo "🛠️  Creating Go Workspace..."
go work init ./proto ./server ./client 2>/dev/null || go work use ./proto ./server ./client

echo "✅ All modules are tidied and linked!"
echo "🐳 You can now run: docker compose up --build"