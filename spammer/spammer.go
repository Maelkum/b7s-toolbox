package spammer

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

	"github.com/Maelkum/b7s/models/bls"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/time/rate"
)

type Config struct {
	Count     uint64
	Frequency uint
	Target    multiaddr.Multiaddr
}

type Spammer struct {
	cfg    Config
	libp2p host.Host

	targetID peer.ID
	conn     network.Conn
}

func New(cfg Config, h host.Host) *Spammer {

	s := &Spammer{
		cfg:    cfg,
		libp2p: h,
	}

	return s
}

func (s *Spammer) Run(ctx context.Context) error {

	err := s.connect(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to the target: %w", err)
	}

	// TODO: Establish stat holder.

	// TODO: Prepare a function that will return input data.

	s.libp2p.SetStreamHandler(bls.ProtocolID, processResponse)

	for i, test := range getTestProfiles() {

		slog.Debug("running test profile",
			"i", i,
			"executions", test.executions,
			"frequency", test.frequency,
		)

		tctx, cancel := context.WithCancel(ctx)
		limiter := rate.NewLimiter(rate.Limit(test.frequency), 1)

		payload := getMessagePayload()

		// Consume responses in a separate goroutine.
		for range test.executions {

			// Test goroutine.
			go func(ctx context.Context) {
				err := limiter.Wait(ctx)
				if err != nil {
					cancel()
				}

				err = s.sendMessage(tctx, payload)
				if err != nil {
					// TODO: Perhaps continue with the test but mark this as a failure.
					panic("could not send message")
				}

			}(ctx)
		}
	}

	time.Sleep(5 * time.Second)

	return nil
}

func (s *Spammer) connect(ctx context.Context) error {

	target := s.cfg.Target
	id, err := peer.IDFromP2PAddr(target)
	if err != nil {
		return fmt.Errorf("could not connect to target: %w", err)
	}

	s.targetID = id

	slog.Debug("parsed target address", "addr", target.String(), "id", id)

	s.libp2p.Peerstore().AddAddr(id, target, 365*24*time.Hour)

	conn, err := s.libp2p.Network().DialPeer(ctx, id)
	if err != nil {
		return fmt.Errorf("could not dial peer: %w", err)
	}

	s.conn = conn

	return nil
}

func (s *Spammer) sendMessage(ctx context.Context, payload []byte) error {

	stream, err := s.libp2p.NewStream(ctx, s.targetID, bls.ProtocolID)
	// stream, err := s.conn.NewStream(ctx)
	if err != nil {
		return fmt.Errorf("could not open stream: %w", err)
	}
	defer func() {
		_ = stream.Close()
	}()

	// err = stream.SetProtocol(bls.ProtocolID)
	// if err != nil {
	// 	stream.Reset()
	// 	return fmt.Errorf("could not set protocol: %w", err)
	// }

	_, err = stream.Write(payload)
	if err != nil {
		stream.Reset()
		return fmt.Errorf("could not write to stream: %w", err)
	}

	return nil
}

func getMessagePayload() []byte {

	rec := Execute{
		Request: Request{
			FunctionID: testFunction.cid,
			Method:     testFunction.method,
			Arguments: []string{
				fmt.Sprint(time.Now().Unix()),
			},
			Config: RequestConfig{
				NodeCount: 1,
			},
		},
	}

	data, err := json.Marshal(rec)
	if err != nil {
		panic("could not marshal message")
	}

	data = append(data, '\n')

	log.Printf("### payload: %s", data)

	return data
}

func processResponse(stream network.Stream) {
	defer stream.Close()

	buf := bufio.NewReader(stream)
	payload, err := buf.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		stream.Reset()
		return
	}

	fmt.Printf("%s\n", payload)

	return
}
