package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

const (
	success = 0
	failure = 1
)

var (
	log zerolog.Logger
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagAddress    string
		flagPort       uint
		flagPrivateKey string
		flagPeer       string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVar(&flagPrivateKey, "private-key", "", "private key that the libp2p host will use")
	pflag.StringVar(&flagPeer, "peer", "", "peer to chat with")

	pflag.Parse()

	log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	host, err := createHost(flagAddress, flagPort, flagPrivateKey)
	if err != nil {
		log.Error().Err(err).Msg("could not create host")
		return failure
	}

	id := host.ID()

	var addrs []string
	for _, addr := range host.Addrs() {
		full := fmt.Sprintf("%s/p2p/%s", addr.String(), id.String())
		addrs = append(addrs, full)
	}

	log.Info().Strs("addr", addrs).Msg("starting...")

	if flagPeer != "" {
		return runWriter(host, flagPeer)
	}

	host.SetStreamHandler(protocol, func(stream network.Stream) {
		defer stream.Close()

		from := stream.Conn().RemotePeer()
		log.Trace().Str("peer", from.String()).Msg("received message")

		buf := bufio.NewReader(stream)

		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			log.Error().Err(err).Msg("error receiving message")
			return
		}

		fmt.Printf("< %s\n", string(msg))
	})

	select {}

	return success
}
