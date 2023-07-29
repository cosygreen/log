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

type SetupOption func(*setupOptions)

func ServiceName(name string) SetupOption {
	return func(opts *setupOptions) {
		opts.serviceName = name
	}
}

func WithFormat[T FormatInput](format T) SetupOption {
	logFormat := getFormat(format)

	return func(opts *setupOptions) {
		opts.format = logFormat
	}
}

func UpdateContext(f func(zerolog.Context) zerolog.Context) SetupOption {
	return func(opts *setupOptions) {
		opts.updateCtx = f
	}
}

func WithOutput(out io.Writer) SetupOption {
	return func(opts *setupOptions) {
		opts.output = out
	}
}
