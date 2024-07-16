package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	topic    = "jibjob"
	protocol = "super-secret-test-protocol"
)

var (
	log            = zerolog.New(os.Stderr).Level(zerolog.TraceLevel)
	publishTimeout = 3 * time.Second
)

func main() {
	// libp2pmain()
	b7smain()
}
