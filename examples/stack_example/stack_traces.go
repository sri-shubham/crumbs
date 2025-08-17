package stack_example

import (
	"context"
	"fmt"

	"github.com/sri-shubham/crumbs"
)

// DemonstrateStackTraces shows how to work with stack traces
func DemonstrateStackTraces() {
	// Enable stack traces globally (disabled by default)
	fmt.Println("Stack Traces Example")
	fmt.Println("===================")

	// Store the original value to restore it later
	originalSetting := crumbs.CaptureStack

	fmt.Println("1. Enabling stack traces globally:")
	crumbs.CaptureStack = true

	ctx := context.Background()

	// Create an error with stack trace
	err1 := level1Function(ctx)

	// Print error with stack trace
	fmt.Println("\nError with stack trace:")
	fmt.Println(crumbs.FormatError(err1, true, true))

	// Disable stack traces
	fmt.Println("\n2. Disabling stack traces globally:")
	crumbs.CaptureStack = false

	// Create an error without stack trace
	err2 := level1Function(ctx)

	// The error won't have a stack trace
	fmt.Println("\nError without stack trace:")
	fmt.Println(crumbs.FormatError(err2, true, true))

	// But we can force a stack trace for critical errors
	fmt.Println("\n3. Forcing a stack trace for a critical error:")
	err3 := level1Function(ctx)

	if cerr, ok := err3.(*crumbs.Error); ok {
		cerr = cerr.ForceStack()
		err3 = cerr
	}

	fmt.Println("\nError with forced stack trace:")
	fmt.Println(crumbs.FormatError(err3, true, true))

	// Configure stack trace depth
	fmt.Println("\n4. Configuring stack trace depth:")
	crumbs.CaptureStack = true

	// Set a smaller stack depth to get less frames
	originalDepth := crumbs.StackTraceDepth
	crumbs.StackTraceDepth = 2

	err4 := level1Function(ctx)

	fmt.Println("\nError with limited stack depth:")
	fmt.Println(crumbs.FormatError(err4, true, true))

	// Restore original settings
	crumbs.CaptureStack = originalSetting
	crumbs.StackTraceDepth = originalDepth
}

// Helper functions to create a call stack
func level1Function(ctx context.Context) error {
	return level2Function(ctx)
}

func level2Function(ctx context.Context) error {
	return level3Function(ctx)
}

func level3Function(ctx context.Context) error {
	return crumbs.New(ctx, "error at level 3", "level", 3)
}
