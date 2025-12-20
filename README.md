# Crumbs

[![Go Report Card](https://goreportcard.com/badge/github.com/sri-shubham/crumbs)](https://goreportcard.com/report/github.com/sri-shubham/crumbs)
[![GoDoc](https://godoc.org/github.com/sri-shubham/crumbs?status.svg)](https://godoc.org/github.com/sri-shubham/crumbs)
[![Coverage Status](https://coveralls.io/repos/github/sri-shubham/crumbs/badge.svg?branch=main)](https://coveralls.io/github/sri-shubham/crumbs?branch=main)
[![GitHub Stars](https://img.shields.io/github/stars/sri-shubham/crumbs.svg)](https://github.com/sri-shubham/crumbs/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/sri-shubham/crumbs.svg)](https://github.com/sri-shubham/crumbs/issues)

Crumbs is a lightweight, flexible error handling library for Go that adds context to your errors. It's designed to enhance error reporting while maintaining full compatibility with the standard library.

## Features

- **Key-Value Context**: Attach key-value pairs ("crumbs") to errors for better debugging and logging
- **Context Integration**: Automatically gather context information from Go's `context.Context`
- **Standard Library Compatible**: Works seamlessly with `errors.Is`, `errors.As`, and `errors.Unwrap`
- **Optional Stack Traces**: Record stack traces when needed
- **Logging Friendly**: Easily extract structured data for your logging system

## Installation

```bash
go get github.com/sri-shubham/crumbs
```

## Quick Start

```go
import (
    "context"
    "errors"
    "fmt"
    "github.com/sri-shubham/crumbs"
)

func main() {
    ctx := context.Background()
    
    // Add crumbs to the context
    ctx = crumbs.AddCrumb(ctx,
        "requestID", "req-12345",
        "userID", "user-abc",
    )
    
    // Create a new error with additional crumbs
    err := crumbs.New(ctx, "operation failed",
        "operation", "getData",
        "status", 500,
    )
    
    // Print detailed error with crumbs
    fmt.Println(crumbs.FormatError(err, false, true))
    
    // Works with standard errors package
    baseErr := errors.New("connection failed")
    wrappedErr := crumbs.Wrap(ctx, baseErr, "database error")
    
    if errors.Is(wrappedErr, baseErr) {
        fmt.Println("Error identity preserved!")
    }
}
```

## Core Concepts

### Creating Errors

```go
// Create a new error
err := crumbs.New(ctx, "something went wrong")

// Create with key-value pairs
err := crumbs.New(ctx, "request failed", 
    "status", 404,
    "path", "/users/123",
)

// Create with formatting
err := crumbs.Errorf(ctx, "failed with code %d", 500)
```

### Wrapping Errors

```go
// Wrap an existing error
baseErr := errors.New("network timeout")
err := crumbs.Wrap(ctx, baseErr, "API request failed")

// Wrap with key-value pairs
err := crumbs.Wrap(ctx, baseErr, "database query failed",
    "query", "SELECT * FROM users",
    "params", []string{"id=123"},
)

// Wrap with formatting
err := crumbs.Wrapf(ctx, baseErr, "operation %s failed", "getData")
```

### Working with Context

```go
// Add crumbs to context
ctx = crumbs.AddCrumb(ctx, "userID", "user-123")

// Add multiple crumbs
ctx = crumbs.AddCrumb(ctx, 
    "requestID", "req-abc",
    "traceID", "trace-xyz",
    "timestamp", time.Now(),
)

// Get crumbs from context
allCrumbs := crumbs.GetCrumbs(ctx)
```

### Stack Traces

Stack traces are disabled by default for performance reasons but can be enabled when needed:

```go
// Enable stack traces globally
crumbs.CaptureStack = true

// Configure stack trace depth (0 for unlimited)
crumbs.StackTraceDepth = 32

// Force a stack trace for a specific error
err := crumbs.New(ctx, "critical error").(*crumbs.Error)
err = err.ForceStack()
```

### Error Formatting

```go
// Format error with crumbs
formatted := crumbs.FormatError(err, false, true)

// Format with stack trace
formatted := crumbs.FormatError(err, true, true)
```

### Extracting Data

```go
if cerr, ok := err.(*crumbs.Error); ok {
    // Get all crumbs (returns []Crumb)
    allCrumbs := cerr.GetCrumbs()
    for _, c := range allCrumbs {
        fmt.Printf("Key: %s, Value: %v\n", c.Key, c.Value)
    }
    
    // Get stack trace
    stack := cerr.GetStack()
}
```

## Logging Integration

Crumbs errors work great with structured logging:

```go
func LogError(logger Logger, err error) {
    if err == nil {
        return
    }
    
    msg := err.Error()
    fields := map[string]interface{}{}
    
    var cerr *crumbs.Error
    if errors.As(err, &cerr) {
        // Add all crumbs as fields
        for _, c := range cerr.GetCrumbs() {
            fields[c.Key] = c.Value
        }
    }
    
    logger.Error(msg, fields)
}
```

## Examples

See the [examples](./examples) directory for comprehensive usage examples:

- Basic usage patterns
- Context integration
- Stack traces
- Logging integration
- Standard library errors compatibility

## Benchmarks

Performance is a key consideration in error handling. Below are benchmark results comparing standard errors with Crumbs:

```
goos: darwin
goarch: arm64
pkg: github.com/sri-shubham/crumbs
cpu: Apple M1
BenchmarkErrorsNew-8                    85283263                13.82 ns/op           16 B/op           1 allocs/op
BenchmarkCrumbsNew-8                    45591188                27.23 ns/op           80 B/op           1 allocs/op
BenchmarkCrumbsNewWithCrumbs-8          22543290                56.04 ns/op          176 B/op           2 allocs/op
BenchmarkErrorsWrap-8                   14096294                82.06 ns/op           56 B/op           2 allocs/op
BenchmarkCrumbsWrap-8                   45918513                26.96 ns/op           80 B/op           1 allocs/op
BenchmarkCrumbsWrapWithCrumbs-8         21582700                55.39 ns/op          176 B/op           2 allocs/op
BenchmarkAddCrumb-8                     19820527                60.57 ns/op          104 B/op           3 allocs/op
BenchmarkAddMultipleCrumbs-8            16305519                73.32 ns/op          168 B/op           3 allocs/op
BenchmarkGetCrumbs-8                    43355341                27.47 ns/op           96 B/op           1 allocs/op
BenchmarkNewWithStackTraceEnabled-8       904336              1322 ns/op             784 B/op           4 allocs/op
BenchmarkNewWithStackTraceDisabled-8    45363045                26.17 ns/op           80 B/op           1 allocs/op
BenchmarkFormatError-8                   4431264               271.9 ns/op           184 B/op           8 allocs/op
BenchmarkFormatErrorWithStack-8          1276100               935.1 ns/op          1872 B/op          24 allocs/op
```

For more detailed benchmark information and analysis, see [BENCHMARKS.md](./BENCHMARKS.md).

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on how to contribute to this project.

## License

[MIT License](./LICENSE) - Copyright (c) 2025 Shubham Srivastava
