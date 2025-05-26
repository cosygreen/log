// Package log provides a global logger for zerolog.
package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cosygreen/errs"
	"github.com/rs/zerolog"
)

// Setup is used to set up logging.
// If a context is passed, the logger will be added to the context and the context returned.
func Setup(ctx context.Context, opts ...SetupOption) context.Context {
	cfg := getOptions(opts)

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(cfg.level)

	//nolint:reassign // zerolog global variables are meant to be reassigned for setup.
	zerolog.ErrorMarshalFunc = zerologErrorMarshalFunc

	//nolint:reassign
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		stack := errs.FormatStack(err)
		if stack == "" {
			return nil
		}
		return stack
	}

	out := getOutput(cfg.output, cfg.format)

	if len(cfg.extraWriters) > 0 {
		out = zerolog.MultiLevelWriter(append(cfg.extraWriters, out)...)
	}

	loggerCtx := zerolog.New(out).With().Timestamp()
	if !cfg.hideCaller {
		loggerCtx = loggerCtx.Caller()
	}
	logger := loggerCtx.Stack().Logger()
	logger.UpdateContext(cfg.updateCtx)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		switch cfg.serviceName {
		case "", "unknown":
		default:
			c = c.Str("service", cfg.serviceName)
		}
		if cfg.hostName != "" {
			c = c.Str("host", cfg.hostName)
		}
		if cfg.region != "" {
			c = c.Str("region", cfg.region)
		}
		if cfg.publicIP != "" {
			c = c.Str("publicIP", cfg.publicIP)
		}
		return c
	})

	Logger = logger
	zerolog.DefaultContextLogger = &logger
	if ctx != nil {
		ctx = logger.WithContext(ctx)
	}

	return ctx
}

// Logger is the global logger.
//
//nolint:gochecknoglobals
var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *zerolog.Event {
	return Logger.Err(err)
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

// Errorf sends a log event using error level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	Logger.Error().CallerSkipFrame(1).Msgf(format, v...)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return Logger.WithLevel(zerolog.FatalLevel)
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return Logger.Panic()
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
