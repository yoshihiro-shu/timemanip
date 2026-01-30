# timemanip

A Go linter that detects usage of `time.Time` manipulation methods.

## Overview

`timemanip` is a static analysis tool that reports usage of the following `time.Time` methods:

- `Add()`
- `AddDate()`
- `Sub()`
- `Truncate()`
- `Round()`

This linter can be useful when you want to enforce consistent time manipulation patterns in your codebase or prevent direct time arithmetic operations.

## Installation

```bash
go install github.com/yoshihiro-shu/timemanip/cmd/timemanip@latest
```

## Usage

### Standalone

```bash
# Analyze current package
timemanip ./...

# Analyze specific package
timemanip ./pkg/...
```

### With go vet

```bash
go vet -vettool=$(which timemanip) ./...
```

### With golangci-lint (v1.57+)

Add to your `.golangci.yml`:

```yaml
linters-settings:
  custom:
    timemanip:
      type: "module"
      description: "Detects usage of time.Time manipulation methods"
      settings:
        # No settings available yet
```

## Examples

The following code will be flagged:

```go
package main

import "time"

func main() {
    t := time.Now()

    // All of these will be reported
    _ = t.Add(time.Hour)           // use of time.Time.Add is not allowed
    _ = t.AddDate(1, 0, 0)         // use of time.Time.AddDate is not allowed
    _ = t.Sub(t)                   // use of time.Time.Sub is not allowed
    _ = t.Truncate(time.Hour)      // use of time.Time.Truncate is not allowed
    _ = t.Round(time.Hour)         // use of time.Time.Round is not allowed
}
```

The following code will NOT be flagged:

```go
package main

import "time"

func main() {
    t := time.Now()

    // These are allowed
    _ = t.Year()
    _ = t.Month()
    _ = t.Format(time.RFC3339)
    _ = t.Before(t)
    _ = t.After(t)
}
```

## Development

### Running tests

```bash
go test ./...
```

### Building

```bash
go build ./cmd/timemanip
```

## License

MIT License
