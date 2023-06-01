package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb/v2"

	"github.com/Maelkum/b7s-toolbox/raft/node"
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

	var (
		flagAddress string
		flagPort    uint
		flagID      string
		flagPeers   []string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVar(&flagID, "id", "", "node id")
	pflag.StringSliceVar(&flagPeers, "peer", []string{}, "peers to add to the cluster")

	pflag.Parse()

	log := zerolog.New(os.Stderr).
		With().Timestamp().
		Str("@module", "main").
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339Nano

	if flagID == "" {
		log.Error().Msg("ID not specified")
		return failure
	}

	log.Info().Str("id", flagID).Msg("starting node")

	address := fmt.Sprintf("%v:%v", flagAddress, flagPort)

	err := os.MkdirAll(flagID, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("could not create node directory")
		return failure
	}

	logDBPath := filepath.Join(flagID, logName)
	logStore, err := boltdb.NewBoltStore(logDBPath)
	if err != nil {
		log.Error().Err(err).Str("path", logDBPath).Msg("could not create log store")
		return failure
	}

	stableDBPath := filepath.Join(flagID, stableName)
	stableDB, err := boltdb.NewBoltStore(stableDBPath)
	if err != nil {
		log.Error().Err(err).Str("path", stableDBPath).Msg("could not create stable store")
		return failure
	}

	conn, err := net.Listen("tcp", address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("could not listen on port")
		return failure
	}

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

		peer := node.Peer{
			ID:      fields[0],
			Address: fields[1],
		}

		peers = append(peers, peer)
	}

	node, err := node.NewNode(log, flagID, address, conn,
		node.WithLogStore(logStore),
		node.WithStableStore(stableDB),
		node.WithSnapshotStore(raft.NewDiscardSnapshotStore()),
		node.WithPeers(peers),
	)
	if err != nil {
		log.Error().Err(err).
			Str("id", flagID).
			Msg("could not create node")
		return failure
	}

	err = node.Run()
	if err != nil {
		log.Error().Err(err).Msg("node main loop failed")
		return failure
	}

	return success
}
