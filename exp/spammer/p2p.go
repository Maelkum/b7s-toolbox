package main

import (
	"context"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

func newP2pHost(address string, port uint, privateKey string) (host.Host, error) {

	hostAddress := fmt.Sprintf("/ip4/%v/tcp/%v", address, port)
	addresses := []string{
		hostAddress,
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(addresses...),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// Read private key, if provided.
	if privateKey != "" {
		key, err := readPrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("could not read private key: %w", err)
		}

		opts = append(opts, libp2p.Identity(key))
	}

	// Create libp2p host.
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p host: %w", err)
	}

	return h, nil
}

func subscribe(ctx context.Context, host host.Host, topic string) (*pubsub.Topic, *pubsub.Subscription, error) {

	// Get a new PubSub object with the default router.
	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create new gossipsub: %w", err)
	}

	// Join the specified topic.
	th, err := pubsub.Join(topic)
	if err != nil {
		return nil, nil, fmt.Errorf("could not join topic: %w", err)
	}

	// Subscribe to the topic.
	subscription, err := th.Subscribe()
	if err != nil {
		return nil, nil, fmt.Errorf("could not subscribe to topic: %w", err)
	}

	return th, subscription, nil
}

func readPrivateKey(filepath string) (crypto.PrivKey, error) {

	payload, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	key, err := crypto.UnmarshalPrivateKey(payload)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal private key: %w", err)
	}

	return key, nil
}
