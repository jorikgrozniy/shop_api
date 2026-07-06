package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	DSN string
}

func NewConfig() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("HTTP_ADDR", ":8080"),
		},

		Database: DatabaseConfig{
			DSN: getEnv("POSTGRES_DSN", "postgres://user:password@host:port/database"),
		},
	}
}

func NewServerConfig(cfg *Config) *ServerConfig {
	return &cfg.Server
}

func NewDatabaseConfig(cfg *Config) *DatabaseConfig {
	return &cfg.Database
}

func getEnv(key string, defVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defVal
}
