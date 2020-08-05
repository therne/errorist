package errorist

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/maruel/panicparse/stack"
	"github.com/pkg/errors"
)

const maxTraces = 300

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Stacktrace returns pretty-formatted stack trace of errors created or wrapped by
// `github.com/pkg/errors` library. Runtime stack traces are skipped for simplicity.
func Stacktrace(err error, opt ...Option) (traceEntries []string) {
	tr, ok := err.(stackTracer)
	if !ok {
		return []string{callerTrace(1)}
	}
	for _, t := range tr.StackTrace() {
		trace := fmt.Sprintf("%+v", t)
		if strings.HasPrefix(trace, "runtime") {
			continue
		}
		fnName := strings.Split(trace, "\n\t")[0]
		paths := strings.Split(trace, "/")
		srcFile := paths[len(paths)-1]

		traceEntries = append(traceEntries, fmt.Sprintf("%s (%s)", fnName, srcFile))
	}
	return traceEntries
}

func stacktrace(skip, limit int, opts Options) []string {
	if opts.DetailedStacktrace {
		return detailedStacktrace(skip+1, limit, opts)
	}
	return simpleStacktrace(skip+1, limit, opts)
}

func simpleStacktrace(skip, limit int, opts Options) (traces []string) {
	goPaths := getGOPATHs()

	pc := make([]uintptr, maxTraces)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		return []string{"unknown"}
	}
	frames := runtime.CallersFrames(pc)
	for i := 0; i < limit; i++ {
		frame, more := frames.Next()
		if !more {
			break
		}
		if opts.SkipNonProjectFiles && isNonProjectFile(goPaths, frame.File) {
			continue
		}
		st := fmt.Sprintf("%s (%s:%d)", frame.Function, filepath.Base(frame.File), frame.Line)
		traces = append(traces, st)
	}
	return traces
}

func detailedStacktrace(skip, limit int, opts Options) (traces []string) {
	st := make([]byte, 1024)
	for {
		n := runtime.Stack(st, false)
		if n < len(st) {
			st = st[:n]
			break
		}
		st = make([]byte, 2*len(st))
	}
	c, err := stack.ParseDump(bytes.NewReader(st), os.Stdout, true)
	if err != nil {
		return append(
			simpleStacktrace(skip+1, limit, opts),
			fmt.Sprintf("warning: error occurred while dumping detailed stacktrace: %v", err),
		)
	}
	goPaths := getGOPATHs()

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
		}
	}

	for i, bucket := range buckets {
		curLine := ""
		panicIndex := -1
		for i, call := range bucket.Stack.Calls {
			if call.Func.Name() == "panic" {
				panicIndex = i
			}
		}
		if i == 0 {
			// remove stacks before main panic
			bucket.Stack.Calls = bucket.Stack.Calls[skip+panicIndex+1:]
		}

		// Print the goroutine header.
		var tags []string
		if s := bucket.SleepString(); s != "" {
			tags = append(tags, s+" sleeping ")
		}
		if bucket.Locked {
			tags = append(tags, "locked ")
		}
		extra := fmt.Sprintf(
			"[%screated by %s (%s:%d)]",
			strings.Join(tags, ", "),
			bucket.CreatedBy.Func.PkgDotName(),
			bucket.CreatedBy.SrcName(),
			bucket.CreatedBy.Line,
		)
		goroutineIDs := ""
		if len(bucket.IDs) < 3 {
			var ids []string
			for _, id := range bucket.IDs {
				ids = append(ids, fmt.Sprintf("#%d", id))
			}
			goroutineIDs = "Goroutine " + strings.Join(ids, ", ")
		} else {
			goroutineIDs = fmt.Sprintf("%d simillar goroutines", len(bucket.IDs))
		}

		curLine += fmt.Sprintf("%s: %s %s", goroutineIDs, bucket.State, extra)
		if len(traces) >= maxTraces {
			traces = append(traces, curLine+" (...)")
			continue
		}
		traces = append(traces, curLine)

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			if opts.SkipNonProjectFiles && isNonProjectFile(goPaths, line.LocalSrcPath) {
				continue
			}
			traces = append(traces, fmt.Sprintf(
				"    %-*s  %s(%s)",
				srcLen, line.SrcLine(),
				line.Func.PkgDotName(), &line.Args,
			))
		}
		if bucket.Stack.Elided {
			traces = append(traces, "    (...)")
		}
	}
	return traces
}

func isNonProjectFile(goPaths []string, absSrcPath string) bool {
	for _, gopath := range goPaths {
		goModRoot := filepath.Join(gopath, "pkg/mod")
		if strings.HasPrefix(absSrcPath, goModRoot) {
			return true
		}
	}
	return false
}

func formatStacktrace(traces []string, opts Options) string {
	if !opts.DetailedStacktrace {
		var indented []string
		for i, trace := range traces {
			if strings.HasPrefix(trace, "runtime.gopanic") {
				indented = indented[i:]
				continue
			}
			indented = append(indented, strings.Repeat(" ", 4) + trace)
		}
		return strings.Join(indented, "\n")
	}
	return "\n" + strings.Join(traces, "\n")
}
