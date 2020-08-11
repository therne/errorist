package errorist

type Stopper interface {
	Stop() error
}

// StopWithErrCapture is used if you want to Stop and fail the function or
// method on a `Stop()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `StopWithErrCapture`
// will append the error caused by `Stop` if any.
func StopWithErrCapture(c Stopper, capture *error, opts ...Option) {
	if err := c.Stop(); err != nil && *capture != nil {
		*capture = maybeWrap(err, applyOptions(opts))
	}
}

// StopWithErrChan is used if you want to Stop and fail the function or
// method on a `Stop()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `StopWithErrChan`
// will send the error to the given channel caused by `Stop` if any.
func StopWithErrChan(c Stopper, errChan chan<- error, opts ...Option) {
	if err := c.Stop(); err != nil {
		errChan <- maybeWrap(err, applyOptions(opts))
	}
}

// StopWithErrLog is used if you want to Stop and fail the function or
// method on a `Stop()` error (make sure the `error` return argument is
// named as `err`). If the error is already present, `StopWithErrLog`
// will log the error caused by `Stop` if any.
func StopWithErrLog(c Stopper, opts ...Option) {
	if err := c.Stop(); err != nil {
		opt := applyOptions(opts)
		opt.Logger(maybeWrap(err, opt).Error())
	}
}
