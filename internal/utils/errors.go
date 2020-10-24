package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	ErrMessageFatal = "fatal error"
)

var (
	ErrFatal = errors.New(ErrMessageFatal)
)

// Must will be panic if one its argument is an error.
//
//	func myFunction(n int) (int, error) {
//		if n === 0 {
//			return 0, fmt.Errorf("myFunction error: %s", "n cannot be zero")
//		}
//		return n, nil
//	}
//	Must(myFunction(0)) // will panic
//	Must(myFunction(1)) // won't panic, execution continues
//
func Must(args ...interface{}) {
	var err error = nil
	if len(args) == 0 {
		return
	}
	if e, ok := args[len(args)-1].(error); ok {
		err = e
		args = args[:len(args)-1]
	}

	if err == nil {
		for _, arg := range args {
			if e, ok := arg.(error); ok {
				err = e
			}
		}
	}
	if err != nil {
		err = errors.Wrap(err, "FATAL ERROR")
		panic(err)
	}
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func FatalErr(err error) {
	log.Fatal(errors.Wrap(err, ErrFatal.Error()))
}

// Fatal starts to panic with error
func Fatal(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		perr := errors.Wrap(err, ErrFatal.Error())
		perr = errors.Wrap(perr, "panic")
		panic(perr)
	}
}

// Catch catches a panic and recovers from it
func Catch(f func(), g func(interface{})) {
	defer func() {
		if v := recover(); v != nil {
			g(v)
		}
	}()
	f()
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok || cause.Cause() == nil {
			break
		}
		err = cause.Cause()
	}
	return err
}

func recovery() {
	recover()
}
