package node

import (
	"context"
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/rs/zerolog"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

type api struct {
	proto.UnimplementedSolveServer

	log zerolog.Logger
}

func newAPI(log zerolog.Logger) *api {

	api := api{
		log: log.With().Str("@module", "api").Logger(),
	}

	return &api
}

func (a *api) SolveExpression(context context.Context, request *proto.SolveRequest) (*proto.SolveResponse, error) {

	a.log.Info().
		Str("expression", request.Expression).
		Interface("params", request.Parameters).
		Msg("received solve request")

	expression, err := govaluate.NewEvaluableExpression(request.Expression)
	if err != nil {
		return nil, fmt.Errorf("could not parse expression: %w", err)
	}

	params := make(map[string]interface{})
	for name, value := range request.Parameters {
		params[name] = value
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		return nil, fmt.Errorf("could not evaluate expression: %w", err)
	}

	a.log.Info().
		Str("expression", request.Expression).
		Interface("params", request.Parameters).
		Interface("result", result).
		Msg("expression solved")

	value, ok := result.(float64)
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", result)
	}

	out := proto.SolveResponse{
		Expression: request.Expression,
		Parameters: request.Parameters,
		Result:     value,
	}

	return &out, nil
}
