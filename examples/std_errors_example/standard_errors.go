package std_errors_example

import (
	"context"
	"errors"
	"fmt"

	"github.com/sri-shubham/crumbs"
)

// Define some sentinel errors for comparison
var (
	ErrNotFound   = errors.New("resource not found")
	ErrPermission = errors.New("permission denied")
	ErrTimeout    = errors.New("operation timed out")
)

// CustomError is a domain-specific error type for examples
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("custom error %d: %s", e.Code, e.Message)
}

// DemonstrateStandardErrorsMethods shows how to use standard library errors methods with crumbs
func DemonstrateStandardErrorsMethods() {
	fmt.Println("Standard Library Errors Integration Example")
	fmt.Println("=========================================")

	ctx := context.Background()

	// Example 1: errors.Is with crumbs
	fmt.Println("\n1. Using errors.Is with crumbs:")

	// Create a chain of wrapped errors
	notFoundErr := crumbs.WrapError(ctx, ErrNotFound, "user profile not found",
		"userID", "12345",
		"source", "database")

	timeoutErr := crumbs.WrapError(ctx, ErrTimeout, "database query timed out",
		"query", "SELECT * FROM users WHERE id = ?",
		"timeout", "5s")

	// Check if errors match the sentinel errors
	fmt.Printf("Is notFoundErr a NotFound error? %v\n", errors.Is(notFoundErr, ErrNotFound))
	fmt.Printf("Is notFoundErr a Permission error? %v\n", errors.Is(notFoundErr, ErrPermission))
	fmt.Printf("Is timeoutErr a Timeout error? %v\n", errors.Is(timeoutErr, ErrTimeout))

	// Example 2: Deep wrapping with errors.Is
	fmt.Println("\n2. Deep wrapping with errors.Is:")

	// Create a deeply nested error chain
	deepErr := crumbs.WrapError(ctx,
		crumbs.WrapError(ctx,
			crumbs.WrapError(ctx, ErrPermission, "access check failed",
				"resource", "file-123"),
			"user verification failed",
			"method", "OAuth"),
		"API request unauthorized",
		"endpoint", "/api/admin")

	// Check if the deeply nested error still matches the sentinel error
	fmt.Printf("Is deepErr a Permission error? %v\n", errors.Is(deepErr, ErrPermission))
	fmt.Printf("Formatted error chain:\n%s\n", crumbs.FormatError(deepErr, false, true))

	// Example 3: Using errors.As with custom error types
	fmt.Println("\n3. Using errors.As with custom error types:")

	// Create a custom error and wrap it with crumbs
	originalErr := &CustomError{Code: 404, Message: "user not found"}
	wrappedCustomErr := crumbs.WrapError(ctx, originalErr, "user lookup failed",
		"username", "alice",
		"method", "ldap")

	// Use errors.As to extract the custom error type
	var customErr *CustomError
	if errors.As(wrappedCustomErr, &customErr) {
		fmt.Printf("Successfully extracted custom error using errors.As\n")
		fmt.Printf("Error code: %d\n", customErr.Code)
		fmt.Printf("Error message: %s\n", customErr.Message)
	} else {
		fmt.Printf("Failed to extract custom error\n")
	}

	// Example 4: Using errors.Unwrap
	fmt.Println("\n4. Using errors.Unwrap:")

	// Create a wrapped error
	baseErr := errors.New("base error")
	wrappedErr := crumbs.WrapError(ctx, baseErr, "wrapped message",
		"key1", "value1")

	// Manually unwrap and print each layer
	fmt.Println("Error chain:")
	currentErr := wrappedErr
	for i := 1; currentErr != nil; i++ {
		fmt.Printf("  Layer %d: %v\n", i, currentErr)
		currentErr = errors.Unwrap(currentErr)
	}

	// Example 5: Combining with Go 1.20+ error joining
	fmt.Println("\n5. Using with error joining (Go 1.20+):")

	err1 := crumbs.NewError(ctx, "first error", "order", 1)
	err2 := crumbs.NewError(ctx, "second error", "order", 2)
	err3 := crumbs.NewError(ctx, "third error", "order", 3)

	// Join errors (Go 1.20+)
	joinedErr := errors.Join(err1, err2, err3)

	// Wrap the joined errors with crumbs
	wrappedJoinedErr := crumbs.WrapError(ctx, joinedErr, "multiple errors occurred",
		"errorCount", 3,
		"operation", "batch processing")

	fmt.Printf("Joined and wrapped errors:\n%v\n", wrappedJoinedErr)

	// Check if the wrapped joined errors contain the original errors
	fmt.Printf("Contains err1? %v\n", errors.Is(wrappedJoinedErr, err1))
	fmt.Printf("Contains err2? %v\n", errors.Is(wrappedJoinedErr, err2))
	fmt.Printf("Contains err3? %v\n", errors.Is(wrappedJoinedErr, err3))
}
