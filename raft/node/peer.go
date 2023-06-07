package node

import (
	"github.com/multiformats/go-multiaddr"
)

type Peer struct {
	ID      string
	Address multiaddr.Multiaddr
}
