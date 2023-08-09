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
		output:      os.Stdout,
		updateCtx:   func(c zerolog.Context) zerolog.Context { return c },
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

type setupOptions struct {
	serviceName string
	format      Format
	output      io.Writer
	updateCtx   func(c zerolog.Context) zerolog.Context
}

// SetupOption defines an option for setting up the logging.
type SetupOption func(*setupOptions)

// ServiceName sets the service name for logging.
func ServiceName(name string) SetupOption {
	return func(opts *setupOptions) {
		opts.serviceName = name
	}
}

// WithFormat sets the format to use for logging.
func WithFormat[T FormatInput](format T) SetupOption {
	logFormat := getFormat(format)

	return func(opts *setupOptions) {
		opts.format = logFormat
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
