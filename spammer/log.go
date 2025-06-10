package spammer

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

func defaultPrefix() string {

	const timeFormat = "0201150405"
	return time.Now().Format(timeFormat)
}

func logName(cfg testConfig) string {
	return fmt.Sprintf("spammer_%v_%v_f%v.log", defaultPrefix(), cfg.executions, cfg.frequency)
}

func mustCreateLogFile(cfg testConfig) io.WriteCloser {
	f, err := createLogFile(cfg)
	if err != nil {
		panic("could not create log file")
	}
	return f
}

func mustCreateOutputFile(cfg testConfig) io.WriteCloser {
	name := fmt.Sprintf("stats_%v_%v_f%v.txt", defaultPrefix(), cfg.executions, cfg.frequency)
	f, err := os.Create(name)
	if err != nil {
		panic("could not create output file")
	}

	slog.Debug("creating output file",
		"name", name,
		"executions", cfg.executions,
		"frequency", cfg.frequency)

	return f
}

func createLogFile(cfg testConfig) (io.WriteCloser, error) {

	name := logName(cfg)
	f, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %w", err)
	}

	slog.Debug("creating log file",
		"name", name,
		"executions", cfg.executions,
		"frequency", cfg.frequency)

	return f, nil
}

// NOTE: Using default logger but I don't care.
func initLogger(w io.WriteCloser) {
	slog.SetDefault(
		slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})),
	)
}
