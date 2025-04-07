package docker

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

// HasImage checks if a Docker image exists locally
func HasImage(ctx context.Context, imageName string) bool {
	filters := filters.NewArgs()
	filters.Add("reference", imageName)

	images, err := sDockerClient.ImageList(ctx, image.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return false
	}

	return len(images) > 0
}

// PullImage pulls a Docker image if it doesn't exist locally
func PullImage(ctx context.Context, imageName string) error {
	if HasImage(ctx, imageName) {
		return nil
	}

	_, err := sDockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	return err
}
