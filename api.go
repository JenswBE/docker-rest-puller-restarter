package main

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/JenswBE/docker-rest-puller-restarter/openapi"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	clients map[string]Client
	docker  *DockerService
}

func NewHandler(dockerService *DockerService, apiClients []Client) *Handler {
	// Build API client map
	clients := make(map[string]Client, len(apiClients))
	for _, client := range apiClients {
		duplicateClient, isDuplicate := clients[client.APIKey]
		if isDuplicate {
			log.Fatal().Str("client_a", client.Name).Str("client_b", duplicateClient.Name).Msg("Duplicate client key detected")
		}
		clients[client.APIKey] = client
	}

	// Build handler
	return &Handler{
		clients: clients,
		docker:  dockerService,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/:container_name/pull_restart/", h.pullRestartContainer)
}

func (h *Handler) pullRestartContainer(c *gin.Context) {
	// Extract params
	containerName := c.Param("container_name")

	// Validate if API key has access
	if err := h.validateAPIKey(c, containerName); err != nil {
		c.JSON(errToResponse(err))
		return
	}

	// Pull and restart container
	if err := h.docker.PullRestartContainer(c.Request.Context(), containerName); err != nil {
		c.JSON(errToResponse(err))
		return
	}

	// Pull and restart successful
	c.Status(http.StatusNoContent)
}

// validateAPIKey checks if the API key exists.
// Next to this, it validates if the API key has access to a container with that name.
func (h *Handler) validateAPIKey(c *gin.Context, containerName string) error {
	// Extact API key
	apiKey := c.GetHeader("API-KEY")

	// Check if API key exists
	client, ok := h.clients[apiKey]
	if !ok {
		log.Info().Str("ip", c.ClientIP()).Str("api_key", apiKey).Msg("Request received with invalid API key")
		return NewError(http.StatusUnauthorized, openapi.APIERRORCODE_INVALID_API_KEY, "", nil)
	}

	// Check if API key has access on container
	if !(sliceContains(client.ContainerNames, "*") || sliceContains(client.ContainerNames, containerName)) {
		log.Info().
			Str("ip", c.ClientIP()).
			Str("client", client.Name).
			Str("container", containerName).
			Strs("allowed_containers", client.ContainerNames).
			Msg("Client requested container for which it has no access")
		return NewError(http.StatusForbidden, openapi.APIERRORCODE_CONTAINER_NAME_NOT_ENABLED_FOR_API_KEY, containerName, nil)
	}

	// Validation successful
	log.Info().
		Str("ip", c.ClientIP()).
		Str("client", client.Name).
		Str("container", containerName).
		Msgf("Client %s requested pull-restart for container %s", client.Name, containerName)
	return nil
}

// errToResponse checks if the provided error is a APIError.
// If yes, status and embedded error message are returned.
// If no, status is 500 and provided error message are returned.
func errToResponse(e error) (int, *APIError) {
	log.Debug().Err(e).Msg("Error received on API level")
	var apiErr *APIError
	if errors.As(e, &apiErr) {
		return apiErr.Status, apiErr
	}
	log.Warn().Err(e).Stringer("error_type", reflect.TypeOf(e)).Msg("API received an non-APIError error")
	return 500, NewError(500, openapi.APIERRORCODE_UNKNOWN_ERROR, "", e)
}
