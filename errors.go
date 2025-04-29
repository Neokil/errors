package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
)

type internalError struct {
	Kind        string         `json:"kind"`
	Message     string         `json:"message"`
	Annotations map[string]any `json:"annotations,omitempty"`
	Stacktrace  string         `json:"stacktrace,omitempty"`
	Parent      error          `json:"parent,omitempty"`
}

// Error implements the error interface for internalError.
func (e *internalError) Error() string {
	return e.Message
}

// Unwrap implements the Unwrap method for internalError.
func (e *internalError) Unwrap() error {
	return e.Parent
}

func (e *internalError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			b, err := json.Marshal(e)
			if err != nil {
				panic(err)
			}
			fmt.Fprint(s, string(b))
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *internalError) LogValue() slog.Value {
	attrs := []slog.Attr{}
	if e.Kind != "" {
		attrs = append(attrs, slog.String("kind", e.Kind))
	}
	if e.Message != "" {
		attrs = append(attrs, slog.String("message", e.Message))
	}
	if len(e.Annotations) > 0 {
		annotations := []slog.Attr{}
		for k, v := range e.Annotations {
			annotations = append(annotations, slog.Any(k, v))
		}
		attrs = append(attrs, slog.Any("annotations", slog.GroupValue(annotations...)))
	}
	if e.Stacktrace != "" {
		attrs = append(attrs, slog.String("stacktrace", e.Stacktrace))
	}
	if e.Parent != nil {
		attrs = append(attrs, slog.String("parent", e.Parent.Error()))
	}

	return slog.GroupValue(attrs...)
}

func internalNew(kind string, msg string) *internalError {
	return &internalError{
		Kind:       kind,
		Message:    msg,
		Stacktrace: stackFromCallers(),
	}
}

func stackFromCallers() string {
	// get callers, skipping the first 4 frames, as those frames contain the errors-internals
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(4, pcs[:])

	// get frames for the callers and build stacktrace
	frames := runtime.CallersFrames(pcs[0:n])
	var stacktrace string
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		stacktrace += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
	}
	return stacktrace
}

// New creates a new error with the given kind and message.
// The stacktrace is captured at the time of creation.
// The kind is used to categorize the error, and the message is used to
// provide additional context about the error.
func New(kind string, msg string) error {
	return internalNew(kind, msg)
}

// Wrap creates a new error that wraps the given parent error.
// The stacktrace is captured at the time of creation.
// The kind is used to categorize the error.
func Wrap(kind string, parent error) error {
	err := internalNew(kind, parent.Error())
	err.Parent = parent

	return err
}

// Annotate adds a key-value pair to the error's annotations.
// Annotations are used to provide additional context about the error.
// If the error is nil, it returns nil.
// If the error is not of type *internalError, it first Wraps the error.
func Annotate(err error, key string, value any) error {
	if err == nil {
		return nil
	}
	internalErr, ok := err.(*internalError)
	if !ok {
		internalErr = internalNew("", err.Error())
		internalErr.Parent = err
	}
	if internalErr.Annotations == nil {
		internalErr.Annotations = make(map[string]any)
	}
	internalErr.Annotations[key] = value

	return internalErr
}
