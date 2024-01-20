package node

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetwork/b7s/models/blockless"
)

func (n *Node) Run(ctx context.Context) error {

	n.log.Info().Msg("starting node main loop")

	go n.establishConnections(ctx)

	// TODO: Just set this on node creation.
	n.host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()
		address := stream.Conn().RemoteMultiaddr()

		buf := bufio.NewReader(stream)
		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			n.log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		n.log.Debug().Str("peer", from.String()).Str("address", address.String()).Msg("received direct message")

		var request SolveRequest
		err = json.Unmarshal(msg, &request)
		if err != nil {
			stream.Reset()
			n.log.Error().Err(err).Msg("could not unmarshal request")
			return
		}

		response, err := n.SolveExpression(ctx, &request)
		if err != nil {
			n.log.Error().Err(err).Msg("could not solve expression")
			response.Error = err.Error()
		}

		payload, err := json.Marshal(response)
		if err != nil {
			stream.Reset()
			n.log.Error().Err(err).Msg("could not marshal response")
			return
		}

		err = n.host.SendMessage(context.Background(), from, payload)
		if err != nil {
			stream.Reset()
			n.log.Error().Err(err).Msg("could not write response")
			return
		}

		n.log.Info().Msg("written response")
	})

	<-ctx.Done()
	n.log.Info().Msg("stopping node main loop")
	return nil
}

func (n *Node) establishConnections(ctx context.Context) {

	count := len(n.cfg.Peers)
	connected := 0

	connectDelay := 2 * time.Second

	for ; ctx.Err() == nil; time.Sleep(connectDelay) {

		for _, p := range n.cfg.Peers {

			peerID := peer.ID(p.ID)

			if p.ID == n.host.ID().String() {
				n.log.Info().Msg("skipping connection to self")
				continue
			}

			if n.host.Host.Network().Connectedness(peerID) == network.Connected {
				n.log.Info().Str("peer", p.ID).Msg("already connected to peer, skipping")
				continue
			}

			addrInfo, err := peer.AddrInfoFromP2pAddr(p.Address)
			if err != nil {
				n.log.Warn().Err(err).Str("address", p.Address.String()).Msg("could not get addrinfo for peer - skipping")
				continue
			}

			err = n.host.Host.Connect(ctx, *addrInfo)
			if err != nil {
				n.log.Error().Err(err).Str("address", p.Address.String()).Msg("could not connect to peer")
				continue
			}

			connected++

			// Connect to all but myself.
			if connected == count-1 {
				n.log.Info().Msg("connected to all cluster peers")
				return
			}
		}
	}
}
