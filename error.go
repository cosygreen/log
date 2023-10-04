package log

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/tehsphinx/errs"
)

// zerologErrorMarshalFunc implements custom error marshalling for the error types in the errs package.
func zerologErrorMarshalFunc(err error) interface{} {
	var info errInfo
	if errors.As(err, &info) {
		return zerologErr{error: err}
	}

	return err
}

// errInfo is used to identify errors that have additional information.
// stackError does not implement this interface as stack is handled differently by zerolog.
type errInfo interface {
	error

	HasErrInfo()
}

type zerologErr struct {
	error
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler.
// It gathers all info the error chain contains and adds it to the log statement.
func (s zerologErr) MarshalZerologObject(event *zerolog.Event) {
	event.Str("msg", s.Error())

	var fields fieldsT
	defer func() {
		event.Object("fields", fields)
	}()

	nextErrs := []error{s.error}
	for {
		if len(nextErrs) == 0 {
			return
		}

		curErrs := nextErrs
		nextErrs = nil

		for _, err := range curErrs {
			//nolint:errorlint // This loop is a customized unwrap function.
			//goland:noinspection GoTypeAssertionOnErrors
			if x, ok := err.(errs.FieldsError); ok {
				fields = fields.Join(x.GetFields())
			}

			//nolint:errorlint // This loop is a customized unwrap function.
			switch x := err.(type) {
			case interface{ Unwrap() error }:
				if e := x.Unwrap(); e != nil {
					nextErrs = append(nextErrs, e)
				}
			case interface{ Unwrap() []error }:
				for _, e := range x.Unwrap() {
					if e != nil {
						nextErrs = append(nextErrs, e)
					}
				}
			}
		}
	}
}

type fieldsT map[string]interface{}

func (f fieldsT) Join(fields fieldsT) fieldsT {
	for k, v := range f {
		fields[k] = v
	}
	return fields
}

func (f fieldsT) MarshalZerologObject(e *zerolog.Event) {
	for k, v := range f {
		e.Any(k, v)
	}
}
