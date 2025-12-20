# Logger Interface

This package defines a generic `Logger` interface for the `crumbs` ecosystem. It allows libraries and applications to log messages with context and structured data without being tied to a specific logging implementation.

## Interface Definition

```go
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, err error, args ...any)
}
```

## Usage

Use this interface in your service or library code to accept any compatible logger.

```go
type MyService struct {
    logger logger.Logger
}

func NewService(l logger.Logger) *MyService {
    return &MyService{logger: l}
}

func (s *MyService) DoWork(ctx context.Context) {
    s.logger.Info(ctx, "work started")
}
```

## Implementations

- **[slog](../integrations/slog)**: Adapter for Go's standard `log/slog` library.
