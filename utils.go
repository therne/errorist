package errorist

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func callerPackageName(skip int) string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return "unknown"
	}
	frame, _ := runtime.CallersFrames(pc).Next()
	frags := strings.Split(frame.Function, ".")
	if len(frags) == 0 {
		if frame.Function != "" {
			return frame.Function
		}
		return frame.File
	}
	return strings.Join(frags[:len(frags)-1], ".")
}

func callerTrace(skip int) string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return "unknown"
	}
	frame, _ := runtime.CallersFrames(pc).Next()
	return fmt.Sprintf("%s (%s:%d)", frame.Function, filepath.Base(frame.File), frame.Line)
}

func maybeWrap(err error, opts Options) error {
	if err == nil {
		return nil
	}
	if len(opts.WrapArguments) == 0 {
		return err
	}
	switch format := opts.WrapArguments[0].(type) {
	case string:
		return errors.Wrapf(err, format, opts.WrapArguments[1:]...)
	default:
		return errors.Wrap(err, fmt.Sprint(opts.WrapArguments...))
	}
}

// getGOPATHs returns parsed GOPATH or its default, using "/" as path separator.
func getGOPATHs() []string {
	var out []string
	if gp := os.Getenv("GOPATH"); gp != "" {
		for _, v := range filepath.SplitList(gp) {
			// Disallow non-absolute paths?
			if v != "" {
				v = strings.Replace(v, "\\", "/", -1)
				// Trim trailing "/".
				if l := len(v); v[l-1] == '/' {
					v = v[:l-1]
				}
				out = append(out, v)
			}
		}
	}
	if len(out) == 0 {
		homeDir := ""
		u, err := user.Current()
		if err != nil {
			homeDir = os.Getenv("HOME")
			if homeDir == "" {
				panic(fmt.Sprintf("Could not get current user or $HOME: %s\n", err.Error()))
			}
		} else {
			homeDir = u.HomeDir
		}
		out = []string{strings.Replace(homeDir+"/go", "\\", "/", -1)}
	}
	return out
}
