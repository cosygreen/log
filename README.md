# log

Log package is a very opinionated replacement for the `zerolog/log` package, building on zerolog.

## Features

In addition to zerolog this package adds:

- Out of the box formats `log.FormatColor` "color", `log.FormatPlain` "plain", `log.FormatJSON` "json"
- Logging and formatting for stack traces of `github.com/tehsphins/errs` package

## Usage

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosygreen/errs"
	"github.com/cosygreen/log"
	"github.com/rs/zerolog"
)

func main() {
	fmt.Println("--- CONSOLE LOGGING ---")
	logPkg(log.FormatColor)

	fmt.Println("--- JSON LOGGING ---")
	logPkg(log.FormatJSON)
}

func logPkg(format log.Format) {
	ctx := log.Setup(context.Background(),
		log.WithFormat(format),
		log.ServiceName("test"),
		log.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("abc", "xyz")
		}),
	)
	err := errors.New("error msg")
	errStack := errs.WithStack(err)

	log.Ctx(ctx).Info().Str("key1", "val1").Msg("this is a log message")
	log.Ctx(ctx).Err(err).Msg("this is a custom err message")
	log.Ctx(ctx).Err(errStack).Msg("this is a custom err message")
}
```
