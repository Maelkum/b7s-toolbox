package main

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/spf13/pflag"
)

func libp2pmain() {

	var (
		flagPort       uint
		flagPrivateKey string
		flagBootNodes  []string
	)

	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVarP(&flagPrivateKey, "private-key", "k", "", "path to key file")
	pflag.StringSliceVarP(&flagBootNodes, "boot-node", "b", nil, "boot node to connect to")

	pflag.Parse()

	// Init libp2p host.
	key, err := readPrivateKey(flagPrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read private key")
	}

	myAddr := fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", flagPort)
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(myAddr),
		libp2p.Identity(key),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}
	host, err := libp2p.New(opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create host")
	}

	host.SetStreamHandler(protocol, messagePrinter)

	log.Info().Msg("#1 - created libp2p host")

	fullAddr := myAddr + fmt.Sprintf("/p2p/%s", host.ID().String())
	log.Info().Str("addr", fullAddr).Msg("my address")

	listenAddrs := formatMultiaddrs(host.Addrs())
	log.Info().Strs("addrs", listenAddrs).Msg("addresses I'm listening on")

	// Connect to boot nodes.

	addrs, err := parseMultiAddresses(flagBootNodes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse multiaddresses")
	}

	err = connectToBootNodes(host, addrs)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to boot nodes")
	}

	log.Info().Msg("#2 - connected to boot nodes")

	// Say hello to peers.
	ids, err := getPeerIDs(flagBootNodes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse peer IDs")
	}

	err = sayHello(host, ids)
	if err != nil {
		log.Fatal().Err(err).Msg("could not say hello to peers")
	}

	log.Info().Msg("#3 - said hello")

	// Pubsub.
	topic, sub, err := subToTopic(host, topic)
	if err != nil {
		log.Fatal().Err(err).Msg("could not subscribe to topic")
	}

	log.Info().Msg("#4 - subscribed to topic - listening for messages")

	go publishMessages(host, topic)

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
