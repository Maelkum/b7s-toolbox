package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/request"
)

var log = zerolog.New(os.Stderr).With().Logger().Level(zerolog.DebugLevel)

func main() {

	var (
		flagAddress string
		flagCount   uint64
		flagRPS     uint
		flagKey     string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "", "address of the b7s node")
	pflag.Uint64VarP(&flagCount, "count", "c", 1, "number of requests to send")
	pflag.UintVar(&flagRPS, "rps", 0, "requests per second limit")
	pflag.StringVarP(&flagKey, "private-key", "k", "", "private key to use")

	pflag.Parse()

	if flagAddress == "" {
		log.Fatal().Msg("node address is required")
	}

	host := createLibp2pHost(flagKey)
	host.SetStreamHandler(blockless.ProtocolID, pongReceiver)

	addrInfo, err := peer.AddrInfoFromP2pAddr(multiaddr.StringCast(flagAddress))
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse multiaddress")
	}

	err = host.Connect(context.Background(), *addrInfo)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to b7s node")
	}

	for i := uint64(0); i < flagCount; i++ {

		data, _ := json.Marshal(request.Ping{ID: i})
		_ = data
	}
}
