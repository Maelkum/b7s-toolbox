package main

import (
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func getZerologLevel() zerolog.Level {

	if len(os.Args) != 2 {
		log.Fatal("usage: hclogger <level>")
	}

	level := strings.TrimSpace(strings.ToLower(os.Args[1]))

	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	}

	log.Fatal("invalid log level")
	return zerolog.Disabled
}
