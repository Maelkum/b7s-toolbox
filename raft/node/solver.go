package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Knetic/govaluate"
	"github.com/hashicorp/raft"
	"github.com/rs/zerolog"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

type solver struct {
	log zerolog.Logger
}

func newSolver(log zerolog.Logger) *solver {

	fsm := solver{
		log: log.With().Str("@module", "solver").Logger(),
	}

	return &fsm
}

func (f *solver) Apply(log *raft.Log) interface{} {

	f.log.Info().Msg("received Apply() call")

	payload := log.Data

	var request *proto.SolveRequest
	err := json.Unmarshal(payload, &request)
	if err != nil {
		return fmt.Errorf("could not unmarshal request: %w", err)
	}

	expression, err := govaluate.NewEvaluableExpression(request.Expression)
	if err != nil {
		return fmt.Errorf("could not parse expression: %w", err)
	}

	params := make(map[string]interface{})
	for name, value := range request.Parameters {
		params[name] = value
	}

	result, err := expression.Evaluate(params)
	if err != nil {
		return fmt.Errorf("could not evaluate expression: %w", err)
	}

	f.log.Info().
		Str("expression", request.Expression).
		Interface("params", request.Parameters).
		Interface("result", result).
		Msg("expression solved")

	value, ok := result.(float64)
	if !ok {
		return fmt.Errorf("unexpected result type: %T", result)
	}

	return value
}

func (f *solver) Snapshot() (raft.FSMSnapshot, error) {
	f.log.Info().Msg("received Snapshot() call")
	return nil, errors.New("TBD: Not implemented")
}

func (f *solver) Restore(snapshot io.ReadCloser) error {
	f.log.Info().Msg("received Restore() call")
	return errors.New("TBD: Not implemented")
}
