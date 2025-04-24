package main

import (
	"context"
)

func runReader(flags cliFlags) {

	for i := 0; i < int(flags.readerCount); i++ {

		mainLogger.Info().Int("reader_no", i).Msg("creating reader thread")

		go func() {

			// TODO: peer list
			host, err := newP2pHost(flags.address, flags.port, flags.privateKey)
			if err != nil {
				mainLogger.Error().Err(err).Msg("could not create host")
				return
			}

			me := host.ID().String()

			peerChan, err := discoverPeers(host)
			if err != nil {
				mainLogger.Error().Err(err).Msg("could not start peer discovery")
				return
			}

			go func() {

				for peer := range peerChan {
					mainLogger.Trace().Str("me", me).Str("peer", peer.ID.String()).Msg("discovered peer")

					err = host.Connect(context.Background(), peer)
					if err != nil {
						mainLogger.Error().Err(err).Str("peer", peer.ID.String()).Msg("could not connect to peer")
						continue
					}

					mainLogger.Trace().Str("me", me).Str("peer", peer.ID.String()).Msg("connected to peer")
				}
			}()

			hostID := host.ID().String()
			mainLogger.Info().Str("me", me).Str("id", hostID).Msg("libp2p host created")

			_, subscription, err := subscribe(context.Background(), host, flags.topic)
			if err != nil {
				mainLogger.Error().Err(err).Msg("could not subscribe to topic")
				return
			}

			for {
				msg, err := subscription.Next(context.Background())
				if err != nil {
					mainLogger.Error().Err(err).Msg("could not receive message")
					return
				}

				peer := msg.ReceivedFrom.String()
				if peer == hostID {
					continue
				}

				mainLogger.Info().Str("me", me).Str("peer", peer).Str("payload", string(msg.Data)).Msg("new message")
			}
		}()
	}

	// Just sit there.
	select {}
}
