package node

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/raft"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

// SolveExpression is required to satisfy the GRPC interface requirement.
func (n *Node) SolveExpression(context context.Context, request *proto.SolveRequest) (*proto.SolveResponse, error) {

	n.log.Info().
		Str("expression", request.Expression).
		Interface("params", request.Parameters).
		Str("state", n.raft.State().String()).
		Msg("received solve request")

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not marshal solve request: %w", err)
	}

	if n.raft.State() != raft.Leader {
		addr, id := n.raft.LeaderWithID()
		return nil, fmt.Errorf("node is not the leader, send requests to %v at %v", id, addr)
	}

	n.log.Info().Msg("node about to apply raft log")

	future := n.raft.Apply(payload, time.Minute)

	n.log.Info().Msg("node called apply for the raft log, waiting for future")

	err = future.Error()
	if err != nil {
		n.log.Error().Err(err).Msg("raft future returned an error")
		return nil, fmt.Errorf("could not apply raft log: %w", err)
	}

	n.log.Info().Msg("future arrived")

	index := future.Index()
	response := future.Response()

	// Response is either the payload, or it's an error returned from the FSM.
	value, ok := response.(float64)
	if !ok {

		fsmErr, ok := response.(error)
		if ok {
			n.log.Error().Err(fsmErr).Msg("fsm returned an error")
			return nil, fmt.Errorf("fsm returned an error: %w", fsmErr)
		}

		n.log.Error().Msg("unexpected raft response format")
		return nil, fmt.Errorf("unexpected raft response format")
	}

	n.log.Info().
		Uint64("index", index).
		Float64("value", value).
		Msg("node applied raft log")

	out := proto.SolveResponse{
		Expression: request.Expression,
		Parameters: request.Parameters,
		Result:     value,
	}

	return &out, nil
}
