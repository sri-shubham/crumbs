package crumbs

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	t.Run("basic error creation", func(t *testing.T) {
		err := New(ctx, "test error")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", err.Error())
		}
	})

	t.Run("with crumbs", func(t *testing.T) {
		err := New(ctx, "test error", "key1", "value1", "key2", 42)

		cerr, ok := err.(*Error)
		if !ok {
			t.Fatal("Expected *Error type")
		}

		crumbs := cerr.GetCrumbs()
		if crumbs["key1"] != "value1" {
			t.Errorf("Expected crumbs['key1'] = 'value1', got '%v'", crumbs["key1"])
		}

		if crumbs["key2"] != 42 {
			t.Errorf("Expected crumbs['key2'] = 42, got '%v'", crumbs["key2"])
		}
	})
}

func TestWrap(t *testing.T) {
	ctx := context.Background()
	baseErr := errors.New("base error")

	t.Run("basic wrapping", func(t *testing.T) {
		err := Wrap(ctx, baseErr, "wrapped error")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "wrapped error") {
			t.Errorf("Expected message to contain 'wrapped error', got '%s'", err.Error())
		}

		if !strings.Contains(err.Error(), "base error") {
			t.Errorf("Expected message to contain 'base error', got '%s'", err.Error())
		}

		if !errors.Is(err, baseErr) {
			t.Error("errors.Is failed to match the base error")
		}
	})

	t.Run("with crumbs", func(t *testing.T) {
		err := Wrap(ctx, baseErr, "wrapped error", "key1", "value1")

		cerr, ok := err.(*Error)
		if !ok {
			t.Fatal("Expected *Error type")
		}

		crumbs := cerr.GetCrumbs()
		if crumbs["key1"] != "value1" {
			t.Errorf("Expected crumbs['key1'] = 'value1', got '%v'", crumbs["key1"])
		}
	})

	t.Run("wrap nil", func(t *testing.T) {
		err := Wrap(ctx, nil, "wrapped nil")
		if err != nil {
			t.Errorf("Expected nil when wrapping nil, got '%v'", err)
		}
	})
}

func TestErrorf(t *testing.T) {
	ctx := context.Background()

	err := Errorf(ctx, "formatted %s %d", "error", 42)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expected := "formatted error 42"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestWrapf(t *testing.T) {
	ctx := context.Background()
	baseErr := errors.New("base error")

	err := Wrapf(ctx, baseErr, "formatted %s %d", "wrapper", 42)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "formatted wrapper 42") {
		t.Errorf("Expected message to contain 'formatted wrapper 42', got '%s'", err.Error())
	}

	if !errors.Is(err, baseErr) {
		t.Error("errors.Is failed to match the base error")
	}
}

func TestErrorsIs(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("sentinel error")

	t.Run("direct wrap", func(t *testing.T) {
		err := Wrap(ctx, sentinel, "wrapped")
		if !errors.Is(err, sentinel) {
			t.Error("errors.Is should find sentinel in direct wrap")
		}
	})

	t.Run("deep wrap", func(t *testing.T) {
		err1 := Wrap(ctx, sentinel, "inner")
		err2 := Wrap(ctx, err1, "middle")
		err3 := Wrap(ctx, err2, "outer")

		if !errors.Is(err3, sentinel) {
			t.Error("errors.Is should find sentinel in deep wrap")
		}
	})

	t.Run("different error", func(t *testing.T) {
		other := errors.New("other error")
		err := Wrap(ctx, sentinel, "wrapped")

		if errors.Is(err, other) {
			t.Error("errors.Is should not match different errors")
		}
	})
}

type customError struct {
	value int
}

func (e *customError) Error() string {
	return "custom error"
}

func TestErrorsAs(t *testing.T) {
	ctx := context.Background()
	custom := &customError{value: 42}

	t.Run("direct wrap", func(t *testing.T) {
		err := Wrap(ctx, custom, "wrapped")

		var ce *customError
		if !errors.As(err, &ce) {
			t.Error("errors.As should extract custom error")
		} else if ce.value != 42 {
			t.Errorf("Expected value 42, got %d", ce.value)
		}
	})

	t.Run("deep wrap", func(t *testing.T) {
		err1 := Wrap(ctx, custom, "inner")
		err2 := Wrap(ctx, err1, "middle")
		err3 := Wrap(ctx, err2, "outer")

		var ce *customError
		if !errors.As(err3, &ce) {
			t.Error("errors.As should extract custom error from deep wrap")
		} else if ce.value != 42 {
			t.Errorf("Expected value 42, got %d", ce.value)
		}
	})
}

func TestContextCrumbs(t *testing.T) {
	ctx := context.Background()

	// Add crumbs to context
	ctx = AddCrumb(ctx, "ctx1", "value1", "ctx2", 42)

	// Check GetCrumbs
	crumbs := GetCrumbs(ctx)
	if crumbs["ctx1"] != "value1" || crumbs["ctx2"] != 42 {
		t.Errorf("GetCrumbs failed, got %v", crumbs)
	}

	// Check that crumbs are included in errors
	err := New(ctx, "test error")
	cerr, ok := err.(*Error)
	if !ok {
		t.Fatal("Expected *Error type")
	}

	errCrumbs := cerr.GetCrumbs()
	if errCrumbs["ctx1"] != "value1" || errCrumbs["ctx2"] != 42 {
		t.Errorf("Context crumbs not included in error, got %v", errCrumbs)
	}
}

func TestAddCrumb(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key", "value")

		crumbs := GetCrumbs(ctx)
		if crumbs["key"] != "value" {
			t.Errorf("Expected crumbs['key'] = 'value', got '%v'", crumbs["key"])
		}
	})

	t.Run("existing crumbs", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key1", "value1")
		ctx = AddCrumb(ctx, "key2", "value2")

		crumbs := GetCrumbs(ctx)
		if crumbs["key1"] != "value1" || crumbs["key2"] != "value2" {
			t.Errorf("Expected both crumbs, got %v", crumbs)
		}
	})

	t.Run("override crumb", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key", "value1")
		ctx = AddCrumb(ctx, "key", "value2")

		crumbs := GetCrumbs(ctx)
		if crumbs["key"] != "value2" {
			t.Errorf("Expected crumbs['key'] = 'value2', got '%v'", crumbs["key"])
		}
	})

	t.Run("multiple crumbs at once", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key1", "value1", "key2", "value2")

		crumbs := GetCrumbs(ctx)
		if crumbs["key1"] != "value1" || crumbs["key2"] != "value2" {
			t.Errorf("Expected both crumbs, got %v", crumbs)
		}
	})

	t.Run("non-string key", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, 123, "value") // Should be ignored

		crumbs := GetCrumbs(ctx)
		if len(crumbs) > 0 {
			t.Errorf("Expected no crumbs, got %v", crumbs)
		}
	})
}

func TestStackTrace(t *testing.T) {
	ctx := context.Background()
	origCaptureStack := CaptureStack

	t.Run("capture disabled", func(t *testing.T) {
		CaptureStack = false
		err := New(ctx, "test error").(*Error)

		if len(err.GetStack()) > 0 {
			t.Error("Stack trace should not be captured when disabled")
		}
	})

	t.Run("capture enabled", func(t *testing.T) {
		CaptureStack = true
		err := New(ctx, "test error").(*Error)

		if len(err.GetStack()) == 0 {
			t.Error("Stack trace should be captured when enabled")
		}
	})

	t.Run("force stack", func(t *testing.T) {
		CaptureStack = false
		err := New(ctx, "test error").(*Error)
		err = err.ForceStack()

		if len(err.GetStack()) == 0 {
			t.Error("Stack trace should be captured when forced")
		}
	})

	t.Run("stack depth", func(t *testing.T) {
		CaptureStack = true
		origDepth := StackTraceDepth
		StackTraceDepth = 2

		err := New(ctx, "test error").(*Error)

		// Check that frames were limited
		if len(err.GetStack()) > 5 {
			t.Errorf("Expected limited stack frames, got %d", len(err.GetStack()))
		}

		StackTraceDepth = origDepth
	})

	// Restore original setting
	CaptureStack = origCaptureStack
}

func TestFormatError(t *testing.T) {
	ctx := context.Background()
	CaptureStack = true
	defer func() { CaptureStack = false }()

	err := New(ctx, "test error", "key1", "value1")

	t.Run("basic format", func(t *testing.T) {
		formatted := FormatError(err, false, false)
		if !strings.Contains(formatted, "test error") {
			t.Errorf("Formatted error should contain message, got: %s", formatted)
		}
	})

	t.Run("with crumbs", func(t *testing.T) {
		formatted := FormatError(err, false, true)
		if !strings.Contains(formatted, "key1") || !strings.Contains(formatted, "value1") {
			t.Errorf("Formatted error should contain crumbs, got: %s", formatted)
		}
	})

	t.Run("with stack", func(t *testing.T) {
		formatted := FormatError(err, true, false)
		if !strings.Contains(formatted, "Stack trace") {
			t.Errorf("Formatted error should contain stack trace, got: %s", formatted)
		}
	})
}
