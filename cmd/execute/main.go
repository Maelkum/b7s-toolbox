package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/executor"
	"github.com/blocklessnetworking/b7s/models/execute"
)

const (
	success = 0
	failure = 1
)

var (
	log zerolog.Logger
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagLogLevel  string
		flagWorkspace string
		flagRuntime   string

		flagFunctionID string
		flagMethod     string
		flagStdin      string
		flagArgs       []string
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "debug", "log level to use")
	pflag.StringVar(&flagWorkspace, "workspace", "./workspace", "directory that the executor can use for file storage")
	pflag.StringVar(&flagRuntime, "runtime", "", "runtime path")

	pflag.StringVarP(&flagFunctionID, "function-id", "f", "", "function id to execute")
	pflag.StringVarP(&flagMethod, "method", "m", "", "function method")
	pflag.StringVarP(&flagStdin, "stdin", "s", "", "stdin of the application to be executed")
	pflag.StringArrayVarP(&flagArgs, "arg", "a", []string{}, "cli argument to the application to be executed")
	// TODO: Add environment variables.

	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	log = zerolog.New(os.Stderr)

	level, err := zerolog.ParseLevel(flagLogLevel)
	if err == nil {
		log = log.Level(level)
	}

	ex, err := executor.New(
		log,
		executor.WithWorkDir(flagWorkspace),
		executor.WithRuntimeDir(flagRuntime),
	)
	if err != nil {
		log.Error().Err(err).Msg("could not create executor")
	}

	requestID := "dummy-request-id"
	request := execute.Request{
		FunctionID: flagFunctionID,
		Method:     flagMethod,
	}

	// Add stdin if specified.
	if flagStdin != "" {
		request.Config = execute.Config{
			Stdin: &flagStdin,
		}
	}

	// Add args if specified.
	if len(flagArgs) > 0 {
		params := make([]execute.Parameter, 0, len(flagArgs))
		for _, arg := range flagArgs {

			p := execute.Parameter{
				Name:  "",
				Value: arg,
			}
			params = append(params, p)
		}

		request.Parameters = params
	}

	log.Info().Interface("request", request).Msg("request to be executed")

	res, err := ex.ExecuteFunction(requestID, request)
	if err != nil {
		log.Error().Err(err).Str("request_id", requestID).Msg("function execution failed")
	}

	log.Info().Interface("response", res).Msg("execution result")

	return success
}
