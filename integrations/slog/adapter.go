package slog

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/sri-shubham/crumbs"
	"github.com/sri-shubham/crumbs/logger"
)

// Adapter implements logger.Logger using log/slog
type Adapter struct {
	logger *slog.Logger
}

// New creates a new slog adapter
func New(l *slog.Logger) *Adapter {
	if l == nil {
		l = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return &Adapter{logger: l}
}

func (l *Adapter) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, nil, args...)
}

func (l *Adapter) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, nil, args...)
}

func (l *Adapter) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, nil, args...)
}

func (l *Adapter) Error(ctx context.Context, msg string, err error, args ...any) {
	l.log(ctx, slog.LevelError, msg, err, args...)
}

func (l *Adapter) log(ctx context.Context, level slog.Level, msg string, err error, args ...any) {
	if !l.logger.Enabled(ctx, level) {
		return
	}

	// Start with provided args
	logArgs := make([]any, 0, len(args)+4) // Pre-allocate for args + potential crumbs
	logArgs = append(logArgs, args...)

	// Add error if present
	if err != nil {
		logArgs = append(logArgs, "error", err.Error())

		// Extract crumbs from error
		var cerr *crumbs.Error
		if errors.As(err, &cerr) {
			for _, c := range cerr.GetCrumbs() {
				logArgs = append(logArgs, c.Key, c.Value)
			}
		}
	}

	// Extract crumbs from context
	if ctxCrumbs := crumbs.GetCrumbs(ctx); len(ctxCrumbs) > 0 {
		for _, c := range ctxCrumbs {
			logArgs = append(logArgs, c.Key, c.Value)
		}
	}

	l.logger.Log(ctx, level, msg, logArgs...)
}

// Ensure Adapter implements logger.Logger
var _ logger.Logger = (*Adapter)(nil)
