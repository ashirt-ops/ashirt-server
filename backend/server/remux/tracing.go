package remux

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// StackTraceEntry is a small structure for holding relevent data pertaining to a single line
// in a stacktrace.
type StackTraceEntry struct {
	File       string
	LineNumber int
}

// watcher provides a mechanism to recover from panics. Start this prior to launching
// questionable code. If a panic occurs, this will catch that panic, log the response,
// and execute the provided cleanup code. a boolean provided to the cleanup code
// indicates if a panic occurred.
func watcher(ctx context.Context, log *slog.Logger, cleanup func(bool)) {
	paniced := false
	if r := recover(); r != nil {
		paniced = true
		strTrace := formatStackTrace(retrace(25))
		log.ErrorContext(ctx, "unexpected panic", "recoveredData", r, "stacktrace", strTrace)
	}
	cleanup(paniced)
}

// retrace follows the call stack backwards to determine where the program has been
// This follows up to the specified number of steps, or when the callstack is unavailable
func retrace(maxSteps int) []StackTraceEntry {
	steps := make([]StackTraceEntry, 0, maxSteps)

	for i := 0; i < maxSteps; i++ {
		// 0 = self, 1 = caller, so start at 2 to go back in time
		_, file, line, ok := runtime.Caller(2 + i)
		if ok {
			steps = append(steps, StackTraceEntry{File: file, LineNumber: line})
		} else {
			break
		}
	}
	return steps
}

// formatStackTrace provides a mechanism for turning a list of StackTraceEntries into a single-line
// structure
func formatStackTrace(points []StackTraceEntry) string {
	output := make([]string, len(points))
	glue := " : "

	for i, step := range points {
		output[i] = fmt.Sprintf("[ %v ] @ line %v", step.File, step.LineNumber)
	}
	return strings.Join(output, glue)
}
