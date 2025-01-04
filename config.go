package log

import (
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel    zerolog.Level `default:"0"     env:"LOG_LEVEL"`
	LogFormat   string        `default:"color" env:"LOG_FORMAT"`
	ServiceName string        `env:"SERVICE_NAME"`
	HostName    string        `env:"HOST_NAME"`
	Region      string        `env:"REGION"`
	PublicIP    string        `env:"PUBLIC_IP"`
	HideTime    bool          `default:"true" env:"HIDE_TIME"`
}
