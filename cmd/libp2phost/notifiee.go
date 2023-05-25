package main

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
)

type connectionNotifiee struct {
	log zerolog.Logger
}

func newConnectionNotifee(log zerolog.Logger) *connectionNotifiee {

	cn := connectionNotifiee{
		log: log.With().Str("component", "notifiee").Logger(),
	}

	return &cn
}

func (n *connectionNotifiee) Connected(network network.Network, conn network.Conn) {

	// Get peer information.
	peerID := conn.RemotePeer()
	maddr := conn.RemoteMultiaddr()
	addr := conn.LocalMultiaddr()

	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_addr", maddr.String()).
		Str("local_addr", addr.String()).
		Msg("peer connected")
}

func (n *connectionNotifiee) Disconnected(_ network.Network, conn network.Conn) {

	peerID := conn.RemotePeer()
	maddr := conn.RemoteMultiaddr()
	addr := conn.LocalMultiaddr()

	n.log.Debug().
		Str("peer", peerID.String()).
		Str("remote_addr", maddr.String()).
		Str("local_addr", addr.String()).
		Msg("peer disconnected")
}

func (n *connectionNotifiee) Listen(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}

func (n *connectionNotifiee) ListenClose(_ network.Network, _ multiaddr.Multiaddr) {
	// Noop
}
