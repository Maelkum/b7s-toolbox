package main

import (
	"context"
	"math"
	"time"
)

func runWriter(flags cliFlags) {

	// TODO: peer list
	host, err := newP2pHost(flags.address, flags.port, flags.privateKey)
	if err != nil {
		mainLogger.Error().Err(err).Msg("could not create host")
		return
	}

	peerChan, err := discoverPeers(host)
	if err != nil {
		mainLogger.Error().Err(err).Msg("could not start peer discovery")
		return
	}

	go func() {

		for peer := range peerChan {
			mainLogger.Info().Str("peer", peer.ID.String()).Msg("discovered peer")

			err = host.Connect(context.Background(), peer)
			if err != nil {
				mainLogger.Error().Err(err).Str("peer", peer.ID.String()).Msg("could not connect to peer")
				continue
			}

			mainLogger.Info().Str("peer", peer.ID.String()).Msg("connected to peer")
		}
	}()

	hostID := host.ID().String()
	mainLogger.Info().Str("id", hostID).Msg("libp2p host created")

	topic, _, err := subscribe(context.Background(), host, flags.topic)
	if err != nil {
		mainLogger.Error().Err(err).Msg("could not subscribe to topic")
		return
	}

	limit := flags.count
	if limit == 0 {
		limit = math.MaxUint64
	}

	for i := 0; uint64(i) < limit; i++ {

		payload := flags.payload
		if payload == "" {
			payload = generatePayload()
		}

		err = topic.Publish(context.Background(), []byte(payload))
		if err != nil {
			mainLogger.Error().Err(err).Msg("could not publish to topic")
			return
		}

		mainLogger.Info().Int("i", i).Str("payload", payload).Msg("publishing message")

		time.Sleep(flags.delay)
	}
}

func generatePayload() string {
	return time.Now().Format(time.RFC3339)
}
