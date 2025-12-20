package crumbs

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Package level configuration settings
var (
	// CaptureStack controls whether stack traces are captured
	CaptureStack = false

	// StackTraceDepth controls how many frames to capture (0 for unlimited)
	StackTraceDepth = 32
)

type crumbsKey struct{}

// StackFrame represents a single frame in the stack trace
type StackFrame struct {
	Function string
	File     string
	Line     int
}

// Crumb represents a key-value pair
type Crumb struct {
	Key   string
	Value interface{}
}

// Empty line to replace the removed category code

// Error is a custom error type that wraps a standard error and supports key-value pairs.
type Error struct {
	Err    error
	Msg    string
	Crumbs []Crumb
	Stack  []StackFrame
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Msg == "" {
		if e.Err != nil {
			return e.Err.Error()
		}
		return "unknown error"
	}

	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}

	return e.Msg
}

// Unwrap returns the underlying error for errors.Is and errors.As compatibility.
func (e *Error) Unwrap() error {
	return e.Err
}

// GetCrumbs returns the key-value pairs associated with the error.
func (e *Error) GetCrumbs() []Crumb {
	return e.Crumbs
}

// GetStack returns the stack trace captured when the error was created.
func (e *Error) GetStack() []StackFrame {
	return e.Stack
}

// FormatStack returns a formatted string representation of the stack trace.
func (e *Error) FormatStack() string {
	if len(e.Stack) == 0 {
		return "no stack trace available"
	}

	var sb strings.Builder
	for i, frame := range e.Stack {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d", frame.Function, frame.File, frame.Line))
	}
	return sb.String()
}

// ForceStack forces a stack trace to be captured even if disabled globally.
// This is useful for critical errors where you always want stack traces.
// Returns the same error for chaining.
func (e *Error) ForceStack() *Error {
	if len(e.Stack) == 0 {
		e.captureStack(3) // Skip this function + caller
	}
	return e
}

// captureStack captures the current stack trace.
func (e *Error) captureStack(skip int) {
	depth := StackTraceDepth
	if depth <= 0 {
		depth = 32 // Reasonable default if set to unlimited
	}

	pcs := make([]uintptr, depth)
	n := runtime.Callers(skip, pcs)

	frames := runtime.CallersFrames(pcs[:n])
	e.Stack = make([]StackFrame, 0, n)

	for {
		frame, more := frames.Next()

		// Skip runtime and standard library frames
		if strings.Contains(frame.File, "runtime/") {
			if more {
				continue
			}
			break
		}

		e.Stack = append(e.Stack, StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}
}

// New creates a new Error with context information and key-value pairs.
func New(ctx context.Context, msg string, kv ...interface{}) error {
	return newError(ctx, nil, msg, kv...)
}

// Wrap creates an Error with a message and key-value pairs.
func Wrap(ctx context.Context, err error, msg string, kv ...interface{}) error {
	if err == nil {
		return nil
	}
	return newError(ctx, err, msg, kv...)
}

// Errorf creates a new Error with formatted message and context information.
func Errorf(ctx context.Context, format string, args ...interface{}) error {
	return newError(ctx, nil, fmt.Sprintf(format, args...), nil)
}

// Wrapf creates an Error with formatted message and context information.
func Wrapf(ctx context.Context, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return newError(ctx, err, fmt.Sprintf(format, args...), nil)
}

// newError is a helper function to create errors with context
func newError(ctx context.Context, err error, msg string, kv ...interface{}) error {
	// Create crumbs slice
	var crumbs []Crumb

	// First, add any crumbs from the context
	if ctx != nil {
		if ctxCrumbs, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
			crumbs = make([]Crumb, len(ctxCrumbs), len(ctxCrumbs)+len(kv)/2)
			copy(crumbs, ctxCrumbs)
		}
	}

	// If crumbs is still nil (no context crumbs), initialize it
	if crumbs == nil {
		crumbs = make([]Crumb, 0, len(kv)/2)
	}

	// Then add any key-value pairs passed directly
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue
		}
		crumbs = append(crumbs, Crumb{Key: key, Value: kv[i+1]})
	}

	// Use the global configuration for stack trace capture
	captureStack := CaptureStack

	// Create the error
	e := &Error{
		Err:    err,
		Msg:    msg,
		Crumbs: crumbs,
	}

	// Capture the stack trace if needed
	if captureStack {
		e.captureStack(3) // Skip newError, Wrap/New, and caller
	}

	return e
}

// AddCrumb adds multiple crumbs (key-value pairs) to the context
func AddCrumb(ctx context.Context, kv ...interface{}) context.Context {
	if len(kv) == 0 {
		return ctx
	}

	var crumbs []Crumb

	// Get existing crumbs if any
	if existing, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
		// Create a copy to avoid modifying the existing slice in the context
		crumbs = make([]Crumb, len(existing), len(existing)+len(kv)/2)
		copy(crumbs, existing)
	} else {
		crumbs = make([]Crumb, 0, len(kv)/2)
	}

	// Add the new crumbs
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue // Skip if key is not a string
		}
		crumbs = append(crumbs, Crumb{Key: key, Value: kv[i+1]})
	}

	// Return new context with updated crumbs
	return context.WithValue(ctx, crumbsKey{}, crumbs)
}

// GetCrumbs retrieves all crumbs from a context
func GetCrumbs(ctx context.Context) []Crumb {
	if ctx == nil {
		return nil
	}

	if crumbs, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
		// Return a copy to prevent modification
		result := make([]Crumb, len(crumbs))
		copy(result, crumbs)
		return result
	}

	return nil
}
// FormatError returns a detailed string representation of the error,
// optionally including stack trace, crumbs, and category
func FormatError(err error, includeStack bool, includeCrumbs bool) string {
	var sb strings.Builder

	// Basic error message
	sb.WriteString(err.Error())

	// Try to get our custom error
	var cerr *Error
	if errors.As(err, &cerr) {
		// Add crumbs if requested
		if includeCrumbs && len(cerr.Crumbs) > 0 {
			sb.WriteString("\nCrumbs:")
			for _, crumb := range cerr.Crumbs {
				sb.WriteString(fmt.Sprintf("\n  %s: %v", crumb.Key, crumb.Value))
			}
		}

		// Add stack trace if requested
		if includeStack && len(cerr.Stack) > 0 {
			sb.WriteString("\nStack trace:")
			for _, frame := range cerr.Stack {
				sb.WriteString(fmt.Sprintf("\n  %s\n    %s:%d", frame.Function, frame.File, frame.Line))
			}
		}
	}

	return sb.String()
}
