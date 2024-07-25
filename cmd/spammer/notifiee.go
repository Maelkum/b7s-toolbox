package main

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
)

type notifiee struct{}

func (n *notifiee) Connected(network network.Network, conn network.Conn) {
	log.Debug().
		Stringer("peer", conn.RemotePeer()).
		Stringer("remote_address", conn.RemoteMultiaddr()).
		Stringer("local_address", conn.LocalMultiaddr()).
		Msg("peer connected")
}

func (n *notifiee) Disconnected(_ network.Network, conn network.Conn) {

	log.Debug().
		Stringer("peer", conn.RemotePeer()).
		Stringer("remote_address", conn.RemoteMultiaddr()).
		Stringer("local_address", conn.LocalMultiaddr()).
		Msg("peer disconnected")
}

func (n *notifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func (n *notifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}
