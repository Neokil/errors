package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/Neokil/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogNewErrorWithSlog(t *testing.T) {
	logOutput := bytes.NewBuffer([]byte{})
	logger := slog.New(slog.NewJSONHandler(logOutput, &slog.HandlerOptions{}))

	err := errors.New("server-side", "an internal server error occurred")
	errors.Annotate(err, "current_user", "demo-user@example.com")
	logger.Error("an error occurred", "err", err)
	t.Log(logOutput.String())

	logMessage := map[string]any{}
	err = json.Unmarshal(logOutput.Bytes(), &logMessage)
	require.NoError(t, err)

	assert.Equal(t, "ERROR", logMessage["level"])
	assert.Equal(t, "an error occurred", logMessage["msg"])
	assert.NotNil(t, logMessage["time"])
	require.NotNil(t, logMessage["err"])
	errMap, ok := logMessage["err"].(map[string]any)
	require.True(t, ok)
	assert.Len(t, errMap, 4)
	assert.Equal(t, "server-side", errMap["kind"])
	assert.Equal(t, "an internal server error occurred", errMap["message"])
	require.NotNil(t, errMap["annotations"])
	annotationsMap, ok := errMap["annotations"].(map[string]any)
	require.True(t, ok)
	assert.Len(t, annotationsMap, 1)
	assert.Equal(t, "demo-user@example.com", annotationsMap["current_user"])
	assert.NotEmpty(t, errMap["stacktrace"])
	assert.Equal(t, 4, strings.Count(errMap["stacktrace"].(string), "\n"))
}

func TestPrintErrorWithFormat(t *testing.T) {
	err := errors.New("server-side", "an internal server error occurred")
	errors.Annotate(err, "current_user", "demo-user@example.com")
	output := fmt.Sprintf("%+v", err)
	t.Log(output)

	errMap := map[string]any{}
	err = json.Unmarshal([]byte(output), &errMap)
	require.NoError(t, err)

	assert.Len(t, errMap, 4)
	assert.Equal(t, "server-side", errMap["kind"])
	assert.Equal(t, "an internal server error occurred", errMap["message"])
	require.NotNil(t, errMap["annotations"])
	annotationsMap, ok := errMap["annotations"].(map[string]any)
	require.True(t, ok)
	assert.Len(t, annotationsMap, 1)
	assert.Equal(t, "demo-user@example.com", annotationsMap["current_user"])
	assert.NotEmpty(t, errMap["stacktrace"])
	assert.Equal(t, 4, strings.Count(errMap["stacktrace"].(string), "\n"))
}
