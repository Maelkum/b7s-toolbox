package main

import (
	"errors"
	"time"

	//	b7shclog "github.com/blocklessnetworking/b7s/log/hclog"
	"github.com/hashicorp/go-hclog"
)

func main() {

	// Hashicorp logger.
	logger := getHCLog(getLevel())

	// Our logger.
	// log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(getZerologLevel())
	// logger := b7shclog.New(log)

	run(logger)
}

func run(log hclog.Logger) {

	log.Trace("here's a trace message")
	log.Trace("tracing with fields",
		"timestampf", time.Now().Format(time.RFC3339),
		"value", 45,
		"str", "whatever",
	)

	log.Debug("here's a debug message")
	log.Debug("debug with fields",
		"username", "aco",
	)

	log.Info("here's an info message")
	log.Info("info with fields",
		"password", "wutwut",
	)

	log.Warn("here's a warn message")
	log.Warn("warn with fields",
		"object", map[string]any{
			"key":       "value",
			"other-key": 45,
		},
	)

	log.Error("here's an error message")
	log.Error("error with fields",
		"error", errors.New("some generic error"),
		"timestampf", time.Now().Format(time.RFC3339),
		"dbname", "something something",
	)

	log.Error("showing levels",
		"is_trace", log.IsTrace(),
		"is_debug", log.IsDebug(),
		"is_info", log.IsInfo(),
		"is_warn", log.IsWarn(),
		"is_error", log.IsError(),
	)

	log.Log(hclog.Info, "this is a log-log message",
		"field", "value",
	)

	sublogger := log.With("subsystem", "dummy")
	sublogger.Info("this is a dummy sublogger that we're using")

	sublogger.Info("this is the dummy sublogger again, with fields",
		"field1", "value1",
		"field2", "value2",
		"field3", "value3",
	)

	sublogger.Info("this is my current log level (I'm info btw)",
		"level_returned", sublogger.GetLevel().String())

	named := log.Named("super-duper-logger")
	named.Info("here's a named logger")

	second := named.Named("another-logger")
	second.Info("here's another named logger")

	rlog := second.ResetNamed("single-logger")

	rlog.Info("we're now logging a message that should have the logger name reset")
}
