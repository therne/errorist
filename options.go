package errorist

import "log"

var (
	globalOptions   []Option
	pkgLocalOptions = map[string][]Option{}
)

type Options struct {
	// Logger specifies logging function used on functions end with "WithErrLog".
	// It uses `log.Println` by default.
	Logger func(err string)

	// DetailedStacktrace specifies verbosity of stacktrace.
	// If it's true, running goroutine and its traces will be dumped.
	// Otherwise, it only dumps package and function name with source lines, similar to Java's.
	DetailedStacktrace bool

	// SkipNonProjectFiles specifies whether to skip stacktrace from non-project sources.
	// true by default.
	//
	// It is not recommended to use this option for production because it refers
	// GOPATH from environment variable for the decision. IncludedPackages is recommended.
	SkipNonProjectFiles bool

	// IncludedPackages specifies allowed list of package names in stacktrace.
	// If set, only packages starting with given names will be included.
	IncludedPackages []string

	// WrapArguments specifies additional context info added on error.
	// If first argument is string, it is used to format message with rest of the arguments and
	// will be passed to errors.Wrapf (by default) or fmt.Errorf (optional).
	WrapArguments []interface{}

	// WrapWithFmtErrorf specifies whether to use fmt.Errorf on error wrapping. false by default.
	WrapWithFmtErrorf bool
}

var DefaultOptions = Options{
	Logger: func(err string) { log.Println(err) },

	DetailedStacktrace:  false,
	SkipNonProjectFiles: true,
}

type Option func(o *Options)

// Wrapf is an option for adding context info to the error by wrapping it.
func Wrapf(format string, args ...interface{}) Option {
	return func(o *Options) {
		o.WrapArguments = []interface{}{format}
		o.WrapArguments = append(o.WrapArguments, args...)
	}
}

// WrapWithFmtErrorf is an option for using fmt.Errorf on error wrapping.
func WrapWithFmtErrorf() Option {
	return func(o *Options) {
		o.WrapWithFmtErrorf = true
	}
}

// WithDetailedTrace is an option for dumping running goroutine and its traces will be dumped.
func WithDetailedTrace() Option {
	return func(o *Options) {
		o.DetailedStacktrace = true
	}
}

// WithDetailedTrace is an option for including non-project files to stacktrace.
// It is not recommended to use this option for production because it refers
// GOPATH from environment variable for the decision. IncludedPackages is recommended.
func IncludeNonProjectFiles() Option {
	return func(o *Options) {
		o.SkipNonProjectFiles = true
	}
}

// IncludedPackages is an option specifying allowed list of package names in stacktrace.
// If set, only packages starting with given names will be included.
func IncludedPackages(pkgs ...string) Option {
	return func(o *Options) {
		o.IncludedPackages = pkgs
	}
}

// LogrusLikeLoggingFunc includes leveled logging functions on Logrus.
// https://github.com/sirupsen/logrus#level-logging
// Other loggers sharing same function signature can be also used.
type LogrusLikeLoggingFunc func(...interface{})

// LogWithLogrus is an option for using logrus on functions end with "WithErrorLog".
func LogWithLogrus(lf LogrusLikeLoggingFunc) Option {
	return func(o *Options) {
		o.Logger = func(err string) { lf(err) }
	}
}

type PrintfFamily func(string, ...interface{})

// LogWithPrintfFamily is an option for using printf-like loggers on functions end with "WithErrorLog".
func LogWithPrintfFamily(lf PrintfFamily) Option {
	return func(o *Options) {
		o.Logger = func(err string) { lf(err) }
	}
}

// WithLogHandler is an option for specifying custom logger on functions end with "WithErrorLog".
func WithLogHandler(handler func(err string)) Option {
	return func(o *Options) {
		o.Logger = handler
	}
}

// SetPackageLevelOptions sets options applied in current package scope.
// It can override global options.
func SetGlobalOptions(opts ...Option) {
	globalOptions = opts
}

// SetPackageLevelOptions sets options applied in current package scope.
// It can override global options.
func SetPackageLevelOptions(opts ...Option) {
	pkg := callerPackageName(1)
	pkgLocalOptions[pkg] = opts
}

func applyOptions(opts []Option) Options {
	var merged []Option
	merged = append(merged, globalOptions...)
	merged = append(merged, pkgLocalOptions[callerPackageName(1)]...)
	merged = append(merged, opts...)

	o := DefaultOptions
	for _, optFn := range merged {
		optFn(&o)
	}
	return o
}
