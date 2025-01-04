package log

import (
	"bytes"
	"io"
	"time"

	"github.com/rs/zerolog"
)

// Format defines the log format.
type Format string

const (
	// FormatPlain uses a zerolog.ConsoleWriter without colors and sends the output to provided output or stdout.
	FormatPlain Format = "plain"
	// FormatPlainWithoutTime uses a zerolog.ConsoleWriter without colors and without time and sends the output to provided output or stdout.
	FormatPlainWithoutTime Format = "plain-notime"
	// FormatColor uses a zerolog.ConsoleWriter with colors and sends the output to provided output or stdout.
	FormatColor Format = "color"
	// FormatColorWithoutTime uses a zerolog.ConsoleWriter with colors and without time and sends the output to provided output or stdout.
	FormatColorWithoutTime Format = "color-notime"
	// FormatJSON writes JSON output to the output provided output or stdout.
	FormatJSON Format = "json"
	// FormatCustom can be used to pass in a custom output writer.
	FormatCustom Format = "custom"
)

// FormatInput allows the input format to be of type string or type Format.
type FormatInput interface {
	string | Format
}

// LevelInput allows the level input to be of type string or type zerolog.Level.
type LevelInput interface {
	string | zerolog.Level
}

func getFormat[T FormatInput](format T) Format {
	var logFormat Format
	switch f := any(format).(type) {
	case Format:
		logFormat = f
	case string:
		logFormat = parseLogFormat(f)
	}
	return logFormat
}

func getLevel[T LevelInput](format T) zerolog.Level {
	logLevel := zerolog.TraceLevel
	switch f := any(format).(type) {
	case zerolog.Level:
		logLevel = f
	case string:
		logLevel, _ = zerolog.ParseLevel(f)
	}
	return logLevel
}

func parseLogFormat(format string) Format {
	switch Format(format) {
	case FormatPlain, FormatColor, FormatJSON, FormatCustom, FormatPlainWithoutTime, FormatColorWithoutTime:
		return Format(format)
	}
	return FormatJSON
}

func getOutput(output io.Writer, format Format) io.Writer {
	switch format {
	case FormatJSON, FormatCustom:
		return output
	case FormatColor:
		return getConsoleOutput(output, true, false)
	case FormatPlain:
		return getConsoleOutput(output, false, false)
	case FormatColorWithoutTime:
		return getConsoleOutput(output, true, true)
	case FormatPlainWithoutTime:
		return getConsoleOutput(output, false, true)
	}
	return output
}

func getConsoleOutput(output io.Writer, color, noTime bool) io.Writer {
	var timeFormat string
	if noTime {
		if color {
			timeFormat = "🌱"
		} else {
			timeFormat = " "
		}
	} else {
		timeFormat = time.Kitchen
	}

	out := zerolog.ConsoleWriter{Out: output, NoColor: !color, TimeFormat: timeFormat}
	out.FieldsExclude = []string{"stack"}
	out.FormatExtra = func(m map[string]interface{}, buf *bytes.Buffer) error {
		if stackI, ok := m["stack"]; ok {
			if stack, ok := stackI.(string); ok {
				buf.WriteByte('\n')
				buf.WriteString(stack)
			}
		}
		return nil
	}
	return out
}
