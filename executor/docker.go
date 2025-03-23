package main

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

func HasImageInstalled(ctx context.Context, imageName string) bool {
	filters := filters.NewArgs()
	filters.Add("reference", imageName)

	images, err := Docker.ImageList(ctx, image.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		log.Fatalf("Failed to list images: %v", err)
	}

	return len(images) > 0
}
