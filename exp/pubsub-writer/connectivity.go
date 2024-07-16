package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/sync/errgroup"
)

func connectToBootNodes(h host.Host, addrs []multiaddr.Multiaddr) error {

	var eg errgroup.Group
	for _, ma := range addrs {

		eg.Go(func() error {

			_, id := peer.SplitAddr(ma)
			if err := id.Validate(); err != nil {
				log.Fatal().Err(err).Msg("empty peerID found")
			}

			info := peer.AddrInfo{
				ID:    id,
				Addrs: []multiaddr.Multiaddr{ma},
			}

			err := h.Connect(context.Background(), info)
			if err != nil {
				log.Error().Err(err).Stringer("addr", ma).Stringer("perr", id).Msg("could not connect to boot node")
			}

			return err
		})
	}

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("could not connect to peers: %w", err)
	}

	return nil
}

func sayHello(h host.Host, peers []peer.ID) error {

	var eg errgroup.Group
	for _, peer := range peers {

		eg.Go(func() error {

			type helloRecord struct {
				From      string `json:"from"`
				Message   string `json:"message"`
				Timestamp string `json:"timestamp"`
			}

			hello := helloRecord{
				From:      h.ID().String(),
				Message:   "hello there!",
				Timestamp: time.Now().String(),
			}

			payload, _ := json.Marshal(hello)

			err := send(h, peer, payload)
			if err != nil {
				log.Error().Err(err).Msg("could not send message")
			}

			return err
		})
	}

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("could not say hello to peers: %w", err)
	}

	return nil
}

func send(h host.Host, peer peer.ID, payload []byte) error {

	stream, err := h.NewStream(context.Background(), peer, protocol)
	if err != nil {
		return fmt.Errorf("could not create stream: %w", err)
	}
	defer stream.Close()

	_, err = stream.Write(payload)
	if err != nil {
		stream.Reset()
		return fmt.Errorf("could not write payload: %w", err)
	}

	return nil

}

func messagePrinter(stream network.Stream) {

	defer stream.Close()

	from := stream.Conn().RemotePeer()

	buf := bufio.NewReader(stream)
	msg, err := buf.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		stream.Reset()
		log.Error().Err(err).Msg("error receiving direct message")
		return
	}

	log.Trace().Str("peer", from.String()).Msg("received direct message")

	log.Info().Str("peer", from.String()).Str("payload", string(msg)).Msg("received direct message")

	fmt.Printf("%s\n", msg)
}

func formatMultiaddrs(maddrs []multiaddr.Multiaddr) []string {
	addrs := make([]string, 0, len(maddrs))
	for _, addr := range maddrs {
		addrs = append(addrs, addr.String())
	}

	return addrs
}
