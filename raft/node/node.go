package node

import (
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/hashicorp/raft"

	"github.com/Maelkum/b7s-toolbox/raft/proto"
)

type Node struct {
	log zerolog.Logger
	cfg Config

	id       raft.ServerID
	address  raft.ServerAddress
	node     *raft.Raft
	listener net.Listener

	grpcServer *grpc.Server
	api        *api
}

func NewNode(log zerolog.Logger, id string, address string, listener net.Listener, options ...Option) (*Node, error) {

	cfg := Config{}
	for _, option := range options {
		option(&cfg)
	}

	// logOpts := hclog.LoggerOptions{
	// 	JSONFormat: true,
	// 	Level:      hclog.Debug,
	// 	Output:     os.Stderr,
	// 	Name:       "raft",
	// }
	// raftLogger := hclog.New(&logOpts)

	// raftCfg := raft.DefaultConfig()
	// raftCfg.LocalID = raft.ServerID(id)
	// raftCfg.Logger = raftLogger

	// fsm := newFSM(log)

	// // NOTE: Using a fixed transport for now.
	// transport := transport.New(raft.ServerAddress(address), []grpc.DialOption{grpc.WithInsecure()})
	// raftNode, err := raft.NewRaft(raftCfg, fsm, cfg.LogStore, cfg.StableStore, cfg.SnapshotStore, transport.Transport())
	// if err != nil {
	// 	return nil, fmt.Errorf("could not create raft node: %w", err)
	// }

	api := newAPI(log)
	server := grpc.NewServer()

	proto.RegisterSolveServer(server, api)
	// transport.Register(server)
	// leaderhealth.Setup(raftNode, server, []string{"Test"})
	// raftadmin.Register(server, raftNode)
	// reflection.Register(server)

	node := Node{
		log: log.With().Str("@module", "node").Logger(),
		cfg: cfg,

		id:      raft.ServerID(id),
		address: raft.ServerAddress(address),
		// node:       raftNode,
		listener:   listener,
		grpcServer: server,
		api:        api,
	}

	// If we're not boostrapping the cluster, we're done.
	if !node.cfg.Bootstrap {
		node.log.Debug().Msg("node created")
		return &node, nil
	}

	node.log.Info().Msg("bootstrapping cluster")

	// clusterCfg := raft.Configuration{
	// 	Servers: []raft.Server{
	// 		{
	// 			Suffrage: raft.Voter,
	// 			ID:       node.id,
	// 			Address:  node.address,
	// 		},
	// 	},
	// }
	//
	// ret := raftNode.BootstrapCluster(clusterCfg)
	// err = ret.Error()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not bootstrap cluster: %w", err)
	// }

	return &node, nil
}
