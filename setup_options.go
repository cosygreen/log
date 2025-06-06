package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func getOptions(opts []SetupOption) setupOptions {
	cfg := setupOptions{
		serviceName: "unknown",
		format:      FormatJSON,
		level:       zerolog.TraceLevel,
		output:      os.Stdout,
		updateCtx:   func(c zerolog.Context) zerolog.Context { return c },
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type setupOptions struct {
	serviceName  string
	hostName     string
	region       string
	publicIP     string
	format       Format
	level        zerolog.Level
	output       io.Writer // os.Stdout is by default but it can be replaced by WithOutput functional option
	updateCtx    func(c zerolog.Context) zerolog.Context
	hideCaller   bool
	extraWriters []io.Writer
}

// SetupOption defines an option for setting up the logging.
type SetupOption func(*setupOptions)

func WithConfig(config Config) SetupOption {
	return func(opts *setupOptions) {
		level, _ := zerolog.ParseLevel(config.LogLevel)
		opts.level = level
		opts.format = parseLogFormat(config.LogFormat)
		opts.serviceName = config.ServiceName
		opts.hostName = config.HostName
		opts.region = config.Region
		opts.publicIP = config.PublicIP
		opts.hideCaller = config.HideCaller
	}
}

// WithExtraWriters adds extra writers to the logger, so you can use any Format option.
func WithExtraWriters(w ...io.Writer) SetupOption {
	return func(opts *setupOptions) {
		opts.extraWriters = append(opts.extraWriters, w...)
	}
}

// ServiceName sets the service name for logging.
func ServiceName(name string) SetupOption {
	return func(opts *setupOptions) {
		opts.serviceName = name
	}
}

// HostName sets the host name for logging.
func HostName(name string) SetupOption {
	return func(opts *setupOptions) {
		opts.hostName = name
	}
}

// Region sets the data center region for logging.
func Region(name string) SetupOption {
	return func(opts *setupOptions) {
		opts.region = name
	}
}

// PublicIP sets the public ip for logging.
func PublicIP(ip string) SetupOption {
	return func(opts *setupOptions) {
		opts.publicIP = ip
	}
}

// WithFormat sets the format to use for logging.
func WithFormat[T FormatInput](format T) SetupOption {
	logFormat := getFormat(format)

	return func(opts *setupOptions) {
		opts.format = logFormat
	}
}

// WithLevel sets the global log level to use for logging.
func WithLevel[T LevelInput](level T) SetupOption {
	logLevel := getLevel(level)

	return func(opts *setupOptions) {
		opts.level = logLevel
	}
}

// UpdateContext can be used to update log context, adding additional information at setup.
func UpdateContext(f func(zerolog.Context) zerolog.Context) SetupOption {
	return func(opts *setupOptions) {
		opts.updateCtx = f
	}
}

// WithOutput sets the writer to write the log output to.
// Can be used to use a customized output writer. Doing so should be done in combination
// with passing FormatCustom to WithFormat.
func WithOutput(out io.Writer) SetupOption {
	return func(opts *setupOptions) {
		opts.output = out
	}
}
