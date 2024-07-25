package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"golang.org/x/time/rate"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

var (
	log = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
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
	pflag.Float64Var(&flagRPS, "rps", math.Inf(1), "requests per second")
	pflag.StringVarP(&flagKey, "private-key", "k", "", "private key to use")

	pflag.Parse()

	if flagAddress == "" {
		log.Fatal().Msg("node address is required")
	}

	host := createLibp2pHost(flagKey)

	// Close channel to signal when we're done.
	done := make(chan struct{})
	pongReceiver := getPongReceiver(flagCount, done)
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

	for i := uint64(0); i < flagCount; i++ {

		_ = limiter.Wait(context.Background())

		opts := []libp2p.Option{
			libp2p.DefaultTransports,
			libp2p.DefaultMuxers,
			libp2p.DefaultSecurity,
			libp2p.NATPortMap(),
		}

		h, err := libp2p.New(opts...)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create libp2p host")
		}

		h.Network().Notify(&notifiee{})

		err = h.Connect(context.Background(), *addrInfo)
		if err != nil {
			log.Error().Err(err).Msg("could not connect to b7s node")
			break
		}

		log.Info().Uint64("i", i).Msg("established connection")
	}

	select {}

	/*

		req := Ping{Type: blockless.MessagePing}

		start := time.Now()
		for i := uint64(0); i < flagCount; i++ {

			_ = limiter.Wait(context.Background())

			// Update request we're sending.
			req.ID = i

			err = sendPing(host, addrInfo.ID, req)
			if err != nil {
				log.Fatal().Err(err).Uint64("i", i).Msg("could not send ping")
			}
		}

		<-done
		end := time.Now()

		log.Info().Msg("all done")

		printStats(start, end, flagCount, flagRPS)
	*/
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
