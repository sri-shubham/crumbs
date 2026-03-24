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
	// captureStack controls whether stack traces are captured globally
	captureStack = false

	// stackTraceDepth controls how many frames to capture (0 for unlimited)
	stackTraceDepth = 32
)

// ConfigureStackTraces sets the global stack trace configuration.
// It is recommended to only call this during application initialization.
func ConfigureStackTraces(capture bool, depth int) {
	captureStack = capture
	stackTraceDepth = depth
}

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
	Value any
}

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
	depth := stackTraceDepth
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

// NewError creates a new error with the given message and key-value pairs.
// It captures the stack trace if ConfigureStackTraces has enabled it.
func NewError(ctx context.Context, msg string, kv ...interface{}) error {
	return newError(ctx, nil, msg, kv...)
}

// WrapError creates an Error with a message and key-value pairs.
func WrapError(ctx context.Context, err error, msg string, kv ...interface{}) error {
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
	var errCrumbs []Crumb
	if err != nil {
		var cerr *Error
		if errors.As(err, &cerr) {
			errCrumbs = cerr.Crumbs
		}
	}

	var ctxCrumbs []Crumb
	if ctx != nil {
		if c, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
			ctxCrumbs = c
		}
	}

	kvLen := (len(kv) + 1) / 2
	totalCap := len(errCrumbs) + len(ctxCrumbs) + kvLen

	var crumbs []Crumb
	if totalCap > 0 {
		crumbs = make([]Crumb, 0, totalCap)
		if len(errCrumbs) > 0 {
			crumbs = append(crumbs, errCrumbs...)
		}
		if len(ctxCrumbs) > 0 {
			crumbs = append(crumbs, ctxCrumbs...)
		}

		for i := 0; i < len(kv); i += 2 {
			if i+1 < len(kv) {
				key, ok := kv[i].(string)
				if !ok {
					continue
				}
				crumbs = append(crumbs, Crumb{Key: key, Value: kv[i+1]})
			} else {
				key, ok := kv[i].(string)
				if ok {
					crumbs = append(crumbs, Crumb{Key: "!BADKEY", Value: key})
				} else {
					crumbs = append(crumbs, Crumb{Key: "!BADKEY", Value: kv[i]})
				}
			}
		}
	}

	e := &Error{
		Err:    err,
		Msg:    msg,
		Crumbs: crumbs,
	}

	// Inherit (keep) the root stack trace from the underlying error if it already has one
	var hasStack bool
	if err != nil {
		var cerr *Error
		if errors.As(err, &cerr) && len(cerr.Stack) > 0 {
			e.Stack = cerr.Stack
			hasStack = true
		}
	}

	if captureStack && !hasStack {
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

	if existing, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
		crumbs = make([]Crumb, len(existing), len(existing)+(len(kv)+1)/2)
		copy(crumbs, existing)
	} else {
		crumbs = make([]Crumb, 0, (len(kv)+1)/2)
	}

	for i := 0; i < len(kv); i += 2 {
		if i+1 < len(kv) {
			key, ok := kv[i].(string)
			if !ok {
				continue // Skip if key is not a string
			}
			crumbs = append(crumbs, Crumb{Key: key, Value: kv[i+1]})
		} else {
			key, ok := kv[i].(string)
			if ok {
				crumbs = append(crumbs, Crumb{Key: "!BADKEY", Value: key})
			} else {
				crumbs = append(crumbs, Crumb{Key: "!BADKEY", Value: kv[i]})
			}
		}
	}

	return context.WithValue(ctx, crumbsKey{}, crumbs)
}

// GetCrumbs retrieves all crumbs from a context
func GetCrumbs(ctx context.Context) []Crumb {
	if ctx == nil {
		return nil
	}

	if crumbs, ok := ctx.Value(crumbsKey{}).([]Crumb); ok {
		result := make([]Crumb, len(crumbs))
		copy(result, crumbs)
		return result
	}

	return nil
}

// FormatError returns a detailed string representation of the error,
// formatted with crumbs and stack trace belonging to the outermost wrapper.
func FormatError(err error, includeStack bool, includeCrumbs bool) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(err.Error())

	if !includeStack && !includeCrumbs {
		return sb.String()
	}

	var cerr *Error
	if errors.As(err, &cerr) {
		if includeCrumbs && len(cerr.Crumbs) > 0 {
			sb.WriteString("\nCrumbs:")
			for _, crumb := range cerr.Crumbs {
				sb.WriteString(fmt.Sprintf("\n  %s: %v", crumb.Key, crumb.Value))
			}
		}

		if includeStack && len(cerr.Stack) > 0 {
			sb.WriteString("\nStack trace:")
			for _, frame := range cerr.Stack {
				sb.WriteString(fmt.Sprintf("\n  %s\n    %s:%d", frame.Function, frame.File, frame.Line))
			}
		}
	}

	return sb.String()
}
