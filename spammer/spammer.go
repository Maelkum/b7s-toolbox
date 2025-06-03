package spammer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Config struct {
	Count     uint64
	Frequency uint
	Target    multiaddr.Multiaddr
}

type Spammer struct {
	cfg    Config
	libp2p host.Host

	conn network.Conn
}

func New(cfg Config, h host.Host) *Spammer {

	s := &Spammer{
		cfg:    cfg,
		libp2p: h,
	}

	return s
}

func (s *Spammer) Run(ctx context.Context) error {

	err := s.connect(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to the target: %w", err)
	}

	return errors.New("TBD")
}

func (s *Spammer) connect(ctx context.Context) error {

	target := s.cfg.Target
	id, err := peer.IDFromP2PAddr(target)
	if err != nil {
		return fmt.Errorf("could not connect to target: %w", err)
	}

	slog.Debug("parsed target address", "addr", target.String(), "id", id)

	s.libp2p.Peerstore().AddAddr(id, target, 365*24*time.Hour)

	conn, err := s.libp2p.Network().DialPeer(ctx, id)
	if err != nil {
		return fmt.Errorf("could not dial peer: %w", err)
	}

	s.conn = conn

	return nil
}
