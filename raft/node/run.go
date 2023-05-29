package node

import (
	"fmt"
)

// TODO: What the node should do:
// Receive a mathematical expression and evaluate it.
// https://github.com/Knetic/govaluate

func (n *Node) Run() error {

	n.log.Info().Msg("starting node main loop")

	err := n.grpcServer.Serve(n.listener)
	if err != nil {
		return fmt.Errorf("could not start node: %w", err)
	}

	return nil
}
