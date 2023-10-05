package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// NewSplitWriter returns an output writer that logs to 2 targets with different output formats and levels.
// Use with WithFormat(FormatCustom).
// Using WithLevel would affect both outputs restricting them more. It doesn't overwrite the levels for the 2 outputs.
func NewSplitWriter[f FormatInput, l LevelInput](userOut, devOut io.Writer, userFormat, devFormat f, userLevel, devLevel l) io.WriteCloser {
	return realtymeOutput{
		userOut:   getOutput(userOut, getFormat(userFormat)),
		userLevel: getLevel(userLevel),
		devOut:    getOutput(devOut, getFormat(devFormat)),
		devLevel:  getLevel(devLevel),
	}
}

// realtymeOutput implements zerolog.LevelWriter and writes to multiple outputs.
type realtymeOutput struct {
	userOut   io.Writer
	userLevel zerolog.Level
	devOut    io.Writer
	devLevel  zerolog.Level
}

// Write should not be called
func (l realtymeOutput) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

// WriteLevel write to the appropriate output
func (l realtymeOutput) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if l.userLevel <= level {
		n, err = l.userOut.Write(p)
	}
	if l.devLevel <= level {
		n, err = l.devOut.Write(p)
	}
	return n, err
}

func (l realtymeOutput) Close() error { return nil }
