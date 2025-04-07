package docker

import (
	"log"

	"github.com/docker/docker/client"
)

// The singleton for Docker Client.
var sDockerClient *client.Client = nil

func Get() *client.Client {
	if sDockerClient == nil {
		initDockerClient()
	}

	return sDockerClient
}

func initDockerClient() {
	cl, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	sDockerClient = cl
}
