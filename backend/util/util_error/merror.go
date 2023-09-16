package util_error

import (
	"github.com/pkg/errors"
	"runtime"
	"strconv"
)

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

type Error struct {
	Code    int
	Desc    string
	Message string
}

func (e *Error) Error() string {
	return e.Desc
}
