package main

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func discoverPeers(host host.Host) (chan peer.AddrInfo, error) {

	notifiee := newNotifiee()

	svc := mdns.NewMdnsService(host, rendezvousString, notifiee)
	err := svc.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start peer discovery: %w", err)
	}

	return notifiee.peers, err
}

type discoveryNotifiee struct {
	peers chan peer.AddrInfo
}

func newNotifiee() *discoveryNotifiee {

	n := discoveryNotifiee{
		peers: make(chan peer.AddrInfo),
	}

	return &n
}

func (n *discoveryNotifiee) HandlePeerFound(pi peer.AddrInfo) {
	n.peers <- pi
}
