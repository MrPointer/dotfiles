---
name: writing-go-code
description: Apply Go coding standards when writing, reviewing, or modifying Go code. Use when implementing functions, writing tests with testify, generating mocks with mockery, using dependency injection, handling errors idiomatically, or working with interfaces. Use this skill for any Go file editing task. For Go language conventions and idioms from official docs, see effective-go skill instead.
---

# Go Development Standards

Project-specific Go coding standards for this codebase.

## Quick Reference

**Coding patterns:** See [Code Style Reference](references/code-style.md)
**Testing patterns:** See [Test Style Reference](references/test-style.md)

## Code Organization

```go
// 1. Struct definition
type MyService struct {
    logger Logger
    fs     FileSystem
}

// 2. Interface verification (immediately after struct)
var _ Service = (*MyService)(nil)

// 3. Constructor with dependency injection
func NewMyService(logger Logger, fs FileSystem) *MyService {
    return &MyService{logger: logger, fs: fs}
}
```

## Testing Patterns

Use `testify/require` for all assertions. Name tests descriptively:

```go
func Test_ServiceReturnsErrorWhenFileNotFound(t *testing.T) {
    // Arrange: create mock
    fs := &MoqFileSystem{
        ReadFileFunc: func(path string) ([]byte, error) {
            return nil, os.ErrNotExist
        },
    }
    svc := NewMyService(logger, fs)

    // Act
    err := svc.LoadConfig("missing.yaml")

    // Assert: check error by keyword, not full message
    require.Error(t, err)
    require.Contains(t, err.Error(), "not found")
}
```

## Table-Driven Tests

```go
func Test_ParseConfig(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Config
        wantErr bool
    }{
        {"valid config", "key: value", Config{Key: "value"}, false},
        {"empty input", "", Config{}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseConfig(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tt.want, got)
        })
    }
}
```

## Mock Generation

Mocks use `mockery` with moq template. To regenerate all mocks:

```bash
mockery
```

Mock types are prefixed with `Moq` (e.g., `MoqLogger`, `MoqFileSystem`).

## Optional Types

Use `samber/mo` for safer nil handling:

```go
import "github.com/samber/mo"

type Config struct {
    Shell mo.Option[string]
}

if shell, ok := config.Shell.Get(); ok {
    // use shell
}
```

## Key Rules

- Line length: 120 characters max
- Pre-allocate slices/maps when size is known
- Wrap OS operations in interfaces for mockability
- End doc comments with a period
- Never edit mock files manually
