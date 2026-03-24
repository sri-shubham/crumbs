package crumbs

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

var (
	benchErr      = errors.New("benchmark error")
	benchResult   error
	benchCtxCrumb context.Context
)

// Benchmarks for error creation
func BenchmarkErrorsNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchResult = errors.New("benchmark error")
	}
}

func BenchmarkCrumbsNewError(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchResult = NewError(ctx, "benchmark error")
	}
}

func BenchmarkCrumbsNewErrorWithCrumbs(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchResult = NewError(ctx, "benchmark error",
			"key1", "value1",
			"key2", 2,
			"key3", true)
	}
}

// Benchmarks for error wrapping
func BenchmarkErrorsWrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchResult = fmt.Errorf("wrapped: %w", benchErr)
	}
}

func BenchmarkCrumbsWrapError(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchResult = WrapError(ctx, benchErr, "wrapped")
	}
}

func BenchmarkCrumbsWrapErrorWithCrumbs(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchResult = WrapError(ctx, benchErr, "wrapped",
			"key1", "value1",
			"key2", 2,
			"key3", true)
	}
}

// Benchmarks for context operations
func BenchmarkAddCrumb(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchCtxCrumb = AddCrumb(ctx, "key", "value")
	}
}

func BenchmarkAddMultipleCrumbs(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		benchCtxCrumb = AddCrumb(ctx,
			"key1", "value1",
			"key2", 2,
			"key3", true)
	}
}

func BenchmarkGetCrumbs(b *testing.B) {
	ctx := context.Background()
	ctx = AddCrumb(ctx,
		"key1", "value1",
		"key2", 2,
		"key3", true)

	b.ResetTimer()
	var result []Crumb
	for i := 0; i < b.N; i++ {
		result = GetCrumbs(ctx)
	}
	_ = result
}

// Benchmarks with stack traces
func BenchmarkNewErrorWithStackTraceEnabled(b *testing.B) {
	ctx := context.Background()
	origSetting := captureStack
	captureStack = true
	defer func() { captureStack = origSetting }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResult = NewError(ctx, "benchmark error")
	}
}

func BenchmarkNewErrorWithStackTraceDisabled(b *testing.B) {
	ctx := context.Background()
	origSetting := captureStack
	captureStack = false
	defer func() { captureStack = origSetting }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResult = NewError(ctx, "benchmark error")
	}
}

// Benchmark error formatting
func BenchmarkFormatError(b *testing.B) {
	ctx := context.Background()
	err := NewError(ctx, "benchmark error", "key1", "value1", "key2", 2)

	b.ResetTimer()
	var result string
	for i := 0; i < b.N; i++ {
		result = FormatError(err, false, true)
	}
	_ = result
}

func BenchmarkFormatErrorWithStack(b *testing.B) {
	ctx := context.Background()
	origSetting := captureStack
	captureStack = true
	err := NewError(ctx, "benchmark error", "key1", "value1", "key2", 2)
	captureStack = origSetting

	b.ResetTimer()
	var result string
	for i := 0; i < b.N; i++ {
		result = FormatError(err, true, true)
	}
	_ = result
}
