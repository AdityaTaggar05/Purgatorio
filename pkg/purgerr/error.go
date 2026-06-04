package purgerr

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// Error wraps a domain sentinel with a low-level cause and a captured stack trace.
//
// .Error() returns only the sentinel message — safe to propagate to API consumers.
// When passed to slog.Any("error", err), LogValue emits a structured group
// {message, cause, stack} automatically, with no replaceAttrs needed.
type Error struct {
	sentinel error
	cause    error
	frames   []frame
}

type frame struct {
	fn   string
	file string
	line int
}

func Wrap(sentinel, cause error) *Error {
	return &Error{
		sentinel: sentinel,
		cause:    cause,
		frames:   capture(2),
	}
}

func (e *Error) Error() string { return e.sentinel.Error() }

func (e *Error) Unwrap() error { return e.cause }

func (e *Error) Is(target error) bool {
	return errors.Is(e.sentinel, target)
}

func (e *Error) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("message", e.sentinel.Error()),
	}
	if e.cause != nil {
		attrs = append(attrs, slog.String("cause", e.cause.Error()))
	}
	if len(e.frames) > 0 {
		attrs = append(attrs, slog.String("stack", e.formatStack()))
	}
	return slog.GroupValue(attrs...)
}

func (e *Error) formatStack() string {
	var b strings.Builder
	for _, f := range e.frames {
		fmt.Fprintf(&b, "\n  %s\n    %s:%d", f.fn, f.file, f.line)
	}
	return b.String()
}

func capture(skip int) []frame {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(skip+1, pcs)
	cf := runtime.CallersFrames(pcs[:n])
	var out []frame
	for {
		f, more := cf.Next()
		if f.Function != "" && !isRuntime(f.File) {
			out = append(out, frame{fn: f.Function, file: f.File, line: f.Line})
		}
		if !more {
			break
		}
	}
	return out
}

func isRuntime(file string) bool {
	return strings.Contains(file, "runtime/") || strings.Contains(file, "testing/")
}
