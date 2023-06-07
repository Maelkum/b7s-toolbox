package node

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"

	libp2praft "github.com/libp2p/go-libp2p-raft"

	"github.com/blocklessnetworking/b7s/host"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

type Node struct {
	proto.UnimplementedSolveServer

	log  zerolog.Logger
	host *host.Host
	cfg  Config

	raft   *raft.Raft
	solver *solver
}

func NewNode(log zerolog.Logger, host *host.Host, options ...Option) (*Node, error) {

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
	raftCfg.LocalID = raft.ServerID(host.ID().String())
	raftCfg.Logger = raftLogger

	addresses := host.Addresses()
	if len(addresses) == 0 {
		return nil, fmt.Errorf("libp2p host has no addresses")
	}

	address := addresses[0]

	solver := newSolver(log)

	transport, err := libp2praft.NewLibp2pTransport(host, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p transport: %w", err)
	}

	raftNode, err := raft.NewRaft(raftCfg, solver, cfg.LogStore, cfg.StableStore, cfg.SnapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("could not create raft node: %w", err)
	}

	node := Node{
		log:  log.With().Str("@module", "node").Logger(),
		host: host,
		cfg:  cfg,

		raft:   raftNode,
		solver: solver,
	}

	// If we have no peers we're useless so return an error.
	if len(cfg.Peers) == 0 {
		return nil, fmt.Errorf("no known peers")
	}

	node.log.Info().Msg("bootstrapping cluster")

	servers := make([]raft.Server, 0, len(cfg.Peers)+1)

	// Add self to cluster.
	self := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(host.ID().String()),
		Address:  raft.ServerAddress(address),
	}

	servers = append(servers, self)

	for _, peer := range cfg.Peers {

		s := raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(peer.ID),
			Address:  raft.ServerAddress(peer.Address.String()),
		}

		// Self is already added.
		if s.Address == self.Address {
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
