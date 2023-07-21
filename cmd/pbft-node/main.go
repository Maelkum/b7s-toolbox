package main

import (
	"fmt"
	"os"

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
	pflag.StringSliceVar(&flagPeers, "peers", []string{}, "peers to connect to")
	pflag.StringVar(&flagKey, "private-key", "", "private key to use")

	pflag.Parse()

	if len(flagPeers) == 0 {
		log.Error().Msg("peer list cannot be empty")
		return 1
	}

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	peers := make([]multiaddr.Multiaddr, 0, len(flagPeers))
	for _, peer := range flagPeers {
		ma, err := multiaddr.NewMultiaddr(peer)
		if err != nil {
			log.Error().Err(err).Msg("could not parse multiaddress")
			return 1
		}

		peers = append(peers, ma)
	}

	host, err := host.New(log, flagAddress, flagPort, host.WithBootNodes(peers))
	if err != nil {
		log.Error().Err(err).Msg("could not create host")
		return 1
	}

	log.Info().Str("id", host.Host.ID().String()).Strs("address", host.Addresses()).Msg("created host")

	ids := make([]peer.ID, 0, len(peers))
	for _, p := range peers {

		addr, err := peer.AddrInfoFromP2pAddr(p)
		if err != nil {
			log.Error().Err(err).Msg("could not extract addrinfo from address")
			return 1
		}

		ids = append(ids, addr.ID)
	}

	if flagKey == "" {
		log.Error().Msg("key is mandatory")
		return 1
	}

	key, err := readPrivateKey(flagKey)
	if err != nil {
		log.Error().Err(err).Msg("could not read private key")
		return 1
	}

	pbft, err := pbft.NewReplica(log, host, ids, key)
	if err != nil {
		log.Error().Err(err).Msg("could not initialize replica")
		return 1
	}

	_ = pbft

	log.Info().Msg("all done")

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
