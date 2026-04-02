package log

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func captureOutput(t *testing.T, format Format, fn func()) string {
	t.Helper()

	var buf bytes.Buffer

	// `log.Setup` mutates global zerolog configuration, so restore the
	// relevant globals to keep tests isolated.
	oldLogger := Logger
	oldDefaultContextLogger := zerolog.DefaultContextLogger
	oldErrorMarshalFunc := zerolog.ErrorMarshalFunc
	oldErrorStackMarshaler := zerolog.ErrorStackMarshaler
	oldFatalExitFunc := zerolog.FatalExitFunc

	defer func() {
		Logger = oldLogger
		zerolog.DefaultContextLogger = oldDefaultContextLogger
		zerolog.ErrorMarshalFunc = oldErrorMarshalFunc
		zerolog.ErrorStackMarshaler = oldErrorStackMarshaler
		zerolog.FatalExitFunc = oldFatalExitFunc
	}()

	Setup(
		context.Background(),
		WithFormat(format),
		WithOutput(&buf),
		WithLevel(zerolog.TraceLevel),
	)

	fn()
	return buf.String()
}

func TestConsolePlain_ErrorField_Present(t *testing.T) {
	testErr := errors.New("test fatal")
	msg := "*** cosy.green a-service crashed ***"

	cases := []struct {
		name string
		run  func()
	}{
		{
			name: "logErr",
			run: func() {
				Err(testErr).Msg(msg)
			},
		},
		{
			name: "logErrorErr",
			run: func() {
				Logger.Error().Err(testErr).Msg(msg)
			},
		},
		{
			name: "logFatalErr",
			run: func() {
				// Prevent terminating the test process.
				oldFatalExitFunc := zerolog.FatalExitFunc
				zerolog.FatalExitFunc = func() {}
				defer func() { zerolog.FatalExitFunc = oldFatalExitFunc }()
				Logger.Fatal().Err(testErr).Msg(msg)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := captureOutput(t, FormatPlain, tc.run)
			if !strings.Contains(out, "error=") || !strings.Contains(out, "test fatal") {
				t.Fatalf("expected console output to include error field; got: %q", out)
			}
			if !strings.Contains(out, msg) {
				t.Fatalf("expected message to be present; got: %q", out)
			}
		})
	}
}

func TestConsolePlain_StackPrintedOnlyWhenString(t *testing.T) {
	out := captureOutput(t, FormatPlain, func() {
		Logger.Error().Str("stack", "STACKTRACE").Msg("hello")
	})
	if !strings.Contains(out, "\nSTACKTRACE") {
		t.Fatalf("expected stack to be printed on its own line; got: %q", out)
	}
	if strings.Contains(out, "stack=") {
		t.Fatalf("expected stack field name to be hidden (printed only as extra); got: %q", out)
	}

	out = captureOutput(t, FormatPlain, func() {
		Logger.Error().Interface("stack", 123).Msg("hello")
	})
	if strings.Contains(out, "123") {
		t.Fatalf("expected non-string stack to be ignored; got: %q", out)
	}
	if strings.Contains(out, "stack=") {
		t.Fatalf("expected stack field name to be hidden even if non-string; got: %q", out)
	}
}

func TestJSON_ErrorField_Present(t *testing.T) {
	testErr := errors.New("test fatal")
	msg := "*** cosy.green a-service crashed ***"

	cases := []struct {
		name string
		run  func()
	}{
		{
			name: "logErr",
			run: func() {
				Err(testErr).Msg(msg)
			},
		},
		{
			name: "logErrorErr",
			run: func() {
				Logger.Error().Err(testErr).Msg(msg)
			},
		},
		{
			name: "logFatalErr",
			run: func() {
				oldFatalExitFunc := zerolog.FatalExitFunc
				zerolog.FatalExitFunc = func() {}
				defer func() { zerolog.FatalExitFunc = oldFatalExitFunc }()
				Logger.Fatal().Err(testErr).Msg(msg)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := captureOutput(t, FormatJSON, tc.run)
			if !strings.Contains(out, `"error":"test fatal"`) {
				t.Fatalf("expected JSON output to include error field; got: %q", out)
			}
		})
	}
}

func TestErrorMarshalFunc_PlainError(t *testing.T) {
	var buf bytes.Buffer

	oldLogger := Logger
	oldDefaultContextLogger := zerolog.DefaultContextLogger
	oldErrorMarshalFunc := zerolog.ErrorMarshalFunc
	oldErrorStackMarshaler := zerolog.ErrorStackMarshaler
	oldFatalExitFunc := zerolog.FatalExitFunc
	defer func() {
		Logger = oldLogger
		zerolog.DefaultContextLogger = oldDefaultContextLogger
		zerolog.ErrorMarshalFunc = oldErrorMarshalFunc
		zerolog.ErrorStackMarshaler = oldErrorStackMarshaler
		zerolog.FatalExitFunc = oldFatalExitFunc
	}()

	Setup(
		context.Background(),
		WithFormat(FormatJSON),
		WithOutput(&buf),
		WithLevel(zerolog.TraceLevel),
	)

	testErr := errors.New("test fatal")
	v := zerolog.ErrorMarshalFunc(testErr)
	if v == nil {
		t.Fatalf("expected ErrorMarshalFunc to return non-nil value for plain error")
	}
	ee, ok := v.(error)
	if !ok {
		t.Fatalf("expected ErrorMarshalFunc to return error for plain error; got %T (%v)", v, v)
	}
	if ee.Error() != "test fatal" {
		t.Fatalf("expected marshaled error to preserve error string; got %q", ee.Error())
	}
}
