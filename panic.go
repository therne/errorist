package errorist

import (
	"fmt"
	"math"
)

func RecoverWithErrCapture(capture *error, opts ...Option) {
	if err := WrapPanic(recover()); err != nil {
		*capture = maybeWrap(err, applyOptions(opts))
	}
}

func RecoverWithErrChan(errChan chan<- error, opts ...Option) {
	if err := WrapPanic(recover()); err != nil {
		errChan <- maybeWrap(err, applyOptions(opts))
	}
}

type PanicError struct {
	Reason  string
	Stack   []string
	Options Options
}

func (pe PanicError) Error() string {
	return pe.Pretty()
}

func (pe PanicError) Pretty() string {
	return fmt.Sprintf("%s\n%s", pe.Reason, formatStacktrace(pe.Stack, pe.Options))
}

func WrapPanic(recovered interface{}, opt ...Option) *PanicError {
	if recovered == nil {
		return nil
	}
	opts := applyOptions(opt)
	return &PanicError{
		Reason:  fmt.Sprintf("panic: %s", recovered),
		Stack:   stacktrace(0, math.MaxInt32, opts),
		Options: opts,
	}
}
