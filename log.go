// Package log provides a global logger for zerolog.
package log

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/tehsphinx/errs"
)

func Setup(ctx context.Context, opts ...SetupOption) context.Context {
	cfg := getOptions(opts)

	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		stack := errs.FormatStack(err)
		if stack == "" {
			return nil
		}
		return stack
	}

	out := getOutput(cfg.output, cfg.format)

	logger := zerolog.New(out).With().Timestamp().Caller().Stack().Logger()
	logger.UpdateContext(cfg.updateCtx)
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("service", cfg.serviceName)
	})

	Logger = logger
	zerolog.DefaultContextLogger = &logger
	if ctx != nil {
		ctx = logger.WithContext(ctx)
	}
	return ctx
}

// Logger is the global logger.
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
