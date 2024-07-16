package main

import (
	"context"
	"fmt"

	"github.com/blocklessnetwork/b7s/host"
	"github.com/spf13/pflag"
)

func b7smain() {

	var (
		flagPort       uint
		flagPrivateKey string
		flagBootNodes  []string
	)

	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVarP(&flagPrivateKey, "private-key", "k", "", "path to key file")
	pflag.StringSliceVarP(&flagBootNodes, "boot-node", "b", nil, "boot node to connect to")

	pflag.Parse()

	if flagPrivateKey == "" {
		log.Fatal().Msg("private key is required")
	}

	addrs, err := parseMultiAddresses(flagBootNodes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse multiaddresses")
	}

	// Create host.

	host, err := host.New(log, "127.0.0.1", flagPort,
		host.WithPrivateKey(flagPrivateKey),
		host.WithBootNodes(addrs),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create host")
	}

	myAddrs := host.Addresses()
	log.Info().Strs("addrs", myAddrs).Msg("addresses I'm listening on")

	log.Info().Msg("#1 - created b7s host")

	// Connect to boot peers.

	err = host.ConnectToKnownPeers(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to known peers")
	}

	log.Info().Msg("#2 - connected to known peers")

	// Init pubsub and subscribe.

	err = host.InitPubSub(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("could not init pubsub")
	}

	topic, sub, err := host.Subscribe(topic)
	if err != nil {
		log.Fatal().Err(err).Msg("could not subscribe to topic")
	}

	log.Info().Msg("#3 - subscribed to topic")

	go publishMessagesb7s(host, topic)

	// Remain alive.
	for {

		msg, err := sub.Next(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("could not receive message")
			break
		}

		if msg.ReceivedFrom == host.ID() {
			log.Debug().Msg("skipping message from self")
			continue
		}

		log.Info().Str("peer", msg.ReceivedFrom.String()).Str("payload", string(msg.GetData())).Msg("received message")

		fmt.Printf("%s\n", string(msg.GetData()))
	}
}
