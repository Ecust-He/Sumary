package i18n

import (
	"fmt"
	"io"
)

type ErrLevel int

const (
	ERROR ErrLevel = iota // 0
	INFO                  // 1
	WARN                  // 2
)

var (
	ErrLevelToString = map[ErrLevel]string{
		ERROR: "ERROR",
		INFO:  "INFO",
		WARN:  "WARN",
	}
)

type (
	Error interface {
		Error() string
		WithCode(string) Error
		Code() string
		WithArgs(...interface{}) Error
		Args() []interface{}
		WithLevel(level ErrLevel) Error
		Level() ErrLevel
		Unwrap() error
		Cause() error
		Format(s fmt.State, verb rune)
	}

	defaultError struct {
		err   error
		C     string `json:"code"`
		M     string `json:"message"`
		level ErrLevel
		args  []interface{}
	}
)

func WrapCodeError(err error, c string) Error {
	return &defaultError{
		err:   err,
		C:     c,
		level: ERROR,
	}
}

func (e *defaultError) Error() string {
	return e.err.Error()
}

func (e *defaultError) WithCode(code string) Error {
	e.C = code
	return e
}

func (e *defaultError) Code() string {
	return e.C
}

func (e *defaultError) WithArgs(args ...interface{}) Error {
	e.args = append(e.args, args...)
	return e
}

func (e *defaultError) Args() []interface{} {
	return e.args
}

func (e *defaultError) WithLevel(level ErrLevel) Error {
	e.level = level
	return e
}

func (e *defaultError) Level() ErrLevel {
	return e.level
}

func (e *defaultError) Unwrap() error {
	return e.err
}

func (e *defaultError) Cause() error {
	return e.err
}

func (e *defaultError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.Cause())
			_, _ = io.WriteString(s, e.M)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
