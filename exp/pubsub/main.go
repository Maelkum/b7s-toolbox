package main

import (
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const (
	topic = "jibjob"
)

var (
	log = zerolog.New(os.Stderr).Level(zerolog.TraceLevel)
)

func main() {

	var (
		flagPort       uint
		flagPrivateKey string
		flagBootNodes  []string
	)

	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVarP(&flagPrivateKey, "private-key", "k", "", "path to key file")
	pflag.StringSliceVarP(&flagBootNodes, "boot-nodes", "b", nil, "boot nodes to connect to")

	pflag.Parse()

	// Init libp2p host.
	key, err := readPrivateKey(flagPrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read private key")
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%v", flagPort)),
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

	log.Info().Msg("#1 - created libp2p host")

	// Connect to boot nodes.
	err = connectToBootNodes(host, flagBootNodes)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to boot nodes")
	}

	log.Info().Msg("#2 - connected to boot nodes")

	// Remain alive.
	select {}
}
