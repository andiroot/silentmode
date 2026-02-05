package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	pb "silentmode/proto"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type cloudServer struct {
	pb.UnimplementedFileServiceServer
	// Map to store active client streams
	mu      sync.Mutex
	clients map[string]pb.FileService_ConnectAndListenServer
}

func (s *cloudServer) ConnectAndListen(stream pb.FileService_ConnectAndListenServer) error {
	// 1. Get initial identity
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	clientID := firstMsg.ClientId
	s.mu.Lock()
	s.clients[clientID] = stream
	s.mu.Unlock()

	log.Printf("[%s] is now online", clientID)

	// Variables to track the current active download
	var currentFile *os.File

	for {
		payload, err := stream.Recv()
		if err != nil {
			if currentFile != nil {
				currentFile.Close()
			}
			break
		}

		// 2. If we receive data, handle the file writing
		if len(payload.Data) > 0 {
			if currentFile == nil {
				// Start of a new unique download session
				dir := fmt.Sprintf("server-data/%s", clientID)
				os.MkdirAll(dir, 0755)

				// We create a temporary "active" file
				tempPath := fmt.Sprintf("%s/transferring.tmp", dir)
				currentFile, err = os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					log.Printf("Error creating temp file: %v", err)
					continue
				}
			}
			currentFile.Write(payload.Data)
		}

		// 3. When upload is finished, make it unique
		if payload.IsComplete && currentFile != nil {
			currentFile.Close()
			currentFile = nil // Reset for next time

			// Rename to a unique timestamped filename
			timestamp := time.Now().Format("20060102-150405")
			finalPath := fmt.Sprintf("server-data/%s/download-%s.txt", clientID, timestamp)

			oldPath := fmt.Sprintf("server-data/%s/transferring.tmp", clientID)
			os.Rename(oldPath, finalPath)

			log.Printf("[%s] Saved unique file: %s", clientID, finalPath)
		}
	}

	s.mu.Lock()
	delete(s.clients, clientID)
	s.mu.Unlock()
	return nil
}

func main() {
	srv := &cloudServer{clients: make(map[string]pb.FileService_ConnectAndListenServer)}

	// 1. Start gRPC Server (Clients connect here)
	go func() {
		lis, _ := net.Listen("tcp", ":50051")
		g := grpc.NewServer()
		pb.RegisterFileServiceServer(g, srv)
		log.Println("gRPC Cloud Server on :50051")
		g.Serve(lis)
	}()

	// 2. Start HTTP Server (Trigger via: curl -X POST http://localhost:8080/trigger?id=client1)
	http.HandleFunc("/trigger", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		srv.mu.Lock()
		stream, ok := srv.clients[id]
		srv.mu.Unlock()

		if !ok {
			http.Error(w, "Client not connected", 404)
			return
		}

		// Send command TO client
		stream.Send(&pb.DownloadCommand{Filename: "data.txt"})
		fmt.Fprintf(w, "Triggered download for %s", id)
	})

	log.Println("API Trigger available on :8080")
	http.ListenAndServe(":8080", nil)
}
