package errorist

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWrapPanicWith(t *testing.T) {
	Convey("Calling errorist.WrapPanic", t, func() {
		Convey("It should catch panic correctly", func() {
			So(func() {
				defer func() {
					if err := WrapPanic(recover()); err != nil {
						fmt.Println(err.Pretty())
					}
				}()
				panicStation()
			}, ShouldNotPanic)
		})

		Convey("With just non-project files", func() {
			Convey("It should catch panic correctly", func() {
				So(func() {
					defer func() {
						if err := WrapPanic(recover(), IncludeNonProjectFiles()); err != nil {
							fmt.Println(err.Pretty())
						}
					}()
					panicStation()
				}, ShouldNotPanic)
			})
		})

		Convey("With DetailedTrace", func() {
			Convey("It should catch panic correctly", func() {
				So(func() {
					defer func() {
						if err := WrapPanic(recover(), WithDetailedTrace()); err != nil {
							fmt.Println(err.Pretty())
						}
					}()
					panicStation()
				}, ShouldNotPanic)
			})
		})

		Convey("With DetailedTrace skipping non-project files", func() {
			Convey("It should catch panic correctly", func() {
				So(func() {
					defer func() {
						if err := WrapPanic(recover(), WithDetailedTrace(), IncludeNonProjectFiles()); err != nil {
							fmt.Println(err.Pretty())
						}
					}()
					panicStation()
				}, ShouldNotPanic)
			})
		})
	})
}

//noinspection ALL
func panicStation() {
	var empty map[string]string
	empty["a"] = "b"
}
