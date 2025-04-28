package errors

import "runtime"

const stackTraceSize = 1024 * 2

type internalError struct {
	Kind        string
	Message     string
	Annotations map[string]any
	Stacktrace  string
	Parent      error
}

func (e *internalError) Error() string {
	return e.Message
}
func (e *internalError) Unwrap() error {
	return e.Parent
}

func internalNew(kind string, msg string) *internalError {
	stacktrace := make([]byte, stackTraceSize)
	n := runtime.Stack(stacktrace, false)

	return &internalError{
		Message:    msg,
		Stacktrace: string(stacktrace[:n]),
	}
}

func New(kind string, msg string) error {
	return internalNew(kind, msg)
}

func Wrap(kind string, parent error) error {
	err := internalNew(kind, parent.Error())
	err.Parent = parent

	return err
}

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
