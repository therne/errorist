package errorist

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStacktrace(t *testing.T) {
	errorsWrapTraces := Stacktrace(errWithPkgErrorsWrap())
	errorsNewTraces := Stacktrace(errWithPkgErrorsNew())
	noStackTraces := Stacktrace(fmt.Errorf("no stacktrace"))

	Convey("Calling errorist.Stacktrace", t, func() {
		Convey("With an error wrapped by errors.Wrap, it should return its stacktrace", func() {
			So(errorsWrapTraces, ShouldHaveLength, 3)
			So(errorsWrapTraces[0], ShouldStartWith, "github.com/therne/errorist.errWithPkgErrorsWrap")
			So(errorsWrapTraces[1], ShouldStartWith, "github.com/therne/errorist.TestStacktrace")
		})

		Convey("With an error created by errors.New, it should return its stacktrace", func() {
			So(errorsNewTraces, ShouldHaveLength, 3)
			So(errorsNewTraces[0], ShouldStartWith, "github.com/therne/errorist.errWithPkgErrorsNew")
			So(errorsNewTraces[1], ShouldStartWith, "github.com/therne/errorist.TestStacktrace")
		})

		Convey("With an error without stacktrace, it should return caller trace instead", func() {
			So(noStackTraces, ShouldHaveLength, 1)
			So(noStackTraces[0], ShouldStartWith, "github.com/therne/errorist.TestStacktrace")
		})
	})
}

func errWithPkgErrorsWrap() error {
	return errors.Wrap(fmt.Errorf("no stacktrace"), "world")
}

func errWithPkgErrorsNew() error {
	return errors.New("hello")
}