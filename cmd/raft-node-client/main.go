package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/blocklessnetworking/b7s/host"
	"github.com/blocklessnetworking/b7s/models/blockless"

	"github.com/Maelkum/b7s-toolbox/raft/node"
)

const (
	success = 0
	failure = 1

	defaultAddress = "127.0.0.1"
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagAddress    string
		flagPort       uint
		flagExpression string
		flagParams     []string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "", "address of the API server")
	pflag.StringVarP(&flagExpression, "expression", "e", "", "expression to be evaluated")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringSliceVar(&flagParams, "param", []string{}, "parameters")

	pflag.Parse()

	log := zerolog.New(os.Stderr).
		With().Timestamp().
		Logger()

	zerolog.TimeFieldFormat = time.RFC3339Nano

	if flagExpression == "" {
		log.Error().Msg("expression cannot be empty")
		return failure
	}

	if flagAddress == "" {
		log.Error().Msg("server address cannot be empty")
		return failure
	}

	params := make(map[string]float64)
	for _, param := range flagParams {

		fields := strings.Split(param, "=")
		if len(fields) != 2 {
			log.Error().Str("param", param).Msg("parameter should be in the 'name=123.4' format")
			return failure
		}

		name := fields[0]
		val := fields[1]

		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Error().Str("value", val).Msg("value has to be a valid number")
			return failure
		}

		params[name] = value
	}

	log.Info().
		Str("expression", flagExpression).
		Interface("params", params).
		Msg("prepared expression")

	host, err := host.New(log, defaultAddress, 0)
	if err != nil {
		log.Error().Err(err).Msg("could not create libp2p host")
		return failure
	}

	hostID := host.ID().String()
	log.Info().Str("id", hostID).Msg("starting client")

	// NOTE: We're reusing the b7s host, which means we'll need to operate on that protocol.

	var wg sync.WaitGroup
	wg.Add(1)

	host.SetStreamHandler(blockless.ProtocolID, func(stream network.Stream) {
		defer stream.Close()
		defer wg.Done()

		from := stream.Conn().RemotePeer()
		log.Debug().Str("peer", from.String()).Msg("received message")

		buf := bufio.NewReader(stream)
		msg, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			log.Error().Err(err).Msg("error receiving direct message")
			return
		}

		var response node.SolveResponse
		err = json.Unmarshal(msg, &response)
		if err != nil {
			stream.Reset()
			log.Error().Err(err).Msg("could not unmarshal message")
			return
		}

		log.Info().Interface("response", response).Msg("received response")
	})

	peerID, err := connectToHost(host, flagAddress)
	if err != nil {
		log.Error().Err(err).Msg("could not connect to server")
		return failure
	}

	request := node.SolveRequest{
		Expression: flagExpression,
		Parameters: params,
	}

	err = sendMessage(host, peerID, request)
	if err != nil {
		log.Error().Err(err).Msg("could not send message to server")
		return failure
	}

	log.Info().Interface("request", request).Msg("sent request")

	wg.Wait()

	return success
}

func connectToHost(host *host.Host, address string) (peer.ID, error) {

	maddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return "", fmt.Errorf("could not parse multiaddress: %w", err)
	}

	addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return "", fmt.Errorf("could not get address info: %w", err)

	}

	err = host.Connect(context.Background(), *addrInfo)
	if err != nil {
		return "", fmt.Errorf("could connect to host: %w", err)

	}

	return addrInfo.ID, nil
}

func sendMessage(host *host.Host, id peer.ID, request node.SolveRequest) error {

	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	err = host.SendMessage(context.Background(), id, payload)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	return nil
}
