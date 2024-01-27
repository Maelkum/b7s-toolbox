package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/multiformats/go-multiaddr"

	"github.com/blocklessnetwork/b7s/models/blockless"
	"github.com/blocklessnetwork/b7s/models/request"
	"github.com/blocklessnetwork/b7s/models/response"
)

func runInstall(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		log.Fatalf("no function CIDs specified")
	}

	// TODO: Add a variant of log.Fatal()
	addr := cmd.Flag("address").Value.String()
	if addr == "" {
		slog.Error("node address cannot be empty")
		return
	}

	err := installFunctions(addr, args)
	if err != nil {
		slog.Error("could not install function(s)", "error", err)
		return
	}
}

func installFunctions(address string, cids []string) error {

	maddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return fmt.Errorf("could not parse multiaddress: %w", err)
	}

	slog.Info("node address", "address", maddr.String())

	host, err := libp2p.New()
	if err != nil {
		return fmt.Errorf("could not create libp2p host: %w", err)
	}
	defer func() {
		err = host.Close()
		if err != nil {
			slog.Error("could not close libp2p host", "error", err)
		}
	}()

	slog.Info("libp2p host we use", "id", host.ID().String())

	addrInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("could not get address info: %w", err)
	}

	err = host.Connect(context.Background(), *addrInfo)
	if err != nil {
		return fmt.Errorf("could not connect to node: %w", err)
	}

	slog.Info("connected to node")

	stream, err := host.NewStream(context.Background(), addrInfo.ID, blockless.ProtocolID)
	if err != nil {
		return fmt.Errorf("could not open stream to node: %w", err)
	}
	defer func() {
		err := stream.Close()
		if err != nil {
			slog.Error("could not close stream", "error", err)
		}
	}()

	responses := make(chan response.InstallFunction, len(cids))
	host.SetStreamHandler(blockless.ProtocolID, getReadResponseHandler(responses))

	for _, cid := range cids {

		slog.Info("installing function", "cid", cid)

		err = sendMessage(stream, cid)
		if err != nil {
			stream.Close()
			return fmt.Errorf("could not send install message: %w", err)
		}
	}

	installed := 0
	count := 0

outerLoop:
	for {
		select {
		case res := <-responses:
			slog.Info("received install function response", "cid", res.CID, "status", res.Message)
			if res.Message == "installed" {
				installed++
			}

			count++
			if count == len(cids) {
				break outerLoop
			}

		case <-time.After(readTimeout):
			slog.Warn("read timeout expired")
			break outerLoop
		}
	}

	if installed == len(cids) {
		slog.Info("all functions installed")
		return nil
	}

	if installed > 0 {
		slog.Warn("some functions were not installed", "requested", len(cids), "installed", installed)
		return nil
	}

	slog.Error("functions failed to install")

	return nil
}

func manifestURLFromCID(cid string) string {
	return fmt.Sprintf("https://%s.ipfs.w3s.link/manifest.json", cid)
}

func sendMessage(stream network.Stream, cid string) error {

	msg := request.InstallFunction{
		Type:        blockless.MessageInstallFunction,
		CID:         cid,
		ManifestURL: manifestURLFromCID(cid),
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}
	payload = append(payload, '\n')

	_, err = stream.Write(payload)
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}

func getReadResponseHandler(responses chan response.InstallFunction) network.StreamHandler {

	return func(stream network.Stream) {
		defer stream.Close()

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			slog.Error("could not read message", "err", err)
			return
		}

		var res response.InstallFunction
		err = json.Unmarshal(payload, &res)
		if err != nil {
			stream.Reset()
			slog.Error("could not decode message", "err", err)
			return
		}

		responses <- res
	}
}
