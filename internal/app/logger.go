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

func initializeLogger(filepath string) (*slog.Logger, func() error, error) {
	var (
		handlers []slog.Handler
		closers  []func() error
	)

	handlers = append(handlers, tint.NewHandler(os.Stderr, &tint.Options{
		Level:   slog.LevelDebug,
		NoColor: !isatty.IsCygwinTerminal(os.Stderr.Fd()) && !isatty.IsTerminal(os.Stderr.Fd()),
	}))

	if filepath != "" {
		logFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot open the log file: %v\n", err)
		}

		bufferedFile := bufio.NewWriter(logFile)

		handlers = append(handlers, slog.NewJSONHandler(bufferedFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		closers = append(closers, func() error {
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
