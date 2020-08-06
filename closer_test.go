package errorist

import (
	stdlibErrors "errors"
	"testing"

	pkgErrors "github.com/pkg/errors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCloseWithErrCapture(t *testing.T) {
	Convey("Calling errorist.CloseWithErrCapture", t, func() {
		expectedErr := pkgErrors.New("test")
		var actualErr error

		Convey("It should capture error caused while closing", func() {
			m := &closerMock{ReturnError: expectedErr}
			CloseWithErrCapture(m, &actualErr)

			So(m.CloseCalled, ShouldEqual, 1)
			So(actualErr, ShouldEqual, expectedErr)
		})

		Convey("It should not capture error if underlying error is already present", func() {
			actualErr = pkgErrors.New("already present")

			m := &closerMock{ReturnError: expectedErr}
			defer CloseWithErrCapture(m, &actualErr)

			So(m.CloseCalled, ShouldEqual, 1)
			So(actualErr, ShouldBeError, "already present")
		})

		Convey("It should add context with Wrapf", func() {
			m := &closerMock{ReturnError: expectedErr}
			CloseWithErrCapture(m, &actualErr, Wrapf("closing some"))

			So(m.CloseCalled, ShouldEqual, 1)
			So(actualErr.Error(), ShouldEqual, "closing some: test")
			So(pkgErrors.Cause(actualErr), ShouldEqual, expectedErr)
			So(stdlibErrors.Is(actualErr, expectedErr), ShouldBeTrue)
		})
	})
}

func TestCloseWithErrChan(t *testing.T) {
	Convey("Calling errorist.CloseWithErrChan", t, func() {
		expectedErr := pkgErrors.New("test")
		actualErrChan := make(chan error, 1)

		Convey("It should send error to given channel if there's error on closing", func() {
			m := &closerMock{ReturnError: expectedErr}
			CloseWithErrChan(m, actualErrChan)

			So(m.CloseCalled, ShouldEqual, 1)
			So(<-actualErrChan, ShouldEqual, expectedErr)
		})
	})
}

func TestCloseWithLogOnErr(t *testing.T) {
	Convey("Calling errorist.CloseWithLogOnErr", t, func() {
		expectedErr := pkgErrors.New("test")
		var actualErr string
		loggerMock := func(err string) { actualErr = err }

		Convey("It should log if there's error on closing", func() {
			m := &closerMock{ReturnError: expectedErr}
			CloseWithLogOnErr(m, WithLogHandler(loggerMock))

			So(m.CloseCalled, ShouldEqual, 1)
			So(actualErr, ShouldEqual, expectedErr.Error())
		})
	})
}

type closerMock struct {
	CloseCalled int
	ReturnError error
}

func (m *closerMock) Close() error {
	m.CloseCalled++
	return m.ReturnError
}
