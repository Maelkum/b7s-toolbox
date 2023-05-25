package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagPort           uint
		flagConnectAddress string
		flagWebsocket      bool
		flagWebsocketPort  uint
	)

	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use for libp2p host")
	pflag.StringVarP(&flagConnectAddress, "connect-address", "c", "", "address of a libp2p host to connect to")
	pflag.BoolVarP(&flagWebsocket, "websocket", "w", false, "use websocket protocol")
	pflag.UintVar(&flagWebsocketPort, "websocket-port", 0, "port to use for libp2p host for websocket protocol")

	pflag.Parse()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	host, err := host.New(log, "127.0.0.1", flagPort, host.WithWebsocket(flagWebsocket), host.WithWebsocketPort(flagWebsocketPort))
	if err != nil {
		log.Error().Err(err).Msg("could not create libp2p host")
		return failure
	}

	log.Info().Strs("addresses", host.Addresses()).Msg("created host")

	cn := newConnectionNotifee(log)
	host.Network().Notify(cn)

	host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()

		buf := bufio.NewReader(stream)
		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			log.Error().Err(err).Msg("could not receive message")
			return
		}

		str := string(msg)
		str = strings.TrimSpace(str)

		fmt.Printf("received message: %s\n", str)
	})

	// If the connect address is empty, just print a listening message and wait.
	if flagConnectAddress == "" {
		fmt.Printf("listening on %+v...\n", host.Addresses())

		in := bufio.NewReader(os.Stdin)
		in.ReadString('\n')

		return success
	}

	// Connect to the listed host.
	addr, err := multiaddr.NewMultiaddr(flagConnectAddress)
	if err != nil {
		log.Error().Err(err).Msg("could not parse multiaddress")
		return failure
	}

	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		log.Error().Err(err).Msg("could not get address info")
		return failure
	}

	ctx := context.Background()
	err = host.Connect(ctx, *info)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to libp2p host")
		return failure
	}

	log.Debug().Str("address", flagConnectAddress).Msg("connected to host")

	fmt.Printf("Enter messages to send:\n")

	in := bufio.NewReader(os.Stdin)
	for {

		fmt.Printf("> ")
		text, _ := in.ReadString('\n')

		err = host.SendMessage(ctx, info.ID, []byte(text))
		if err != nil {
			log.Error().Err(err).Msg("could not send message")
			return failure
		}
	}

	return success
}
