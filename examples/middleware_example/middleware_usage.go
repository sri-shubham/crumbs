package middleware_example

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sri-shubham/crumbs"
)

// TraceMiddleware is an HTTP middleware that adds a TraceID to the context
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate or extract TraceID (simulated here)
		traceID := fmt.Sprintf("trace-%d", time.Now().UnixNano())

		// Add TraceID and other request details to the context using crumbs
		// This context will now carry these values to all downstream functions
		ctx := crumbs.AddCrumb(r.Context(),
			"trace_id", traceID,
			"path", r.URL.Path,
			"method", r.Method,
		)

		// Update request with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ServiceLayer simulates a business logic layer
func ServiceLayer(ctx context.Context) error {
	// Simulate some work
	// ...

	// Call the database layer
	if err := DatabaseLayer(ctx); err != nil {
		// Wrap the error with more context
		return crumbs.WrapError(ctx, err, "service layer failed")
	}
	return nil
}

// DatabaseLayer simulates a database operation
func DatabaseLayer(ctx context.Context) error {
	// Simulate a failure
	// We can add more specific context here
	// Simulate an error
	return crumbs.NewError(ctx, "connection timeout",
		"service", "db",
		"retry_count", 3)
}

// Handler simulates an HTTP handler
func Handler(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := r.Context()

	// Call service layer
	if err := ServiceLayer(ctx); err != nil {
		// When logging the error, we can extract all the crumbs.
		// This log entry will contain:
		// - trace_id (from middleware)
		// - path (from middleware)
		// - method (from middleware)
		// - db_host (from database layer)
		// - db_port (from database layer)
		// - The error message chain

		LogError(logger, "request failed", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

// LogError helper (as defined in README)
func LogError(logger *slog.Logger, msg string, err error) {
	args := []any{"error", err.Error()}

	// Extract crumbs if available
	if cerr, ok := err.(*crumbs.Error); ok {
		for _, c := range cerr.GetCrumbs() {
			args = append(args, slog.Any(c.Key, c.Value))
		}
	}

	logger.Error(msg, args...)
}

// RunExample starts a server to demonstrate the middleware
func RunExample() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data", Handler)

	// Wrap mux with middleware
	handler := TraceMiddleware(mux)

	fmt.Println("Starting server on :8080...")
	fmt.Println("Try: curl http://localhost:8080/api/data")

	// In a real app: http.ListenAndServe(":8080", handler)

	// For demonstration, let's simulate a request directly
	req, _ := http.NewRequest("GET", "/api/data", nil)
	w := &mockResponseWriter{}
	handler.ServeHTTP(w, req)
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() http.Header        { return http.Header{} }
func (m *mockResponseWriter) Write([]byte) (int, error)  { return 0, nil }
func (m *mockResponseWriter) WriteHeader(statusCode int) {}
