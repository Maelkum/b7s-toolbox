package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"

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
		flagAddress   string
		flagPort      uint
		flagID        string
		flagBootstrap bool
	)

	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVar(&flagID, "id", "", "node id")
	pflag.BoolVar(&flagBootstrap, "bootstrap", false, "should node cluster be bootstrapped")

	pflag.Parse()

	log := zerolog.New(os.Stderr).
		With().Timestamp().
		Str("@module", "main").
		Logger()

	if flagID == "" {
		log.Error().Msg("ID not specified")
		return failure
	}

	log.Info().Str("id", flagID).Msg("starting node")

	address := fmt.Sprintf("%v:%v", flagAddress, flagPort)

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

	node, err := node.NewNode(log, flagID, address,
		node.BootstrapCluster(flagBootstrap),
		node.WithLogStore(logStore),
		node.WithStableStore(stableDB),
		node.WithSnapshotStore(raft.NewDiscardSnapshotStore()),
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
