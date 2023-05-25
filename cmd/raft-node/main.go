package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
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
		flagAddress string
		flagPort    uint
	)

	pflag.StringVarP(&flagAddress, "address", "a", "127.0.0.1", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")

	pflag.Parse()

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	log.Info().Msg("starting node")

	return success
}
