package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/consensus/pbft"
	"github.com/blocklessnetworking/b7s/host"
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagAddress string
		flagPort    uint
		flagPeers   []string
		flagKey     string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "0.0.0.0", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringSliceVar(&flagPeers, "peer", []string{}, "peers to connect to")
	pflag.StringVar(&flagKey, "private-key", "", "private key to use")

	pflag.Parse()

	if len(flagPeers) == 0 {
		log.Error().Msg("peer list cannot be empty")
		return 1
	}

	if flagKey == "" {
		log.Error().Msg("key is mandatory")
		return 1
	}

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	host, err := host.New(log, flagAddress, flagPort, host.WithPrivateKey(flagKey))
	if err != nil {
		log.Error().Err(err).Msg("could not create host")
		return 1
	}

	log.Info().Str("id", host.Host.ID().String()).Strs("address", host.Addresses()).Msg("created host")

	// Wait for other replicas to boot up and try to connect.
	// Could be done via peer discovery.
	time.Sleep(1 * time.Second)

	var replicas []peer.ID
	for _, addr := range flagPeers {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Error().Err(err).Msg("could not parse multiaddress")
			return 1
		}

		addrInfo, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			log.Error().Err(err).Msg("could not extract addrinfo from address")
			return 1
		}

		replicas = append(replicas, addrInfo.ID)

		// Skip self.
		if addrInfo.ID == host.ID() {
			continue
		}

		err = host.Host.Connect(context.Background(), *addrInfo)
		if err != nil {
			log.Error().Err(err).Msg("could not connect to peer")
			return 1
		}
	}

	log.Info().Msg("connected to other replicas")

	key, err := readPrivateKey(flagKey)
	if err != nil {
		log.Error().Err(err).Msg("could not read private key")
		return 1
	}

	pbft, err := pbft.NewReplica(log, host, DummyExecutor{}, replicas, key)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize replica")
		return 1
	}

	_ = pbft

	log.Info().Msg("all done")

	select {}

	return 0
}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}
