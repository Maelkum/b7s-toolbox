package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/Maelkum/b7s/models/bls"
	"github.com/Maelkum/b7s/models/execute"
	"github.com/Maelkum/b7s/models/request"
	"github.com/Maelkum/b7s/models/response"
)

func runExecute(cmd *cobra.Command, args []string) {

	flags := ExecuteCmdFlags

	// TODO: Add a variant of log.Fatal()
	addr := cmd.Flag("address").Value.String()
	if addr == "" {
		slog.Error("node address cannot be empty")
		return
	}

	if flags.FunctionID == "" {
		slog.Error("function CID can't be empty")
		return
	}

	if flags.Method == "" {
		slog.Error("function method can't be empty")
		return
	}

	req := request.Execute{
		Request: execute.Request{
			FunctionID: flags.FunctionID,
			Method:     flags.Method,

			Config: execute.Config{
				Permissions: flags.Permissions,
				Runtime: execute.BLSRuntimeConfig{
					Entry: flags.Entry,
				},
			},
		},
	}

	if flags.Nodes > 1 {
		req.Config.NodeCount = int(flags.Nodes)
	}

	if flags.Stdin != "" {
		req.Config.Stdin = &flags.Stdin
	}

	// Append environment variables, if any.
	if len(flags.EnvironmentVars) > 0 {
		envVars := make([]execute.EnvVar, 0, len(flags.EnvironmentVars))
		// Ensure correct format for the environment variables.
		for _, env := range flags.EnvironmentVars {
			fields := strings.Split(env, "=")
			if len(fields) != 2 {
				slog.Error("invalid env var value, expected name=value format", "env", env)
				return
			}

			ev := execute.EnvVar{
				Name:  fields[0],
				Value: fields[1],
			}
			envVars = append(envVars, ev)
		}

		req.Config.Environment = envVars
	}

	// Append parameters.
	if len(flags.Parameters) > 0 {
		params := make([]execute.Parameter, 0, len(flags.Parameters))
		for _, p := range flags.Parameters {

			param := execute.Parameter{
				Value: p,
			}

			params = append(params, param)
		}
		req.Parameters = params
	}

	res, err := executeFunction(addr, req)
	if err != nil {
		slog.Error("execution failed", "error", err)
		return
	}

	payload, err := json.Marshal(res)
	if err != nil {
		slog.Error("could not marshal payload", "error", err)
		return
	}

	fmt.Printf("%s\n", payload)
}

func executeFunction(address string, req request.Execute) (response.Execute, error) {

	// TODO: ALL of this is duplicated with the install function.

	maddr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		return response.Execute{}, fmt.Errorf("could not parse multiaddress: %w", err)
	}

	slog.Info("node address", "address", maddr.String())

	host, err := libp2p.New()
	if err != nil {
		return response.Execute{}, fmt.Errorf("could not create libp2p host: %w", err)
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
		return response.Execute{}, fmt.Errorf("could not get address info: %w", err)
	}

	err = host.Connect(context.Background(), *addrInfo)
	if err != nil {
		return response.Execute{}, fmt.Errorf("could not connect to node: %w", err)
	}

	slog.Info("connected to node")

	stream, err := host.NewStream(context.Background(), addrInfo.ID, bls.ProtocolID)
	if err != nil {
		return response.Execute{}, fmt.Errorf("could not open stream to node: %w", err)
	}
	defer func() {
		err := stream.Close()
		if err != nil {
			slog.Error("could not close stream", "error", err)
		}
	}()

	responseCh := make(chan response.Execute, 1)
	host.SetStreamHandler(bls.ProtocolID, getReadExecuteResponseHandler(responseCh))

	err = sendExecuteMessage(stream, req)
	if err != nil {
		return response.Execute{}, fmt.Errorf("could not send execution message: %w", err)
	}

	select {
	case res := <-responseCh:
		slog.Info("received execution response")
		return res, nil

	case <-time.After(readTimeout):
		break
	}

	slog.Warn("read timeout expired")
	return response.Execute{}, errors.New("no response received")
}

func getReadExecuteResponseHandler(responseCh chan response.Execute) network.StreamHandler {

	return func(stream network.Stream) {
		defer stream.Close()

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			stream.Reset()
			slog.Error("could not read message", "err", err)
			return
		}

		var res response.Execute
		err = json.Unmarshal(payload, &res)
		if err != nil {
			stream.Reset()
			slog.Error("could not decode message", "err", err)
			return
		}

		responseCh <- res
	}
}

func sendExecuteMessage(stream network.Stream, req request.Execute) error {

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	payload = append(payload, '\n')

	_, err = stream.Write(payload)
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}
