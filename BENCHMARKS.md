# Benchmarks for Crumbs

This document explains the benchmarks available in the Crumbs library and how to run them.

## Available Benchmarks

The benchmarks in `benchmark_test.go` compare the performance of:

1. **Standard Go errors vs Crumbs errors**
   - Creating new errors
   - Wrapping existing errors

2. **Context operations**
   - Adding single crumbs to context
   - Adding multiple crumbs at once
   - Retrieving crumbs from context

3. **Stack trace impact**
   - Error creation with stack traces enabled
   - Error creation with stack traces disabled

4. **Error formatting**
   - Formatting errors with and without stack traces

## Running the Benchmarks

You can run benchmarks using the following methods:

### Using Go directly

```bash
# Run all benchmarks
go test -bench=. -benchmem ./...

# Run specific benchmarks (using regex pattern)
go test -bench=BenchmarkNew -benchmem ./...
```

### Using Makefile

```bash
# Run all benchmarks
make benchmark
```

## Interpreting Results

The benchmark results show:

- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes of memory allocated per operation (lower is better)
- **allocs/op**: Number of heap allocations per operation (lower is better)

Example output:
```
BenchmarkErrorsNew-8                  14602931        82.00 ns/op        48 B/op        1 allocs/op
BenchmarkCrumbsNew-8                   4672731       256.4 ns/op        208 B/op        3 allocs/op
```

## Performance Considerations

- The Crumbs library adds additional context and capabilities to errors, which comes with some overhead
- Stack trace capture is particularly expensive - only enable it when needed
- For high-performance applications, consider the tradeoff between rich error information and raw performance
