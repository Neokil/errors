# Errors Library

A Go library for creating, wrapping, and annotating errors with additional context, stack traces, and structured logging support.

## Features

- **Error Categorization**: Use `kind` to categorize errors.
- **Stack Traces**: Automatically captures stack traces when errors are created or wrapped.
- **Annotations**: Add key-value pairs to errors for additional context.
- **Structured Logging**: Integrates with `log/slog` for structured error logging.
- **Error Wrapping**: Wrap existing errors while preserving their context.

## Installation

To install the library, run:

```bash
go get github.com/Neokil/errors
```

## Usage
### Creating a New Error
Use the `New` function to create a new error with a specific kind and message:
```golang
import "github.com/Neokil/errors"

err := errors.New("server-side", "an internal server error occurred")
```

### Wrapping an Existing Error
Use the `Wrap` function to wrap an existing error with a new kind:
```golang
import "github.com/Neokil/errors"

wrappedErr := errors.Wrap("database", err)
```

Annotating an Error
Use the `Annotate` function to add key-value pairs to an error for additional context:
```golang
import "github.com/Neokil/errors"

annotatedErr := errors.Annotate(err, "current_user", "demo-user@example.com")
```

### Logging an Error
The library integrates with `log/slog` for structured logging:
```golang
import (
    "log/slog"
    "github.com/Neokil/errors"
)

err := errors.New("server-side", "an internal server error occurred")
errors.Annotate(err, "current_user", "demo-user@example.com")

slog.Error("an error occurred", "err", err)
```
```json
{
    "time": "2025-01-01T00:00:00.000000+00:00",
    "level": "ERROR",
    "msg": "an error occurred",
    "err": {
        "kind": "server-side",
        "message": "an internal server error occurred",
        "annotations": {
            "current_user": "demo-user@example.com"
        },
        "stacktrace": "github.com/Neokil/errors_test.TestLogNewErrorWithSlog\n\t/.../errors/errors_test.go:20\ntesting.tRunner\n\t/.../src/testing/testing.go:1792\n"
    }
}
```

### Formatting an Error
The library supports custom formatting for errors:
- `%s` or `%q`: Prints the error message.
- `%+v`: Prints the error as a JSON object, including stack trace and annotations.
```golang
fmt.Printf("%+v\n", err)
```
```json
{
    "kind": "server-side",
    "message": "an internal server error occurred",
    "annotations": {
        "current_user": "demo-user@example.com"
    },
    "stacktrace": "github.com/Neokil/errors_test.TestPrintErrorWithFormat\n\t/.../errors/errors_test.go:48\ntesting.tRunner\n\t/.../src/testing/testing.go:1792\n"
}
```


### Testing
Run the tests using:
```bash
go test ./...
```

### License
This project is licensed under the MIT License. See the LICENSE file for details.
