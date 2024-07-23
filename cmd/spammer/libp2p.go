package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"

	"github.com/Maelkum/b7s-toolbox/util"
	"github.com/blocklessnetwork/b7s/models/response"
)

func createLibp2pHost(keyfile string) host.Host {

	key, err := util.ReadLibp2pKey(keyfile)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read key")
	}

	addr := fmt.Sprintf("/ip4/0.0.0.0/tcp/0")

	opts := []libp2p.Option{
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(addr),
		libp2p.Identity(key),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create libp2p host")
	}

	return host
}

var pongReceiver = func(stream network.Stream) {
	defer stream.Close()

	buf := bufio.NewReader(stream)
	msg, err := buf.ReadBytes('\n')
	if err != nil {
		stream.Reset()
		log.Error().Err(err).Msg("could not read message")
		return
	}

	var pong response.Pong
	err = json.Unmarshal(msg, &pong)
	if err != nil {
		stream.Reset()
		log.Error().Err(err).Msg("could not process response")
		return
	}

	processPong(pong)
}
