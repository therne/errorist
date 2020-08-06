package errorist

import "io"

// CloseWithErrCapture is used if you want to close and fail the function or
// method on a `io.Closer.Close()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `CloseWithErrCapture`
// will append the error caused by `Close` if any.
func CloseWithErrCapture(c io.Closer, capture *error, opts ...Option) {
	if err := c.Close(); err != nil && *capture != nil {
		*capture = maybeWrap(err, applyOptions(opts))
	}
}

// CloseWithErrCapture is used if you want to close and fail the function or
// method on a `io.Closer.Close()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `CloseWithErrChan`
// will send the error to the given channel caused by `Close` if any.
func CloseWithErrChan(c io.Closer, errChan chan<- error, opts ...Option) {
	if err := c.Close(); err != nil {
		errChan <- maybeWrap(err, applyOptions(opts))
	}
}

// CloseWithLogOnErr is used if you want to close and fail the function or
// method on a `io.Closer.Close()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `CloseWithLogOnErr`
// will log the error caused by `Close` if any.
func CloseWithLogOnErr(c io.Closer, opts ...Option) {
	if err := c.Close(); err != nil {
		opt := applyOptions(opts)
		opt.Logger(maybeWrap(err, opt).Error())
	}
}
