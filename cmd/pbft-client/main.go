package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/consensus/pbft"
	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
	"github.com/blocklessnetworking/b7s/models/execute"
	"github.com/blocklessnetworking/b7s/models/response"
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

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	host, err := host.New(log, flagAddress, flagPort)
	if err != nil {
		log.Error().Err(err).Msg("could not create host")
		return 1
	}

	// Delay start for cluster to start up.
	time.Sleep(3 * time.Second)

	log.Info().Str("id", host.Host.ID().String()).Msg("created host")

	var peers []peer.ID
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

		err = host.Host.Connect(context.Background(), *addrInfo)
		if err != nil {
			log.Error().Err(err).Msg("could not connect to peer")
			return 1
		}

		peers = append(peers, addrInfo.ID)
	}

	log.Info().Msg("connected to peers")

	// Send request to all members of the cluster.
	id := "dummy-request-id"
	request := pbft.Request{
		ID:        id,
		Timestamp: time.Now(),
		Origin:    host.ID(),
		Execute:   execute.Request{},
	}

	payload, _ := json.Marshal(request)

	for _, peer := range peers {
		err = host.SendMessage(context.Background(), peer, payload)
		if err != nil {
			log.Error().Err(err).Str("peer", peer.String()).Msg("could not send message to peer")
			return 1
		}

		log.Info().Str("peer", peer.String()).Msg("sent request to replica")
	}

	host.Host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		log.Debug().Str("peer", from.String()).Msg("received message")

		var response response.Execute
		err = json.Unmarshal(payload, &response)
		if err != nil {
			log.Error().Err(err).Msg("could not unpack message")
			return
		}

		log.Info().Str("peer", from.String()).Interface("response", response).Msg("received response from a replica peer")
	})

	select {}

	return 0
}
