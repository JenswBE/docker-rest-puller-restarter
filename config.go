package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Clients []Client
	Server  struct {
		Debug          bool
		Port           int
		TrustedProxies []string
	}
}

type Client struct {
	Name           string
	APIKey         string
	ContainerNames []string
}

func parseConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("Server.Debug", false)
	viper.SetDefault("Server.Port", 8080)
	viper.SetDefault("Server.TrustedProxies", []string{"172.16.0.0/16"}) // Default Docker IP range

	// Parse config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed reading config file: %w", err)
		}
		log.Warn().Err(err).Msg("No config file found, expecting configuration through ENV variables")
	}

	// Bind ENV variables
	err = bindEnvs([]envBinding{
		{"Server.Debug", "SERVER_DEBUG"},
		{"Server.Port", "SERVER_PORT"},
		{"Server.TrustedProxies", "SERVER_TRUSTED_PROXIES"},
	})
	if err != nil {
		return nil, err
	}

	// Unmarshal to Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return &config, nil
}

type envBinding struct {
	configPath string
	envName    string
}

func bindEnvs(bindings []envBinding) error {
	for _, binding := range bindings {
		err := viper.BindEnv(binding.configPath, binding.envName)
		if err != nil {
			return fmt.Errorf("failed to bind env var %s to %s: %w", binding.envName, binding.configPath, err)
		}
	}
	return nil
}

func validateConfig(config Config) error {
	// Validate config
	if len(config.Clients) == 0 {
		return fmt.Errorf("no clients defined, please define at least 1 client")
	}

	// Validate clients
	for i, client := range config.Clients {
		if err := validateClient(client, i); err != nil {
			return err
		}
	}

	// Validation successful
	return nil
}

func validateClient(client Client, index int) error {
	switch {
	case client.Name == "":
		return fmt.Errorf("client name is required, but missing for client %d", index+1)
	case client.APIKey == "":
		return fmt.Errorf("client API key is required, but missing for client %s", client.Name)
	case len(client.ContainerNames) == 0:
		return fmt.Errorf("client container names are required, but missing for client %s. Please use * if you want to allow all container names", client.Name)
	default:
		// Validation successful
		return nil
	}
}
