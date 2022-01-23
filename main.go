package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	gin.SetMode(gin.ReleaseMode)

	// Parse config
	config, err := parseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	// Setup Debug logging if enabled
	if config.Server.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.Debug().Msg("Debug logging enabled")
	}

	// Setup Gin
	router := gin.Default()
	err = router.SetTrustedProxies(config.Server.TrustedProxies)
	if err != nil {
		log.Fatal().Err(err).Strs("trusted_proxies", config.Server.TrustedProxies).Msg("Failed to set trusted proxies")
	}
	router.StaticFile("/", "docs/index.html")
	router.StaticFile("/index.html", "docs/index.html")
	router.StaticFile("/openapi.yml", "docs/openapi.yml")

	// Setup Docker service
	dockerService, err := NewDockerService()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Docker service")
	}

	// Setup handler
	handler := NewHandler(dockerService, config.Clients)
	handler.RegisterRoutes(router.Group("/"))

	// Start Gin
	port := strconv.Itoa(config.Server.Port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal().Err(err).Int("port", config.Server.Port).Msg("Failed to start Gin server")
	}
}
