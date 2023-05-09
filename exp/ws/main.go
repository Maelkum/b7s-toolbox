package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"

	"github.com/lxzan/gws"
)

const (
	success = 0
	failure = 1
)

var (
	log zerolog.Logger
)

func main() {
	os.Exit(run())
}

func run() int {

	var (
		flagAddress string
		flagPort    uint
		flagConnect string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "localhost", "address to use")
	pflag.UintVarP(&flagPort, "port", "p", 0, "port to use")
	pflag.StringVarP(&flagConnect, "connect", "c", "", "address to connect to")

	pflag.Parse()

	log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	addr := fmt.Sprintf("%v:%v", flagAddress, flagPort)

	// Start server.
	if flagConnect == "" {

		srv := gws.NewServer(&Handler{}, nil)

		log.Info().Str("address", addr).Msg("starting server")

		err := srv.Run(addr)
		if err != nil {
			log.Error().Err(err).Msg("could not start websocket server")
			return failure
		}

		return success
	}

	// Start client.
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("could not resolve TCP address")
		return failure
	}

	dialer := net.Dialer{
		LocalAddr: tcpAddr,
	}

	opt := &gws.ClientOption{
		Addr:   flagConnect,
		Dialer: &dialer,
	}
	conn, _, err := gws.NewClient(&Handler{}, opt)
	if err != nil {
		log.Error().Err(err).Msg("could not create websocket client")
		return failure
	}

	fmt.Printf("Enter messages to send:\n")

	in := bufio.NewReader(os.Stdin)
	for {

		fmt.Printf("msg> ")
		text, _ := in.ReadString('\n')

		err = conn.WriteString(text)
		if err != nil {
			log.Error().Err(err).Msg("could not send message")
			return failure
		}
	}

	return success
}
