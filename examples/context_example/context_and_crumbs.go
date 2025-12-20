package context_example

import (
	"context"
	"fmt"

	"github.com/sri-shubham/crumbs"
)

// simulateRequest simulates an API request with context
func simulateRequest(ctx context.Context, userID string, action string) error {
	// Add some crumbs to the context for the entire request lifecycle
	ctx = crumbs.AddCrumb(ctx,
		"requestID", "req-12345",
		"userID", userID,
		"action", action,
		"timestamp", "2025-08-17T10:30:00Z",
	)

	// First operation
	if err := validateRequest(ctx, userID); err != nil {
		return err
	}

	// Second operation
	if err := processData(ctx); err != nil {
		return err
	}

	// Third operation
	if err := saveResults(ctx); err != nil {
		return err
	}

	return nil
}

func validateRequest(ctx context.Context, userID string) error {
	if userID == "invalid" {
		// Create an error that will automatically include all crumbs from ctx
		return crumbs.New(ctx, "invalid user",
			"validationTime", "2025-08-17T10:30:01Z")
	}
	return nil
}

func processData(ctx context.Context) error {
	// Add another crumb to the context for downstream operations
	ctx = crumbs.AddCrumb(ctx, "dataSize", 1024)

	// Simulate an error
	// The error will include all crumbs from the context
	return crumbs.New(ctx, "processing failed",
		"processingTime", "2025-08-17T10:30:02Z")
}

func saveResults(ctx context.Context) error {
	// This wouldn't execute due to earlier error,
	// but included for completeness
	return crumbs.New(ctx, "save succeeded")
}

// RunExample demonstrates the context and crumbs example
func RunExample() {
	fmt.Println("Context and Crumbs Example")
	fmt.Println("=========================")

	// Create a base context
	ctx := context.Background()

	// Simulate a request that will fail in processing
	err := simulateRequest(ctx, "user123", "getData")

	if err != nil {
		fmt.Println("Error occurred during request:")
		fmt.Println(err)

		fmt.Println("\nDetailed error with crumbs:")
		fmt.Println(crumbs.FormatError(err, false, true))
	}

	// Show crumbs explicitly
	fmt.Println("\nExtracting and using crumbs directly:")
	if cerr, ok := err.(*crumbs.Error); ok {
		crumbsSlice := cerr.GetCrumbs()
		for _, c := range crumbsSlice {
			switch c.Key {
			case "requestID":
				fmt.Println("Request ID:", c.Value)
			case "userID":
				fmt.Println("User ID:", c.Value)
			case "action":
				fmt.Println("Action:", c.Value)
			case "dataSize":
				fmt.Println("Data Size:", c.Value)
			}
		}
	}

	// Demonstrate how to get crumbs from a context
	fmt.Println("\nGetting crumbs from a context:")
	ctx = crumbs.AddCrumb(ctx, "key1", "value1", "key2", "value2")
	contextCrumbs := crumbs.GetCrumbs(ctx)
	fmt.Printf("Context crumbs: %+v\n", contextCrumbs)
}
