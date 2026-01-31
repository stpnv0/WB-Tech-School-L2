package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	defaultConfigPath = "config/config.yaml"
	envConfigPath     = "CONFIG_PATH"
)

type Config struct {
	HTTP   HTTPConfig   `yaml:"http"`
	Logger LoggerConfig `yaml:"logger"`
}

type HTTPConfig struct {
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"HTTP_PORT" env-default:"8080"`

	Timeout struct {
		Server time.Duration `yaml:"server" env-default:"5s"`
		Write  time.Duration `yaml:"write"  env-default:"5s"`
		Read   time.Duration `yaml:"read"   env-default:"5s"`
		Idle   time.Duration `yaml:"idle"   env-default:"30s"`
	} `yaml:"timeout"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"  env:"LOG_LEVEL"  env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`   // json | text
	Output string `yaml:"output" env:"LOG_OUTPUT" env-default:"stdout"` // stdout | file
	Path   string `yaml:"path"   env:"LOG_PATH"   env-default:""`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	cfg, err := load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func load() (*Config, error) {
	path := getConfigPath()

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file error: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func getConfigPath() string {
	if path := os.Getenv(envConfigPath); path != "" {
		return path
	}
	return defaultConfigPath
}

func (h HTTPConfig) Addr() string {
	return net.JoinHostPort(h.Host, strconv.Itoa(h.Port))
}
