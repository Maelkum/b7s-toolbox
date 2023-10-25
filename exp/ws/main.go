package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"nhooyr.io/websocket"
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
		flagConnect string
	)

	pflag.StringVarP(&flagAddress, "address", "a", "localhost:9000", "address to use")
	pflag.StringVarP(&flagConnect, "connect", "c", "", "address to connect to")

	pflag.Parse()

	log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Start server.
	if flagConnect == "" {

		log.Info().Str("address", flagAddress).Msg("serving on address")

		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			conn, err := websocket.Accept(w, req, nil)
			if err != nil {
				log.Error().Err(err).Msg("could not accept incoming connection")
				return
			}
			defer conn.CloseNow()

			for {
				_, payload, err := conn.Read(context.Background())
				if err != nil {
					log.Error().Err(err).Msg("could not read message")
					break
				}

				fmt.Printf("> %s\n", payload)
			}

			conn.Close(websocket.StatusNormalClosure, "")
		})

		err := http.ListenAndServe(flagAddress, handler)
		if err != nil {
			log.Error().Err(err).Msg("could not start server")
			return failure
		}

		return success
	}

	opts := websocket.DialOptions{}
	conn, _, err := websocket.Dial(context.Background(), flagConnect, &opts)
	if err != nil {
		log.Error().Err(err).Msg("could not dial client")
		return failure
	}

	fmt.Printf("Enter messages to send:\n")

	in := bufio.NewReader(os.Stdin)
	for {

		fmt.Printf("msg> ")
		text, _ := in.ReadString('\n')
		text = strings.TrimSpace(text)

		err = conn.Write(context.Background(), websocket.MessageText, []byte(text))
		if err != nil {
			log.Error().Err(err).Msg("could not send message")
			return failure
		}
	}

	return success
}
