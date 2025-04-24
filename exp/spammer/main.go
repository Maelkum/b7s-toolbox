package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

var (
	mainLogger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
)

const (
	roleWriter   = "writer"
	roleReader   = "reader"
	defaultTopic = "spammer"

	rendezvousString = "spammer-rendezvous"
)

type cliFlags struct {
	role    string
	payload string
	count   uint64
	delay   time.Duration
	topic   string

	address     string
	port        uint
	privateKey  string
	readerCount uint

	peerListFile string
}

func main() {

	var cliFlags cliFlags

	pflag.StringVarP(&cliFlags.role, "role", "r", roleReader, "role to use")
	pflag.StringVar(&cliFlags.payload, "payload", "", "payload to publish, current time by default")
	pflag.Uint64VarP(&cliFlags.count, "count", "c", 1, "how many messages to publish")
	pflag.DurationVarP(&cliFlags.delay, "delay", "d", time.Second, "delay between two messages")
	pflag.StringVarP(&cliFlags.topic, "topic", "t", defaultTopic, "topic to publish/subscribe to")

	pflag.StringVarP(&cliFlags.address, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&cliFlags.port, "port", "p", 0, "port to use")
	pflag.StringVar(&cliFlags.privateKey, "private-key", "", "private key to use")
	pflag.UintVar(&cliFlags.readerCount, "reader-count", 1, "how many readers to create")

	pflag.StringVar(&cliFlags.peerListFile, "peer-list-file", "", "path to file with peer list")

	pflag.Parse()

	switch cliFlags.role {
	case roleWriter:
		runWriter(cliFlags)
		return
	case roleReader:
		runReader(cliFlags)
		return
	default:
		mainLogger.Fatal().Msg("invalid role specified")
	}
}
