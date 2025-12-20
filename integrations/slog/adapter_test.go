package slog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/sri-shubham/crumbs"
	crumbslog "github.com/sri-shubham/crumbs/integrations/slog"
	"github.com/sri-shubham/crumbs/logger"
)

func TestAdapter_ImplementsLogger(t *testing.T) {
	var _ logger.Logger = (*crumbslog.Adapter)(nil)
}

func TestAdapter_Logging(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(h)
	adapter := crumbslog.New(l)

	ctx := context.Background()

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		adapter.Info(ctx, "info message", "key", "value")

		var logEntry map[string]any
		if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
			t.Fatalf("failed to unmarshal log entry: %v", err)
		}

		if logEntry["msg"] != "info message" {
			t.Errorf("expected msg 'info message', got %v", logEntry["msg"])
		}
		if logEntry["level"] != "INFO" {
			t.Errorf("expected level 'INFO', got %v", logEntry["level"])
		}
		if logEntry["key"] != "value" {
			t.Errorf("expected key 'value', got %v", logEntry["key"])
		}
	})

	t.Run("Context Crumbs", func(t *testing.T) {
		buf.Reset()
		ctxWithCrumbs := crumbs.AddCrumb(ctx, "request_id", "12345")
		adapter.Info(ctxWithCrumbs, "context crumbs")

		var logEntry map[string]any
		if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
			t.Fatalf("failed to unmarshal log entry: %v", err)
		}

		if logEntry["request_id"] != "12345" {
			t.Errorf("expected request_id '12345', got %v", logEntry["request_id"])
		}
	})

	t.Run("Error Crumbs", func(t *testing.T) {
		buf.Reset()
		err := crumbs.New(ctx, "something went wrong", "user_id", "u-999")
		adapter.Error(ctx, "error occurred", err)

		var logEntry map[string]any
		if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
			t.Fatalf("failed to unmarshal log entry: %v", err)
		}

		if logEntry["msg"] != "error occurred" {
			t.Errorf("expected msg 'error occurred', got %v", logEntry["msg"])
		}
		if logEntry["error"] != "something went wrong" {
			t.Errorf("expected error 'something went wrong', got %v", logEntry["error"])
		}
		if logEntry["user_id"] != "u-999" {
			t.Errorf("expected user_id 'u-999', got %v", logEntry["user_id"])
		}
	})
}
