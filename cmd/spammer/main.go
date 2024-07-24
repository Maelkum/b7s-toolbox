package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"golang.org/x/time/rate"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

var (
	log  = zerolog.New(os.Stderr).With().Logger().Level(zerolog.DebugLevel)
	done = make(chan struct{})
)

type Ping struct {
	Type string `json:"type"`
	ID   uint64 `json:"id"`
}

func main() {

	var (
		flagAddress string
		flagCount   uint64
		flagRPS     float64
		flagKey     string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "", "address of the b7s node")
	pflag.Uint64VarP(&flagCount, "count", "c", 1, "number of requests to send")
	pflag.Float64Var(&flagRPS, "rps", 100, "requests per second")
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

	limiter := rate.NewLimiter(rate.Limit(flagRPS), 1)

	req := Ping{Type: blockless.MessagePing}
	for i := uint64(0); i < flagCount; i++ {

		_ = limiter.Wait(context.Background())

		// Update request we're sending.
		req.ID = i
		err = sendPing(host, addrInfo.ID, req)
		if err != nil {
			log.Fatal().Err(err).Uint64("i", i).Msg("could not send ping")
		}

	}

	log.Info().Msg("all done")
	<-done
	printStats()
}

func sendPing(host host.Host, peer peer.ID, msg Ping) error {

	stream, err := host.NewStream(context.Background(), peer, blockless.ProtocolID)
	if err != nil {
		return fmt.Errorf("could not create new stream: %w", err)
	}
	defer func() {
		err = stream.Close()
		if err != nil {
			log.Error().Err(err).Uint64("id", msg.ID).Msg("could not close stream")
		}
	}()

	data, _ := json.Marshal(msg)
	data = append(data, '\n')

	_, err = stream.Write(data)
	if err != nil {
		return fmt.Errorf("could not write to stream: %w", err)
	}

	return nil
}
