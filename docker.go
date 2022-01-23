package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JenswBE/docker-rest-puller-restarter/openapi"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

type DockerService struct {
	client *docker.Client
}

func NewDockerService() (*DockerService, error) {
	// Create Docker client
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	// Build Docker service
	return &DockerService{client: client}, nil
}

func (s *DockerService) PullRestartContainer(ctx context.Context, containerName string) error {
	// Find containers
	container, err := s.findContainer(ctx, containerName)
	if err != nil {
		return err
	}

	// Pull image
	imageRC, err := s.client.ImagePull(ctx, container., types.ImagePullOptions{})
	if err != nil {
		log.Warn().Err(err).Str("image", container.Image).Str("container", containerName).Msgf("Failed to pull image")
		err = fmt.Errorf("failed to pull image %s: %w", container.Image, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	imageRC.Close()

	// Check if image was updated
	timout := 5 * time.Minute
	s.client.ContainerCr
	err = s.client.ContainerRestart(ctx, container.ID, &timout)
	if err != nil {
		err = fmt.Errorf("failed to restart container %s: %w", containerName, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}

	// Pull-Restart successful
	return nil
}

func (s *DockerService) findContainer(ctx context.Context, containerName string) (types.ContainerJSON, error) {
	// List containers
	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		err = fmt.Errorf("failed to list Docker containers: %w", err)
		return types.Container{}, NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}

	// Find container
	dockerContainerName := "/" + containerName
	var container types.Container
	for _, cont := range containers {
		if sliceContains(cont.Names, dockerContainerName) {
			container = cont
			break
		}
	}
	if container.ID == "" {
		log.Debug().
			Interface("containers", containers).
			Str("container_name", containerName).
			Msg("Container with name not found in Docker host")
		return types.Container{}, NewError(http.StatusNotFound, openapi.APIERRORCODE_UNKNOWN_CONTAINER, containerName, nil)
	}

	// Inspect container
	return container, nil
}
