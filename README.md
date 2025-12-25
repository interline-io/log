# log

A lightweight logging wrapper for [zerolog](https://github.com/rs/zerolog), providing convenience functions and sensible defaults for [Interline](https://www.interline.io/) and [Transitland](https://www.transit.land/) projects.

> **Note:** This library is primarily intended for internal use across Interline/Transitland repositories—mainly for our convenience and to avoid having `zerolog.Xyz` scattered throughout our codebase and copy-pasting logging middlewares. It's published publicly for our convenience; external users should probably just use zerolog directly.

## Features

- Simple printf-style logging functions (`Infof`, `Debugf`, `Tracef`, `Errorf`)
- Structured logging via zerolog's fluent API
- Pretty console output with colors (automatically disabled when piping/redirecting)
- JSON output mode for production environments
- Context-aware logging with request ID support
- HTTP middleware for request logging

## Installation

```bash
go get github.com/interline-io/log
```

## Usage

### Basic Logging

```go
package main

import "github.com/interline-io/log"

func main() {
    // Printf-style logging
    log.Infof("Processing %d items", 42)
    log.Debugf("Debug info: %s", "details")
    log.Errorf("Something went wrong: %v", err)
    log.Tracef("Verbose trace: %+v", obj)

    // Structured logging (zerolog style)
    log.Info().Str("user", "alice").Int("count", 5).Msg("User action")
    log.Error().Err(err).Str("op", "save").Msg("Operation failed")

    // Simple print (no timestamp, ignores log level)
    log.Print("Plain output: %s", "hello")
}
```

### Context-Aware Logging

```go
// Get logger from context
logger := log.For(ctx)
logger.Info().Msg("Using context logger")

// Add logger to context
ctx = log.WithLogger(ctx, customLogger)
```

### HTTP Middleware

```go
import "github.com/interline-io/log"

// Add request ID to each request
r.Use(log.RequestIDMiddleware)

// Add request ID to context logger
r.Use(log.RequestIDLoggingMiddleware)

// Log request duration and details
r.Use(log.LoggingMiddleware(1000, getUserNameFunc))
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TL_LOG` | Log level: `TRACE`, `DEBUG`, `INFO`, `ERROR`, `FATAL` | `INFO` |
| `TL_LOG_JSON` | Set to `true` for JSON output (no colors) | `false` |

### Examples

```bash
# Enable debug logging
TL_LOG=DEBUG ./myapp

# Enable trace logging (most verbose)
TL_LOG=TRACE ./myapp

# JSON output for production/log aggregation
TL_LOG_JSON=true ./myapp
```

### Terminal Colors

Console output includes ANSI colors when writing to a terminal. Colors are automatically disabled when output is piped or redirected:

```bash
# Colors enabled (terminal)
./myapp

# Colors disabled (piped)
./myapp | tee output.log

# Colors disabled (redirected)
./myapp > output.log 2>&1
```

## Log Levels

| Level | Function | Use Case |
|-------|----------|----------|
| `TRACE` | `Tracef()`, `Trace()` | Verbose debugging, performance tracing |
| `DEBUG` | `Debugf()`, `Debug()` | Development debugging |
| `INFO` | `Infof()`, `Info()` | General operational messages |
| `ERROR` | `Errorf()`, `Error()` | Error conditions |

## License

See [LICENSE](LICENSE) file.