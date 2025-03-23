package main

import (
	"encoding/json"
	"flag"
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
	sampleRun := flag.Bool("sample-run", false, "Run a sample code execution test")
	flag.Parse()

	client, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	Docker = client

	if *sampleRun {
		fmt.Println("Running sample code execution test...")
		err = ExecuteCode("123", "python:3.9.21-slim", map[string]File{
			"main.py": {
				Encoding: "utf-8",
				Content: `
a = input()
print(a)
`,
			},
			"runner.sh": {
				Encoding: "utf-8",
				Content: `#!/bin/sh
input="{{ .Input }}"
echo "{{ .Boundary }}"
echo "$input" | python3 -u main.py
echo "{{ .Boundary }}"
echo "::Status Code:: $?"`,
			},
		}, []Test{
			{
				Input:  "Hello, World!",
				Output: "Hello, World!",
			},
			{
				Input:  "Heyyyyyy",
				Output: "Heyyyyyy",
			},
		}, Configuration{})
		fmt.Println("Error: ", err)
		return
	}

	http.HandleFunc("/submit", handleSubmit)

	// Initialize poller background function.
	go RunPollerLoop()

	port := ":8080"
	fmt.Printf("Starting executor service on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
