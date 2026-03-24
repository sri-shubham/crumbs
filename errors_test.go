package crumbs

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func crumbsToMap(crumbs []Crumb) map[string]interface{} {
	m := make(map[string]interface{})
	for _, c := range crumbs {
		m[c.Key] = c.Value
	}
	return m
}

func TestNewError(t *testing.T) {
	ctx := context.Background()

	t.Run("basic error creation", func(t *testing.T) {
		err := NewError(ctx, "test error")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if err.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", err.Error())
		}
	})

	t.Run("with crumbs", func(t *testing.T) {
		err := NewError(ctx, "test error", "key1", "value1", "key2", 42)

		cerr, ok := err.(*Error)
		if !ok {
			t.Fatal("Expected *Error type")
		}

		crumbs := crumbsToMap(cerr.GetCrumbs())
		if crumbs["key1"] != "value1" {
			t.Errorf("Expected crumbs['key1'] = 'value1', got '%v'", crumbs["key1"])
		}

		if crumbs["key2"] != 42 {
			t.Errorf("Expected crumbs['key2'] = 42, got '%v'", crumbs["key2"])
		}
	})
}

func TestWrapError(t *testing.T) {
	ctx := context.Background()
	baseErr := errors.New("base error")

	t.Run("basic wrapping", func(t *testing.T) {
		err := WrapError(ctx, baseErr, "wrapped error")

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
		err := WrapError(ctx, baseErr, "wrapped error", "key1", "value1")

		cerr, ok := err.(*Error)
		if !ok {
			t.Fatal("Expected *Error type")
		}

		crumbs := crumbsToMap(cerr.GetCrumbs())
		if crumbs["key1"] != "value1" {
			t.Errorf("Expected crumbs['key1'] = 'value1', got '%v'", crumbs["key1"])
		}
	})

	t.Run("wrap nil", func(t *testing.T) {
		err := WrapError(ctx, nil, "wrapped nil")
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
		err := WrapError(ctx, sentinel, "wrapped")
		if !errors.Is(err, sentinel) {
			t.Error("errors.Is should find sentinel in direct wrap")
		}
	})

	t.Run("deep wrap", func(t *testing.T) {
		err1 := WrapError(ctx, sentinel, "inner")
		err2 := WrapError(ctx, err1, "middle")
		err3 := WrapError(ctx, err2, "outer")

		if !errors.Is(err3, sentinel) {
			t.Error("errors.Is should find sentinel in deep wrap")
		}
	})

	t.Run("different error", func(t *testing.T) {
		other := errors.New("other error")
		err := WrapError(ctx, sentinel, "wrapped")

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
		err := WrapError(ctx, custom, "wrapped")

		var ce *customError
		if !errors.As(err, &ce) {
			t.Error("errors.As should extract custom error")
		} else if ce.value != 42 {
			t.Errorf("Expected value 42, got %d", ce.value)
		}
	})

	t.Run("deep wrap", func(t *testing.T) {
		err1 := WrapError(ctx, custom, "inner")
		err2 := WrapError(ctx, err1, "middle")
		err3 := WrapError(ctx, err2, "outer")

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
	crumbs := crumbsToMap(GetCrumbs(ctx))
	if crumbs["ctx1"] != "value1" || crumbs["ctx2"] != 42 {
		t.Errorf("GetCrumbs failed, got %v", crumbs)
	}

	// Check that crumbs are included in errors
	err := NewError(ctx, "test error")
	cerr, ok := err.(*Error)
	if !ok {
		t.Fatal("Expected *Error type")
	}

	errCrumbs := crumbsToMap(cerr.GetCrumbs())
	if errCrumbs["ctx1"] != "value1" || errCrumbs["ctx2"] != 42 {
		t.Errorf("Context crumbs not included in error, got %v", errCrumbs)
	}
}

func TestAddCrumb(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key", "value")

		crumbs := crumbsToMap(GetCrumbs(ctx))
		if crumbs["key"] != "value" {
			t.Errorf("Expected crumbs['key'] = 'value', got '%v'", crumbs["key"])
		}
	})

	t.Run("existing crumbs", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key1", "value1")
		ctx = AddCrumb(ctx, "key2", "value2")

		crumbs := crumbsToMap(GetCrumbs(ctx))
		if crumbs["key1"] != "value1" || crumbs["key2"] != "value2" {
			t.Errorf("Expected both crumbs, got %v", crumbs)
		}
	})

	t.Run("override crumb", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key", "value1")
		ctx = AddCrumb(ctx, "key", "value2")

		crumbs := crumbsToMap(GetCrumbs(ctx))
		if crumbs["key"] != "value2" {
			t.Errorf("Expected crumbs['key'] = 'value2', got '%v'", crumbs["key"])
		}
	})

	t.Run("multiple crumbs at once", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, "key1", "value1", "key2", "value2")

		crumbs := crumbsToMap(GetCrumbs(ctx))
		if crumbs["key1"] != "value1" || crumbs["key2"] != "value2" {
			t.Errorf("Expected both crumbs, got %v", crumbs)
		}
	})

	t.Run("non-string key", func(t *testing.T) {
		ctx := context.Background()
		ctx = AddCrumb(ctx, 123, "value") // Should be ignored

		crumbs := crumbsToMap(GetCrumbs(ctx))
		if len(crumbs) > 0 {
			t.Errorf("Expected no crumbs, got %v", crumbs)
		}
	})
}

func TestStackTrace(t *testing.T) {
	ctx := context.Background()
	origcaptureStack := captureStack

	t.Run("capture disabled", func(t *testing.T) {
		captureStack = false
		err := NewError(ctx, "test error").(*Error)

		if len(err.GetStack()) > 0 {
			t.Error("Stack trace should not be captured when disabled")
		}
	})

	t.Run("capture enabled", func(t *testing.T) {
		captureStack = true
		err := NewError(ctx, "test error").(*Error)

		if len(err.GetStack()) == 0 {
			t.Error("Stack trace should be captured when enabled")
		}
	})

	t.Run("force stack", func(t *testing.T) {
		captureStack = false
		err := NewError(ctx, "test error").(*Error)
		err = err.ForceStack()

		if len(err.GetStack()) == 0 {
			t.Error("Stack trace should be captured when forced")
		}
	})

	t.Run("stack depth", func(t *testing.T) {
		captureStack = true
		origDepth := stackTraceDepth
		stackTraceDepth = 2

		err := NewError(ctx, "test error").(*Error)

		// Check that frames were limited
		if len(err.GetStack()) > 5 {
			t.Errorf("Expected limited stack frames, got %d", len(err.GetStack()))
		}

		stackTraceDepth = origDepth
	})

	// Restore original setting
	captureStack = origcaptureStack
}

func TestFormatError(t *testing.T) {
	ctx := context.Background()
	captureStack = true
	defer func() { captureStack = false }()

	err := NewError(ctx, "test error", "key1", "value1")

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

func TestConfigureStackTraces(t *testing.T) {
	ConfigureStackTraces(true, 10)
	if !captureStack || stackTraceDepth != 10 {
		t.Errorf("ConfigureStackTraces failed")
	}
	ConfigureStackTraces(false, 32)
}

func TestBadKeysAndCoverage(t *testing.T) {
	ctx := context.Background()

	// 1. Wrapf with nil err
	if Wrapf(ctx, nil, "fmt %s", "a") != nil {
		t.Error("Wrapf should handle nil err")
	}

	// 2. FormatStack with no stack
	err := NewError(ctx, "msg").(*Error)
	if err.FormatStack() != "no stack trace available" {
		t.Error("FormatStack should return 'no stack trace available'")
	}

	// 3. newError dangling key
	err2 := NewError(ctx, "msg", "key", "val", "dangling").(*Error)
	if len(err2.Crumbs) != 2 || err2.Crumbs[1].Key != "!BADKEY" || err2.Crumbs[1].Value != "dangling" {
		t.Error("Dangling key not mapped to !BADKEY")
	}

	// 4. newError odd key that is not string
	err3 := NewError(ctx, "msg", "key", "val", 123).(*Error)
	if len(err3.Crumbs) != 2 || err3.Crumbs[1].Key != "!BADKEY" || err3.Crumbs[1].Value != 123 {
		t.Error("Odd non-string key not mapped to !BADKEY")
	}

	// 4b. newError even non-string key ignored
	err3b := NewError(ctx, "msg", 123, "val").(*Error)
	if len(err3b.Crumbs) != 0 {
		t.Error("Even non-string key not ignored")
	}

	// 5. AddCrumb dangling key
	ctx2 := AddCrumb(ctx, "key", "val", "dangling")
	crumbs := GetCrumbs(ctx2)
	if len(crumbs) != 2 || crumbs[1].Key != "!BADKEY" || crumbs[1].Value != "dangling" {
		t.Error("AddCrumb dangling key failed")
	}

	// 6. AddCrumb odd non-string key
	ctx3 := AddCrumb(ctx, "key", "val", 123)
	crumbs = GetCrumbs(ctx3)
	if len(crumbs) != 2 || crumbs[1].Key != "!BADKEY" || crumbs[1].Value != 123 {
		t.Error("AddCrumb non-string dangling key failed")
	}

	// 7. AddCrumb even non-string key
	ctx4 := AddCrumb(ctx, 123, "val")
	if len(GetCrumbs(ctx4)) != 0 {
		t.Error("AddCrumb even non-string key not ignored")
	}

	// 8. Error method fallbacks
	var emptyErr *Error = &Error{}
	if emptyErr.Error() != "unknown error" {
		t.Error("Empty Error.Error() failed")
	}
	wrapErr := &Error{Err: errors.New("base")}
	if wrapErr.Error() != "base" {
		t.Error("Error() falling back to base failed")
	}
}

func TestGetCrumbsNil(t *testing.T) {
	//nolint:staticcheck // deliberately testing nil safeguard
	if GetCrumbs(nil) != nil {
		t.Error("GetCrumbs(nil) should be nil")
	}
	if GetCrumbs(context.Background()) != nil {
		t.Error("GetCrumbs(emptyCtx) should be nil")
	}
}
