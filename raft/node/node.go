package node

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"

	transport "github.com/Jille/raft-grpc-transport"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

type Node struct {
	proto.UnimplementedSolveServer

	log zerolog.Logger
	cfg Config

	id       raft.ServerID
	address  raft.ServerAddress
	raft     *raft.Raft
	listener net.Listener
	solver   *solver

	grpcServer *grpc.Server
}

func NewNode(log zerolog.Logger, id string, address string, listener net.Listener, options ...Option) (*Node, error) {

	cfg := Config{}
	for _, option := range options {
		option(&cfg)
	}

	logOpts := hclog.LoggerOptions{
		JSONFormat: true,
		Level:      hclog.Debug,
		Output:     os.Stderr,
		Name:       "raft",
	}
	raftLogger := hclog.New(&logOpts)

	raftCfg := raft.DefaultConfig()
	raftCfg.LocalID = raft.ServerID(id)
	raftCfg.Logger = raftLogger

	solver := newSolver(log)
	// NOTE: Using a fixed transport for now.
	transport := transport.New(raft.ServerAddress(address), []grpc.DialOption{grpc.WithInsecure()})
	raftNode, err := raft.NewRaft(raftCfg, solver, cfg.LogStore, cfg.StableStore, cfg.SnapshotStore, transport.Transport())
	if err != nil {
		return nil, fmt.Errorf("could not create raft node: %w", err)
	}

	server := grpc.NewServer()

	node := Node{
		log: log.With().Str("@module", "node").Logger(),
		cfg: cfg,

		id:         raft.ServerID(id),
		address:    raft.ServerAddress(address),
		raft:       raftNode,
		listener:   listener,
		solver:     solver,
		grpcServer: server,
	}

	proto.RegisterSolveServer(server, &node)
	transport.Register(server)
	reflection.Register(server)

	// If we have no peers we're useless so return an error.
	if len(cfg.Peers) == 0 {
		return nil, fmt.Errorf("no known peers")
	}

	node.log.Info().Msg("bootstrapping cluster")

	servers := make([]raft.Server, 0, len(cfg.Peers)+1)

	// Add self to cluster.
	self := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(id),
		Address:  raft.ServerAddress(address),
	}

	servers = append(servers, self)

	for _, peer := range cfg.Peers {

		s := raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(peer.ID),
			Address:  raft.ServerAddress(peer.Address),
		}

		// Self is already added.
		if s.ID == self.ID {
			continue
		}

		servers = append(servers, s)
	}

	clusterCfg := raft.Configuration{
		Servers: servers,
	}

	node.log.Info().Interface("cluster", clusterCfg).Msg("bootstrapping cluster")

	ret := raftNode.BootstrapCluster(clusterCfg)
	err = ret.Error()
	if err != nil && !errors.Is(err, raft.ErrCantBootstrap) {
		return nil, fmt.Errorf("could not bootstrap cluster: %w", err)
	}

	return &node, nil
}
