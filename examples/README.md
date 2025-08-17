# Crumbs Library Examples

This directory contains examples demonstrating how to use the crumbs error handling library.

## Running the Examples

To run all examples:

```bash
go run main.go
```

## Example Categories

### Basic Usage

Basic examples showing how to create and wrap errors with key-value pairs:
- Creating simple errors
- Adding key-value pairs to errors
- Wrapping existing errors
- Using `errors.Is` with wrapped errors
- Formatted error creation with `Errorf` and `Wrapf`

### Context and Crumbs

Shows how to work with context:
- Adding crumbs to a context
- Propagating crumbs through functions
- Automatically capturing context crumbs in errors
- Extracting crumbs from errors and contexts

### Stack Traces

Demonstrates stack trace capabilities:
- Enabling/disabling stack traces globally
- Configuring stack trace depth
- Forcing stack traces for specific errors
- Formatting and displaying stack traces

### Logging Integration

Shows how to integrate with logging libraries:
- Extracting crumbs from errors for structured logging
- Including stack traces in logs
- Working with error chains and wrapped errors

### Standard Library Errors Integration

Demonstrates compatibility with the standard library errors package:
- Using `errors.Is` to compare with sentinel errors
- Using `errors.As` to extract custom error types
- Working with deeply nested error chains
- Using `errors.Unwrap` to walk the error chain
- Combining with `errors.Join` (Go 1.20+)

## Key Concepts

- **Crumbs**: Key-value pairs that provide context to errors
- **Context Integration**: Automatically including context information in errors
- **Stack Traces**: Optional capturing of stack traces for debugging
- **Standard Library Compatibility**: Works with `errors.Is` and `errors.As`
