package app

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func replaceAttrs(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "error" {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}

		// Unwrap the error tree (if present) generated using %w or errors.Join
		if errs, ok := err.(interface{ Unwrap() []error }); ok {
			var errAttrs []slog.Attr

			for i, err := range errs.Unwrap() {
				errAttrs = append(errAttrs, slog.GroupAttrs(fmt.Sprintf("error_%d", i+1), slog.Any("message", err.Error())))
			}

			return slog.GroupAttrs("errors", errAttrs...)
		}

		return a
	}

	return a
}

func initializeLogger(filepath string) (*slog.Logger, func() error, error) {
	var (
		handlers []slog.Handler
		closers  []func() error
	)

	handlers = append(handlers, tint.NewHandler(os.Stderr, &tint.Options{
		Level:       slog.LevelDebug,
		NoColor:     !isatty.IsCygwinTerminal(os.Stderr.Fd()) && !isatty.IsTerminal(os.Stderr.Fd()),
		ReplaceAttr: replaceAttrs,
	}))

	if filepath != "" {
		logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot open the log file: %v\n", err)
		}

		// Create a buffered writer so that each log doesn't involve an I/O operation
		bufferedFile := bufio.NewWriter(logFile)

		handlers = append(handlers, slog.NewJSONHandler(bufferedFile, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: replaceAttrs,
		}))

		closers = append(closers, func() error {
			// Flush the buffer for the remaining logs before closing
			if err := bufferedFile.Flush(); err != nil {
				return fmt.Errorf("failed to flush the log file: %v\n", err)
			}

			if err := logFile.Close(); err != nil {
				return fmt.Errorf("failed to close the log file: %v\n", err)
			}

			return nil
		})
	}

	close := func() error {
		var errs []error

		for _, closer := range closers {
			if err := closer(); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}

	return slog.New(slog.NewMultiHandler(handlers...)), close, nil
}
