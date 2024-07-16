package main

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// parse list of strings with multiaddresses
func parseMultiAddresses(addrs []string) ([]multiaddr.Multiaddr, error) {

	var out []multiaddr.Multiaddr
	for _, addr := range addrs {

		addr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("could not parse multiaddress (addr: %s): %w", addr, err)
		}

		out = append(out, addr)
	}

	return out, nil
}

func getPeerIDs(addrs []string) ([]peer.ID, error) {

	parsed, err := parseMultiAddresses(addrs)
	if err != nil {
		return nil, fmt.Errorf("could not parse multiaddresses: %w", err)
	}

	var ids []peer.ID
	for _, ma := range parsed {

		_, id := peer.SplitAddr(ma)
		if err := id.Validate(); err != nil {
			return nil, fmt.Errorf("invalid peer ID found: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}
