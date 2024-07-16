package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/sync/errgroup"
)

func connectToBootNodes(h host.Host, addrs []string) error {

	maddrs, err := parseMultiAddresses(addrs)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse multiaddresses")
	}

	var eg errgroup.Group
	for _, ma := range maddrs {

		eg.Go(func() error {

			_, id := peer.SplitAddr(ma)
			if id.Validate() != nil {
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

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("could not connect to peers: %w", err)
	}

	return nil

}
