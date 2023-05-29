package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagAddress    string
		flagExpression string
		flagParams     []string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "", "address of the GRPC server")
	pflag.StringVarP(&flagExpression, "expression", "e", "", "expression to be evaluated")
	pflag.StringSliceVarP(&flagParams, "param", "p", []string{}, "parameters")

	pflag.Parse()

	log := zerolog.New(os.Stderr).
		With().Timestamp().
		Logger()

	if flagExpression == "" {
		log.Error().Msg("expression cannot be empty")
		return failure
	}

	if flagAddress == "" {
		log.Error().Msg("server address cannot be empty")
		return failure
	}

	params := make(map[string]float64)
	for _, param := range flagParams {

		fields := strings.Split(param, "=")
		if len(fields) != 2 {
			log.Error().Str("param", param).Msg("parameter should be in the 'name=123.4' format")
			return failure
		}

		name := fields[0]
		val := fields[1]

		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Error().Str("value", val).Msg("value has to be a valid number")
			return failure
		}

		params[name] = value
	}

	log.Info().
		Str("expression", flagExpression).
		Interface("params", params).
		Msg("prepared expression")

	conn, err := grpc.Dial(flagAddress, grpc.WithInsecure())
	if err != nil {
		log.Error().Err(err).Msg("could not connect to server")
		return failure
	}

	client := proto.NewSolveClient(conn)

	request := proto.SolveRequest{
		Expression: flagExpression,
		Parameters: params,
	}

	response, err := client.SolveExpression(context.Background(), &request)
	if err != nil {
		log.Error().Err(err).Msg("solve request failed")
		return failure
	}

	out, _ := json.Marshal(response)
	fmt.Printf("%s\n", out)

	return success
}
