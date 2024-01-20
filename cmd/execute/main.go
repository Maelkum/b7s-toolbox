package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s/executor"
	"github.com/blocklessnetwork/b7s/executor/limits"
	"github.com/blocklessnetwork/b7s/models/execute"
)

const (
	success = 0
	failure = 1
)

const (
	sleepEnvVar = "B7S_SLEEP"
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
		flagEnv        []string

		flagCPURate     float64
		flagMemoryMaxKB int64

		flagCLIDuration int

		flagCfg execute.RuntimeConfig
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "debug", "log level to use")
	pflag.StringVar(&flagWorkspace, "workspace", "./workspace", "directory that the executor can use for file storage")
	pflag.StringVar(&flagRuntime, "runtime", "", "runtime path")

	pflag.StringVarP(&flagFunctionID, "function-id", "f", "", "function id to execute")
	pflag.StringVarP(&flagMethod, "method", "m", "", "function method")
	pflag.StringVarP(&flagStdin, "stdin", "s", "", "stdin of the application to be executed")
	pflag.StringArrayVarP(&flagArgs, "arg", "a", []string{}, "cli argument to the application to be executed")
	pflag.StringArrayVarP(&flagEnv, "env", "e", []string{}, "environment variables to pass, in the format 'name=value'")

	pflag.StringVar(&flagCfg.Entry, "runtime-entry", "", "runtime entry")
	pflag.StringVar(&flagCfg.Logger, "runtime-logger", "", "runtime logger")
	pflag.Uint64Var(&flagCfg.Fuel, "runtime-fuel", 0, "runtime fuel")
	pflag.Uint64Var(&flagCfg.Memory, "runtime-memory", 0, "runtime memory")
	pflag.Uint64Var(&flagCfg.ExecutionTime, "runtime-execution-time", 0, "runtime execution time")
	pflag.BoolVar(&flagCfg.DebugInfo, "runtime-debug", false, "runtime debug")

	pflag.Float64Var(&flagCPURate, "cpu-rate", 1.0, "cpu rate")
	pflag.Int64Var(&flagMemoryMaxKB, "memory-limit", 0, "memory limit (kB)")

	pflag.IntVar(&flagCLIDuration, "cli-duration", 0, "value for the B7S_SLEEP environment variable")

	// Runtime flags

	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	log = zerolog.New(os.Stderr)

	level, err := zerolog.ParseLevel(flagLogLevel)
	if err == nil {
		log = log.Level(level)
	}

	options := []executor.Option{
		executor.WithWorkDir(flagWorkspace),
		executor.WithRuntimeDir(flagRuntime),
	}

	// Create a limiter if needed.
	if flagCPURate != 1.0 || flagMemoryMaxKB > 0 {

		log.Info().Msg("creating limiter")

		limiter, err := limits.New(limits.WithCPUPercentage(flagCPURate), limits.WithMemoryKB(flagMemoryMaxKB))
		if err != nil {
			log.Error().Err(err).Msg("could not create limiter")
			return failure
		}

		defer limiter.Shutdown()

		options = append(options, executor.WithLimiter(limiter))
	}

	ex, err := executor.New(
		log,
		options...,
	)
	if err != nil {
		log.Error().Err(err).Msg("could not create executor")
		return failure
	}

	requestID := "dummy-request-id"
	request := execute.Request{
		FunctionID: flagFunctionID,
		Method:     flagMethod,
		Config: execute.Config{
			Runtime: flagCfg,
		},
	}

	// Add stdin if specified.
	if flagStdin != "" {
		request.Config.Stdin = &flagStdin
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

	// Set environment variables if needed.
	if len(flagEnv) > 0 {
		vars := make([]execute.EnvVar, 0, len(flagEnv))
		for _, env := range flagEnv {

			fields := strings.Split(env, "=")
			if len(fields) != 2 {
				log.Error().Str("input", env).Msg("bad environment variable format")
				return failure
			}

			v := execute.EnvVar{
				Name:  fields[0],
				Value: fields[1],
			}
			vars = append(vars, v)
		}

		request.Config.Environment = vars
	}

	// Set the environment variable for the dummy CLI, if specified.
	if flagCLIDuration != 0 {
		sleepVar := execute.EnvVar{
			Name:  sleepEnvVar,
			Value: fmt.Sprint(flagCLIDuration),
		}
		request.Config.Environment = append(request.Config.Environment, sleepVar)
	}

	log.Info().Interface("request", request).Msg("request to be executed")

	res, err := ex.ExecuteFunction(requestID, request)
	if err != nil {
		log.Error().Err(err).Str("request_id", requestID).Msg("function execution failed")
	}

	log.Info().Interface("response", res).Msg("execution result")

	return success
}
