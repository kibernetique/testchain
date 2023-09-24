package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	// Server
	ServerPort            string        `env:"SERVER_PORT,default=8080"`
	ServerIdleTimeout     time.Duration `env:"SERVER_IDLE_TIMEOUT,default=30s"`
	ServerReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT,default=10s"`
	ServerWriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT,default=10s"`
	ServerShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT,default=30s"`
	ServerMaxUploadSize   int           `env:"SERVER_MAX_UPLOAD_SIZE,default=128"` // in megabytes
}

func GetFromEnv(ctx context.Context) (*Config, error) {
	config := Config{}
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
