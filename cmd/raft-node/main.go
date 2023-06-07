package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb/v2"

	"github.com/Maelkum/b7s-toolbox/raft/node"
	"github.com/blocklessnetworking/b7s/host"
)

const (
	success = 0
	failure = 1
)

const (
	logName    = "logs.dat"
	stableName = "stable.dat"
)

func main() {
	os.Exit(run())
}

func run() int {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagAddress    string
		flagPort       uint
		flagPrivateKey string
		flagPeers      []string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVar(&flagPrivateKey, "private-key", "", "private key that the libp2p host will use")
	pflag.StringSliceVar(&flagPeers, "peer", []string{}, "peers to add to the cluster")

	pflag.Parse()

	log := zerolog.New(os.Stderr).
		With().Timestamp().
		Str("@module", "main").
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339Nano

	if len(flagPeers) == 0 {
		log.Error().Msg("peer list cannot be empty")
		return failure
	}

	var peers []node.Peer
	for _, p := range flagPeers {

		fields := strings.Split(p, "=")
		if len(fields) != 2 {
			log.Error().Str("peer", p).Msg("invalid peer specification, should be <node_id>=<address>")
			return failure
		}

		addr, err := multiaddr.NewMultiaddr(fields[1])
		if err != nil {
			log.Error().Err(err).Str("addr", fields[1]).Msg("could not parse multiaddress")
			return failure
		}

		peer := node.Peer{
			ID:      fields[0],
			Address: addr,
		}

		peers = append(peers, peer)
	}

	addrs := make([]multiaddr.Multiaddr, 0, len(peers))
	for _, peer := range peers {
		addrs = append(addrs, peer.Address)
	}

	host, err := host.New(log, flagAddress, flagPort,
		host.WithPrivateKey(flagPrivateKey),
		host.WithBootNodes(addrs),
	)
	if err != nil {
		log.Error().Err(err).Msg("could not create libp2p host")
		return failure
	}

	hostID := host.ID().String()

	log.Info().Str("id", hostID).Msg("starting node")

	for i, addr := range host.Addresses() {
		log.Info().Int("i", i).Str("address", addr).Msg("host address")
	}

	err = os.MkdirAll(hostID, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("could not create node directory")
		return failure
	}

	logDBPath := filepath.Join(hostID, logName)
	logStore, err := boltdb.NewBoltStore(logDBPath)
	if err != nil {
		log.Error().Err(err).Str("path", logDBPath).Msg("could not create log store")
		return failure
	}

	stableDBPath := filepath.Join(hostID, stableName)
	stableDB, err := boltdb.NewBoltStore(stableDBPath)
	if err != nil {
		log.Error().Err(err).Str("path", stableDBPath).Msg("could not create stable store")
		return failure
	}

	node, err := node.NewNode(log, host,
		node.WithLogStore(logStore),
		node.WithStableStore(stableDB),
		node.WithSnapshotStore(raft.NewDiscardSnapshotStore()),
		node.WithPeers(peers),
	)
	if err != nil {
		log.Error().Err(err).
			Str("id", hostID).
			Msg("could not create node")
		return failure
	}

	// Create the main context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = node.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("node main loop failed")
		}
	}()

	<-sig
	log.Info().Msg("node stopping")

	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	return success
}
