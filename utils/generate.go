package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. 
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi 
ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit 
in voluptate velit esse cillum dolore eu fugiat nulla pariatur. 
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia 
deserunt mollit anim id est laborum.
`

func generateLoremFile(path string, sizeMB int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	targetSize := int64(sizeMB) * 1024 * 1024
	var written int64

	for written < targetSize {
		n, err := f.WriteString(lorem + "\n")
		if err != nil {
			return err
		}
		written += int64(n)
	}

	fmt.Printf("Generated %s (~%d MB)\n", path, sizeMB)
	return nil
}

func main() {
	// Create directories for 5 clients
	for i := 1; i <= 5; i++ {
		clientDir := fmt.Sprintf("client-data/client%d", i)
		os.MkdirAll(clientDir, 0755)

		// Generate unique file for each client
		filePath := filepath.Join(clientDir, "data.txt")
		err := generateLoremFile(filePath, 100) // 100MB files
		if err != nil {
			log.Printf("Error generating file for client %d: %v", i, err)
		} else {
			// Add unique identifier to each file
			f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("\n=== This file belongs to CLIENT-%d ===\n", i))
				f.Close()
			}
		}
	}

	// Create server directory
	os.MkdirAll("server-data", 0755)

	fmt.Println("Data generation complete!")
	fmt.Println("Generated files for 5 clients in client-data/ directory")
}
