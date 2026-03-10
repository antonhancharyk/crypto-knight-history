package config

import (
	"fmt"
	"os"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type GRPC struct {
	Host string
}

type Server struct {
	Port string
}

type Config struct {
	DB     DB
	GRPC   GRPC
	Server Server
}

func Load() (*Config, error) {
	cfg := &Config{
		DB: DB{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		GRPC: GRPC{
			Host: os.Getenv("GRPC_HOST"),
		},
		Server: Server{
			Port: getEnv("APP_SERVER_PORT", "8080"),
		},
	}

	if cfg.DB.User == "" || cfg.DB.Name == "" {
		return nil, fmt.Errorf("config: DB_USER and DB_NAME are required")
	}
	if cfg.GRPC.Host == "" {
		return nil, fmt.Errorf("config: GRPC_HOST is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
