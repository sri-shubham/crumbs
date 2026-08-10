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

	// withArgs are key-value pairs accumulated via With, merged into every
	// subsequent log call's args (rather than baked into logger via
	// slog.Logger.With) so they still go through the same error-stringify
	// and crumb-splat logic as call-site args. See adapter_test.go's
	// TestAdapter_WithPreservesErrorCrumbs for the regression this guards
	// against.
	withArgs []any
}

// New creates a new slog adapter
func New(l *slog.Logger) *Adapter {
	if l == nil {
		l = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return &Adapter{logger: l}
}

// Debug logs at DEBUG level, attaching crumbs from ctx and any error args.
func (l *Adapter) Debug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

// Info logs at INFO level, attaching crumbs from ctx and any error args.
func (l *Adapter) Info(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs at WARN level, attaching crumbs from ctx and any error args.
func (l *Adapter) Warn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at ERROR level, attaching crumbs from ctx and any error args.
func (l *Adapter) Error(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

// log forwards to the underlying slog.Logger while enriching args with crumbs.
// Crumb-extraction rules:
//   - Any value in args (in key/value positions) that is itself an error gets
//     scanned with errors.As(&*crumbs.Error). The first such *crumbs.Error
//     contributes its crumbs to the log line, and the context lookup is
//     skipped (the error already snapshotted the ctx crumbs).
//   - Otherwise, ctx crumbs are appended.
func (l *Adapter) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.logger.Enabled(ctx, level) {
		return
	}

	all := args
	if len(l.withArgs) > 0 {
		all = make([]any, 0, len(l.withArgs)+len(args))
		all = append(all, l.withArgs...)
		all = append(all, args...)
	}

	logArgs := make([]any, 0, len(all)+4)

	var cerr *crumbs.Error
	for _, a := range all {
		// Replace error values with their string form so JSON/text handlers
		// emit the message rather than an opaque struct representation. Also
		// remember the first *crumbs.Error so we can splat its crumbs.
		if e, ok := a.(error); ok && e != nil {
			if cerr == nil {
				_ = errors.As(e, &cerr)
			}
			logArgs = append(logArgs, e.Error())
			continue
		}
		logArgs = append(logArgs, a)
	}

	if cerr != nil {
		for _, c := range cerr.GetCrumbs() {
			logArgs = append(logArgs, c.Key, c.Value)
		}
	} else if ctxCrumbs := crumbs.GetCrumbs(ctx); len(ctxCrumbs) > 0 {
		for _, c := range ctxCrumbs {
			logArgs = append(logArgs, c.Key, c.Value)
		}
	}

	l.logger.Log(ctx, level, msg, logArgs...)
}

// With returns a derived Adapter that includes the supplied key-value args
// on every subsequent log call. Args are merged at log time (not baked into
// the underlying slog.Logger via slog.Logger.With) so an error passed here
// still gets stringified and, if it is a *crumbs.Error, has its crumbs
// splatted onto every subsequent log line exactly as if passed at the call
// site.
func (l *Adapter) With(args ...any) logger.Logger {
	merged := make([]any, 0, len(l.withArgs)+len(args))
	merged = append(merged, l.withArgs...)
	merged = append(merged, args...)
	return &Adapter{logger: l.logger, withArgs: merged}
}

// Ensure Adapter implements logger.Logger
var _ logger.Logger = (*Adapter)(nil)
