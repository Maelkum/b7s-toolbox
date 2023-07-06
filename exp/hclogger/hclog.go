package main

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func getHCLog(level hclog.Level) hclog.Logger {

	opts := hclog.LoggerOptions{
		Level:      level,
		Output:     os.Stderr,
		JSONFormat: true,
	}

	logger := hclog.New(&opts)
	return logger
}

func getLevel() hclog.Level {
	if len(os.Args) != 2 {
		log.Fatal("usage: hclogger <level>")
	}

	level := strings.TrimSpace(strings.ToLower(os.Args[1]))

	switch level {
	case "trace":
		return hclog.Trace
	case "debug":
		return hclog.Debug
	case "info":
		return hclog.Info
	case "warn":
		return hclog.Warn
	case "error":
		return hclog.Error
	}

	log.Fatal("invalid log level")
	return hclog.Off
}
