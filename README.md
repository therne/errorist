<img src="https://raw.githubusercontent.com/therne/errorist/master/docs/logo.png" alt="errorist" height="300px" />


[![Godoc Reference](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/github.com/therne/errorist)
![MIT License Badge](https://img.shields.io/github/license/therne/errorist?style=for-the-badge)

Package errorist provides useful error handling utilities inspired by
[Thanos Coding Style Guide](https://thanos.io/contributing/coding-style-guide.md/#defers-don-t-forget-to-check-returned-errors)
and [Uber go style guide](https://github.com/uber-go/guide/blob/master/style.md).

## Closing Resources

errorist provides hassle-free functions closing `io.Closer` with `defer`.

### With Error Capture

The most recommended way is using error capture. An error caused by `Close` will be captured on the given `error` pointer
 unless it already has an error as the value.

```go
func printFile(path string) (err error) {
    f, err := os.OpenFile(path)
    if err != nil {
        return errors.Wrapf(err, "open %s", path)
    }
    defer errorist.CloseWithErrCapture(f, &err)
    ...
}
```

### With Error Channel

An error also can be captured and sent to an error channel (`chan error`). It is a good fit with resources in goroutines.

```go
errChan := make(chan error)
go func() {
    f, err := os.OpenFile(path)
    defer errorist.CloseWithErrChan(f, errChan)
}()
```


### With Error Log

Otherwise, why don't we just log and ignore it? Default logger is `log.Println` but you are able to customize it with options.


```go
defer errorist.CloseWithLogOnErr(f)

// with logrus
defer errorist.CloseWithLogOnErr(f, errorist.LogWithLogrus(logger.Warn))
```

### Adding Contexts with Error Wrapping

If you're familiar with `errors.Wrap` or `fmt.Errorf`, you may want to do the same error handling with errorist.
errorist provides the option wrapping errors with context.

```go
func printFile(path string) (err error) {
    f, err := os.OpenFile(path)
    if err != nil {
        return errors.Wrapf(err, "open %s", path)
    }
    defer errorist.CloseWithErrCapture(f, &err, errorist.Wrapf("close %s", path))
    ...
}
```

> Note that `errorist.Wrapf` is just an option specifier for errorist; it cannot be used solely.
If you want to just wrap errors, you can use [pkg/errors](http://github.com/pkg/errors) or `fmt.Errorf` with `%w` pattern added in Go 1.13.


## Recovering from Panics

errorist provides panic recovery functions that can be used with `defer`.

```go
wg, ctx := errgroup.WithContext(context.Background())
wg.Go(func() (err error) {
    defer errorist.RecoverWithErrCapture(&err)
    ...
})
```

Stacktrace is prettified, and calls from non-project source will be filtered by default.
You can customize stacktrace format with options. For details, please refer
[options.go](https://github.com/therne/errorist/blob/master/options.go).

```go
fmt.Println("%+v", err)
```

```
panic: assignment to entry in nil map
    github.com/therne/errorist.panicStation (panic_test.go:67)
    github.com/therne/errorist.TestWrapPanicWith.func1.1.1 (panic_test.go:19)
    github.com/therne/errorist.TestWrapPanicWith.func1.1 (panic_test.go:13)
    github.com/therne/errorist.TestWrapPanicWith.func1 (panic_test.go:12)
    github.com/therne/errorist.TestWrapPanicWith (panic_test.go:11)
```

## Prettifying Stacktraces on Errors

[pkg/errors](http://github.com/pkg/errors) is the most popular and powerful tool for handling and wrapping errors.
errorist provides extracting and prettifying a stacktrace from errors created or wrapped by pkg/errors.

```go
errorist.Stacktrace(err) ==
    []string{
        "github.com/some/app/api/Server.Handle (api.go:84)",
        "github.com/some/app/controller/Controller.Get (controller.go:11)",
    }
```

The below example shows how to use it in API responses for debugging purposes.

```go
func handleGinError(c *gin.Context) {
    c.Next()
    if len(c.Errors) > 0 {
        err := c.Errors.Last().Err
        c.JSON(500, gin.H{
            "error":       err.Error(),
            "stacktraces": errorist.Stacktrace(err)
        })
    }
}
```

## Options

You can use global options, package-wide, or with call arguments. Options set on a smaller scope can override options set on a wider scope.

```go
errorist.SetGlobalOptions(
    errorist.IncludeNonProjectFiles(),
    errorist.WithDetailedTrace(),
)

// package-local options can override global options
errorist.SetPackageLevelOptions(
    errorist.IncludedPackages("github.com/my/pkg"),
    errorist.LogWithLogrus(pkgLogger),
)

// call options can override options above
errorist.CloseWithErrCapture(f, &err, errorist.Wrapf("close %s", path))
```

For detailed options, please refer [Godoc](https://pkg.go.dev/github.com/therne/errorist?tab=doc#Options) or [options.go](https://github.com/therne/errorist/blob/master/options.go).

###### License: MIT
