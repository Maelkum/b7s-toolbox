package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func runWriter(host host.Host, peeraddr string) int {

	addr, err := multiaddr.NewMultiaddr(peeraddr)
	if err != nil {
		log.Error().Err(err).Msg("could not parse peer address")
		return failure
	}

	pi, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Error().Err(err).Msg("could not get p2p addr info")
		return failure
	}

	err = host.Connect(context.Background(), *pi)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to peer")
		return failure
	}

	log.Info().Msg("connected to peer")

	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("msg> ")
		text, _ := in.ReadString('\n')

		stream, err := host.NewStream(context.Background(), pi.ID, protocol)
		if err != nil {
			log.Error().Err(err).Msg("could not create stream")
			return failure
		}

		_, err = stream.Write([]byte(text))
		if err != nil {
			stream.Close()
			log.Error().Err(err).Msg("could not write to stream")
			return failure
		}

		log.Info().Msg("sent message")
		stream.Close()
	}

	return success
}
