package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "silentmode/proto"

	"google.golang.org/grpc"
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "server:50051"
	}

	var conn *grpc.ClientConn
	var err error

	log.Printf("[%s] Starting connection attempts to %s...", clientID, serverAddr)

	for {
		// Create a context with a 5-second timeout for the dial attempt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// grpc.DialContext with WithBlock will wait until the connection is ready or the context expires
		conn, err = grpc.DialContext(ctx, serverAddr,
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)

		if err == nil {
			cancel()
			break
		}

		log.Printf("[%s] Connection failed (%v). Retrying in 2s...", clientID, err)
		cancel()
		time.Sleep(2 * time.Second)
	}

	defer conn.Close()
	log.Printf("[%s] Connected to Server!", clientID)

	client := pb.NewFileServiceClient(conn)

	// Open the stream
	stream, err := client.ConnectAndListen(context.Background())
	if err != nil {
		log.Fatalf("Fatal stream error: %v", err)
	}

	// Send the first ID message
	stream.Send(&pb.FilePayload{ClientId: clientID})
	log.Printf("[%s] Connected and listening!", clientID)
	for {
		cmd, err := stream.Recv()
		if err != nil {
			log.Printf("[%s] Connection lost: %v", clientID, err)
			break
		}

		log.Printf("[%s] Server requested: %s. Uploading...", clientID, cmd.Filename)

		// 3. Open and stream the file
		filePath := filepath.Join("client-data", "data.txt")
		file, err := os.Open(filePath)
		if err != nil {
			// Create a dummy if missing
			os.MkdirAll("client-data", 0755)
			os.WriteFile(filePath, []byte("Hello from "+clientID), 0644)
			file, _ = os.Open(filePath)
		}

		buf := make([]byte, 64*1024)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			// Send chunk with our unique ID
			stream.Send(&pb.FilePayload{
				ClientId: clientID,
				Data:     buf[:n],
			})
		}
		file.Close()
		stream.Send(&pb.FilePayload{ClientId: clientID, IsComplete: true})
		log.Printf("[%s] Upload complete", clientID)
	}
}
