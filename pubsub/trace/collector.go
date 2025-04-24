package trace

import (
	"compress/gzip"
	"encoding/json"
	"fmt"

	ggio "github.com/gogo/protobuf/io"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/rs/zerolog/log"
)

type Collector struct{}

func NewTraceCollector(h host.Host) *Collector {
	h.SetStreamHandler(pubsub.RemoteTracerProtoID, messagePrinter)
	return &Collector{}
}

func (c *Collector) Flush() {}

func (c *Collector) Stop() {}

func messagePrinter(stream network.Stream) {

	defer stream.Close()

	from := stream.Conn().RemotePeer()
	log.Debug().Stringer("from", from).Msg("new stream")

	gzipR, err := gzip.NewReader(stream)
	if err != nil {
		log.Debug().Err(err).Stringer("from", from).Msg("error opening compressed stream")
		stream.Reset()
		return
	}

	reader := ggio.NewDelimitedReader(gzipR, 1<<22)
	var msg pb.TraceEventBatch

	err = reader.ReadMsg(&msg)
	if err != nil {
		log.Error().Err(err).Stringer("from", from).Msg("could not read message")
		stream.Reset()
		return
	}

	//	log.Info().Msg("received a message")

	for i, event := range msg.Batch {
		_ = i
		// 	log.Info().Int("no", i).Msg("processing event")

		ev := createEventRecord(event)
		log.Info().Any("event", ev).Msg("received event")

		payload, _ := json.Marshal(ev)
		fmt.Println(string(payload))
	}
}
