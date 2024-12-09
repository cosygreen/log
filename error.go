package log

import (
	"errors"

	"github.com/cosygreen/errs"
	"github.com/rs/zerolog"
)

// zerologErrorMarshalFunc implements custom error marshalling for the error types in the errs package.
func zerologErrorMarshalFunc(err error) interface{} {
	var info errInfo
	if !errors.As(err, &info) {
		return err
	}

	zErr := newZerologErr(err)
	zErr.evaluate()
	if len(zErr.objects) == 0 {
		return err
	}

	return zErr
}

// errInfo is used to identify errors that have additional information.
// stackError does not implement this interface as stack is handled differently by zerolog.
type errInfo interface {
	error

	HasErrInfo()
}

func newZerologErr(err error) zerologErr {
	return zerologErr{error: err, objects: map[string]zerolog.LogObjectMarshaler{}}
}

type zerologErr struct {
	error
	objects map[string]zerolog.LogObjectMarshaler
}

func (s *zerologErr) evaluate() {
	var (
		fields fieldsT
	)
	defer func() {
		if fields != nil {
			s.objects["fields"] = fields
		}
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
			switch x := err.(type) {
			case errs.FieldsError:
				fields = fields.Join(x.GetFields())
			default:
			}

			nextErrs = unwrap(err, nextErrs)
		}
	}
}

func unwrap(err error, unwrapped []error) []error {
	//nolint:errorlint // This loop is a customized unwrap function.
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		if e := x.Unwrap(); e != nil {
			unwrapped = append(unwrapped, e)
		}
	case interface{ Unwrap() []error }:
		for _, e := range x.Unwrap() {
			if e != nil {
				unwrapped = append(unwrapped, e)
			}
		}
	}

	return unwrapped
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler.
// It gathers all info the error chain contains and adds it to the log statement.
func (s zerologErr) MarshalZerologObject(event *zerolog.Event) {
	event.Str("msg", s.Error())
	for key, object := range s.objects {
		event.Object(key, object)
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
