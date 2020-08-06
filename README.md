errorist
============

[![Godoc Reference](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/github.com/therne/errorist)
![MIT License Badge](https://img.shields.io/github/license/therne/errorist?style=for-the-badge)

Package errorist provides useful error handling utilities, including ones recommended in
[Thanos Coding Style Guide](https://thanos.io/contributing/coding-style-guide.md/#defers-don-t-forget-to-check-returned-errors)
or [Uber go style guide](https://github.com/uber-go/guide/blob/master/style.md).

## Closing Resources

errorist provides resource management utilities for easily closing `io.Closer` with `defer`.

### With Error Capture

Most recommended way is using error capture. Error caused by `Close` will be captured on given `error` pointer
 if there's no error already present on it.

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

Error also can be captured to error channel (`chan error`). It has a good fit with resources managed in goroutines.

```go
errChan := make(chan error)
go func() {
    f, err := os.OpenFile(path)
    defer errorist.CloseWithErrChan(f, errChan)
}()
```


### With Error Log

Otherwise, why can't we just log and ignore it? Default logger is `log.Println` but you can customize it with options.


```go
defer errorist.CloseWithLogOnErr(f)

// with logrus
defer errorist.CloseWithLogOnErr(f, errorist.LogWithLogrus(logger.Warn))
```

### Adding Contexts with Error Wrapping

If you're familiar with `errors.Wrap` or `fmt.Errorf`, you might want to do same thing on errors handling with errorist.
errorist provides option for wrapping errors with context.

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

> Note that `errorist.Wrapf` is just option specifier for errorist; it cannot be used solely.
If you want to just wrap errors, you may use [pkg/errors](http://github.com/pkg/errors) or `fmt.Errorf` with `%w` pattern added in Go 1.13.


## Recovering from Panics

errorist provides panic recovery functions that can be used together with `defer`.

```go
wg, ctx := errgroup.WithContext(context.Background())
wg.Go(func() (err error) {
    defer errorist.RecoverWithErrCapture(&err)
    ...
})
```

Stacktrace is prettified and calls from non-project source will be filtered by default.
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

## Prettifying Stacktraces on Error

[pkg/errors](http://github.com/pkg/errors) is most popular and powerful tool for handling and wrapping errors.
errorist provides extracting and prettifying a stacktrace from errors created and wrapped by pkg/errors.

```go
errorist.Stacktrace(err) ==
    []string{
        "github.com/some/app/api/Server.Handle (api.go:84)",
        "github.com/some/app/controller/Controller.Get (controller.go:11)",
    }
```

An example usecase is using it with API response for debugging errors.

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

You can apply options globally, package-wide, or with call arguments. Options set on smaller scope can override options set on wider scope.

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

For detailed options reference, please refer [Godoc](https://pkg.go.dev/github.com/therne/errorist?tab=doc#Options) or [options.go](https://github.com/therne/errorist/blob/master/options.go).

###### License: MIT
