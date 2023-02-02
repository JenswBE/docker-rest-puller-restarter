package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/JenswBE/docker-rest-puller-restarter/openapi"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DockerService struct {
	client *docker.Client
}

func NewDockerService() (*DockerService, error) {
	// Create Docker client
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create new Docker client: %w", err)
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
	imageName := container.Config.Image
	log.Debug().Str("image", imageName).Msg("Pulling image ...")
	imagePullLogs, err := s.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		log.Warn().Err(err).Str("image", imageName).Str("container", containerName).Msgf("Failed to pull image")
		err = fmt.Errorf("failed to pull image %s: %w", imageName, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	defer imagePullLogs.Close()
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		logs, err := io.ReadAll(imagePullLogs)
		if err != nil {
			log.Warn().Err(err).Str("image", imageName).Msg("Latest image pulled, but failed to read logs")
		} else {
			log.Debug().Bytes("logs", logs).Str("image", imageName).Msg("Latest image pulled")
		}
	}

	// Re-create container
	return s.recreateContainer(ctx, container)
}

func (s *DockerService) findContainer(ctx context.Context, containerName string) (types.ContainerJSON, error) {
	// List containers
	log.Debug().Msg("Listing all containers ...")
	containers, err := s.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		err = fmt.Errorf("failed to list Docker containers: %w", err)
		return types.ContainerJSON{}, NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	log.Debug().Msgf("%d containers received", len(containers))

	// Find container
	log.Debug().Str("container_name", containerName).Msg("Finding container in list ...")
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
		return types.ContainerJSON{}, NewError(http.StatusNotFound, openapi.APIERRORCODE_UNKNOWN_CONTAINER, containerName, nil)
	}
	log.Debug().Str("container_name", containerName).Str("container_id", container.ID).Msg("Container found in list")

	// Inspect container
	log.Debug().Str("container_id", container.ID).Msg("Inspecting container ...")
	containerJSON, err := s.client.ContainerInspect(ctx, container.ID)
	if err != nil {
		err = fmt.Errorf("failed to inspect Docker container %s: - %w", container.ID, err)
		return types.ContainerJSON{}, NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	return containerJSON, nil
}

func (s *DockerService) recreateContainer(ctx context.Context, container types.ContainerJSON) error {
	// Stop old container
	log.Debug().Str("container_id", container.ID).Msg("Stopping old container ...")
	err := s.client.ContainerStop(ctx, container.ID, nil)
	if err != nil {
		err = fmt.Errorf("failed to stop old Docker container %s: - %w", container.ID, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	log.Debug().Str("container_id", container.ID).Msg("Old container stopped")

	// Remove old container
	log.Debug().Str("container_id", container.ID).Msg("Removing old container ...")
	err = s.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
	if err != nil {
		err = fmt.Errorf("failed to remove old Docker container %s: - %w", container.ID, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	log.Debug().Str("container_id", container.ID).Msg("Old container removed")

	// Create new container
	log.Debug().Msg("Creating new container ...")
	newContainer, err := s.client.ContainerCreate(
		ctx,
		container.Config,
		container.HostConfig,
		&network.NetworkingConfig{EndpointsConfig: container.NetworkSettings.Networks},
		nil,
		container.Name,
	)
	if err != nil {
		err = fmt.Errorf("failed to create new Docker container %s: %w", container.Name, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}
	log.Debug().Str("container_id", newContainer.ID).Msg("New container created")

	// Start new container
	log.Debug().Str("container_id", newContainer.ID).Msg("Starting new container ...")
	err = s.client.ContainerStart(ctx, newContainer.ID, types.ContainerStartOptions{})
	if err != nil {
		err = fmt.Errorf("failed to start new Docker container %s: %w", container.Name, err)
		return NewError(http.StatusInternalServerError, openapi.APIERRORCODE_UNKNOWN_ERROR, "", err)
	}

	// Recreate successful
	log.Debug().Str("container_id", newContainer.ID).Msg("New container started")
	return nil
}
