package errors

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

const (
	System   = 50001
	Business = 40001
	NoLogin  = 1007

	Signature = 6003
)

const (
	RPCNoAuthorization = iota + 8000
	RPCNoAuthorizationStep1
	RPCNoAuthorizationStep2
	RPCNoAuthorizationStep3
	RPCNoAuthorizationStep4
)

type Error struct {
	error error
	code  int
}

func New(info interface{}, code ...int) error {
	err := &Error{}
	if e, ok := info.(error); ok {
		err.error = fmt.Errorf("%w", e)
	} else {
		err.error = fmt.Errorf("%v", info)
	}
	if len(code) == 0 {
		err.code = Business
	} else {
		err.code = code[0]
	}
	return err
}

func As(err error) (target Error, ok bool) {
	ok = errors.As(err, &target)
	return
}

func Wrap(info interface{}, e error, code ...int) error {
	err := Error{}
	err.error = fmt.Errorf("%v, %w", info, e)
	if len(code) == 0 {
		err.code = Business
	} else {
		err.code = code[0]
	}
	return err
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func (e Error) Error() string {
	return e.error.Error()
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Is(target error) bool {
	return errors.Is(e.error, target)
}

func Background(err error) error {
	_, file, line, _ := runtime.Caller(1)
	dir, f := filepath.Split(file)
	e := fmt.Errorf("<%s/%s:%d> [%w]", filepath.Base(dir), f, line, err)
	return e
}
