package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
	DB     DBConfig     `mapstructure:"db"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type DBConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func Load() (*Config, error) {
	// Get environment from ENV variable, default to "dev"
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	return LoadForEnvironment(env)
}

func LoadForEnvironment(env string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("logger.level", "info")
	v.SetDefault("db.enabled", false)

	// Configure file reading
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// Try to read the config file
	if err := v.ReadInConfig(); err != nil {
		// If config file is not found, continue with defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Environment variables override config files
	v.SetEnvPrefix("BLOG_API")
	v.AutomaticEnv()

	v.BindEnv("server.port", "PORT")
	v.BindEnv("logger.level", "LOG_LEVEL")
	v.BindEnv("db.enabled", "DB_ENABLED")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *LoggerConfig) GetSlogLevel() slog.Level {
	switch c.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
