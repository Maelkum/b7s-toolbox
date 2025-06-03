package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Maelkum/b7s-toolbox/spammer"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/pflag"
)

type config struct {
	address string
	port    uint
	node    string
	key     string

	count     uint64
	frequency uint
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	var cfg config

	pflag.StringVarP(&cfg.address, "address", "a", "127.0.0.1", "address to listen on")
	pflag.UintVarP(&cfg.port, "port", "p", 0, "port to use")
	pflag.StringVarP(&cfg.node, "node", "n", "", "multiaddress of the node to connect to")
	pflag.StringVarP(&cfg.key, "key", "k", "", "path to the private key to use to determine identity")

	pflag.Uint64VarP(&cfg.count, "count", "c", 1, "how many requests to send")
	pflag.UintVarP(&cfg.frequency, "frequency", "f", 10, "how many requests per second should we send")

	pflag.Parse()

	// TODO: Parse multiaddress
	// TODO: Generate key if not found
	host, err := createLibp2pHost(cfg)
	if err != nil {
		return fmt.Errorf("could not create host: %w", err)
	}

	target, err := multiaddr.NewMultiaddr(cfg.node)
	if err != nil {
		return fmt.Errorf("could not parse node multiaddress: %w", err)
	}

	scfg := spammer.Config{
		Count:     cfg.count,
		Frequency: cfg.frequency,
		Target:    target,
	}

	spammer := spammer.New(scfg, host)

	err = spammer.Run(context.TODO())
	if err != nil {
		return fmt.Errorf("could not run spammer: %w", err)
	}

	return nil
}

func createLibp2pHost(cfg config) (host.Host, error) {

	key, err := readPrivateKey(cfg.key)
	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/%v/tcp/%v", cfg.address, cfg.port),
		),
		libp2p.Identity(key),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create libp2p key: %w", err)
	}

	return host, nil
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
