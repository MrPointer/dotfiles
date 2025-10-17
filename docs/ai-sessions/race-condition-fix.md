# Race Condition Fix: Progress Display Thread Safety

## Problem Description

The progress display implementation was experiencing race conditions when running tests with Go's race detector (`go test -race`). The race occurred between:

1. **Writer goroutines**: `huh` spinner library (via `bubbletea`) writing terminal control sequences and content to the output buffer
2. **Reader goroutines**: Test code reading buffer content using `buffer.String()` to verify output

## Root Cause Analysis

The core issue was that `bytes.Buffer` is **not thread-safe** for concurrent reads and writes. Even though writes were happening through the spinner library, tests were directly reading from the same buffer concurrently, causing data races at the memory level.

### Race Condition Details

```
WARNING: DATA RACE
Read at 0x... by goroutine 7:
  bytes.(*Buffer).String()
Previous write at 0x... by goroutine 8:
  bytes.(*Buffer).Write()
```

This occurred because:
- Multiple goroutines in `bubbletea`'s renderer were writing to the buffer
- Test goroutines were reading via `buffer.String()` simultaneously
- No synchronization existed between these operations

## Solution Implementation

### 1. Thread-Safe Buffer Wrapper

Created `safeBytesBuffer` that wraps `*bytes.Buffer` with `sync.RWMutex`:

```go
type safeBytesBuffer struct {
    buf   *bytes.Buffer
    mutex sync.RWMutex
}

func (sbb *safeBytesBuffer) Write(p []byte) (n int, err error) {
    sbb.mutex.Lock()
    defer sbb.mutex.Unlock()
    return sbb.buf.Write(p)
}

func (sbb *safeBytesBuffer) SafeString() string {
    sbb.mutex.RLock()
    defer sbb.mutex.RUnlock()
    return sbb.buf.String()
}
```

### 2. Smart Output Writer Detection

Modified `NewProgressDisplay` to automatically detect `*bytes.Buffer` and wrap it:

```go
// Special handling for *bytes.Buffer to ensure thread safety
if buf, ok := output.(*bytes.Buffer); ok {
    safeBuffer = &safeBytesBuffer{buf: buf}
    outputWriter = safeBuffer
} else {
    outputWriter = &synchronizedWriter{writer: output}
}
```

### 3. Safe Testing Interface

Added `GetOutputSafely()` method for tests to safely read buffer content:

```go
func (p *ProgressDisplay) GetOutputSafely() string {
    if p.safeBuffer != nil {
        return p.safeBuffer.SafeString()
    }
    return ""
}
```

## Testing Best Practices

### ❌ Unsafe (Race Condition)
```go
func TestExample(t *testing.T) {
    var output bytes.Buffer
    display := NewProgressDisplay(&output)
    
    display.Start("Test")
    // Race condition: concurrent read while spinner writes
    content := output.String()
}
```

### ✅ Safe (Thread-Safe)
```go
func TestExample(t *testing.T) {
    var output bytes.Buffer
    display := NewProgressDisplay(&output)
    
    display.Start("Test")
    // Thread-safe read using synchronized access
    content := display.GetOutputSafely()
}
```

## Key Insights

1. **Library Assumption**: The issue wasn't in the `huh` or `bubbletea` libraries themselves, but in **our usage pattern** where we exposed the same buffer to both the library (for writes) and tests (for reads) without proper synchronization.

2. **Buffer Thread Safety**: `bytes.Buffer` is designed for single-threaded use. Concurrent access requires explicit synchronization.

3. **Testing vs Production**: The race only manifested in tests because production code typically doesn't read from the output buffer while the progress display is running.

## Resolution Verification

- ✅ All tests pass with race detection enabled (`go test -race`)
- ✅ All tests pass in normal mode (`go test`)
- ✅ No performance degradation in production usage
- ✅ Backward compatibility maintained for existing code

## Lessons Learned

- Always consider concurrent access patterns when sharing mutable state between goroutines
- Race conditions can be subtle and may only appear in testing scenarios
- Thread-safe wrappers should protect both reads and writes, not just writes
- Go's race detector is invaluable for catching these issues early