package log

import (
	"bytes"
	"io"

	"github.com/rs/zerolog"
)

// Format defines the log format.
type Format string

const (
	// FormatPlain uses a zerolog.ConsoleWriter without colors and sends the output to provided output or stdout.
	FormatPlain Format = "plain"
	// FormatColor uses a zerolog.ConsoleWriter with colors and sends the output to provided output or stdout.
	FormatColor Format = "color"
	// FormatJSON writes JSON output to the output provided output or stdout.
	FormatJSON Format = "json"
	// FormatCustom can be used to pass in a custom output writer.
	FormatCustom Format = "custom"
)

// FormatInput allows the input format to be of type string or type Format.
type FormatInput interface {
	string | Format
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

func parseLogFormat(format string) Format {
	switch Format(format) {
	case FormatPlain, FormatColor, FormatJSON, FormatCustom:
		return Format(format)
	}
	return FormatJSON
}

func getOutput(output io.Writer, format Format) io.Writer {
	switch format {
	case FormatJSON:
		return output
	case FormatColor:
		return getConsoleOutput(output, true)
	case FormatPlain:
		return getConsoleOutput(output, false)
	case FormatCustom:
		return output
	}
	return output
}

func getConsoleOutput(output io.Writer, color bool) io.Writer {
	out := zerolog.ConsoleWriter{Out: output, NoColor: !color}
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
