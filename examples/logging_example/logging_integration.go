package logging_example

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sri-shubham/crumbs"
)

// Logger is a simple logger interface
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

// SimpleLogger is a basic logger implementation
type SimpleLogger struct {
	logger *log.Logger
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

// Debug logs debug messages
func (l *SimpleLogger) Debug(msg string, fields map[string]interface{}) {
	l.log("DEBUG", msg, fields)
}

// Info logs info messages
func (l *SimpleLogger) Info(msg string, fields map[string]interface{}) {
	l.log("INFO", msg, fields)
}

// Warn logs warning messages
func (l *SimpleLogger) Warn(msg string, fields map[string]interface{}) {
	l.log("WARN", msg, fields)
}

// Error logs error messages
func (l *SimpleLogger) Error(msg string, fields map[string]interface{}) {
	l.log("ERROR", msg, fields)
}

func (l *SimpleLogger) log(level string, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["msg"] = msg
	fields["level"] = level

	// Convert to JSON
	jsonData, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		l.logger.Printf("[%s] %s (failed to marshal fields)", level, msg)
		return
	}

	l.logger.Println(string(jsonData))
}

// LogError logs an error with crumbs
func LogError(logger Logger, err error) {
	if err == nil {
		return
	}

	msg := err.Error()
	fields := map[string]interface{}{}

	// Check if it's a crumbs error
	var cerr *crumbs.Error
	if errors.As(err, &cerr) {
		// Add all crumbs as fields
		for k, v := range cerr.GetCrumbs() {
			fields[k] = v
		}

		// Add stack trace if available
		stack := cerr.GetStack()
		if len(stack) > 0 {
			frames := make([]string, 0, len(stack))
			for _, frame := range stack {
				frames = append(frames, fmt.Sprintf("%s:%d", frame.File, frame.Line))
			}
			fields["stack"] = frames
		}
	}

	logger.Error(msg, fields)
}

// DemonstrateLoggingIntegration shows how to integrate with logging libraries
func DemonstrateLoggingIntegration() {
	fmt.Println("Logging Integration Example")
	fmt.Println("==========================")

	logger := NewSimpleLogger()
	ctx := context.Background()

	// Add some crumbs to the context
	ctx = crumbs.AddCrumb(ctx,
		"requestID", "req-abcd",
		"userID", "user-1234",
	)

	// Enable stack traces for this example
	crumbs.CaptureStack = true
	defer func() { crumbs.CaptureStack = false }()

	// Create an error with additional crumbs
	err := crumbs.New(ctx, "operation failed",
		"operation", "getData",
		"status", 500,
	)

	fmt.Println("\nLogging an error with crumbs:")
	LogError(logger, err)

	// Create a more complex error chain
	baseErr := errors.New("network timeout")
	wrappedErr := crumbs.Wrap(ctx, baseErr, "API request failed",
		"endpoint", "/api/v1/users",
		"method", "GET",
	)

	fmt.Println("\nLogging a wrapped error with crumbs:")
	LogError(logger, wrappedErr)
}
