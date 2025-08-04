package spammer

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Maelkum/b7s/models/bls"
	"github.com/blessnetwork/b7s/models/response"
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

	start := time.Now()

	for profile_no, test := range getTestProfiles() {

		log().Debug("running test profile",
			"i", profile_no,
			"executions", test.executions,
			"frequency", test.frequency,
		)

		logfile := mustCreateLogFile(test)
		defer logfile.Close()

		initActiveSublogger(logfile)

		// TODO: Hardcoded limit for each individual execution.
		const testTimeout = 5 * time.Minute

		tctx, cancel := context.WithTimeout(ctx, testTimeout)
		defer cancel()
		limiter := rate.NewLimiter(rate.Limit(test.frequency), 1)

		var (
			stats = testProfileResult{
				Map: &sync.Map{},
			}
			wg     sync.WaitGroup
			fireWG sync.WaitGroup
		)

		s.libp2p.SetStreamHandler(bls.ProtocolID, responseHandler(&wg, stats.Map))

		fireWG.Add(int(test.executions))

		stats.start = time.Now()

		// Consume responses in a separate goroutine.
		for i := range test.executions {

			// Test goroutine.
			go func(ctx context.Context) {
				defer fireWG.Done()

				err := limiter.Wait(ctx)
				if err != nil {
					panic("could not wait on limiter slot: " + err.Error())
				}

				key := executionMapKey(i, test.executions, test.frequency)
				payload := getMessagePayload(key)

				// Record timestamp.
				ts := time.Now()

				stats.Store(key, runResponse{
					start: ts,
				})

				err = s.sendMessage(tctx, payload)
				if err != nil {
					// TODO: Perhaps continue with the test but mark this as a failure.
					panic("could not send message" + err.Error())
				}

				wg.Add(1)

				log().Debug("sent execution request",
					"key", key)

			}(tctx)
		}

		fireWG.Wait()

		done := make(chan struct{}, 1)

		go func() {
			log().Info("waiting for all stream handlers to be called")
			wg.Wait()
			close(done)
		}()

		log().Info("waiting on stop condition")

		select {
		case <-tctx.Done():
			log().Info("context is done")
		case <-done:
			log().Info("all stream handlers ran")
		}

		log().Info("stop condition met")

		stats.end = time.Now()

		var keys []string
		stats.Range(func(k, v any) bool {
			keys = append(keys, k.(string))
			return true
		})

		if uint(len(keys)) != test.executions {
			log().Warn("not all executions accounted for",
				"want", test.executions,
				"got", len(keys),
			)
			for i, key := range keys {
				log().Debug("key accounted for",
					"i", i,
					"key", key,
				)
			}

			return fmt.Errorf("not all execution responses collected: %w", err)
		}

		processResults(&stats, mustCreateOutputFile(test), detailedTable)

		deactivateSublogger()

		log().Info("completed test profile",
			"duration", stats.end.Sub(stats.start).String(),
		)
	}

	log().Info("completed all profiles",
		"duration", time.Since(start).String())

	return nil
}

func (s *Spammer) connect(ctx context.Context) error {

	target := s.cfg.Target
	id, err := peer.IDFromP2PAddr(target)
	if err != nil {
		return fmt.Errorf("could not connect to target: %w", err)
	}

	s.targetID = id

	log().Debug("parsed target address", "addr", target.String(), "id", id)

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

	// TODO: Figure out why stream.SetProtocol(bls.ProtocolID) does not work

	_, err = stream.Write(payload)
	if err != nil {
		stream.Reset()
		return fmt.Errorf("could not write to stream: %w", err)
	}

	return nil
}

func responseHandler(wg *sync.WaitGroup, stats *sync.Map) network.StreamHandler {

	var n atomic.Int64

	return func(stream network.Stream) {
		defer stream.Close()
		defer func() {
			wg.Done()

			updated := n.Add(1)
			if updated%100 == 0 {
				log().Debug("number of responses recorded",
					"n", updated,
				)
			}
		}()

		// Record timestamp early on.
		ts := time.Now()

		buf := bufio.NewReader(stream)
		payload, err := buf.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			log().Error("could not read response", "err", err)
			stream.Reset()
			return
		}

		var response response.Execute
		err = json.Unmarshal(payload, &response)
		if err != nil {
			log().Error("could not unmarshal response", "err", err)
			stream.Reset()
			return
		}

		if len(response.Results) == 0 {
			log().Warn("execution failed",
				"code", response.Code.String())
		}

		for peer, res := range response.Results {

			out := res.Result.Result.Stdout

			// We expect a single result.
			log().Debug("processing execution response",
				"peer", peer.String(),
				"out", out,
				"exit_code", res.Result.Result.ExitCode,
			)

			s, _ := stats.Load(out)
			stat := s.(runResponse)

			stat.end = ts
			stat.success = res.Result.Result.ExitCode == 0

			stats.Store(out, stat)

			break
		}
	}
}
