package basic_example

import (
	"context"
	"errors"
	"fmt"

	"github.com/sri-shubham/crumbs"
)

// DemonstrateBasicUsage shows basic usage of the crumbs library
func DemonstrateBasicUsage() {
	// Create a context
	ctx := context.Background()

	// Example 1: Creating a simple error
	fmt.Println("Example 1: Basic Error Creation")
	err1 := crumbs.New(ctx, "something went wrong")
	fmt.Println(err1)
	fmt.Println()

	// Example 2: Creating an error with key-value pairs
	fmt.Println("Example 2: Error with Key-Value Pairs")
	err2 := crumbs.New(ctx, "failed to process request",
		"requestID", "123456",
		"user", "alice",
		"timestamp", "2025-08-17T10:15:30Z")
	fmt.Println(err2)

	// Print with details using FormatError
	fmt.Println("\nFormatted error with crumbs:")
	fmt.Println(crumbs.FormatError(err2, false, true))
	fmt.Println()

	// Example 3: Wrapping an existing error
	fmt.Println("Example 3: Wrapping Errors")
	baseErr := errors.New("connection refused")
	wrappedErr := crumbs.Wrap(ctx, baseErr, "database connection failed",
		"host", "db.example.com",
		"port", 5432,
		"attempts", 3)
	fmt.Println(wrappedErr)

	// Print with details
	fmt.Println("\nFormatted wrapped error:")
	fmt.Println(crumbs.FormatError(wrappedErr, false, true))
	fmt.Println()

	// Example 4: Using errors.Is with wrapped errors
	fmt.Println("Example 4: Using errors.Is")
	if errors.Is(wrappedErr, baseErr) {
		fmt.Println("wrappedErr contains baseErr - errors.Is works!")
	} else {
		fmt.Println("errors.Is failed")
	}
	fmt.Println()

	// Example 5: Formatted errors
	fmt.Println("Example 5: Formatted Errors")
	err3 := crumbs.Errorf(ctx, "failed with code %d: %s", 404, "not found")
	fmt.Println(err3)

	wrappedErr2 := crumbs.Wrapf(ctx, baseErr, "server %s returned error %d", "api.example.com", 500)
	fmt.Println(wrappedErr2)
	fmt.Println()
}
