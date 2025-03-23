package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

var Docker *client.Client

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate execution group ID and individual execution IDs
	groupID := uuid.New().String()
	executionIDs := make(map[string]string)

	response := SubmitResponse{
		ExecutionGroupID: groupID,
		Executions:       executionIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// http.HandleFunc("/submit", handleSubmit)

	// client, err := client.NewClientWithOpts(client.WithHost("host.docker.internal:2375"), client.WithAPIVersionNegotiation())
	// if err != nil {
	// 	log.Fatalf("Failed to create Docker client: %v", err)
	// }
	// Docker = client

	client, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	Docker = client

	// // Initialize poller background function.
	// go RunPollerLoop()

	// port := ":8080"
	// fmt.Printf("Starting executor service on port %s\n", port)
	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
	fmt.Println("Starting executor service")
	stdout, stderr, err := ExecuteCode("123", "python:3.9.21-slim", map[string]File{
		"main.py": {
			Encoding: "utf-8",
			Content:  `print("Hello world!", flush=True)`,
		},
		"runner.sh": {
			Encoding: "utf-8",
			Content: `#!/bin/sh
echo "==== Running main.py ===="
python3 -u main.py
echo "==== Status Code: $?"`,
		},
	}, []Test{
		// {
		// 	Input:  "Hello, World!",
		// 	Output: "Hello, World!",
		// },
	}, Configuration{})
	fmt.Println("Stdout: ", stdout)
	fmt.Println("Stderr: ", stderr)
	fmt.Println("Error: ", err)
}

func RunPollerLoop() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		runExecutePoller()
	}
}

func runExecutePoller() {

}
