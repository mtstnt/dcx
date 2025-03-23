package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func ExecuteCode(executionGroupID string, imageName string, files map[string]File, tests []Test, configurations Configuration) (string, string, error) {
	if !HasImageInstalled(context.Background(), imageName) {
		return "", "", fmt.Errorf("image %s not found", imageName)
	}

	fmt.Println("Creating container")

	containerCreateResponse, err := Docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:           imageName,
			AttachStdin:     true,
			AttachStdout:    true,
			AttachStderr:    true,
			WorkingDir:      "/code",
			Cmd:             []string{"sh", "runner.sh"},
			NetworkDisabled: true, // By default, we don't want any networking.
		},
		&container.HostConfig{
			AutoRemove: true,
			Privileged: true,
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		executionGroupID,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create container: %v", err)
	}

	containerID := containerCreateResponse.ID

	tarArchive, err := CreateTarArchiveFromJSON(files)
	if err != nil {
		return "", "", fmt.Errorf("failed to create tar archive: %v", err)
	}

	if err := Docker.CopyToContainer(context.Background(), containerID, "/", tarArchive, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}); err != nil {
		return "", "", fmt.Errorf("failed to copy files to container: %v", err)
	}

	fmt.Println("Starting container")
	if err := Docker.ContainerStart(context.Background(), containerID, container.StartOptions{}); err != nil {
		return "", "", fmt.Errorf("failed to start container: %v", err)
	}

	logs, err := Docker.ContainerLogs(context.Background(), containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer logs.Close()

	// Copy the logs to buffers
	var (
		bufStdout = bytes.NewBuffer(nil)
		bufStderr = bytes.NewBuffer(nil)
	)

	if _, err := stdcopy.StdCopy(bufStdout, bufStderr, logs); err != nil {
		return "", "", fmt.Errorf("failed to copy stdout: %v", err)
	}

	fmt.Println("Waiting for container")
	waitResponseCh, errCh := Docker.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		return "", "", fmt.Errorf("failed to wait for container 1: %v", err)
	case waitResponse := <-waitResponseCh:
		if waitResponse.Error != nil {
			return "", "", fmt.Errorf("failed to wait for container 2: %v", waitResponse.Error)
		}
	}

	return bufStdout.String(), bufStderr.String(), nil
}

// CreateTarArchiveFromJSON creates a TAR archive from a JSON representation of files
// and returns it as an io.ReadCloser for use with Docker's CopyToContainer
func CreateTarArchiveFromJSON(files map[string]File) (io.ReadCloser, error) {
	// Create a new TAR writer
	var buffer bytes.Buffer
	{
		tw := tar.NewWriter(&buffer)
		defer tw.Close()

		// Create a base directory for all files
		const baseDir = "/code/"

		// Write the files to the TAR archive in a separate goroutine
		for filename, file := range files {
			// Create a header for the file
			header := &tar.Header{
				Name: baseDir + filename,
				Mode: 0777,
				Size: int64(len(file.Content)),
			}

			// Write the header to the TAR archive
			if err := tw.WriteHeader(header); err != nil {
				return nil, fmt.Errorf("failed to write tar header: %v", err)
			}

			content := []byte(file.Content)

			// Write the content to the TAR archive
			if _, err := tw.Write(content); err != nil {
				return nil, fmt.Errorf("failed to write file content: %v", err)
			}
		}
	}

	return io.NopCloser(bytes.NewReader(buffer.Bytes())), nil
}
