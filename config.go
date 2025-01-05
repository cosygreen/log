package log

type Config struct {
	LogLevel    string `default:"debug" env:"LOG_LEVEL"`
	LogFormat   string `default:"color" env:"LOG_FORMAT"`
	ServiceName string `env:"SERVICE_NAME"`
	HostName    string `env:"HOST_NAME"`
	Region      string `env:"REGION"`
	PublicIP    string `env:"PUBLIC_IP"`
	HideCaller  bool   `env:"LOG_HIDE_CALLER"`
}
