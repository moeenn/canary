package config

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Port    int
	Timeout time.Duration
}

func newServerConfig() (*ServerConfig, error) {
	defaultServerPort := 3000
	port, err := readInt("SERVER_PORT", &defaultServerPort)
	if err != nil {
		return nil, err
	}

	cfg := ServerConfig{
		Port:    port,
		Timeout: time.Second * 10,
	}

	return &cfg, nil
}

func (cfg ServerConfig) Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

type DatabaseConfig struct {
	Filepath string
}

func newDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Filepath: "./service.db",
	}
}

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
}

func NewConfig() (*Config, error) {
	server, err := newServerConfig()
	if err != nil {
		return nil, fmt.Errorf("invalid server config: %w", err)
	}

	db := newDatabaseConfig()

	cfg := Config{
		Server:   server,
		Database: &db,
	}

	return &cfg, nil
}
