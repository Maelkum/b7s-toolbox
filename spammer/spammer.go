package spammer

import (
	"errors"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	Count     uint64
	Frequency uint
}

type Spammer struct {
	cfg    Config
	libp2p host.Host
	target multiaddr.Multiaddr
}

func New(cfg Config, h host.Host, target multiaddr.Multiaddr) *Spammer {

	s := &Spammer{
		cfg:    cfg,
		libp2p: h,
		target: target,
	}

	return s
}

func (s *Spammer) Run() error {
	return errors.New("TBD")
}
